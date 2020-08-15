package model

import (
	"fmt"
	"image"
	"regexp"
	"strings"

	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	gorm.Model
	CartID   string           `json:"cart_id"`
	Name     string           `json:"name"`
	Email    string           `json:"email"`
	Password string           `json:"password"`
	Orders   []ordering.Order `json:"orders" gorm:"foreignkey:UserID"`
	Reviews  []Review         `json:"reviews" gorm:"foreignkey:UserID"`
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
	qr, err := qrcode.New(fmt.Sprintf("http://127.0.0.1/users/%d", u.ID), qrcode.Medium)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the qrcode")
	}

	img := qr.Image(256)

	return img, nil
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

		if err := u.ValidateEmail(u.Email); err != nil {
			return err
		}

	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}

		if err := u.ValidateEmail(u.Email); err != nil {
			return err
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

		if err := u.ValidateEmail(u.Email); err != nil {
			return err
		}
	}

	return nil
}

// ValidateEmail checks if the email is valid.
func (u *User) ValidateEmail(email string) error {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}
