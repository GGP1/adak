package user

import (
	"time"

	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shopping/ordering"
)

// User represents platform customers.
// Each user has a unique cart.
type User struct {
	ID               string           `json:"id,omitempty" validate:"unique"`
	CartID           string           `json:"cart_id,omitempty" db:"cart_id"`
	Username         string           `json:"username,omitempty"`
	Email            string           `json:"email,omitempty" validate:"email"`
	Password         string           `json:"password,omitempty"`
	VerifiedEmail    bool             `json:"verified_email,omitempty" db:"verified_email"`
	IsAdmin          bool             `json:"is_admin,omitempty" db:"is_admin"`
	ConfirmationCode string           `json:"confirmation_code,omitempty" db:"confirmation_code"`
	Orders           []ordering.Order `json:"orders,omitempty"`
	Reviews          []review.Review  `json:"reviews,omitempty"`
	CreatedAt        time.Time        `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty" db:"updated_at"`
}

// AddUser is used to create new users.
type AddUser struct {
	ID        string    `json:"id,omitempty"`
	CartID    string    `json:"cart_id,omitempty" db:"cart_id"`
	Username  string    `json:"username,omitempty" validate:"required,max=25"`
	Email     string    `json:"email,omitempty" validate:"email,required"`
	Password  string    `json:"password,omitempty" validate:"required,min=6"`
	IsAdmin   bool      `json:"is_admin,omitempty" db:"is_admin"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

// ListUser is the structure used to list users.
type ListUser struct {
	ID       string          `json:"id,omitempty"`
	CartID   string          `json:"cart_id,omitempty" db:"cart_id"`
	Username string          `json:"username,omitempty"`
	Email    string          `json:"email,omitempty" validate:"email"`
	IsAdmin  bool            `json:"is_admin,omitempty" db:"is_admin"`
	Reviews  []review.Review `json:"reviews,omitempty"`
}

// UpdateUser is the structure used to update users.
type UpdateUser struct {
	Username string `json:"username,omitempty" validate:"required"`
}
