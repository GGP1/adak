package cart_test

import (
	"context"
	"os"
	"testing"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4/zero"
)

const cartID = "1234"

var service cart.Service

func TestMain(m *testing.M) {
	poolMc, resourceMc, mc, err := test.RunMemcached()
	if err != nil {
		logger.Fatal(err)
	}
	poolPg, resourcePg, db, err := test.RunPostgres()
	if err != nil {
		logger.Fatal(err)
	}

	service = cart.NewService(db, mc)
	if err := service.Create(context.Background(), cartID); err != nil {
		logger.Fatal(err)
	}

	code := m.Run()

	if err := poolMc.Purge(resourceMc); err != nil {
		logger.Fatal(err)
	}
	if err := poolPg.Purge(resourcePg); err != nil {
		logger.Fatal(err)
	}

	os.Exit(code)
}

func TestAdd(t *testing.T) {
	ctx := context.Background()
	quantity := zero.IntFrom(5)
	product := &cart.Product{
		ID:       zero.StringFrom("1"),
		Brand:    zero.StringFrom("Mónaco"),
		Quantity: quantity,
	}

	err := service.Add(ctx, cartID, product)
	assert.NoError(t, err)

	products, err := service.FilterBy(ctx, cartID, "brand", "Mónaco")
	assert.NoError(t, err)

	assert.Equal(t, quantity, products[0].Quantity)
}

func TestCheckout(t *testing.T) {
	total, err := service.Checkout(context.Background(), cartID)
	assert.NoError(t, err)

	assert.Equal(t, int64(0), total)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	err := service.Delete(ctx, cartID)
	assert.NoError(t, err)

	c, err := service.Get(ctx, cartID)
	assert.NoError(t, err)
	assert.Equal(t, "", c.ID)

	// Recreate
	err = service.Create(ctx, cartID)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	_, err := service.Get(context.Background(), cartID)
	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	ctx := context.Background()
	product := &cart.Product{
		ID:       zero.StringFrom("2"),
		Type:     zero.StringFrom("deleted"),
		Quantity: zero.IntFrom(1),
	}

	err := service.Add(ctx, cartID, product)
	assert.NoError(t, err)

	err = service.Remove(ctx, cartID, "2", 1)
	assert.NoError(t, err)

	_, err = service.FilterBy(ctx, cartID, "type", "deleted")
	assert.Error(t, err, "Expected no products found error")
}

func TestReset(t *testing.T) {
	ctx := context.Background()
	err := service.Reset(ctx, cartID)
	assert.NoError(t, err)

	expected := cart.New(cartID)
	got, err := service.Get(ctx, cartID)
	assert.NoError(t, err)

	assert.Equal(t, expected.Counter.Int64, got.Counter.Int64)
	assert.Equal(t, expected.Discount.Int64, got.Discount.Int64)
	assert.Equal(t, expected.Subtotal.Int64, got.Subtotal.Int64)
	assert.Equal(t, expected.Taxes.Int64, got.Taxes.Int64)
	assert.Equal(t, expected.Total.Int64, got.Total.Int64)
}

func TestSize(t *testing.T) {
	size, err := service.Size(context.Background(), cartID)
	assert.NoError(t, err)

	assert.Equal(t, int64(0), size)
}
