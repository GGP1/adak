package email

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage.
type Repository interface {
	Add(email, token string) error
	Read() (map[string]string, error)
	Remove(key string) error
	Seek(email string) error
}

// Service provides email lists operations.
type Service interface {
	Add(email, token string) error
	Read() (map[string]string, error)
	Remove(key string) error
	Seek(email string) error
}

// List represents an email lists and provides the methods to make
// queries, each list may contain a unique table.
type List struct {
	DB *gorm.DB
	// Used to distinguish between tables of the same struct
	tableName string
	r         Repository
	Email     string `json:"email"`
	Token     string `json:"token"`
}

// NewList creates the email list service.
func NewList(db *gorm.DB, tableName string, r Repository) Service {
	return &List{
		DB:        db,
		tableName: tableName,
		r:         r,
		Email:     "",
		Token:     "",
	}
}

// Add a user to the list.
func (l *List) Add(email, token string) error {
	l.Email = email
	l.Token = token

	err := l.DB.Table(l.tableName).Create(l).Error
	if err != nil {
		return fmt.Errorf("couldn't create the pending list")
	}

	return nil
}

// Read returns a map with the email list or an error.
func (l *List) Read() (map[string]string, error) {
	err := l.DB.Table(l.tableName).Find(l).Error
	if err != nil {
		return nil, fmt.Errorf("list not found")
	}
	emailList := make(map[string]string)
	emailList[l.Email] = l.Token

	return emailList, nil
}

// Remove deletes an email from the list.
func (l *List) Remove(email string) error {
	err := l.DB.Table(l.tableName).Where("email=?", email).Delete(l).Error
	if err != nil {
		return fmt.Errorf("couldn't delete the email from the list")
	}

	return nil
}

// Seek looks for the specified email in the list.
func (l *List) Seek(email string) error {
	err := l.DB.Table(l.tableName).First(l, "email = ?", email).Error
	if err != nil {
		return fmt.Errorf("email not found")
	}

	return nil
}
