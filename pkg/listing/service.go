/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides access to the storage
type Repository interface {
	GetProducts(*[]model.Product) error
	GetProductByID(*model.Product, string) error

	GetReviews(*[]model.Review) error
	GetReviewByID(*model.Review, string) error

	GetShops(*[]model.Shop) error
	GetShopByID(*model.Shop, string) error

	GetUsers(*[]model.User) error
	GetUserByID(*model.User, string) error
}

// Service provides models listing operations.
type Service interface {
	GetProducts(*[]model.Product) error
	GetProductByID(*model.Product, string) error

	GetReviews(*[]model.Review) error
	GetReviewByID(*model.Review, string) error

	GetShops(*[]model.Shop) error
	GetShopByID(*model.Shop, string) error

	GetUsers(*[]model.User) error
	GetUserByID(*model.User, string) error
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetProducts lists all the products stored in the database
func (s *service) GetProducts(product *[]model.Product) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").Find(product).Error
	if err != nil {
		return errors.Wrap(err, "error: products not found")
	}
	return nil
}

// GetProductByID lists the product requested from the database
func (s *service) GetProductByID(product *model.Product, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").First(product, id).Error
	if err != nil {
		return errors.Wrap(err, "error: product not found")
	}
	return nil
}

// GetReviews lists all the reviews stored in the database
func (s *service) GetReviews(review *[]model.Review) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Find(review).Error
	if err != nil {
		return errors.Wrap(err, "error: reviews not found")
	}
	return nil
}

// GetReviewByID lists the review requested from the database
func (s *service) GetReviewByID(review *model.Review, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.First(review, id).Error
	if err != nil {
		return errors.Wrap(err, "error: review not found")
	}
	return nil
}

// GetShops lists all the shops stored in the database
func (s *service) GetShops(shop *[]model.Shop) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Location").Preload("Reviews").Preload("Products").Find(shop).Error
	if err != nil {
		return errors.Wrap(err, "error: shops not found")
	}
	return nil
}

// GetShopByID lists the shop requested from the database
func (s *service) GetShopByID(shop *model.Shop, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Location").Preload("Reviews").Preload("Products").First(shop, id).Error
	if err != nil {
		return errors.Wrap(err, "error: shop not found")
	}
	return nil
}

// GetUsers lists all the users stored in the database
func (s *service) GetUsers(user *[]model.User) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").Find(user).Error
	if err != nil {
		return errors.Wrap(err, "error: users not found")
	}
	return nil
}

// GetUserByID lists the user requested from the database
func (s *service) GetUserByID(user *model.User, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").First(user, id).Error
	if err != nil {
		return errors.Wrap(err, "error: user not found")
	}
	return nil
}
