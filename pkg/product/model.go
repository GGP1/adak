package product

import (
	"time"

	"github.com/GGP1/adak/pkg/review"
)

// Product represents a market commodity.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type Product struct {
	ID          string `json:"id,omitempty"`
	ShopID      string `json:"shop_id,omitempty" db:"shop_id" validate:"required"`
	Stock       uint   `json:"stock,omitempty"`
	Brand       string `json:"brand,omitempty" validate:"required"`
	Category    string `json:"category,omitempty" validate:"required"`
	Type        string `json:"type,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
	// 1000 = 1kg
	Weight int64 `json:"weight,omitempty" validate:"required,min=1"`
	// This field should be used as a percentage
	Discount int64 `json:"discount,omitempty" validate:"min=0"`
	// This field should be used as a percentage
	Taxes     int64           `json:"taxes,omitempty" validate:"min=0"`
	Subtotal  int64           `json:"subtotal,omitempty" validate:"required"`
	Total     int64           `json:"total,omitempty" validate:"min=0"`
	Reviews   []review.Review `json:"reviews,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at,omitempty" db:"updated_at"`
}
