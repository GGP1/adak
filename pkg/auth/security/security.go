/*
Package security provides password hashing and comparison between hash/password functions
*/
package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the user password
func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return hash, err
}

// ComparePasswords takes the user password and its hash and compares it
func ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err
}
