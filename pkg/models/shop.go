package models

import (
	"github.com/jinzhu/gorm"
)

// Shop model
type Shop struct {
	gorm.Model
	Name     string    `json:"name,omitempty"`
	Location Location  `json:"location,omitempty"`
	Products []Product `json:"products,omitempty"`
}

// Location of the shop
type Location struct {
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
}
