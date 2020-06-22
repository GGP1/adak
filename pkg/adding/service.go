/*
Package adding includes database adding operations
*/
package adding

import (
	"errors"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Service provides models adding operations.
type Service interface {
	Add(interface{}) error
	AddUser(*model.User) error
}

// Add takes the input model and appends it to the database
func Add(model interface{}) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}

	if err := db.Create(model).Error; err != nil {
		return err
	}

	return nil
}

// AddUser takes a new user, hashes its password, sends
// a verification email and appends it to the database
func AddUser(user *model.User) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}

	if user.Email == "" || user.Password == "" {
		return errors.New("email or password missing")
	}

	// Hash password
	hash, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create user
	if err := db.Create(user).Error; err != nil {
		return err
	}

	// Send confirmation email to the user
	email.Confirmation(user)

	return nil
}
