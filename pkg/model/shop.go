package model

import (
	"github.com/jinzhu/gorm"
)

// Shop model
type Shop struct {
	gorm.Model
	Name     string    `json:"name;omitempty"`
	Location Location  `json:"location;omitempty"`
	Reviews  []Review  `json:"reviews;omitempty" gorm:"foreignkey:ShopID"`
	Products []Product `json:"products;omitempty" gorm:"foreignkey:ShopID"`
}

// Location of the shop
type Location struct {
	Country string `json:"country;omitempty"`
	City    string `json:"city;omitempty"`
	Address string `json:"address;omitempty"`
	ShopID  uint   `json:"shop_id;omitempty"`
}
