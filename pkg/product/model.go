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
	ID          string `json:"id"`
	ShopID      string `json:"shop_id" db:"shop_id" validate:"required"`
	Stock       uint   `json:"stock"`
	Brand       string `json:"brand" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	// 1000 = 1kg
	Weight int64 `json:"weight" validate:"required,min=1"`
	// This field should be used as a percentage
	Discount int64 `json:"discount" validate:"min=0"`
	// This field should be used as a percentage
	Taxes     int64           `json:"taxes" validate:"min=0"`
	Subtotal  int64           `json:"subtotal" validate:"required"`
	Total     int64           `json:"total" validate:"min=0"`
	Reviews   []review.Review `json:"reviews,omitempty"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}
