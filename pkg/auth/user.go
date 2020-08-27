package auth

import (
	"github.com/GGP1/palo/pkg/email"

	"github.com/pkg/errors"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	ID       string `json:"id"`
	CartID   string `json:"cart_id" db:"cart_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate checks if the user inputs are correct.
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if err := email.Validate(u.Email); err != nil {
		return err
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
