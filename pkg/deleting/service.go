/*
Package deleting gives us the methods to delete objects from the database
*/
package deleting

import (
	"github.com/GGP1/palo/pkg/model"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides user and product deleting operations.
type Service interface {
	DeleteUser(model.User, string) error
	DeleteProduct(model.Product, string) error
	DeleteReview(model.Review, string) error
	DeleteShop(model.Shop, string) error
}

// DeleteUser returns nil and deletes a user
func DeleteUser(user *model.User, id string) error {
	stg.DB.Delete(user, id)
	return nil
}

// DeleteProduct returns nil and deletes a product
func DeleteProduct(product *model.Product, id string) error {
	stg.DB.Delete(product, id)
	return nil
}

// DeleteReview returns nil and deletes a review
func DeleteReview(review *model.Review, id string) error {
	stg.DB.Delete(review, id)
	return nil
}

// DeleteShop returns nil and deletes a shop
func DeleteShop(shop *model.Shop, id string) error {
	stg.DB.Delete(shop, id)
	return nil
}
