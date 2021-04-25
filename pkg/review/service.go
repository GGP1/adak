package review

import (
	"context"
	"time"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/token"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Service provides review operations.
type Service interface {
	Create(ctx context.Context, r *Review, userID string) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]*Review, error)
	GetByID(ctx context.Context, id string) (Review, error)
}

type service struct {
	db *sqlx.DB
}

// NewService returns a new review service.
func NewService(db *sqlx.DB) Service {
	return &service{db}
}

// Create a review.
func (s *service) Create(ctx context.Context, r *Review, userID string) error {
	q := `INSERT INTO reviews
	(id, stars, comment, user_id, product_id, shop_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	id := token.RandString(30)
	r.CreatedAt = time.Now()

	_, err := s.db.ExecContext(ctx, q, id, r.Stars, r.Comment, userID, r.ProductID, r.ShopID, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		logger.Log.Errorf("failed creating review: %v", err)
		return errors.Wrap(err, "couldn't create the review")
	}

	return nil
}

// Delete permanently deletes a review from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		logger.Log.Errorf("failed deleting review: %v", err)
		return errors.Wrap(err, "couldn't delete the review")
	}

	return nil
}

// Get returns a list with all the reviews stored in the database.
func (s *service) Get(ctx context.Context) ([]*Review, error) {
	var reviews []*Review

	if err := s.db.SelectContext(ctx, &reviews, "SELECT * FROM reviews"); err != nil {
		logger.Log.Errorf("failed listing reviews: %v", err)
		return nil, errors.Wrap(err, "couldn't find the reviews")
	}

	return reviews, nil
}

// GetByID retrieves the review requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (Review, error) {
	var review Review

	if err := s.db.GetContext(ctx, &review, "SELECT * FROM reviews WHERE id=$1", id); err != nil {
		return Review{}, errors.Wrap(err, "couldn't find the review")
	}

	return review, nil
}
