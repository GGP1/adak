package model

import (
	"errors"
	"strings"

	"github.com/GGP1/palo/pkg/ordering"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	gorm.Model
	Name     string           `json:"name"`
	Email    string           `json:"email"`
	Password string           `json:"password"`
	CartID   string           `json:"cart_id"`
	Orders   []ordering.Order `json:"orders" gorm:"foreignkey:UserID"`
	Reviews  []Review         `json:"reviews" gorm:"foreignkey:UserID"`
}

// Validate checks if the inputs are correct.
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Name == "" {
			return errors.New("username is required")
		}

		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}

	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}

		if u.Password == "" {
			return errors.New("password is required")
		}

	default:
		if u.Name == "" {
			return errors.New("username is required")
		}

		if u.Password == "" {
			return errors.New("password is required")
		}

		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
	}

	return nil
}
