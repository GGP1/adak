// Package listing includes database listing operations
package listing

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to the storage.
type Repository interface {
	GetProducts() ([]model.Product, error)
	GetProductByID(id string) (model.Product, error)

	GetReviews() ([]model.Review, error)
	GetReviewByID(id string) (model.Review, error)

	GetShops() ([]model.Shop, error)
	GetShopByID(id string) (model.Shop, error)

	GetUsers() ([]model.User, error)
	GetUserByID(id string) (model.User, error)
}

// Service provides models listing operations.
type Service interface {
	GetProducts() ([]model.Product, error)
	GetProductByID(id string) (model.Product, error)

	GetReviews() ([]model.Review, error)
	GetReviewByID(id string) (model.Review, error)

	GetShops() ([]model.Shop, error)
	GetShopByID(id string) (model.Shop, error)

	GetUsers() ([]model.User, error)
	GetUserByID(id string) (model.User, error)
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
func (s *service) GetProducts() ([]model.Product, error) {
	var (
		products []model.Product
		result   []model.Product
	)

	if err := s.DB.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, fmt.Errorf("products not found: %v", err)
	}

	for _, product := range products {
		var reviews []model.Review

		if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
			return nil, fmt.Errorf("error fetching reviews: %v", err)
		}

		product.Reviews = reviews

		result = append(result, product)
	}

	return result, nil
}

// GetProductByID lists the product requested from the database.
func (s *service) GetProductByID(id string) (model.Product, error) {
	var (
		product model.Product
		reviews []model.Review
	)

	if err := s.DB.Get(&product, "SELECT * FROM products WHERE id=$1", id); err != nil {
		return model.Product{}, fmt.Errorf("product not found: %v", err)
	}

	if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE product_id=$1", id); err != nil {
		return model.Product{}, fmt.Errorf("error fetching reviews: %v", err)
	}

	product.Reviews = reviews

	return product, nil
}

// GetReviews lists all the reviews stored in the database.
func (s *service) GetReviews() ([]model.Review, error) {
	var reviews []model.Review

	if err := s.DB.Select(&reviews, "SELECT * FROM reviews"); err != nil {
		return nil, fmt.Errorf("reviews not found: %v", err)
	}

	return reviews, nil
}

// GetReviewByID lists the review requested from the database.
func (s *service) GetReviewByID(id string) (model.Review, error) {
	var review model.Review

	if err := s.DB.Get(&review, "SELECT * FROM reviews WHERE id=$1", id); err != nil {
		return model.Review{}, fmt.Errorf("review not found: %v", err)
	}

	return review, nil
}

// GetShops lists all the shops stored in the database.
func (s *service) GetShops() ([]model.Shop, error) {
	var (
		shops  []model.Shop
		result []model.Shop
	)

	if err := s.DB.Select(&shops, "SELECT * FROM shops"); err != nil {
		return nil, fmt.Errorf("shops not found: %v", err)
	}

	for _, shop := range shops {
		var (
			location model.Location
			reviews  []model.Review
			products []model.Product
		)

		if err := s.DB.Get(&location, "SELECT * FROM locations WHERE shop_id=$1", shop.ID); err != nil {
			return nil, fmt.Errorf("location not found: %v", err)
		}

		if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE shop_id=$1", shop.ID); err != nil {
			return nil, fmt.Errorf("reviews not found: %v", err)
		}

		if err := s.DB.Select(&products, "SELECT * FROM products WHERE shop_id=$1", shop.ID); err != nil {
			return nil, fmt.Errorf("products not found: %v", err)
		}

		shop.Location = location
		shop.Reviews = reviews
		shop.Products = products

		result = append(result, shop)
	}

	return result, nil
}

// GetShopByID lists the shop requested from the database.
func (s *service) GetShopByID(id string) (model.Shop, error) {
	var (
		shop     model.Shop
		location model.Location
		reviews  []model.Review
		products []model.Product
	)

	if err := s.DB.Get(&shop, "SELECT * FROM shops WHERE id=$1", id); err != nil {
		return model.Shop{}, fmt.Errorf("shop not found: %v", err)
	}

	if err := s.DB.Get(&location, "SELECT * FROM locations WHERE shop_id=$1", id); err != nil {
		return model.Shop{}, fmt.Errorf("location not found: %v", err)
	}

	if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE shop_id=$1", id); err != nil {
		return model.Shop{}, fmt.Errorf("reviews not found: %v", err)
	}

	if err := s.DB.Select(&products, "SELECT * FROM products WHERE shop_id=$1", id); err != nil {
		return model.Shop{}, fmt.Errorf("products not found: %v", err)
	}

	shop.Location = location
	shop.Reviews = reviews
	shop.Products = products

	return shop, nil
}

// GetUsers lists all the users stored in the database.
func (s *service) GetUsers() ([]model.User, error) {
	var (
		users  []model.User
		result []model.User
	)

	if err := s.DB.Select(&users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("users not found: %v", err)
	}

	for _, user := range users {
		var reviews []model.Review

		orders, err := ordering.GetByUserID(s.DB, user.ID)
		if err != nil {
			return nil, err
		}

		if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
			return nil, fmt.Errorf("error fetching reviews: %v", err)
		}

		user.Orders = orders
		user.Reviews = reviews

		result = append(result, user)
	}

	return result, nil
}

// GetUserByID lists the user requested from the database.
func (s *service) GetUserByID(id string) (model.User, error) {
	var (
		user    model.User
		reviews []model.Review
	)

	if err := s.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return model.User{}, fmt.Errorf("user not found: %v", err)
	}

	if err := s.DB.Select(&reviews, "SELECT * FROM reviews WHERE user_id=$1", id); err != nil {
		return model.User{}, fmt.Errorf("error fetching reviews: %v", err)
	}

	orders, err := ordering.GetByUserID(s.DB, id)
	if err != nil {
		return model.User{}, err
	}

	user.Orders = orders

	return user, nil
}
