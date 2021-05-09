package review

import (
	"gopkg.in/guregu/null.v4/zero"
)

// Review represents users critics over a shop or product.
type Review struct {
	ID        zero.String `json:"id,omitempty"`
	Stars     zero.Int    `json:"stars,omitempty" validate:"min=0,max=5"`
	Comment   zero.String `json:"comment,omitempty"`
	UserID    zero.String `json:"user_id,omitempty" db:"user_id" validate:"required"`
	ProductID zero.String `json:"product_id,omitempty" db:"product_id" validate:"required_without=ShopID"`
	ShopID    zero.String `json:"shop_id,omitempty" db:"shop_id" validate:"required_without=ProductID"`
	CreatedAt zero.Time   `json:"created_at,omitempty" db:"created_at"`
}
