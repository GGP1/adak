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
	Firstname string `json:"firstname;not null"`
	Lastname  string `json:"lastname;not null"`
	Email     string `json:"email;unique;not null"`
	Password  string `json:"password;not null"`
}

// Validate checks if the inputs are correct
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Firstname == "" {
			return errors.New("Firstname is required")
		}

		if u.Lastname == "" {
			return errors.New("Lastname is required")
		}

		if u.Email == "" {
			return errors.New("Email is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}

	case "login":
		if u.Email == "" {
			return errors.New("Email is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}

		if u.Password == "" {
			return errors.New("Password is required")
		}

	default:
		if u.Firstname == "" {
			return errors.New("Firstname is required")
		}

		if u.Lastname == "" {
			return errors.New("Lastname is required")
		}

		if u.Password == "" {
			return errors.New("Password is required")
		}

		if u.Email == "" {
			return errors.New("Email is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
	}

	return nil
}
