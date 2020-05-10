package deleting

import (
	"palo/pkg/models"
	stg "palo/pkg/storage"
)

// Deleter provides user and product deleting operations.
type Deleter interface {
	DeleteUser(models.User, string) error
	DeleteProduct(models.Product) error
	DeleteReview(models.Review) error
}

// DeleteUser returns nil and deletes a user
func DeleteUser(user *models.User, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(user)
	return nil
}

// DeleteProduct returns nil and deletes a user
func DeleteProduct(product *models.Product, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(product)
	return nil
}

// DeleteReview returns nil and deletes a user
func DeleteReview(review *models.Review, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(review)
	return nil
}
