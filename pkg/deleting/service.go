/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides access to the storage
type Repository interface {
	DeleteProduct(*model.Product, string) error
	DeleteReview(*model.Review, string) error
	DeleteShop(*model.Shop, string) error
	DeleteUser(*model.User, string) error
}

// Service provides models deleting operations.
type Service interface {
	DeleteProduct(*model.Product, string) error
	DeleteReview(*model.Review, string) error
	DeleteShop(*model.Shop, string) error
	DeleteUser(*model.User, string) error
}

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// DeleteProduct takes a product from the database and permanently deletes it
func (s *service) DeleteProduct(product *model.Product, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(product, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the product")
	}
	return nil
}

// DeleteReview takes a review from the database and permanently deletes it
func (s *service) DeleteReview(review *model.Review, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(review, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the review")
	}
	return nil
}

// DeleteShop takes a shop from the database and permanently deletes it
func (s *service) DeleteShop(shop *model.Shop, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(shop, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the shop")
	}
	return nil
}

// DeleteUser takes a user from the database and permanently deletes it
func (s *service) DeleteUser(user *model.User, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(user, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the user")
	}
	return nil
}
