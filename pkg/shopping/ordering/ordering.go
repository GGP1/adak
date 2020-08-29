package ordering

import "time"

// Order represents a user purchase request.
type Order struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id" db:"user_id"`
	Currency     string         `json:"currency"`
	Address      string         `json:"address"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	ZipCode      string         `json:"zip_code" db:"zip_code"`
	Country      string         `json:"country"`
	Status       string         `json:"status"`
	OrderedAt    time.Time      `json:"ordered_at" db:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date" db:"delivery_date"`
	CartID       string         `json:"cart_id" db:"cart_id"`
	Cart         OrderCart      `json:"cart"`
	Products     []OrderProduct `json:"products"`
}

// OrderCart represents the cart ordered by the user.
type OrderCart struct {
	OrderID  string  `json:"order_id" db:"order_id"`
	Counter  int     `json:"counter"`
	Weight   float64 `json:"weight"`
	Discount float64 `json:"discount"`
	Taxes    float64 `json:"taxes"`
	Subtotal float64 `json:"subtotal"`
	Total    float64 `json:"total"`
}

// OrderProduct represents a product placed into the cart ordered by the user.
type OrderProduct struct {
	ProductID   string  `json:"product_id" db:"product_id"`
	OrderID     string  `json:"order_id" db:"order_id"`
	Quantity    int     `json:"quantity"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Description string  `json:"description,omitempty"`
	Weight      float64 `json:"weight"`
	Discount    float64 `json:"discount"`
	Taxes       float64 `json:"taxes"`
	Subtotal    float64 `json:"subtotal"`
	Total       float64 `json:"total"`
}
