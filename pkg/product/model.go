package product

import (
	"errors"
	"time"

	"github.com/GGP1/palo/pkg/review"
)

// Product represents a market commodity.
type Product struct {
	ID          string          `json:"id"`
	ShopID      string          `json:"shop_id" db:"shop_id"`
	Stock       uint            `json:"stock"`
	Brand       string          `json:"brand"`
	Category    string          `json:"category"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Weight      float64         `json:"weight"`
	Discount    float64         `json:"discount"`
	Taxes       float64         `json:"taxes"`
	Subtotal    float64         `json:"subtotal"`
	Total       float64         `json:"total"`
	Reviews     []review.Review `json:"reviews,omitempty"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// Validate checks that there is no missing fields.
func (p *Product) Validate() error {
	if p.ShopID == "" {
		return errors.New("shop id is required")
	}

	if p.Brand == "" {
		return errors.New("brand is required")
	}

	if p.Category == "" {
		return errors.New("category is required")
	}

	if p.Type == "" {
		return errors.New("type is required")
	}

	if p.Weight == 0 {
		return errors.New("weight is required")
	}

	if p.Weight < 0 {
		return errors.New("invalid weight")
	}

	if p.Subtotal == 0 {
		return errors.New("subtotal is required")
	}

	if p.Subtotal < 0 {
		return errors.New("invalid subtotal")
	}

	if p.Discount < 0 {
		return errors.New("invalid discount")
	}

	if p.Taxes < 0 {
		return errors.New("invalid taxes")
	}

	return nil
}
