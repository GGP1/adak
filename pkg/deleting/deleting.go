// Package deleting includes database deleting operations.
package deleting

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository provides access to the storage.
type Repository interface {
	DeleteProduct(db *sqlx.DB, id string) error
	DeleteReview(db *sqlx.DB, id string) error
	DeleteShop(db *sqlx.DB, id string) error
	DeleteUser(db *sqlx.DB, id string) error
}

// Service provides models deleting operations.
type Service interface {
	DeleteProduct(db *sqlx.DB, id string) error
	DeleteReview(db *sqlx.DB, id string) error
	DeleteShop(db *sqlx.DB, id string) error
	DeleteUser(db *sqlx.DB, id string) error
}

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// DeleteProduct takes a product from the database and permanently deletes it.
func (s *service) DeleteProduct(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("couldn't delete the product: %v", err)
	}

	return nil
}

// DeleteReview takes a review from the database and permanently deletes it.
func (s *service) DeleteReview(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("couldn't delete the review: %v", err)
	}

	return nil
}

// DeleteShop takes a shop from the database and permanently deletes it.
func (s *service) DeleteShop(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM shops WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("couldn't delete the shop: %v", err)
	}

	return nil
}

// DeleteUser takes a user from the database and permanently deletes it.
func (s *service) DeleteUser(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("couldn't delete the user: %v", err)
	}

	return nil
}
