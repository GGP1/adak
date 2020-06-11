/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
)

// Service provides models deleting operations.
type Service interface {
	DeleteUser() error
	DeleteProduct() error
	DeleteReview() error
	DeleteShop() error
}

// DeleteUser deletes a user, returns an error
func DeleteUser(user *model.User, id string, db *gorm.DB) error {
	db.Delete(user, id)
	return nil
}

// DeleteProduct deletes a product, returns an error
func DeleteProduct(product *model.Product, id string, db *gorm.DB) error {
	db.Delete(product, id)
	return nil
}

// DeleteReview deletes a review, return an error
func DeleteReview(review *model.Review, id string, db *gorm.DB) error {
	db.Delete(review, id)
	return nil
}

// DeleteShop deletes a shop, returns an error
func DeleteShop(shop *model.Shop, id string, db *gorm.DB) error {
	db.Delete(shop, id)
	return nil
}
