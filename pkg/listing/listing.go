// Package listing includes database listing operations
package listing

import (
	"context"
	"time"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	GetProducts(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, id string) (model.Product, error)

	GetReviews(ctx context.Context) ([]model.Review, error)
	GetReviewByID(ctx context.Context, id string) (model.Review, error)

	GetShops(ctx context.Context) ([]model.Shop, error)
	GetShopByID(ctx context.Context, id string) (model.Shop, error)

	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
}

// Service provides models listing operations.
type Service interface {
	GetProducts(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, id string) (model.Product, error)

	GetReviews(ctx context.Context) ([]model.Review, error)
	GetReviewByID(ctx context.Context, id string) (model.Review, error)

	GetShops(ctx context.Context) ([]model.Shop, error)
	GetShopByID(ctx context.Context, id string) (model.Shop, error)

	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a listing service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// GetProducts lists all the products stored in the database.
func (s *service) GetProducts(ctx context.Context) ([]model.Product, error) {
	var (
		products []model.Product
		result   []model.Product
	)

	if err := s.DB.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "products not found")
	}

	for _, product := range products {
		var reviews []model.Review

		if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
			return nil, errors.Wrap(err, "error fetching reviews")
		}

		product.Reviews = reviews

		result = append(result, product)
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetProductByID lists the product requested from the database.
func (s *service) GetProductByID(ctx context.Context, id string) (model.Product, error) {
	var (
		product model.Product
		reviews []model.Review
	)

	if err := s.DB.GetContext(ctx, &product, "SELECT * FROM products WHERE id=$1", id); err != nil {
		return model.Product{}, errors.Wrap(err, "product not found")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", id); err != nil {
		return model.Product{}, errors.Wrap(err, "error fetching reviews")
	}

	product.Reviews = reviews

	select {
	case <-time.After(0 * time.Nanosecond):
		return product, nil
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	}
}

// GetReviews lists all the reviews stored in the database.
func (s *service) GetReviews(ctx context.Context) ([]model.Review, error) {
	var reviews []model.Review

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews"); err != nil {
		return nil, errors.Wrap(err, "reviews not found")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return reviews, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetReviewByID lists the review requested from the database.
func (s *service) GetReviewByID(ctx context.Context, id string) (model.Review, error) {
	var review model.Review

	if err := s.DB.GetContext(ctx, &review, "SELECT * FROM reviews WHERE id=$1", id); err != nil {
		return model.Review{}, errors.Wrap(err, "review not found")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return review, nil
	case <-ctx.Done():
		return model.Review{}, ctx.Err()
	}
}

// GetShops lists all the shops stored in the database.
func (s *service) GetShops(ctx context.Context) ([]model.Shop, error) {
	var (
		shops  []model.Shop
		result []model.Shop
	)

	if err := s.DB.SelectContext(ctx, &shops, "SELECT * FROM shops"); err != nil {
		return nil, errors.Wrap(err, "shops not found")
	}

	for _, shop := range shops {
		var (
			location model.Location
			reviews  []model.Review
			products []model.Product
		)

		if err := s.DB.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", shop.ID); err != nil {
			return nil, errors.Wrap(err, "location not found")
		}

		if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE shop_id=$1", shop.ID); err != nil {
			return nil, errors.Wrap(err, "reviews not found")
		}

		if err := s.DB.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id=$1", shop.ID); err != nil {
			return nil, errors.Wrap(err, "products not found")
		}

		shop.Location = location
		shop.Reviews = reviews
		shop.Products = products

		result = append(result, shop)
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetShopByID lists the shop requested from the database.
func (s *service) GetShopByID(ctx context.Context, id string) (model.Shop, error) {
	var (
		shop     model.Shop
		location model.Location
		reviews  []model.Review
		products []model.Product
	)

	if err := s.DB.GetContext(ctx, &shop, "SELECT * FROM shops WHERE id=$1", id); err != nil {
		return model.Shop{}, errors.Wrap(err, "shop not found")
	}

	if err := s.DB.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", id); err != nil {
		return model.Shop{}, errors.Wrap(err, "location not found")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE shop_id=$1", id); err != nil {
		return model.Shop{}, errors.Wrap(err, "reviews not found")
	}

	if err := s.DB.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id=$1", id); err != nil {
		return model.Shop{}, errors.Wrap(err, "products not found")
	}

	shop.Location = location
	shop.Reviews = reviews
	shop.Products = products

	select {
	case <-time.After(0 * time.Nanosecond):
		return shop, nil
	case <-ctx.Done():
		return model.Shop{}, ctx.Err()
	}
}

// GetUsers lists all the users stored in the database.
func (s *service) GetUsers(ctx context.Context) ([]model.User, error) {
	var (
		users  []model.User
		result []model.User
	)

	if err := s.DB.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, errors.Wrap(err, "users not found")
	}

	for _, user := range users {
		var reviews []model.Review

		orders, err := ordering.GetByUserID(ctx, s.DB, user.ID)
		if err != nil {
			return nil, err
		}

		if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
			return nil, errors.Wrap(err, "error fetching reviews")
		}

		user.Orders = orders
		user.Reviews = reviews

		result = append(result, user)
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetUserByID lists the user requested from the database.
func (s *service) GetUserByID(ctx context.Context, id string) (model.User, error) {
	var (
		user    model.User
		reviews []model.Review
	)

	if err := s.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return model.User{}, errors.Wrap(err, "user not found")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", id); err != nil {
		return model.User{}, errors.Wrap(err, "error fetching reviews")
	}

	orders, err := ordering.GetByUserID(ctx, s.DB, id)
	if err != nil {
		return model.User{}, err
	}

	user.Orders = orders

	select {
	case <-time.After(0 * time.Nanosecond):
		return user, nil
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	}
}
