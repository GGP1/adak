// Package searching provides a service for searching specific information
// in the database related to the core api models.
package searching

import (
	"context"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	SearchProducts(ctx context.Context, search string) ([]model.Product, error)
	SearchShops(ctx context.Context, search string) ([]model.Shop, error)
	SearchUsers(ctx context.Context, search string) ([]model.User, error)
}

// Service provides models searching operations.
type Service interface {
	SearchProducts(ctx context.Context, search string) ([]model.Product, error)
	SearchShops(ctx context.Context, search string) ([]model.Shop, error)
	SearchUsers(ctx context.Context, search string) ([]model.User, error)
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a searching service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// SearchProducts looks for the products that contain the value specified. (Only text fields)
func (s *service) SearchProducts(ctx context.Context, search string) ([]model.Product, error) {
	var (
		products []model.Product
		list     []model.Product
	)

	ch := make(chan model.Product)
	errCh := make(chan error)

	q := `SELECT * FROM products WHERE
	to_tsvector(id || ' ' || shop_id || ' ' || brand || ' ' || type || ' ' || category || ' ' || description)
	@@ to_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &products, q, search); err != nil {
		return nil, errors.Wrap(err, "couldn't find products")
	}

	for _, product := range products {
		go func(product model.Product) {
			var reviews []model.Review

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

// SearchShops looks for the shops that contain the value specified. (Only text fields)
func (s *service) SearchShops(ctx context.Context, search string) ([]model.Shop, error) {
	var (
		shops []model.Shop
		list  []model.Shop
	)

	ch := make(chan model.Shop)
	errCh := make(chan error)

	q := `SELECT * FROM shops WHERE
	to_tsvector(id || ' ' || name) @@ to_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &shops, q, search); err != nil {
		return nil, errors.Wrap(err, "couldn't find shops")
	}

	for _, shop := range shops {
		go func(shop model.Shop) {
			var (
				location model.Location
				reviews  []model.Review
				products []model.Product
			)

			if err := s.DB.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", shop.ID); err != nil {
				errCh <- errors.Wrap(err, "location not found")
			}

			if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE shop_id=$1", shop.ID); err != nil {
				errCh <- errors.Wrap(err, "reviews not found")
			}

			if err := s.DB.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id=$1", shop.ID); err != nil {
				errCh <- errors.Wrap(err, "products not found")
			}

			shop.Location = location
			shop.Reviews = reviews
			shop.Products = products

			ch <- shop
		}(shop)
	}

	select {
	case s := <-ch:
		list = append(list, s)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return list, nil
}

// SearchUsers looks for the users that contain the value specified. (Only text fields)
func (s *service) SearchUsers(ctx context.Context, search string) ([]model.User, error) {
	var (
		users []model.User
		list  []model.User
	)

	ch := make(chan model.User)
	errCh := make(chan error)

	q := `SELECT * FROM users WHERE
	to_tsvector(id || ' ' || username || ' ' || email) 
	@@ to_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &users, q, search); err != nil {
		return nil, errors.Wrap(err, "couldn't find users")
	}

	for _, user := range users {
		go func(user model.User) {
			var reviews []model.Review

			orders, err := ordering.GetByUserID(ctx, s.DB, user.ID)
			if err != nil {
				errCh <- err
			}

			if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
				errCh <- errors.Wrap(err, "error fetching reviews")
			}

			user.Orders = orders
			user.Reviews = reviews

			ch <- user
		}(user)
	}
	select {
	case u := <-ch:
		list = append(list, u)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return list, nil
}
