package ordering_test

import (
	"context"
	"testing"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"
	"github.com/GGP1/adak/pkg/user"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4/zero"
)

const (
	cartID  = "5"
	orderID = "27"
	userID  = "95"
)

func NewOrderingService(t *testing.T) (context.Context, ordering.Service, cart.Service) {
	t.Helper()
	logger.Disable()
	ctx, cancel := context.WithCancel(context.Background())

	db := test.StartPostgres(t)
	service := ordering.NewService(db)

	mc := test.StartMemcached(t)
	cartService := cart.NewService(db, mc)
	err := cartService.Create(ctx, cartID)
	assert.NoError(t, err)
	userService := user.NewService(db, mc)
	err = userService.Create(ctx, user.AddUser{ID: userID})
	assert.NoError(t, err)

	t.Cleanup(func() {
		cancel()
	})

	return ctx, service, cartService
}

func TestOrderingService(t *testing.T) {
	ctx, s, cartService := NewOrderingService(t)

	t.Run("New", new(ctx, s, cartService))
	t.Run("Get", get(ctx, s))
	t.Run("Get by ID", getByID(ctx, s))
	t.Run("Get by user ID", getByUserID(ctx, s))
	t.Run("Get cart by ID", getCartByID(ctx, s))
	t.Run("Get products by ID", getProductsByID(ctx, s))
	t.Run("Update status", updateStatus(ctx, s))
	t.Run("Delete", delete(ctx, s))
}

func new(ctx context.Context, s ordering.Service, cartService cart.Service) func(*testing.T) {
	return func(t *testing.T) {
		p := cart.Product{ID: zero.StringFrom("test"), Quantity: zero.IntFrom(1)}
		err := cartService.Add(ctx, p)
		assert.NoError(t, err)

		params := ordering.OrderParams{
			Date: ordering.Date{
				Year:    2150,
				Month:   8,
				Day:     14,
				Hour:    1,
				Minutes: 0,
			},
		}
		_, err = s.New(ctx, orderID, userID, cartID, params, cartService)
		assert.NoError(t, err)
	}
}

func delete(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		err := s.Delete(ctx, orderID)
		assert.NoError(t, err)

		order, err := s.GetByID(ctx, orderID)
		assert.NoError(t, err)

		assert.Equal(t, "", order.ID.String)
	}
}

func get(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		params := params.Query{}
		orders, err := s.Get(ctx, params)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(orders))
		assert.Equal(t, orderID, orders[0].ID.String)
	}
}

func getByID(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		_, err := s.GetByID(ctx, orderID)
		assert.NoError(t, err)
	}
}

func getByUserID(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		_, err := s.GetByUserID(ctx, userID)
		assert.NoError(t, err)
	}
}

func getCartByID(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		_, err := s.GetCartByID(ctx, orderID)
		assert.NoError(t, err)
	}
}

func getProductsByID(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		_, err := s.GetProductsByID(ctx, orderID)
		assert.NoError(t, err)
	}
}

func updateStatus(ctx context.Context, s ordering.Service) func(*testing.T) {
	return func(t *testing.T) {
		status := ordering.Shipped
		err := s.UpdateStatus(ctx, orderID, status)
		assert.NoError(t, err)

		order, err := s.GetByID(ctx, orderID)
		assert.NoError(t, err)

		assert.Equal(t, int64(status), order.Status.Int64)
	}
}
