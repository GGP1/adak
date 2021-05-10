package product

import (
	"context"

	"github.com/GGP1/adak/pkg/review"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Service provides product operations.
type Service interface {
	Create(ctx context.Context, p Product) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Search(ctx context.Context, query string) ([]Product, error)
	Update(ctx context.Context, id string, p UpdateProduct) error
}

type service struct {
	db      *sqlx.DB
	mc      *memcache.Client
	metrics metrics
}

// NewService returns a new product service.
func NewService(db *sqlx.DB, mc *memcache.Client) Service {
	return &service{db, mc, initMetrics()}
}

// Create a product.
func (s *service) Create(ctx context.Context, p Product) error {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Create"}).Inc()

	q := `INSERT INTO products 
	(id, shop_id, stock, brand, category, type, description, 
	weight, discount, taxes, subtotal, total, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := s.db.ExecContext(ctx, q, p.ID, p.ShopID, p.Stock, p.Brand,
		p.Category, p.Type, p.Description, p.Weight, p.Discount, p.Taxes,
		p.Subtotal, p.Total, p.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the product")
	}

	s.metrics.totalProducts.Inc()
	return nil
}

// Delete permanently deletes a product from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Delete"}).Inc()
	_, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete product from the database")
	}
	s.metrics.totalProducts.Dec()

	if err := s.mc.Delete(id); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "couldn't delete product from cache")
	}

	return nil
}

// Get returns a list with all the products stored in the database.
func (s *service) Get(ctx context.Context) ([]Product, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Get"}).Inc()

	q := `SELECT p.*, r.*
	FROM products AS p
	LEFT JOIN reviews AS r ON p.id=r.product_id`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		p := Product{}
		r := review.Review{}
		err := rows.Scan(
			&p.ID, &p.ShopID, &p.Stock, &p.Brand, &p.Category, &p.Type,
			&p.Description, &p.Weight, &p.Discount, &p.Taxes, &p.Subtotal,
			&p.Total, &p.CreatedAt, &p.UpdatedAt,
			&r.ID, &r.Stars, &r.Comment, &r.UserID, &r.ProductID, &r.ShopID,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't scan product")
		}

		p.Reviews = append(p.Reviews, r)
		products = append(products, p)
	}

	return products, nil
}

// GetByID retrieves the product requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (Product, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "GetByID"}).Inc()

	q := `SELECT p.*, r.*
	FROM products p
	LEFT JOIN reviews r ON p.id=r.product_id
	WHERE p.id=$1`
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return Product{}, errors.Wrap(err, "couldn't find the product")
	}
	defer rows.Close()

	var p Product
	for rows.Next() {
		r := review.Review{}
		err := rows.Scan(
			&p.ID, &p.ShopID, &p.Stock, &p.Brand, &p.Category, &p.Type,
			&p.Description, &p.Weight, &p.Discount, &p.Taxes, &p.Subtotal,
			&p.Total, &p.CreatedAt, &p.UpdatedAt,
			&r.ID, &r.Stars, &r.Comment, &r.UserID, &r.ProductID, &r.ShopID,
			&r.CreatedAt,
		)
		if err != nil {
			return Product{}, errors.Wrap(err, "couldn't scan product")
		}

		p.Reviews = append(p.Reviews, r)
	}

	return p, nil
}

// Search looks for the products that contain the value specified. (Only text fields)
func (s *service) Search(ctx context.Context, query string) ([]Product, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Search"}).Inc()

	var products []Product
	q := `SELECT * FROM products WHERE
	to_tsvector(id || ' ' ||  shop_id || ' ' || brand || ' ' || type || ' ' || category)
	@@ plainto_tsquery($1)`
	if err := s.db.SelectContext(ctx, &products, q, query); err != nil {
		return nil, errors.Wrap(err, "couldn't find products")
	}

	return products, nil
}

// Update updates product fields.
func (s *service) Update(ctx context.Context, id string, p UpdateProduct) error {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Update"}).Inc()

	q := `UPDATE products SET stock=$2, brand=$3, category=$4, type=$5,
	description=$6, weight=$7, discount=$8, taxes=$9, subtotal=$10, total=$11
	WHERE id=$1`
	_, err := s.db.ExecContext(ctx, q, id, p.Stock, p.Brand, p.Category, p.Type,
		p.Description, p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't update the product")
	}

	if err := s.mc.Delete(id); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "couldn't delete product from cache")
	}

	return nil
}
