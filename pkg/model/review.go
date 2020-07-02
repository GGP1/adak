package model

import (
	"github.com/jinzhu/gorm"
)

// Review model
type Review struct {
	gorm.Model
	Stars     uint8  `json:"stars;omitempty"`
	Comment   string `json:"comment;omitempty"`
	UserID    uint   `json:"user_id;omitempty"`
	ProductID uint   `json:"product_id;omitempty"`
	ShopID    uint   `json:"shop_id;omitempty"`
}
