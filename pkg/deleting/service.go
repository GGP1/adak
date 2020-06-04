/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/GGP1/palo/pkg/models"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides models deleting operations.
type Service interface {
	DeleteUser(models.User, string) error
	DeleteProduct(models.Product, string) error
	DeleteReview(models.Review, string) error
	DeleteShop(models.Shop, string) error
}

// DeleteUser returns nil and deletes a user
func DeleteUser(user *models.User, id string) error {
	stg.DB.Delete(user, id)
	return nil
}

// DeleteProduct returns nil and deletes a product
func DeleteProduct(product *models.Product, id string) error {
	stg.DB.Delete(product, id)
	return nil
}

// DeleteReview returns nil and deletes a review
func DeleteReview(review *models.Review, id string) error {
	stg.DB.Delete(review, id)
	return nil
}

// DeleteShop returns nil and deletes a shop
func DeleteShop(shop *models.Shop, id string) error {
	stg.DB.Delete(shop, id)
	return nil
}
