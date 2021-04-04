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
	ID           string         `json:"id"`
	UserID       string         `json:"user_id" db:"user_id"`
	Currency     string         `json:"currency"`
	Address      string         `json:"address"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	ZipCode      string         `json:"zip_code" db:"zip_code"`
	Country      string         `json:"country"`
	Status       status         `json:"status"`
	OrderedAt    time.Time      `json:"ordered_at" db:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date" db:"delivery_date"`
	CartID       string         `json:"cart_id" db:"cart_id"`
	Cart         OrderCart      `json:"cart"`
	Products     []OrderProduct `json:"products"`
}

// OrderCart represents the cart ordered by the user.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type OrderCart struct {
	OrderID  string `json:"order_id" db:"order_id"`
	Counter  int    `json:"counter"`
	Weight   int64  `json:"weight"`
	Discount int64  `json:"discount"`
	Taxes    int64  `json:"taxes"`
	Subtotal int64  `json:"subtotal"`
	Total    int64  `json:"total"`
}

// OrderProduct represents a product placed into the cart ordered by the user.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type OrderProduct struct {
	ProductID   string `json:"product_id" db:"product_id"`
	OrderID     string `json:"order_id" db:"order_id"`
	Quantity    int    `json:"quantity"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Weight      int64  `json:"weight"`
	Discount    int64  `json:"discount"`
	Taxes       int64  `json:"taxes"`
	Subtotal    int64  `json:"subtotal"`
	Total       int64  `json:"total"`
}
