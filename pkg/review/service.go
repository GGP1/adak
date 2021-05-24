package review

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4/zero"
)

// Service provides review operations.
type Service interface {
	Create(ctx context.Context, r Review) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Review, error)
	GetByID(ctx context.Context, id string) (Review, error)
}

type service struct {
	db      *sqlx.DB
	mc      *memcache.Client
	metrics metrics
}

// NewService returns a new review service.
func NewService(db *sqlx.DB, mc *memcache.Client) Service {
	return &service{db, mc, initMetrics()}
}

// Create a review.
func (s *service) Create(ctx context.Context, r Review) error {
	s.metrics.incMethodCalls("Create")

	q := `INSERT INTO reviews
	(id, stars, comment, user_id, product_id, shop_id, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, q, r.ID, r.Stars, r.Comment, r.UserID, r.ProductID,
		r.ShopID, zero.TimeFrom(time.Now()))
	if err != nil {
		return errors.Wrap(err, "couldn't create the review")
	}

	s.metrics.totalReviews.Inc()
	return nil
}

// Delete permanently deletes a review from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	s.metrics.incMethodCalls("Delete")
	_, err := s.db.ExecContext(ctx, "DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the review")
	}
	s.metrics.totalReviews.Dec()

	if err := s.mc.Delete(id); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting review from cache")
	}

	return nil
}

// Get returns a list with all the reviews stored in the database.
func (s *service) Get(ctx context.Context) ([]Review, error) {
	s.metrics.incMethodCalls("Get")

	var reviews []Review
	if err := s.db.SelectContext(ctx, &reviews, "SELECT * FROM reviews"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the reviews")
	}

	return reviews, nil
}

// GetByID retrieves the review requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (Review, error) {
	s.metrics.incMethodCalls("GetByID")

	var review Review
	row := s.db.QueryRowContext(ctx, "SELECT * FROM reviews WHERE id=$1", id)
	err := row.Scan(
		&review.ID, &review.Stars, &review.Comment, &review.UserID,
		&review.ShopID, &review.ProductID, &review.CreatedAt,
	)
	if err != nil {
		return Review{}, errors.Wrap(err, "couldn't scan review")
	}

	return review, nil
}
