package model

import (
	"errors"
	"time"
)

// Review represents users critics over a shop or a product.
type Review struct {
	ID        string    `json:"id"`
	Stars     uint8     `json:"stars"`
	Comment   string    `json:"comment"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id,omitempty" db:"product_id"`
	ShopID    string    `json:"shop_id,omitempty" db:"shop_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate checks the correction of the reviews sended by the user.
func (r *Review) Validate() error {
	if r.Stars < 0 || r.Stars > 5 {
		return errors.New("invalid stars")
	}

	if r.UserID == "" {
		return errors.New("user id required")
	}

	if r.ProductID == "" && r.ShopID == "" {
		return errors.New("product/shop id required")
	}

	return nil
}
