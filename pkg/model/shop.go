package model

import (
	"errors"

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

// Validate checks shop input correctness
func (s *Shop) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}

	if s.Location.Country == "" {
		return errors.New("country is required")
	}

	if s.Location.City == "" {
		return errors.New("country is required")
	}

	if s.Location.Address == "" {
		return errors.New("address is required")
	}

	if s.Location.ShopID == 0 {
		return errors.New("shop id is required")
	}

	return nil
}
