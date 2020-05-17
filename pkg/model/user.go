package model

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Firstname string `json:"firstname;not null"`
	Lastname  string `json:"lastname;not null"`
	Email     string `json:"email;unique;not null"`
	Password  string `json:"password;not null"`
}
