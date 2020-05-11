package updating

import (
	"palo/pkg/model"
	stg "palo/pkg/storage"
)

// Service provides user and product updating operations.
type Service interface {
	UpdateUser(model.User, string) error
	UpdateProduct(model.Product, string) error
}

// UpdateUser returns nil and updates a user
func UpdateUser(user *model.User, id string) (err error) {
	stg.DB.Where("id=?", id).Update(user)
	return nil
}

// UpdateProduct returns nil and updates a user
func UpdateProduct(product *model.Product, id string) (err error) {
	stg.DB.Where("id=?", id).Update(product)
	return nil
}
