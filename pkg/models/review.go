package models

import (
	"github.com/jinzhu/gorm"
)

// Review model
type Review struct {
	gorm.Model
	Stars     uint8  `json:"stars"`
	Comment   string `json:"comment"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
}
