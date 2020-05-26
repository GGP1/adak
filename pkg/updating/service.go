/*
Package updating includes the functions to update an object from the database
*/
package updating

import (
	"github.com/GGP1/palo/pkg/model"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides user and product updating operations.
type Service interface {
	UpdateUser(*model.User, string) error
	UpdateProduct(*model.Product, string) error
	UpdateShop(*model.Shop, string) error
}

// UpdateUser returns nil and updates a user
func UpdateUser(user *model.User, id string) error {
	stg.DB.Where("id=?", id).Update(user)
	return nil
}

// UpdateProduct returns nil and updates a product
func UpdateProduct(product *model.Product, id string) error {
	stg.DB.Where("id=?", id).Update(product)
	return nil
}

// UpdateShop returns nil and updates a shop
func UpdateShop(shop *model.Shop, id string) error {
	stg.DB.Where("id=?", id).Update(shop)
	return nil
}
