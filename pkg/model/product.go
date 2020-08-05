// Package model contains the objects used in the api.
package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Product represents a market commodity.
// It contains its properties, reviews and belongs to one shop.
type Product struct {
	gorm.Model
	ShopID      uint     `json:"shop_id"`
	Stock       uint     `json:"stock"`
	Brand       string   `json:"brand"`
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Weight      float32  `json:"weight"`
	Taxes       float32  `json:"taxes"`
	Discount    float32  `json:"discount"`
	Subtotal    float32  `json:"subtotal"`
	Total       float32  `json:"total"`
	Reviews     []Review `json:"reviews" gorm:"foreignkey:ProductID"`
}

// Validate checks that there is no missing fields.
func (p *Product) Validate() error {
	if p.ShopID == 0 {
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
