package product_test

import (
	"context"
	"testing"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/shop"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4/zero"
)

var p = product.Product{
	ID:       zero.StringFrom("4"),
	ShopID:   zero.StringFrom("6"),
	Stock:    zero.IntFrom(1),
	Brand:    zero.StringFrom("brand"),
	Category: zero.StringFrom("category"),
	Type:     zero.StringFrom("type"),
	Weight:   zero.IntFrom(1),
	Subtotal: zero.IntFrom(1),
	Total:    zero.IntFrom(1),
}

// TestMain failed when creating the product service.
func NewProductService(t *testing.T) (context.Context, product.Service) {
	t.Helper()
	logger.Disable()
	ctx, cancel := context.WithCancel(context.Background())

	db := test.StartPostgres(t)
	mc := test.StartMemcached(t)
	service := product.NewService(db, mc)
	createRelationship(ctx, t, db, mc)

	t.Cleanup(func() {
		cancel()
	})

	return ctx, service
}

func TestProductService(t *testing.T) {
	ctx, s := NewProductService(t)

	t.Run("Create", create(ctx, s))
	t.Run("Get", get(ctx, s))
	t.Run("Get by id", getByID(ctx, s))
	t.Run("Update", update(ctx, s))
	t.Run("Search", search(ctx, s))
	t.Run("Delete", delete(ctx, s))
}

func create(ctx context.Context, s product.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Create(ctx, p))

		product, err := s.GetByID(ctx, p.ID.String)
		assert.NoError(t, err)

		assert.Equal(t, p.Brand, product.Brand)
	}
}

func delete(ctx context.Context, s product.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Delete(ctx, p.ID.String))

		shop, err := s.GetByID(ctx, p.ID.String)
		assert.NoError(t, err)

		assert.Equal(t, "", shop.ID.String)
	}
}

func get(ctx context.Context, s product.Service) func(t *testing.T) {
	return func(t *testing.T) {
		params := params.Query{}
		products, err := s.Get(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, p.Category, products[0].Category)
	}
}

func getByID(ctx context.Context, s product.Service) func(t *testing.T) {
	return func(t *testing.T) {
		product, err := s.GetByID(ctx, p.ID.String)
		assert.NoError(t, err)
		assert.Equal(t, p.Type, product.Type)
	}
}

func update(ctx context.Context, s product.Service) func(t *testing.T) {
	return func(t *testing.T) {
		pr := product.UpdateProduct{
			Stock:    zero.IntFrom(1),
			Brand:    zero.StringFrom("brand"),
			Category: zero.StringFrom("category"),
			Type:     zero.StringFrom("type"),
			Weight:   zero.IntFrom(1),
			Subtotal: zero.IntFrom(1),
			Total:    zero.IntFrom(10),
		}

		assert.NoError(t, s.Update(ctx, p.ID.String, pr))

		uptProduct, err := s.GetByID(ctx, p.ID.String)
		assert.NoError(t, err)
		assert.Equal(t, pr.Total, uptProduct.Total)
	}
}

func search(ctx context.Context, s product.Service) func(t *testing.T) {
	return func(t *testing.T) {
		products, err := s.Search(ctx, "brand")
		assert.NoError(t, err)

		t.Log(products)

		var found bool
		for _, pr := range products {
			if pr.ID.String == p.ID.String {
				found = true
				break
			}
		}
		assert.Equal(t, true, found)
	}
}

func createRelationship(ctx context.Context, t *testing.T, db *sqlx.DB, mc *memcache.Client) {
	t.Helper()

	shopService := shop.NewService(db, mc)
	err := shopService.Create(ctx, shop.Shop{
		ID:   "6",
		Name: "test",
	})
	assert.NoError(t, err)
}
