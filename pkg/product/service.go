package product

import (
	"context"
	"time"

	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/review"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	Create(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Search(ctx context.Context, search string) ([]Product, error)
	Update(ctx context.Context, p *Product, id string) error
}

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
	r  Repository
	DB *sqlx.DB
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// CreateProduct creates a product.
func (s *service) Create(ctx context.Context, p *Product) error {
	q := `INSERT INTO products 
	(id, shop_id, stock, brand, category, type, description, weight, discount, taxes, subtotal, total, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	if err := p.Validate(); err != nil {
		return err
	}

	id := token.GenerateRunes(35)
	p.CreatedAt = time.Now()

	taxes := ((p.Subtotal / 100) * p.Taxes)
	discount := ((p.Subtotal / 100) * p.Discount)
	p.Total = p.Subtotal + taxes - discount

	_, err := s.DB.ExecContext(ctx, q, id, p.ShopID, p.Stock, p.Brand, p.Category, p.Type, p.Description,
		p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the product")
	}

	return nil
}

// Delete permanently deletes a product from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the product")
	}

	return nil
}

// Get returns a list with all the products stored in the database.
func (s *service) Get(ctx context.Context) ([]Product, error) {
	var (
		products []Product
		list     []Product
	)

	ch := make(chan Product)
	errCh := make(chan error)

	if err := s.DB.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "products not found")
	}

	for _, product := range products {
		go func(product Product) {
			var reviews []review.Review

			if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
				errCh <- errors.Wrap(err, "error fetching reviews")
			}

			product.Reviews = reviews

			ch <- product
		}(product)
	}

	select {
	case p := <-ch:
		list = append(list, p)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
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
		return Product{}, errors.Wrap(err, "product not found")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", id); err != nil {
		return Product{}, errors.Wrap(err, "error fetching reviews")
	}

	product.Reviews = reviews

	return product, nil
}

// Search looks for the products that contain the value specified. (Only text fields)
func (s *service) Search(ctx context.Context, search string) ([]Product, error) {
	var (
		products []Product
		list     []Product
	)

	ch := make(chan Product)
	errCh := make(chan error)

	q := `SELECT * FROM products WHERE
	to_tsvector(id || ' ' || shop_id || ' ' || brand || ' ' || type || ' ' || category || ' ' || description)
	@@ to_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &products, q, search); err != nil {
		return nil, errors.Wrap(err, "couldn't find products")
	}

	for _, product := range products {
		go func(product Product) {
			var reviews []review.Review

			if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
				errCh <- errors.Wrap(err, "error fetching reviews")
			}

			product.Reviews = reviews

			ch <- product
		}(product)
	}

	select {
	case p := <-ch:
		list = append(list, p)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
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
