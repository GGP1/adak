/*
Package updating includes database udpating operations
*/
package updating

import (
	"errors"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Updater provides models updating operations.
type Updater interface {
	UpdateUser(*model.User, string) error
	UpdateProduct(*model.Product, string) error
	UpdateShop(*model.Shop, string) error
}

// UpdateUser returns updates a user, returns an error
func UpdateUser(user *model.User, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Model(&user).Where("id=?", id).Update(
		"firstname", user.Firstname,
		"lastname", user.Lastname,
		"email", user.Email).
		Error
	if err != nil {
		return errors.New("error: couldn't update the user")
	}

	return nil
}

// UpdateProduct updates a product, returns an error
func UpdateProduct(product *model.Product, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Model(&product).Where("id=?", id).Update(
		"brand", product.Brand,
		"category", product.Category,
		"type", product.Type,
		"description", product.Description,
		"weight", product.Weight,
		"discount", product.Discount,
		"taxes", product.Taxes,
		"subtotal", product.Subtotal,
		"total", product.Total).
		Error
	if err != nil {
		return errors.New("error: couldn't update the product")
	}

	return nil
}

// UpdateShop updates a shop, returns an error
func UpdateShop(shop *model.Shop, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Model(&shop).Where("id=?", id).Update(
		"name", shop.Name).
		Error
	if err != nil {
		return errors.New("error: couldn't update the shop")
	}

	return nil
}
