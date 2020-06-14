/*
Package adding includes database adding operations
*/
package adding

import (
	"github.com/GGP1/palo/internal/utils/email"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Service provides models adding operations.
type Service interface {
	Add(interface{}, *gorm.DB) error
	AddUser(*model.User, *gorm.DB) error
}

// Add takes the input model and appends it to the database
func Add(model interface{}, db *gorm.DB) error {
	if err := db.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// AddUser takes a new user, hashes its password, sends
// a verification email and appends it to the database
func AddUser(user *model.User, db *gorm.DB) error {
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
