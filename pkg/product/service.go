package product

import (
	"context"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/review"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Service provides product operations.
type Service interface {
	Create(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Search(ctx context.Context, search string) ([]Product, error)
	Update(ctx context.Context, p *Product, id string) error
}

type service struct {
	DB *sqlx.DB
}

// NewService returns a new product service.
func NewService(db *sqlx.DB) Service {
	return &service{db}
}

// Create a product.
func (s *service) Create(ctx context.Context, p *Product) error {
	q := `INSERT INTO products 
	(id, shop_id, stock, brand, category, type, description, weight, discount, taxes, subtotal, total, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	id := token.RandString(35)
	p.CreatedAt = time.Now()

	// percentages -> numeric values
	taxes := ((p.Subtotal / 100) * p.Taxes)
	discount := ((p.Subtotal / 100) * p.Discount)

	p.Discount = discount
	p.Taxes = taxes
	p.Total = p.Subtotal + p.Taxes - p.Discount

	_, err := s.DB.ExecContext(ctx, q, id, p.ShopID, p.Stock, p.Brand, p.Category, p.Type, p.Description,
		p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		logger.Log.Errorf("failed creating product: %v", err)
		return errors.Wrap(err, "couldn't create the product")
	}

	return nil
}

// Delete permanently deletes a product from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		logger.Log.Errorf("failed deleting product: %v", err)
		return errors.Wrap(err, "couldn't delete the product")
	}

	return nil
}

// Get returns a list with all the products stored in the database.
func (s *service) Get(ctx context.Context) ([]Product, error) {
	var products []Product

	if err := s.DB.Select(&products, "SELECT * FROM products"); err != nil {
		logger.Log.Errorf("failed listing products: %v", err)
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	list, err := s.getRelationships(ctx, products)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByID retrieves the product requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (Product, error) {
	var (
		product Product
		reviews []review.Review
	)

	if err := s.DB.GetContext(ctx, &product, "SELECT * FROM products WHERE id=$1", id); err != nil {
		return Product{}, errors.Wrap(err, "couldn't find the product")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", id); err != nil {
		logger.Log.Errorf("failed listing product's reviews: %v", err)
		return Product{}, errors.Wrap(err, "couldn't find the reviews")
	}

	product.Reviews = reviews

	return product, nil
}

// Search looks for the products that contain the value specified. (Only text fields)
func (s *service) Search(ctx context.Context, search string) ([]Product, error) {
	if strings.ContainsAny(search, ";-\\|@#~€¬<>_()[]}{¡^'") {
		return nil, errors.New("invalid search")
	}

	products := []Product{}
	q := `SELECT * FROM products WHERE
	to_tsvector(id || ' ' || shop_id || ' ' || brand || ' ' || type || ' ' || category || ' ' || description)
	@@ plainto_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &products, q, search); err != nil {
		logger.Log.Errorf("failed searching products: %v", err)
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	list, err := s.getRelationships(ctx, products)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Update updates product fields.
func (s *service) Update(ctx context.Context, p *Product, id string) error {
	q := `UPDATE products SET stock=$2, brand=$3, category=$4, type=$5,
	description=$6, weight=$7, discount=$8, taxes=$9, subtotal=$10, total=$11
	WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, q, id, p.Stock, p.Brand, p.Category, p.Type,
		p.Description, p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't update the product")
	}

	return nil
}

func (s *service) getRelationships(ctx context.Context, products []Product) ([]Product, error) {
	ch, errCh := make(chan Product), make(chan error, 1)

	for _, product := range products {
		go func(product Product) {
			var reviews []review.Review

			if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
				logger.Log.Errorf("failed listing product's reviews: %v", err)
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			product.Reviews = reviews

			ch <- product
		}(product)
	}

	list := make([]Product, len(products))
	for i := range products {
		select {
		case product := <-ch:
			list[i] = product
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil

}
