/*
Package model contains all the objects used in the api
*/
package model

import (
	"github.com/jinzhu/gorm"
)

// Product model
type Product struct {
	gorm.Model
	ShopID      uint     `json:"shop_id;omitempty"`
	Brand       string   `json:"brand;omitempty"`
	Category    string   `json:"category;omitempty"`
	Type        string   `json:"type;omitempty"`
	Description string   `json:"description"`
	Weight      float32  `json:"weight;omitempty"`
	Taxes       float32  `json:"taxes;omitempty"`
	Discount    float32  `json:"discount;omitempty"`
	Subtotal    float32  `json:"subtotal;omitempty"`
	Total       float32  `json:"total;omitempty"`
	Reviews     []Review `json:"reviews;omitempty" gorm:"foreignkey:ProductID"`
}
