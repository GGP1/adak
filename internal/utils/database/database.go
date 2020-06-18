/*
Package database is used to simplify and improve tests services
*/
package database

import (
	"os"

	"github.com/GGP1/palo/internal/utils/env"
	"github.com/jinzhu/gorm"
)

// Connect creates a connection to the database
func Connect() (*gorm.DB, error) {
	env.LoadEnv()
	connect := os.Getenv("PQ_URL")

	db, err := gorm.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	return db, nil
}
