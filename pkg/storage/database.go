/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"log"
	"os"

	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// DB global variable
var DB *gorm.DB

// Connect database
func Connect() {
	var err error

	// Load env variables
	osErr := godotenv.Load("../../.env")
	if osErr != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("PQ_URL")

	// Connection
	DB, err = gorm.Open("postgres", connStr)
	CheckErr(err)

	err = DB.DB().Ping()
	CheckErr(err)

	// Auto-migrate models to the db
	DB.AutoMigrate(&model.Product{}, &model.User{}, &model.Review{})

	// Check if database tables are already created
	hasP := DB.HasTable(model.Product{})
	hasU := DB.HasTable(model.User{})
	hasR := DB.HasTable(model.Review{})

	// If a database does not exist, create it
	if hasP != true {
		DB.CreateTable(model.Product{})
	}

	if hasU != true {
		DB.CreateTable(model.User{})
	}

	if hasR != true {
		DB.CreateTable(model.Review{})
	}
}
