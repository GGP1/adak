package ordering

import "time"

type status int

// Order statuses
const (
	pending status = iota
	paid
	shipping
	shipped
	failed
)

// Order represents a user purchase request.
type Order struct {
	ID           string         `json:"id,omitempty"`
	UserID       string         `json:"user_id,omitempty" db:"user_id"`
	Currency     string         `json:"currency,omitempty"`
	Address      string         `json:"address,omitempty"`
	City         string         `json:"city,omitempty"`
	State        string         `json:"state,omitempty"`
	ZipCode      string         `json:"zip_code,omitempty" db:"zip_code"`
	Country      string         `json:"country,omitempty"`
	Status       status         `json:"status,omitempty"`
	OrderedAt    time.Time      `json:"ordered_at,omitempty" db:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date,omitempty" db:"delivery_date"`
	CartID       string         `json:"cart_id,omitempty" db:"cart_id"`
	Cart         OrderCart      `json:"cart,omitempty"`
	Products     []OrderProduct `json:"products,omitempty"`
}

// OrderCart represents the cart ordered by the user.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type OrderCart struct {
	OrderID  string `json:"order_id,omitempty" db:"order_id"`
	Counter  int    `json:"counter,omitempty"`
	Weight   int64  `json:"weight,omitempty"`
	Discount int64  `json:"discount,omitempty"`
	Taxes    int64  `json:"taxes,omitempty"`
	Subtotal int64  `json:"subtotal,omitempty"`
	Total    int64  `json:"total,omitempty"`
}

// OrderProduct represents a product placed into the cart ordered by the user.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type OrderProduct struct {
	ProductID   string `json:"product_id,omitempty" db:"product_id"`
	OrderID     string `json:"order_id,omitempty" db:"order_id"`
	Quantity    int    `json:"quantity,omitempty"`
	Brand       string `json:"brand,omitempty"`
	Category    string `json:"category,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Weight      int64  `json:"weight,omitempty"`
	Discount    int64  `json:"discount,omitempty"`
	Taxes       int64  `json:"taxes,omitempty"`
	Subtotal    int64  `json:"subtotal,omitempty"`
	Total       int64  `json:"total,omitempty"`
}
