package product

import (
	"github.com/GGP1/adak/pkg/review"

	"gopkg.in/guregu/null.v4/zero"
)

// Product represents a market commodity.
//
// Amounts to be provided in a currency’s smallest unit.
// 100 = 1 USD.
type Product struct {
	ID          zero.String `json:"id,omitempty"`
	ShopID      zero.String `json:"shop_id,omitempty" db:"shop_id" validate:"required"`
	Stock       zero.Int    `json:"stock,omitempty"`
	Brand       zero.String `json:"brand,omitempty" validate:"required"`
	Category    zero.String `json:"category,omitempty" validate:"required"`
	Type        zero.String `json:"type,omitempty" validate:"required"`
	Description zero.String `json:"description,omitempty"`
	// 1000 = 1kg
	Weight    zero.Int        `json:"weight,omitempty" validate:"required,min=1"`
	Discount  zero.Int        `json:"discount,omitempty" validate:"min=0"`
	Taxes     zero.Int        `json:"taxes,omitempty" validate:"min=0"`
	Subtotal  zero.Int        `json:"subtotal,omitempty" validate:"required"`
	Total     zero.Int        `json:"total,omitempty" validate:"min=0"`
	Reviews   []review.Review `json:"reviews,omitempty"`
	CreatedAt zero.Time       `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt zero.Time       `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateProduct is the structure used to update products.
type UpdateProduct struct {
	Stock       zero.Int    `json:"stock,omitempty"`
	Brand       zero.String `json:"brand,omitempty" validate:"required"`
	Category    zero.String `json:"category,omitempty" validate:"required"`
	Type        zero.String `json:"type,omitempty" validate:"required"`
	Description zero.String `json:"description,omitempty"`
	Weight      zero.Int    `json:"weight,omitempty" validate:"required,min=1"`
	Discount    zero.Int    `json:"discount,omitempty" validate:"min=0"`
	Taxes       zero.Int    `json:"taxes,omitempty" validate:"min=0"`
	Subtotal    zero.Int    `json:"subtotal,omitempty" validate:"required"`
	Total       zero.Int    `json:"total,omitempty" validate:"min=0"`
}
