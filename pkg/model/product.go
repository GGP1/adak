/*
Package model contains all the objects used in the api
*/
package model

import (
	"github.com/jinzhu/gorm"
)

// Product model
type Product struct {
	gorm.Model
	Category    string `json:"category"`
	Brand       string `json:"brand"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Weight      string `json:"weight"`
	Cost        int    `json:"cost"`
}
