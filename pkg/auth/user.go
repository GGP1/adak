package auth

// User represents a customer trying to log in.
type User struct {
	ID            string `json:"id"`
	CartID        string `json:"cart_id" db:"cart_id"`
	Username      string `json:"username"`
	Email         string `json:"email" validate:"email,required"`
	Password      string `json:"password" validate:"required,min=6"`
	VerifiedEmail bool   `json:"-" db:"verified_email"`
}

// UserAuth is the login request used to authenticate users.
type UserAuth struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}
