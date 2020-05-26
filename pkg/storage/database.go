/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"log"
	"os"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/utils/env"

	"github.com/jinzhu/gorm"
)

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// DB global variable
var DB *gorm.DB

// Connect to the database
func Connect() {
	var err error

	// Load env file
	env.LoadEnv()

	connStr := os.Getenv("PQ_URL")

	// Connection
	DB, err = gorm.Open("postgres", connStr)
	CheckErr(err)

	err = DB.DB().Ping()
	CheckErr(err)

	// Auto-migrate models to the db
	DB.AutoMigrate(&model.Product{}, &model.User{}, &model.Review{}, &model.Shop{})

	// Check if database tables are already created
	productTable := DB.HasTable(model.Product{})
	userTable := DB.HasTable(model.User{})
	reviewTable := DB.HasTable(model.Review{})
	shopTable := DB.HasTable(model.Shop{})

	// If a database does not exist, create it
	if productTable != true {
		DB.CreateTable(model.Product{})
	}

	if userTable != true {
		DB.CreateTable(model.User{})
	}

	if reviewTable != true {
		DB.CreateTable(model.Review{})
	}

	if shopTable != true {
		DB.CreateTable(model.Shop{})
	}
}
