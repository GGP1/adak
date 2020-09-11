package product

import (
	"time"

	"github.com/GGP1/palo/pkg/review"
)

// Product represents a market commodity.
type Product struct {
	ID          string          `json:"id"`
	ShopID      string          `json:"shop_id" db:"shop_id" validate:"required"`
	Stock       uint            `json:"stock"`
	Brand       string          `json:"brand" validate:"required"`
	Category    string          `json:"category" validate:"required"`
	Type        string          `json:"type" validate:"required"`
	Description string          `json:"description"`
	Weight      float64         `json:"weight" validate:"required,min=0.1"`
	Discount    float64         `json:"discount" validate:"min=0"`
	Taxes       float64         `json:"taxes" validate:"min=0"`
	Subtotal    float64         `json:"subtotal" validate:"required"`
	Total       float64         `json:"total" validate:"min=0"`
	Reviews     []review.Review `json:"reviews,omitempty"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
