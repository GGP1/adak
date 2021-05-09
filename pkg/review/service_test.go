package review_test

import (
	"context"
	"testing"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shop"
	"github.com/GGP1/adak/pkg/user"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4/zero"
)

var r = review.Review{
	ID:        zero.StringFrom("test"),
	Stars:     zero.IntFrom(5),
	Comment:   zero.StringFrom("testing"),
	UserID:    zero.StringFrom("1"),
	ShopID:    zero.StringFrom("5"),
	ProductID: zero.StringFrom("3"),
}

// TestMain failed when creating the review service.
func NewReviewService(t *testing.T) (context.Context, review.Service) {
	t.Helper()
	logger.Disable()
	ctx, cancel := context.WithCancel(context.Background())

	db := test.StartPostgres(t)
	mc := test.StartMemcached(t)
	service := review.NewService(db, mc)
	createRelations(ctx, t, db, mc)

	t.Cleanup(func() {
		cancel()
	})

	return ctx, service
}

func TestReviewService(t *testing.T) {
	ctx, s := NewReviewService(t)

	t.Run("Create", create(ctx, s))
	t.Run("Get", get(ctx, s))
	t.Run("Get by id", getByID(ctx, s))
	t.Run("Delete", delete(ctx, s))
}

func create(ctx context.Context, s review.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Create(ctx, r))

		review, err := s.GetByID(ctx, r.ID.String)
		assert.NoError(t, err)

		assert.Equal(t, r.Comment, review.Comment)
	}
}

func delete(ctx context.Context, s review.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Delete(ctx, r.ID.String))

		_, err := s.GetByID(ctx, r.ID.String)
		assert.Error(t, err)
	}
}

func get(ctx context.Context, s review.Service) func(t *testing.T) {
	return func(t *testing.T) {
		reviews, err := s.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, r.ShopID, reviews[0].ShopID)
	}
}

func getByID(ctx context.Context, s review.Service) func(t *testing.T) {
	return func(t *testing.T) {
		review, err := s.GetByID(ctx, r.ID.String)
		assert.NoError(t, err)
		assert.Equal(t, r.ID, review.ID)
	}
}

func createRelations(ctx context.Context, t *testing.T, db *sqlx.DB, mc *memcache.Client) {
	t.Helper()
	userService := user.NewService(db, mc)
	err := userService.Create(ctx, user.AddUser{
		ID:       "1",
		CartID:   "test",
		Email:    "test",
		Username: "test",
		Password: "test",
	})
	assert.NoError(t, err)

	shopService := shop.NewService(db, mc)
	err = shopService.Create(ctx, shop.Shop{
		ID:   "5",
		Name: "test",
	})
	assert.NoError(t, err)

	productService := product.NewService(db, mc)
	err = productService.Create(ctx, product.Product{
		ID:       zero.StringFrom("3"),
		ShopID:   zero.StringFrom("5"),
		Stock:    zero.IntFrom(1),
		Brand:    zero.StringFrom("test"),
		Category: zero.StringFrom("test"),
		Type:     zero.StringFrom("test"),
		Weight:   zero.IntFrom(1),
		Subtotal: zero.IntFrom(1),
		Total:    zero.IntFrom(1),
	})
	assert.NoError(t, err)
}
