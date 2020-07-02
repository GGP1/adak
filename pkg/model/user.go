package model

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Firstname string   `json:"firstname;omitempty"`
	Lastname  string   `json:"lastname;omitempty"`
	Email     string   `json:"email;unique;omitempty"`
	Password  string   `json:"password;omitempty"`
	Role      string   `json:"role;omitempty"`
	Reviews   []Review `json:"reviews" gorm:"foreignkey:UserID"`
}

// Validate checks if the inputs are correct
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Firstname == "" {
			return errors.New("firstname is required")
		}

		if u.Lastname == "" {
			return errors.New("lastname is required")
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
		if u.Firstname == "" {
			return errors.New("firstname is required")
		}

		if u.Lastname == "" {
			return errors.New("lastname is required")
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
