/*
Package env let us manage the environment variables
*/
package env

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment file
func LoadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
