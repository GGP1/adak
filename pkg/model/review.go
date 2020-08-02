package model

import (
	"github.com/jinzhu/gorm"
)

// Review represents users critics over a shop or a product.
type Review struct {
	gorm.Model
	Stars     uint8  `json:"stars"`
	Comment   string `json:"comment"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	ShopID    uint   `json:"shop_id"`
}
