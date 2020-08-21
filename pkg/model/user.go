package model

import (
	"errors"
	"fmt"
	"image"
	"regexp"
	"strings"
	"time"

	"github.com/GGP1/palo/pkg/shopping/ordering"

	qrcode "github.com/skip2/go-qrcode"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	ID        string           `json:"id"`
	CartID    string           `json:"cart_id" db:"cart_id"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	Orders    []ordering.Order `json:"orders"`
	Reviews   []Review         `json:"reviews"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
}

// Card respresents a user card.
type Card struct {
	Number   string `json:"number"`
	ExpMonth string `json:"exp_month"`
	ExpYear  string `json:"exp_year"`
	CVC      string `json:"cvc"`
}

// QRCode creates a QRCode with the link to the user profile.
func (u *User) QRCode() (image.Image, error) {
	qr, err := qrcode.New(fmt.Sprintf("http://127.0.0.1/users/%s", u.ID), qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("qrcode: %v", err)
	}

	img := qr.Image(256)

	return img, nil
}

// Validate checks if the user inputs are correct.
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("username is required")
		}

		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := ValidateEmail(u.Email); err != nil {
			return err
		}

	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := ValidateEmail(u.Email); err != nil {
			return err
		}

		if u.Password == "" {
			return errors.New("password is required")
		}

	default:
		if u.Username == "" {
			return errors.New("username is required")
		}

		if u.Password == "" {
			return errors.New("password is required")
		}

		if len(u.Password) < 6 {
			return errors.New("password must be equal or greater than 6 characters")
		}

		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := ValidateEmail(u.Email); err != nil {
			return err
		}
	}

	return nil
}

// ValidateEmail checks if the email is valid.
func ValidateEmail(email string) error {
	emailRegexp, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		return err
	}

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}

	return nil
}
