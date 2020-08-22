// Package updating includes database updating operations.
package updating

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to the storage.
type Repository interface {
	UpdateProduct(product *model.Product, id string) error
	UpdateShop(shop *model.Shop, id string) error
	UpdateUser(user *model.User, id string) error
}

// Service provides models updating operations.
type Service interface {
	UpdateProduct(product *model.Product, id string) error
	UpdateShop(shop *model.Shop, id string) error
	UpdateUser(user *model.User, id string) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a updating service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// UpdateProduct updates a product, returns an error.
func (s *service) UpdateProduct(p *model.Product, id string) error {
	query := `UPDATE products SET stock=$2, brand=$3, category=$4, type=$5,
	description=$6, weight=$7, discount=$8, taxes=$9, subtotal=$10, total=$11
	WHERE id=$1`

	_, err := s.DB.Exec(query, id, p.Stock, p.Brand, p.Category, p.Type,
		p.Description, p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return fmt.Errorf("couldn't update the product: %v", err)
	}

	return nil
}

// UpdateShop updates a shop, returns an error.
func (s *service) UpdateShop(shop *model.Shop, id string) error {
	query := `UPDATE shops SET name=$2, country=$3, city=$4, address=$5
	WHERE id=$1`

	_, err := s.DB.Exec(query, id, shop.Name, shop.Location.Country,
		shop.Location.City, shop.Location.Address)
	if err != nil {
		return fmt.Errorf("couldn't update the shop: %v", err)
	}

	return nil
}

// UpdateUser returns updates a user, returns an error.
func (s *service) UpdateUser(u *model.User, id string) error {
	query := `UPDATE users SET name=$2, email=$3 WHERE id=$1`

	_, err := s.DB.Exec(query, id, u.Username, u.Email)
	if err != nil {
		return fmt.Errorf("couldn't update the user: %v", err)
	}

	return nil
}
