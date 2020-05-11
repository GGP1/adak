/*
Package storage saves all the data and connects to the database
*/
package storage

import (
	"log"
	"os"

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

	// connection
	DB, err = gorm.Open("postgres", connStr)
	CheckErr(err)

	err = DB.DB().Ping()
	CheckErr(err)
}
