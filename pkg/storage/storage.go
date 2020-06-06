/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"log"
	"os"

	"github.com/GGP1/palo/internal/utils/env"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
)

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
	DB.AutoMigrate(&model.Product{}, &model.User{}, &model.Review{}, &model.Shop{})

	// Check if database tables are already created
	productsTable(model.Product{})
	usersTable(model.User{})
	reviewsTable(model.Review{})
	shopsTable(model.Shop{})
}

// CheckErr test if there's an error and returns it
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Check if database tables are already
// created if not, create them
func productsTable(m model.Product) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}

func usersTable(m model.User) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}

func reviewsTable(m model.Review) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}

func shopsTable(m model.Shop) {
	h := DB.HasTable(m)
	if h != true {
		DB.CreateTable(m)
	}
}
