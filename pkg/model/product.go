/*
Package model contains all the objects used in the api
*/
package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Product model
type Product struct {
	gorm.Model
	ShopID      uint     `json:"shop_id;omitempty"`
	Quantity    int      `json:"amount;omitempty"`
	Brand       string   `json:"brand;omitempty"`
	Category    string   `json:"category;omitempty"`
	Type        string   `json:"type;omitempty"`
	Description string   `json:"description;omitempty"`
	Weight      float32  `json:"weight;omitempty"`
	Taxes       float32  `json:"taxes;omitempty"`
	Discount    float32  `json:"discount;omitempty"`
	Subtotal    float32  `json:"subtotal;omitempty"`
	Total       float32  `json:"total;omitempty"`
	Reviews     []Review `json:"reviews;omitempty" gorm:"foreignkey:ProductID"`
}

// Validate checks that there is no missing fields
func (p *Product) Validate() error {
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

	if p.Subtotal == 0 {
		return errors.New("subtotal is required")
	}

	if p.Total == 0 {
		return errors.New("total is required")
	}

	return nil
}
