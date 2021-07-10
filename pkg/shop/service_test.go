package shop_test

import (
	"context"
	"testing"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/shop"

	"github.com/stretchr/testify/assert"
)

var sh = shop.Shop{
	ID:   "test",
	Name: "Adak",
	Location: shop.Location{
		ShopID:  "test",
		Country: "New Zealand",
		State:   "Auckland",
		ZipCode: "1023",
		City:    "Auckland",
		Address: "8 Hopetoun St",
	},
}

// TestMain failed when creating the shop service.
func NewShopService(t *testing.T) (context.Context, shop.Service) {
	t.Helper()
	logger.Disable()
	ctx, cancel := context.WithCancel(context.Background())

	db := test.StartPostgres(t)
	mc := test.StartMemcached(t)
	service := shop.NewService(db, mc)

	t.Cleanup(func() {
		cancel()
	})

	return ctx, service
}

func TestShopService(t *testing.T) {
	ctx, s := NewShopService(t)

	t.Run("Create", create(ctx, s))
	t.Run("Get", get(ctx, s))
	t.Run("Get by id", getByID(ctx, s))
	t.Run("Update", update(ctx, s))
	t.Run("Search", search(ctx, s))
	t.Run("Delete", delete(ctx, s))
}

func create(ctx context.Context, s shop.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Create(ctx, sh))

		shop, err := s.GetByID(ctx, sh.ID)
		assert.NoError(t, err)

		assert.Equal(t, sh.Name, shop.Name)
	}
}

func delete(ctx context.Context, s shop.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Delete(ctx, sh.ID))

		shop, err := s.GetByID(ctx, sh.ID)
		assert.NoError(t, err)

		assert.Equal(t, "", shop.ID)
	}
}

func get(ctx context.Context, s shop.Service) func(t *testing.T) {
	return func(t *testing.T) {
		params := params.Query{}
		shops, err := s.Get(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, sh.Name, shops[0].Name)
	}
}

func getByID(ctx context.Context, s shop.Service) func(t *testing.T) {
	return func(t *testing.T) {
		shop, err := s.GetByID(ctx, sh.ID)
		assert.NoError(t, err)
		assert.Equal(t, sh.Name, shop.Name)
	}
}

func update(ctx context.Context, s shop.Service) func(t *testing.T) {
	return func(t *testing.T) {
		name := "updated_name"
		assert.NoError(t, s.Update(ctx, sh.ID, shop.UpdateShop{Name: name}))

		uptShop, err := s.GetByID(ctx, sh.ID)
		assert.NoError(t, err)

		assert.Equal(t, name, uptShop.Name)
	}
}

func search(ctx context.Context, s shop.Service) func(t *testing.T) {
	return func(t *testing.T) {
		shops, err := s.Search(ctx, sh.ID)
		assert.NoError(t, err)

		var found bool
		for _, sp := range shops {
			if sp.ID == sh.ID {
				found = true
				break
			}
		}
		assert.Equal(t, true, found)
	}
}
