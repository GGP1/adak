/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"log"
	"os"

	"github.com/GGP1/palo/pkg/models"
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

	// Auto-migrate modelss to the db
	DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Review{}, &models.Shop{})

	// Check if database tables are already created
	productTable := DB.HasTable(models.Product{})
	userTable := DB.HasTable(models.User{})
	reviewTable := DB.HasTable(models.Review{})
	shopTable := DB.HasTable(models.Shop{})

	// If a database does not exist, create it
	if productTable != true {
		DB.CreateTable(models.Product{})
	}

	if userTable != true {
		DB.CreateTable(models.User{})
	}

	if reviewTable != true {
		DB.CreateTable(models.Review{})
	}

	if shopTable != true {
		DB.CreateTable(models.Shop{})
	}
}
