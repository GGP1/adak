// Package deleting includes database deleting operations.
package deleting

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	DeleteProduct(ctx context.Context, id string) error
	DeleteReview(ctx context.Context, id string) error
	DeleteShop(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
}

// Service provides models deleting operations.
type Service interface {
	DeleteProduct(ctx context.Context, id string) error
	DeleteReview(ctx context.Context, id string) error
	DeleteShop(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// DeleteProduct permanently deletes a product from the database.
func (s *service) DeleteProduct(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the product")
	}

	return nil
}

// DeleteReview permanently deletes a review from the database.
func (s *service) DeleteReview(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the review")
	}

	return nil
}

// DeleteShop permanently deletes a shop from the database.
func (s *service) DeleteShop(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM shops WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the shop")
	}

	return nil
}

// DeleteUser permanently deletes a user from the database.
func (s *service) DeleteUser(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the user")
	}

	return nil
}
