/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides access to the storage
type Repository interface {
	DeleteProduct(db *gorm.DB, product *model.Product, id string) error
	DeleteReview(db *gorm.DB, review *model.Review, id string) error
	DeleteShop(db *gorm.DB, shop *model.Shop, id string) error
	DeleteUser(db *gorm.DB, user *model.User, id string) error
}

// Service provides models deleting operations
type Service interface {
	DeleteProduct(db *gorm.DB, product *model.Product, id string) error
	DeleteReview(db *gorm.DB, review *model.Review, id string) error
	DeleteShop(db *gorm.DB, shop *model.Shop, id string) error
	DeleteUser(db *gorm.DB, user *model.User, id string) error
}

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// DeleteProduct takes a product from the database and permanently deletes it
func (s *service) DeleteProduct(db *gorm.DB, product *model.Product, id string) error {
	if err := db.Delete(product, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the product")
	}
	return nil
}

// DeleteReview takes a review from the database and permanently deletes it
func (s *service) DeleteReview(db *gorm.DB, review *model.Review, id string) error {
	if err := db.Delete(review, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the review")
	}
	return nil
}

// DeleteShop takes a shop from the database and permanently deletes it
func (s *service) DeleteShop(db *gorm.DB, shop *model.Shop, id string) error {
	if err := db.Delete(shop, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the shop")
	}
	return nil
}

// DeleteUser takes a user from the database and permanently deletes it
func (s *service) DeleteUser(db *gorm.DB, user *model.User, id string) error {
	if err := db.Delete(user, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the user")
	}
	return nil
}
