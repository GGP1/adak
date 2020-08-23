// Package deleting includes database deleting operations.
package deleting

import (
	"context"
	"time"

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

// DeleteProduct takes a product from the database and permanently deletes it.
func (s *service) DeleteProduct(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the product")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DeleteReview takes a review from the database and permanently deletes it.
func (s *service) DeleteReview(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the review")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DeleteShop takes a shop from the database and permanently deletes it.
func (s *service) DeleteShop(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM shops WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the shop")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DeleteUser takes a user from the database and permanently deletes it.
func (s *service) DeleteUser(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the user")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
