package user

import (
	"fmt"
	"image"
	"time"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/pkg/review"
	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	ID               string           `json:"id"`
	CartID           string           `json:"cart_id" db:"cart_id"`
	Username         string           `json:"username"`
	Email            string           `json:"email"`
	Password         string           `json:"password"`
	EmailVerifiedAt  time.Time        `json:"-" db:"email_verified_at"`
	ConfirmationCode string           `json:"-" db:"confirmation_code"`
	Orders           []ordering.Order `json:"orders,omitempty"`
	Reviews          []review.Review  `json:"reviews,omitempty"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// AddUser is used to create new users.
type AddUser struct {
	ID        string    `json:"id"`
	CartID    string    `json:"cart_id" db:"cart_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ListUser is the structure used to list users.
type ListUser struct {
	ID       string           `json:"id"`
	CartID   string           `json:"cart_id" db:"cart_id"`
	Username string           `json:"username"`
	Email    string           `json:"email"`
	Orders   []ordering.Order `json:"orders,omitempty"`
	Reviews  []review.Review  `json:"reviews,omitempty"`
}

// UpdateUser is the structure used to update users.
type UpdateUser struct {
	Username string `json:"username"`
}

// QRCode creates a QRCode with the link to the user profile.
func (u *ListUser) QRCode() (image.Image, error) {
	qr, err := qrcode.New(fmt.Sprintf("http://127.0.0.1/users/%s", u.ID), qrcode.Medium)
	if err != nil {
		return nil, errors.Wrap(err, "qrcode")
	}

	img := qr.Image(256)

	return img, nil
}

// Validate checks if inputs when creating a user are correct.
func (u *AddUser) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	if len(u.Username) >= 20 {
		return errors.New("username must be less than or equal to 20 characters long")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be greater than 6 characters long")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if err := email.Validate(u.Email); err != nil {
		return err
	}

	return nil
}

// Validate checks if inputs when updating a user are correct.
func (u *UpdateUser) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	return nil
}
