/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"log"
	"os"

	"github.com/GGP1/palo/internal/utils/env"
	"github.com/GGP1/palo/pkg/models"

	"github.com/jinzhu/gorm"
)

// CheckErr test if there's an error and returns it
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
	productsTable(models.Product{})
	usersTable(models.User{})
	reviewsTable(models.Review{})
	shopsTable(models.Shop{})
}

// Check if database tables are already
// created if not, create them
func productsTable(m models.Product) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}

func usersTable(m models.User) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}

func reviewsTable(m models.Review) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}

func shopsTable(m models.Shop) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}
