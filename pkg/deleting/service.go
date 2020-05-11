package deleting

import (
	"palo/pkg/model"
	stg "palo/pkg/storage"
)

// Service provides user and product deleting operations.
type Service interface {
	DeleteUser(model.User, string) error
	DeleteProduct(model.Product) error
	DeleteReview(model.Review) error
}

// DeleteUser returns nil and deletes a user
func DeleteUser(user *model.User, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(user)
	return nil
}

// DeleteProduct returns nil and deletes a user
func DeleteProduct(product *model.Product, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(product)
	return nil
}

// DeleteReview returns nil and deletes a user
func DeleteReview(review *model.Review, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(review)
	return nil
}
