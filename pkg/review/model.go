package review

import (
	"time"
)

// Review represents users critics over a shop or product.
type Review struct {
	ID        string    `json:"id,omitempty"`
	Stars     uint8     `json:"stars,omitempty" validate:"min=0,max=5"`
	Comment   string    `json:"comment,omitempty"`
	UserID    string    `json:"user_id,omitempty" db:"user_id" validate:"required"`
	ProductID string    `json:"product_id,omitempty" db:"product_id" validate:"required_without=ShopID"`
	ShopID    string    `json:"shop_id,omitempty" db:"shop_id" validate:"required_without=ProductID"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
