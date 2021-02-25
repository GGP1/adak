package shop

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Shops implements the shops interface.
type Shops struct {
	db *sqlx.DB

	productsClient product.ProductsClient
	reviewsClient  review.ReviewsClient
}

// NewService returns a new shops server.
func NewService(db *sqlx.DB, productsConn, reviewsConn *grpc.ClientConn) *Shops {
	return &Shops{
		db:             db,
		productsClient: product.NewProductsClient(productsConn),
		reviewsClient:  review.NewReviewsClient(reviewsConn),
	}
}

// Run starts the server.
func (s *Shops) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errors.Wrapf(err, "shops: failed listening on port %d", port)
	}

	srv := grpc.NewServer()
	RegisterShopsServer(srv, s)

	return srv.Serve(lis)
}

// Create creates a shop.
func (s *Shops) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	sQuery := `INSERT INTO shops
	(id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4)`

	lQuery := `INSERT INTO locations
	(shop_id, country, state, zip_code, city, address)
	VALUES ($1, $2, $3, $4, $5, $6)`

	id := token.GenerateRunes(30)
	req.Shop.CreatedAt = timestamppb.Now()

	_, err := s.db.ExecContext(ctx, sQuery, id, req.Shop.Name, req.Shop.CreatedAt, req.Shop.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the shop")
	}

	_, err = s.db.ExecContext(ctx, lQuery, id, req.Shop.Location.Country, req.Shop.Location.State,
		req.Shop.Location.ZipCode, req.Shop.Location.City, req.Shop.Location.Address)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the location")
	}

	return nil, nil
}

// Delete permanently deletes a shop from the database.
func (s *Shops) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	_, err := s.db.ExecContext(ctx, "DELETE FROM shops WHERE id=$1", req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete the shop")
	}

	return nil, nil
}

// Get returns a list with all the shops stored in the database.
func (s *Shops) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	var shops []*Shop

	if err := s.db.SelectContext(ctx, &shops, "SELECT * FROM shops"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the shops")
	}

	list, err := getRelationships(ctx, s.db, shops)
	if err != nil {
		return nil, err
	}

	return &GetResponse{Shops: list}, nil
}

// GetByID retrieves the shop requested from the database.
func (s *Shops) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	var (
		shop     *Shop
		location *Location
		reviews  []*review.Review
		products []*product.Product
	)

	if err := s.db.GetContext(ctx, &shop, "SELECT * FROM shops WHERE id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the shop")
	}

	if err := s.db.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the location")
	}

	if err := s.db.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE shop_id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the reviews")
	}

	if err := s.db.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	shop.Location = location
	shop.Reviews = reviews
	shop.Products = products

	return &GetByIDResponse{Shop: shop}, nil
}

// Search looks for the shops that contain the value specified. (Only text fields)
func (s *Shops) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	var shops []*Shop

	q := `SELECT * FROM shops WHERE
	to_tsvector(id || ' ' || name) @@ to_tsquery($1)`

	if strings.ContainsAny(req.Search, ";-\\|@#~€¬<>_()[]}{¡'") {
		return nil, errors.New("invalid search")
	}

	if err := s.db.SelectContext(ctx, &shops, q, req.Search); err != nil {
		return nil, errors.Wrap(err, "couldn't find shops")
	}

	list, err := getRelationships(ctx, s.db, shops)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{Shops: list}, nil
}

// Update updates shop fields.
func (s *Shops) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	q := `UPDATE shops SET name=$2, country=$3, city=$4, address=$5
	WHERE id=$1`

	_, err := s.db.ExecContext(ctx, q, req.ID, req.Shop.Name, req.Shop.Location.Country,
		req.Shop.Location.City, req.Shop.Location.Address)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the shop")
	}

	return nil, nil
}

func getRelationships(ctx context.Context, db *sqlx.DB, shops []*Shop) ([]*Shop, error) {
	var list []*Shop

	ch, errCh := make(chan *Shop), make(chan error, 1)

	for _, shop := range shops {
		go func(shop *Shop) {
			var (
				location *Location
				reviews  []*review.Review
				products []*product.Product
			)

			if err := db.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", shop.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the location")
			}

			if err := db.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE shop_id=$1", shop.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			if err := db.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id=$1", shop.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the products")
			}

			shop.Location = location
			shop.Products = products
			shop.Reviews = reviews

			ch <- shop
		}(shop)
	}

	for i := 0; i < len(shops); i++ {
		select {
		case shop := <-ch:
			list = append(list, shop)
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil
}
