package model

import "time"

// Review represents users critics over a shop or a product.
type Review struct {
	ID        string    `json:"id"`
	Stars     uint8     `json:"stars"`
	Comment   string    `json:"comment"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	ShopID    string    `json:"shop_id" db:"shop_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
