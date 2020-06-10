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

// DeleteUser returns nil and deletes a user
func DeleteUser(user *model.User, id string, db *gorm.DB) error {
	db.Delete(user, id)
	return nil
}

// DeleteProduct returns nil and deletes a product
func DeleteProduct(product *model.Product, id string, db *gorm.DB) error {
	db.Delete(product, id)
	return nil
}

// DeleteReview returns nil and deletes a review
func DeleteReview(review *model.Review, id string, db *gorm.DB) error {
	db.Delete(review, id)
	return nil
}

// DeleteShop returns nil and deletes a shop
func DeleteShop(shop *model.Shop, id string, db *gorm.DB) error {
	db.Delete(shop, id)
	return nil
}
