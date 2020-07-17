/*
Package storage makes the database connection
*/
package storage

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
)

// NewDatabase creates a database and returns gorm.DB and an error
//
// Return the close function so it's not avoided in the future
func NewDatabase() (*gorm.DB, func() error, error) {
	// Connection
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return nil, nil, err
	}

	// Check connectivity
	err = db.DB().Ping()
	if err != nil {
		return nil, nil, err
	}

	// Auto-migrate models
	db.AutoMigrate(&model.Product{}, &model.User{}, &model.Review{}, &model.Shop{}, &model.Location{})
	// Create tables
	tableExists(db, model.Product{}, model.User{}, model.Review{}, model.Shop{}, model.Location{})

	if db.Table("pending_list").HasTable(&email.List{}) != true && db.Table("validated_list").HasTable(&email.List{}) != true {
		db.Table("pending_list").CreateTable(&email.List{}).AutoMigrate(&email.List{})
		db.Table("validated_list").CreateTable(&email.List{}).AutoMigrate(&email.List{})
	}

	return db, db.Close, nil
}

// Check if a table is already created, if not, create it
func tableExists(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		if db.HasTable(model) != true {
			db.CreateTable(model)
		}
	}
}
