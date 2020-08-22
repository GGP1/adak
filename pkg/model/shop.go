package model

import (
	"errors"
	"time"
)

// Shop represents a market with its name and location.
// Each shop has multiple reviews and products.
type Shop struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  Location  `json:"location"`
	Reviews   []Review  `json:"reviews,omitempty"`
	Products  []Product `json:"products"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Location of the shop.
type Location struct {
	ShopID  string `json:"shop_id" db:"shop_id"`
	Country string `json:"country"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code" db:"zip_code"`
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

	return nil
}
