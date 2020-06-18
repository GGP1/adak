/*
Package database is used to simplify and improve tests services
*/
package database

import (
	"github.com/jinzhu/gorm"
)

// Connect creates a connection to the database
func Connect(string) (*gorm.DB, error) {
	return gorm.Open("postgres", URL)
}
