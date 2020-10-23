package cart_test

import (
	"testing"

	"github.com/GGP1/palo/pkg/shopping/cart"
)

func TestCartAddProduct(t *testing.T) {
	testCases := []struct {
		product  *cart.Product
		quantity int64
	}{
		{
			product: &cart.Product{
				Weight:   1.5,
				Discount: 6,
				Taxes:    13,
				Subtotal: 4,
			},
			quantity: 3,
		},
		{
			product: &cart.Product{
				Weight:   3,
				Discount: 8,
				Taxes:    21,
				Subtotal: 1.6,
			},
			quantity: 1,
		},
		{
			product: &cart.Product{
				Weight:   1.5,
				Discount: 3,
				Taxes:    7,
				Subtotal: 2,
			},
			quantity: 2,
		},
	}

	c := &cart.Cart{}
	var counter int64

	for _, tC := range testCases {
		cart.AddProduct(c, tC.product, tC.quantity)

		counter += tC.quantity
		taxes := ((tC.product.Subtotal / 100) * tC.product.Taxes)
		discount := ((tC.product.Subtotal / 100) * tC.product.Discount)

		expected := (tC.product.Subtotal + taxes - discount) * float64(tC.quantity)

		if c.Total != expected {
			t.Errorf("Failed computing total, expected: %f, got: %f", expected, c.Total)
		}

		c.Total -= expected
	}

	if c.Counter != counter {
		t.Errorf("Cart counter failed, expected: %d, got: %d", counter, c.Counter)
	}
}

func TestCartRemoveProduct(t *testing.T) {
	testCases := []struct {
		product  *cart.Product
		quantity int64
	}{
		{
			product: &cart.Product{
				Weight:   1.5,
				Discount: 6,
				Taxes:    13,
				Subtotal: 4,
			},
			quantity: 3,
		},
		{
			product: &cart.Product{
				Weight:   3,
				Discount: 8,
				Taxes:    21,
				Subtotal: 1.6,
			},
			quantity: 1,
		},
		{
			product: &cart.Product{
				Weight:   1.5,
				Discount: 0,
				Taxes:    0,
				Subtotal: 2,
			},
			quantity: 2,
		},
	}

	var counter int64
	c := &cart.Cart{}
	for _, tc := range testCases {
		cart.RemoveProduct(c, tc.product, tc.quantity)

		counter -= tc.quantity

		taxes := ((tc.product.Subtotal / 100) * tc.product.Taxes)
		discount := ((tc.product.Subtotal / 100) * tc.product.Discount)

		expected := (-tc.product.Subtotal - taxes + discount) * float64(tc.quantity)

		if c.Total != expected {
			t.Errorf("Failed computing total, expected: %f, got: %f", expected, c.Total)
		}

		c.Total -= expected
	}

	if c.Counter != counter {
		t.Errorf("Cart counter failed, expected: %d, got: %d", counter, c.Counter)
	}
}
