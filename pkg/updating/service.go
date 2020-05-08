package updating

import (
	stg "palo/pkg/storage"
)

// Service provides user and product updating operations.
type Service interface {
	UpdateUser(User, string) error
	UpdateProduct(Product, string) error
}

// UpdateUser returns nil and updates a user
func UpdateUser(user *User, id string) (err error) {
	stg.DB.Where("id=?", id).Update(user)
	return nil
}

// UpdateProduct returns nil and updates a user
func UpdateProduct(product *Product, id string) (err error) {
	stg.DB.Where("id=?", id).Update(product)
	return nil
}
