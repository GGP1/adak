/*
Package models contains all the objects used in the api
*/
package models

import (
	"github.com/jinzhu/gorm"
)

// Product model
type Product struct {
	gorm.Model
	Category    string   `json:"category,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description string   `json:"description"`
	Weight      string   `json:"weight,omitempty"`
	Price       int      `json:"price,omitempty"`
	Shop        Shop     `json:"shop,omitempty"`
	Reviews     []Review `json:"reviews,omitempty"`
}
