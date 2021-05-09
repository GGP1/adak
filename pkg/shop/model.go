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
	ID        string            `json:"id"`
	Name      string            `json:"name" validate:"required"`
	Location  Location          `json:"location"`
	Reviews   []review.Review   `json:"reviews,omitempty"`
	Products  []product.Product `json:"products"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt zero.Time         `json:"updated_at" db:"updated_at"`
}

// UpdateShop is the structure used to update shops.
type UpdateShop struct {
	Name string `json:"name" validate:"required"`
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
