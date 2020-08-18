package email

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Emailer provides email operations.
type Emailer interface {
	Add(email, token string) error
	Exists(email string) bool
	Read() (map[string]string, error)
	Remove(email string) error
	Seek(email string) error
}

// List represents a list of emails and provides the methods to make
// queries, each list may contain a unique table.
type List struct {
	DB *gorm.DB
	// Used to distinguish between tables of the same struct
	tableName string
	Email     string `json:"email"`
	Token     string `json:"token"`
}

// NewList creates the email list service.
func NewList(db *gorm.DB, tableName string) Emailer {
	return &List{
		DB:        db,
		tableName: tableName,
		Email:     "",
		Token:     "",
	}
}

// Add a user to the list.
func (l *List) Add(email, token string) error {
	l.Email = email
	l.Token = token

	if err := l.DB.Table(l.tableName).Create(l).Error; err != nil {
		return errors.New("couldn't create the pending list")
	}

	return nil
}

// Exists checks if the email is already stored in the database.
func (l *List) Exists(email string) bool {
	rows := l.DB.Table(l.tableName).First(l, "email = ?", email).RowsAffected
	if rows == 0 {
		return false
	}

	return true
}

// Read returns a map with the email list or an error.
func (l *List) Read() (map[string]string, error) {
	if err := l.DB.Table(l.tableName).Find(l).Error; err != nil {
		return nil, errors.New("list not found")
	}

	emailList := make(map[string]string)
	emailList[l.Email] = l.Token

	return emailList, nil
}

// Remove deletes an email from the list.
func (l *List) Remove(email string) error {
	if err := l.DB.Table(l.tableName).Where("email=?", email).Delete(l).Error; err != nil {
		return errors.New("couldn't delete the email from the list")
	}

	return nil
}

// Seek looks for the specified email in the list.
func (l *List) Seek(email string) error {
	if err := l.DB.Table(l.tableName).First(l, "email = ?", email).Error; err != nil {
		return errors.New("email not found")
	}

	return nil
}
