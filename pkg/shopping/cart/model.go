package cart

import "sync"

// Cart represents a temporary record of items that the customer selected for purchase.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type Cart struct {
	mu sync.Mutex

	ID string `json:"id"`
	// Counter contains the quantity of products placed in the cart
	Counter int `json:"counter"`
	// 1000 = 1kg
	Weight int64 `json:"weight"`
	// This field should be used as a percentage
	Discount int64 `json:"discount"`
	// This field should be used as a percentage
	Taxes    int64     `json:"taxes"`
	Subtotal int64     `json:"subtotal"`
	Total    int64     `json:"total"`
	Products []Product `json:"products"`
}

// Product represents a product that has been appended to the cart.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type Product struct {
	ID          string `json:"id"`
	CartID      string `json:"cart_id" db:"cart_id"`
	Quantity    int    `json:"quantity"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	Type        string `json:"type"`
	Description string `json:"description"`
	// 1000 = 1kg
	Weight int64 `json:"weight"`
	// This field should be used as a percentage
	Discount int64 `json:"discount"`
	// This field should be used as a percentage
	Taxes    int64 `json:"taxes"`
	Subtotal int64 `json:"subtotal"`
	Total    int64 `json:"total"`
}
