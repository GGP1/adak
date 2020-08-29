package cart

// Cart represents a temporary record of items that the customer
// selected for purchase.
type Cart struct {
	ID       string     `json:"id"`
	Counter  int        `json:"counter"` // Counter contains the quantity of products placed in the cart
	Weight   float64    `json:"weight"`
	Discount float64    `json:"discount"`
	Taxes    float64    `json:"taxes"`
	Subtotal float64    `json:"subtotal"`
	Total    float64    `json:"total"`
	Products []*Product `json:"products"`
}

// Product represents a product that has been appended to the cart.
type Product struct {
	ID          string  `json:"id"`
	CartID      string  `json:"cart_id" db:"cart_id"`
	Quantity    int     `json:"quantity"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Discount    float64 `json:"discount"`
	Taxes       float64 `json:"taxes"`
	Subtotal    float64 `json:"subtotal"`
	Total       float64 `json:"total"`
}
