package model

import (
	"fmt"
	"image"
	"strings"

	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
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
