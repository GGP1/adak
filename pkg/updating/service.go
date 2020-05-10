package updating

import (
	"palo/pkg/models"
	stg "palo/pkg/storage"
)

// Updater provides user and product updating operations.
type Updater interface {
	UpdateUser(models.User, string) error
	UpdateProduct(models.Product, string) error
}

// UpdateUser returns nil and updates a user
func UpdateUser(user *models.User, id string) (err error) {
	stg.DB.Where("id=?", id).Update(user)
	return nil
}

// UpdateProduct returns nil and updates a user
func UpdateProduct(product *models.Product, id string) (err error) {
	stg.DB.Where("id=?", id).Update(product)
	return nil
}
