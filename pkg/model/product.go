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
	Brand       string   `json:"brand;omitempty"`
	Category    string   `json:"category;omitempty"`
	Type        string   `json:"type;omitempty"`
	Description string   `json:"description"`
	Weight      string   `json:"weight;omitempty"`
	Price       float32  `json:"price;omitempty"`
	ShopID      uint     `json:"shop_id;omitempty"`
	Reviews     []Review `json:"reviews;omitempty" gorm:"foreignkey:ProductID"`
}
