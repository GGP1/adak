// Package updating includes database updating operations.
package updating

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage.
type Repository interface {
	UpdateProduct(db *gorm.DB, product *model.Product, id string) error
	UpdateShop(db *gorm.DB, shop *model.Shop, id string) error
	UpdateUser(db *gorm.DB, user *model.User, id string) error
}

// Service provides models updating operations.
type Service interface {
	UpdateProduct(db *gorm.DB, product *model.Product, id string) error
	UpdateShop(db *gorm.DB, shop *model.Shop, id string) error
	UpdateUser(db *gorm.DB, user *model.User, id string) error
}

type service struct {
	r Repository
}

// NewService creates a updating service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// UpdateProduct updates a product, returns an error.
func (s *service) UpdateProduct(db *gorm.DB, product *model.Product, id string) error {
	err := db.Model(&product).Where("id=?", id).Update(
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
		return fmt.Errorf("couldn't update the product")
	}

	return nil
}

// UpdateShop updates a shop, returns an error.
func (s *service) UpdateShop(db *gorm.DB, shop *model.Shop, id string) error {
	err := db.Model(&shop).Where("id=?", id).Update(
		"name", shop.Name,
		"country", shop.Location.Country,
		"city", shop.Location.City,
		"address", shop.Location.Address).
		Error
	if err != nil {
		return fmt.Errorf("couldn't update the shop")
	}

	return nil
}

// UpdateUser returns updates a user, returns an error.
func (s *service) UpdateUser(db *gorm.DB, user *model.User, id string) error {
	err := db.Model(&user).Where("id=?", id).Update(
		"username", user.Name,
		"email", user.Email).
		Error
	if err != nil {
		return fmt.Errorf("couldn't update the user")
	}

	return nil
}
