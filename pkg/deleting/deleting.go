// Package deleting includes database deleting operations.
package deleting

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage.
type Repository interface {
	DeleteProduct(db *gorm.DB, id string) error
	DeleteReview(db *gorm.DB, id string) error
	DeleteShop(db *gorm.DB, id string) error
	DeleteUser(db *gorm.DB, id string) error
}

// Service provides models deleting operations.
type Service interface {
	DeleteProduct(db *gorm.DB, id string) error
	DeleteReview(db *gorm.DB, id string) error
	DeleteShop(db *gorm.DB, id string) error
	DeleteUser(db *gorm.DB, id string) error
}

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// DeleteProduct takes a product from the database and permanently deletes it.
func (s *service) DeleteProduct(db *gorm.DB, id string) error {
	var product model.Product

	if err := db.Delete(product, id).Error; err != nil {
		return fmt.Errorf("couldn't delete the product")
	}
	return nil
}

// DeleteReview takes a review from the database and permanently deletes it.
func (s *service) DeleteReview(db *gorm.DB, id string) error {
	var review model.Review

	if err := db.Delete(review, id).Error; err != nil {
		return fmt.Errorf("couldn't delete the review")
	}
	return nil
}

// DeleteShop takes a shop from the database and permanently deletes it.
func (s *service) DeleteShop(db *gorm.DB, id string) error {
	var shop model.Shop

	if err := db.Delete(shop, id).Error; err != nil {
		return fmt.Errorf("couldn't delete the shop")
	}
	return nil
}

// DeleteUser takes a user from the database and permanently deletes it.
func (s *service) DeleteUser(db *gorm.DB, id string) error {
	var user model.User
	if err := db.Delete(user, id).Error; err != nil {
		return fmt.Errorf("couldn't delete the user")
	}

	return nil
}
