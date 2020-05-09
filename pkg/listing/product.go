/*
Package listing: lists products/users taking them from the database
*/
package listing

import (
	"github.com/jinzhu/gorm"
)

// Product model
type Product struct {
	gorm.Model
	Category string `json:"category"`
	Brand    string `json:"brand"`
	Name     string `json:"name"`
	Weight   string `json:"weight"`
	Cost     int    `json:"cost"`
}
