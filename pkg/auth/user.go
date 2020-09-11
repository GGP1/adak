package auth

import (
	"time"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	ID              string    `json:"id"`
	CartID          string    `json:"cart_id" db:"cart_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email" validate:"email,required"`
	Password        string    `json:"password" validate:"required,min=6"`
	EmailVerifiedAt time.Time `json:"-" db:"email_verified_at"`
}
