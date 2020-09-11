package shop

import (
	"time"

	"github.com/GGP1/palo/pkg/product"
	"github.com/GGP1/palo/pkg/review"
)

// Shop represents a market with its name and location.
// Each shop has multiple reviews and products.
type Shop struct {
	ID        string            `json:"id"`
	Name      string            `json:"name" validate:"required"`
	Location  Location          `json:"location"`
	Reviews   []review.Review   `json:"reviews,omitempty"`
	Products  []product.Product `json:"products"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}

// Location of the shop.
type Location struct {
	ShopID  string `json:"shop_id" db:"shop_id"`
	Country string `json:"country" validate:"required"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code" db:"zip_code"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
}
