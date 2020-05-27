package models

import (
	"github.com/jinzhu/gorm"
)

// Review model
type Review struct {
	gorm.Model
	Stars   uint8  `json:"stars"`
	Comment string `json:"comment"`
	User    User   `json:"user"`
}
