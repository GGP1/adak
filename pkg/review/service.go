package review

import (
	"context"
	"time"

	"github.com/GGP1/palo/internal/token"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	Create(ctx context.Context, r *Review, userID string) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Review, error)
	GetByID(ctx context.Context, id string) (Review, error)
}

// Service provides review operations.
type Service interface {
	Create(ctx context.Context, r *Review, userID string) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Review, error)
	GetByID(ctx context.Context, id string) (Review, error)
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// Create creates a review.
func (s *service) Create(ctx context.Context, r *Review, userID string) error {
	q := `INSERT INTO reviews
	(id, stars, comment, user_id, product_id, shop_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	err := r.Validate()
	if err != nil {
		return err
	}

	id := token.GenerateRunes(30)
	r.CreatedAt = time.Now()

	_, err = s.DB.ExecContext(ctx, q, id, r.Stars, r.Comment, userID, r.ProductID, r.ShopID, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the review")
	}

	return nil
}

// Delete permanently deletes a review from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the review")
	}

	return nil
}

// Get returns a list with all the reviews stored in the database.
func (s *service) Get(ctx context.Context) ([]Review, error) {
	var reviews []Review

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews"); err != nil {
		return nil, errors.Wrap(err, "reviews not found")
	}

	return reviews, nil
}

// GetByID retrieves the review requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (Review, error) {
	var review Review

	if err := s.DB.GetContext(ctx, &review, "SELECT * FROM reviews WHERE id=$1", id); err != nil {
		return Review{}, errors.Wrap(err, "review not found")
	}

	return review, nil
}
