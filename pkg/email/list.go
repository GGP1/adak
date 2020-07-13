package email

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// List contains a list of users emails and their tokens
type List struct {
	DB        *gorm.DB
	tableName string
	sync.RWMutex

	Email string `json:"email"`
	Token string `json:"token"`
	r     Repository
}

// NewList creates the email list service
func NewList(db *gorm.DB, tableName string, r Repository) Service {
	return &List{
		DB:        db,
		tableName: tableName,
		Email:     "",
		Token:     "",
		r:         r,
	}
}

// Add a user to the list
func (l *List) Add(email, token string) error {

	l.Lock()
	l.Email = email
	l.Token = token
	l.Unlock()

	err := l.DB.Table(l.tableName).Create(l).Error
	if err != nil {
		return errors.Wrap(err, "error: couldn't create the pending list")
	}

	return nil
}

// Read returns a map with the email list or an error
func (l *List) Read() (map[string]string, error) {
	err := l.DB.Table(l.tableName).Find(l).Error
	if err != nil {
		return nil, errors.Wrap(err, "error: pending list not found")
	}
	maps := make(map[string]string)
	maps[l.Email] = l.Token

	return maps, nil
}

// Remove deletes a key from the map
func (l *List) Remove(key string) error {
	err := l.DB.Table(l.tableName).Delete(l, key).Error
	if err != nil {
		return errors.Wrap(err, "error: couldn't delete the email from the list")
	}

	return nil
}

// Seek looks for the specified email in the database
func (l *List) Seek(email string) error {
	err := l.DB.Table(l.tableName).First(l, "email = ?", email).Error
	if err != nil {
		return errors.Wrap(err, "error: email not found")
	}

	return nil
}
