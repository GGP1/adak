package model

import (
	"github.com/jinzhu/gorm"
)

// Review model
type Review struct {
	gorm.Model
	Stars     int    `json:"stars"`
	Comment   string `json:"comment"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
}
