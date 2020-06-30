/*
Package env manages the environment variables
*/
package env

import (
	"log"

	"github.com/joho/godotenv"
)

// Load loads the environment file
func Load() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
