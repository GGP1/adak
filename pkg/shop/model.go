package shop

import (
	"time"

	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"
	"gopkg.in/guregu/null.v4/zero"
)

// Shop represents a market with its name and location.
// Each shop has multiple reviews and products.
type Shop struct {
	ID        string            `json:"id,omitempty"`
	Name      string            `json:"name,omitempty" validate:"required"`
	Location  Location          `json:"location,omitempty"`
	Reviews   []review.Review   `json:"reviews,omitempty"`
	Products  []product.Product `json:"products,omitempty"`
	CreatedAt time.Time         `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt zero.Time         `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateShop is the structure used to update shops.
type UpdateShop struct {
	Name string `json:"name,omitempty" validate:"required"`
}

// Location of the shop.
type Location struct {
	ShopID  string `json:"shop_id,omitempty" db:"shop_id"`
	Country string `json:"country,omitempty" validate:"required"`
	State   string `json:"state,omitempty"`
	ZipCode string `json:"zip_code,omitempty" db:"zip_code"`
	City    string `json:"city,omitempty" validate:"required"`
	Address string `json:"address,omitempty" validate:"required"`
}
