package deleting

import (
	stg "palo/pkg/storage"
)

// Service provides user and product deleting operations.
type Service interface {
	DeleteUser(User, string) error
	DeleteProduct(Product) error
}

// DeleteUser returns nil and deletes a user
func DeleteUser(user *User, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(user)
	return nil
}

// DeleteProduct returns nil and deletes a user
func DeleteProduct(product *Product, id string) (err error) {
	stg.DB.Where("id=?", id).Delete(product)
	return nil
}
