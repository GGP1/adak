/*
Package updating includes database udpating operations
*/
package updating

import (
	"github.com/GGP1/palo/pkg/models"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides user and product updating operations.
type Service interface {
	UpdateUser(*models.User, string) error
	UpdateProduct(*models.Product, string) error
	UpdateShop(*models.Shop, string) error
}

// UpdateUser returns nil and updates a user
func UpdateUser(user *models.User, id string) error {
	stg.DB.Where("id=?", id).Update(user)
	return nil
}

// UpdateProduct returns nil and updates a product
func UpdateProduct(product *models.Product, id string) error {
	stg.DB.Where("id=?", id).Update(product)
	return nil
}

// UpdateShop returns nil and updates a shop
func UpdateShop(shop *models.Shop, id string) error {
	stg.DB.Where("id=?", id).Update(shop)
	return nil
}
