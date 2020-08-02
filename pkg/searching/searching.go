// Package searching provides a service for searching specific information
// in the database related to the core api models.
package searching

import (
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	SearchProducts(db *gorm.DB, products *[]model.Product, search string) error
	SearchShops(db *gorm.DB, shops *[]model.Shop, search string) error
	SearchUsers(db *gorm.DB, users *[]model.User, search string) error
}

// Service provides models searching operations.
type Service interface {
	SearchProducts(db *gorm.DB, products *[]model.Product, search string) error
	SearchShops(db *gorm.DB, shops *[]model.Shop, search string) error
	SearchUsers(db *gorm.DB, users *[]model.User, search string) error
}

type service struct {
	r Repository
}

// NewService creates a searching service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// SearchProducts looks for the products that contain the value specified. (Only text fields)
func (s *service) SearchProducts(db *gorm.DB, products *[]model.Product, search string) error {
	err := db.Preload("Reviews").
		Where("deleted_at IS NULL AND to_tsvector(brand || ' ' || type || ' ' || category || ' ' || description) @@ to_tsquery(?)", search).
		Find(&products).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find products")
	}

	return nil
}

// SearchShops looks for the shops that contain the value specified. (Only text fields)
func (s *service) SearchShops(db *gorm.DB, shops *[]model.Shop, search string) error {

	err := db.Preload("Location").Preload("Products").Preload("Reviews").
		Where("deleted_at IS NULL AND to_tsvector(name) @@ to_tsquery(?)", search).Find(&shops).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find shops")
	}

	return nil
}

// SearchUsers looks for the users that contain the value specified. (Only text fields)
func (s *service) SearchUsers(db *gorm.DB, users *[]model.User, search string) error {
	err := db.Preload("Reviews").
		Where("deleted_at IS NULL AND to_tsvector(firstname || ' ' || lastname || ' ' || email) @@ to_tsquery(?)", search).
		Find(&users).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find users")
	}

	return nil
}
