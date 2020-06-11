/*
Package updating includes database udpating operations
*/
package updating

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Service provides models updating operations.
type Service interface {
	UpdateUser() error
	UpdateProduct() error
	UpdateShop() error
}

// UpdateUser returns updates a user, returns an error
func UpdateUser(user *model.User, id string, db *gorm.DB) error {
	err := db.Model(&user).Where("id=?", id).Update(
		"firstname", user.Firstname,
		"lastname", user.Lastname,
		"email", user.Email).
		Error

	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct updates a product, returns an error
func UpdateProduct(product *model.Product, id string, db *gorm.DB) error {
	err := db.Model(&product).Where("id=?", id).Update(
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

// UpdateShop updates a shop, returns an error
func UpdateShop(shop *model.Shop, id string, db *gorm.DB) error {
	err := db.Model(&shop).Where("id=?", id).Update(
		"name", shop.Name).
		Error

	if err != nil {
		return err
	}

	return nil
}
