/*
Package updating includes database udpating operations
*/
package updating

import (
	"github.com/GGP1/palo/pkg/model"
	storage "github.com/GGP1/palo/pkg/storage"
)

// Service provides models updating operations.
type Service interface {
	UpdateUser(*model.User, string) error
	UpdateProduct(*model.Product, string) error
	UpdateShop(*model.Shop, string) error
}

// UpdateUser returns nil and updates a user
func UpdateUser(user *model.User, id string) error {
	err := storage.DB.Model(&user).Where("id=?", id).Update(
		"firstname", user.Firstname,
		"lastname", user.Lastname,
		"email", user.Email).
		Error

	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct returns nil and updates a product
func UpdateProduct(product *model.Product, id string) error {
	err := storage.DB.Model(&product).Where("id=?", id).Update(
		"brand", product.Brand,
		"category", product.Category,
		"type", product.Type,
		"description", product.Description,
		"weight", product.Weight,
		"price", product.Price).
		Error

	if err != nil {
		return err
	}

	return nil
}

// UpdateShop returns nil and updates a shop
func UpdateShop(shop *model.Shop, id string) error {
	err := storage.DB.Model(&shop).Where("id=?", id).Update(
		"name", shop.Name).
		Error

	if err != nil {
		return err
	}

	return nil
}
