package auth

// User represents a customer. Each user has a unique cart.
type User struct {
	ID           string `json:"id"`
	CartID       string `json:"cart_id" db:"cart_id"`
	Username     string `json:"username"`
	Email        string `json:"email" validate:"email,required"`
	Password     string `json:"password" validate:"required,min=6"`
	VerfiedEmail bool   `json:"-" db:"verified_email"`
	IsAdmin      bool   `json:"-" db:"is_admin"`
}

// UserAuth is the login request used to authenticate users.
type UserAuth struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}
