package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Shop represents a market with its name and location.
// Each shop has multiple reviews and products.
type Shop struct {
	gorm.Model
	Name     string    `json:"name"`
	Location Location  `json:"location"`
	Reviews  []Review  `json:"reviews" gorm:"foreignkey:ShopID"`
	Products []Product `json:"products" gorm:"foreignkey:ShopID"`
}

// Location of the shop.
type Location struct {
	ShopID  uint   `json:"shop_id"`
	Country string `json:"country"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	City    string `json:"city"`
	Address string `json:"address"`
}

// Validate checks shop input correctness.
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
