/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"os"

	"github.com/GGP1/palo/internal/utils/env"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
)

// Connect to the database
func Connect() (*gorm.DB, error) {
	var err error

	// Load env file
	env.LoadEnv()
	connStr := os.Getenv("PQ_URL")

	// Connection
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	// Auto-migrate modelss to the db
	db.AutoMigrate(&model.Product{}, &model.User{}, &model.Review{}, &model.Shop{})

	// Check if database tables are already created
	productsTable(model.Product{}, db)
	usersTable(model.User{}, db)
	reviewsTable(model.Review{}, db)
	shopsTable(model.Shop{}, db)

	return db, nil
}

// Check if database tables are already
// created if not, create them
func productsTable(m model.Product, db *gorm.DB) {
	h := db.HasTable(m)
	if h != true {
		db.CreateTable(m)
	}
}

func usersTable(m model.User, db *gorm.DB) {
	h := db.HasTable(m)
	if h != true {
		db.CreateTable(m)
	}
}

func reviewsTable(m model.Review, db *gorm.DB) {
	h := db.HasTable(m)
	if h != true {
		db.CreateTable(m)
	}
}

func shopsTable(m model.Shop, db *gorm.DB) {
	h := db.HasTable(m)
	if h != true {
		db.CreateTable(m)
	}
}
