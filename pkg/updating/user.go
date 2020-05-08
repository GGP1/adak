package updating

import (
	"time"
)

// User model
type User struct {
	ID        string    `json:"id" `
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
