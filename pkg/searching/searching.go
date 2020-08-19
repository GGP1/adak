// Package searching provides a service for searching specific information
// in the database related to the core api models.
package searching

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to the storage.
type Repository interface {
	SearchProducts(db *sqlx.DB, products *[]model.Product, search string) error
	SearchShops(db *sqlx.DB, shops *[]model.Shop, search string) error
	SearchUsers(db *sqlx.DB, users *[]model.User, search string) error
}

// Service provides models searching operations.
type Service interface {
	SearchProducts(db *sqlx.DB, products *[]model.Product, search string) error
	SearchShops(db *sqlx.DB, shops *[]model.Shop, search string) error
	SearchUsers(db *sqlx.DB, users *[]model.User, search string) error
}

type service struct {
	r Repository
}

// NewService creates a searching service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// SearchProducts looks for the products that contain the value specified. (Only text fields)
func (s *service) SearchProducts(db *sqlx.DB, products *[]model.Product, search string) error {
	query := `SELECT * FROM products WHERE deleted_at IS NULL AND
	to_tsvector(brand || ' ' || type || ' ' || category || ' ' || description)
	@@ to_tsquery($1)`

	if err := db.Select(&products, query, search); err != nil {
		return fmt.Errorf("couldn't find products: %v", err)
	}

	return nil
}

// SearchShops looks for the shops that contain the value specified. (Only text fields)
func (s *service) SearchShops(db *sqlx.DB, shops *[]model.Shop, search string) error {
	query := `SELECT * FROM shops WHERE deleted_at IS NULL AND
	to_tsvector(name) @@ to_tsquery($1)`

	if err := db.Select(&shops, query, search); err != nil {
		return fmt.Errorf("couldn't find shops: %v", err)
	}

	return nil
}

// SearchUsers looks for the users that contain the value specified. (Only text fields)
func (s *service) SearchUsers(db *sqlx.DB, users *[]model.User, search string) error {
	query := `SELECT * FROM users WHERE deleted_at IS NULL AND
	to_tsvector(firstname || ' ' || lastname || ' ' || email) 
	@@ to_tsquery($1)`

	if err := db.Select(&users, query, search); err != nil {
		return fmt.Errorf("couldn't find users: %v", err)
	}

	return nil
}
