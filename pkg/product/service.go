package product

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/review"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Products implements the product service.
type Products struct {
	db *sqlx.DB

	reviewClient review.ReviewsClient
}

// NewService returns a new products server.
func NewService(db *sqlx.DB, reviewConn *grpc.ClientConn) *Products {
	return &Products{
		db:           db,
		reviewClient: review.NewReviewsClient(reviewConn),
	}
}

// Run starts the server.
func (p *Products) Run(port int) error {
	srv := grpc.NewServer()
	RegisterProductsServer(srv, p)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "products: failed listening on port %d", port)
	}

	return srv.Serve(lis)
}

// Create creates a product.
func (p *Products) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	q := `INSERT INTO products 
	(id, shop_id, stock, brand, category, type, description, weight, discount, taxes, subtotal, total, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	id := token.GenerateRunes(35)
	req.Product.CreatedAt.Seconds = time.Now().Unix()

	// percentages -> numeric values
	taxes := ((req.Product.Subtotal / 100) * req.Product.Taxes)
	discount := ((req.Product.Subtotal / 100) * req.Product.Discount)

	req.Product.Discount = discount
	req.Product.Taxes = taxes
	req.Product.Total = req.Product.Subtotal + req.Product.Taxes - req.Product.Discount

	_, err := p.db.ExecContext(ctx, q, id, req.Product.ShopID, req.Product.Stock, req.Product.Brand,
		req.Product.Category, req.Product.Type, req.Product.Description, req.Product.Weight, req.Product.Discount,
		req.Product.Taxes, req.Product.Subtotal, req.Product.Total, req.Product.CreatedAt, req.Product.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the product")
	}

	return nil, nil
}

// Delete permanently deletes a product from the database.
func (p *Products) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	_, err := p.db.ExecContext(ctx, "DELETE FROM products WHERE id=$1", req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete the product")
	}

	return nil, nil
}

// Get returns a list with all the products stored in the database.
func (p *Products) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	var products []*Product

	if err := p.db.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	list, err := getRelationships(ctx, p.db, products)
	if err != nil {
		return nil, err
	}

	return &GetResponse{Products: list}, nil
}

// GetByID retrieves the product requested from the database.
func (p *Products) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	var (
		product *Product
		reviews []*review.Review
	)

	if err := p.db.GetContext(ctx, &product, "SELECT * FROM products WHERE id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the product")
	}

	if err := p.db.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the reviews")
	}

	product.Reviews = reviews

	return &GetByIDResponse{Product: product}, nil
}

// Search looks for the products that contain the value specified. (Only text fields)
func (p *Products) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	var products []*Product

	q := `SELECT * FROM products WHERE
	to_tsvector(id || ' ' || shop_id || ' ' || brand || ' ' || type || ' ' || category || ' ' || description)
	@@ to_tsquery($1)`

	if strings.ContainsAny(req.Search, ";-\\|@#~€¬<>_()[]}{¡'") {
		return nil, errors.New("invalid search")
	}

	if err := p.db.SelectContext(ctx, &products, q, req.Search); err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	list, err := getRelationships(ctx, p.db, products)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{Products: list}, nil
}

// Update updates product fields.
func (p *Products) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	q := `UPDATE products SET stock=$2, brand=$3, category=$4, type=$5,
	description=$6, weight=$7, discount=$8, taxes=$9, subtotal=$10, total=$11
	WHERE id=$1`

	_, err := p.db.ExecContext(ctx, q, req.ID, req.Product.Stock, req.Product.Brand, req.Product.Category, req.Product.Type,
		req.Product.Description, req.Product.Weight, req.Product.Discount, req.Product.Taxes, req.Product.Subtotal, req.Product.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the product")
	}

	return nil, nil
}

func getRelationships(ctx context.Context, db *sqlx.DB, products []*Product) ([]*Product, error) {
	var list []*Product

	ch, errCh := make(chan *Product), make(chan error, 1)

	for _, product := range products {
		go func(product *Product) {
			var reviews []*review.Review

			if err := db.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			product.Reviews = reviews

			ch <- product
		}(product)
	}

	for i := 0; i < len(products); i++ {
		select {
		case p := <-ch:
			list = append(list, p)
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil

}
