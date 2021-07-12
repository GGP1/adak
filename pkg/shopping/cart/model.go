package cart

import (
	"gopkg.in/guregu/null.v4/zero"
)

// Cart represents a temporary record of items that the customer selected for purchase.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type Cart struct {
	ID string `json:"id,omitempty"`
	// Counter contains the quantity of products placed in the cart
	Counter zero.Int `json:"counter,omitempty"`
	// 1000 = 1kg
	Weight zero.Int `json:"weight,omitempty"`
	// This field should be used as a percentage
	Discount zero.Int `json:"discount,omitempty"`
	// This field should be used as a percentage
	Taxes    zero.Int  `json:"taxes,omitempty"`
	Subtotal zero.Int  `json:"subtotal,omitempty"`
	Total    zero.Int  `json:"total,omitempty"`
	Products []Product `json:"products,omitempty"`
}

// Product represents a product that has been added to the cart.
type Product struct {
	ID       zero.String `json:"id,omitempty" validate:"uuid4_rfc4122"`
	CartID   zero.String `json:"cart_id,omitempty" db:"cart_id"`
	Quantity zero.Int    `json:"quantity,omitempty" validate:"required,min=1"`
}
