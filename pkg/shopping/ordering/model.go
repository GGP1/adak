package ordering

import (
	"gopkg.in/guregu/null.v4/zero"
)

type status int64

// Order statuses
const (
	Pending status = iota
	Paid
	Shipping
	Shipped
	Failed
)

// Order represents a user purchase request.
type Order struct {
	ID           zero.String    `json:"id,omitempty"`
	UserID       zero.String    `json:"user_id,omitempty" db:"user_id"`
	Currency     zero.String    `json:"currency,omitempty"`
	Address      zero.String    `json:"address,omitempty"`
	City         zero.String    `json:"city,omitempty"`
	State        zero.String    `json:"state,omitempty"`
	ZipCode      zero.String    `json:"zip_code,omitempty" db:"zip_code"`
	Country      zero.String    `json:"country,omitempty"`
	Status       zero.Int       `json:"status,omitempty"`
	OrderedAt    zero.Time      `json:"ordered_at,omitempty" db:"ordered_at"`
	DeliveryDate zero.Time      `json:"delivery_date,omitempty" db:"delivery_date"`
	CartID       zero.String    `json:"cart_id,omitempty" db:"cart_id"`
	Cart         OrderCart      `json:"cart,omitempty"`
	Products     []OrderProduct `json:"products,omitempty"`
	CreatedAt    zero.Time      `json:"created_at,omitempty" db:"created_at"`
}

// OrderCart represents the cart ordered by the user.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type OrderCart struct {
	OrderID  zero.String `json:"order_id,omitempty" db:"order_id"`
	Counter  zero.Int    `json:"counter,omitempty"`
	Weight   zero.Int    `json:"weight,omitempty"`
	Discount zero.Int    `json:"discount,omitempty"`
	Taxes    zero.Int    `json:"taxes,omitempty"`
	Subtotal zero.Int    `json:"subtotal,omitempty"`
	Total    zero.Int    `json:"total,omitempty"`
}

// OrderProduct represents a product placed into the cart ordered by the user.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type OrderProduct struct {
	ProductID   zero.String `json:"product_id,omitempty" db:"product_id"`
	OrderID     zero.String `json:"order_id,omitempty" db:"order_id"`
	Quantity    zero.Int    `json:"quantity,omitempty"`
	Brand       zero.String `json:"brand,omitempty"`
	Category    zero.String `json:"category,omitempty"`
	Type        zero.String `json:"type,omitempty"`
	Description zero.String `json:"description,omitempty"`
	Weight      zero.Int    `json:"weight,omitempty"`
	Discount    zero.Int    `json:"discount,omitempty"`
	Taxes       zero.Int    `json:"taxes,omitempty"`
	Subtotal    zero.Int    `json:"subtotal,omitempty"`
	Total       zero.Int    `json:"total,omitempty"`
}
