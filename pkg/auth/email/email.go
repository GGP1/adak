package email

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Emailer provides email operations.
type Emailer interface {
	Add(email, token string) error
	Exists(email string) bool
	Read() ([]List, error)
	Remove(email string) error
	Seek(email string) error
}

// List represents a list of emails and provides the methods to make
// queries, each list may contain a unique table.
type List struct {
	DB *sqlx.DB
	// Used to distinguish between tables of the same struct
	tableName string
	Email     string `json:"email"`
	Token     string `json:"token"`
}

// NewList creates the email list service.
func NewList(db *sqlx.DB, tableName string) Emailer {
	return &List{
		DB:        db,
		tableName: tableName,
		Email:     "",
		Token:     "",
	}
}

// Add a user to the list.
func (l *List) Add(email, token string) error {
	query := `INSERT INTO ` + l.tableName + `
	(email, token)
	VALUES ($1, $2)`

	l.Email = email
	l.Token = token

	_, err := l.DB.Exec(query, l.Email, l.Token)
	if err != nil {
		return fmt.Errorf("couldn't create the %s list", l.tableName)
	}

	return nil
}

// Exists checks if the email is already stored in the database.
func (l *List) Exists(email string) bool {
	row := l.DB.QueryRow("SELECT * FROM "+l.tableName+" WHERE email=$1", l.Email)
	if row != nil {
		return true
	}

	return false
}

// Read returns a map with the email list or an error.
func (l *List) Read() ([]List, error) {
	var list []List

	if err := l.DB.Select(&list, "SELECT * FROM "+l.tableName+""); err != nil {
		return nil, errors.New("list not found")
	}

	return list, nil
}

// Remove deletes an email from the list.
func (l *List) Remove(email string) error {
	_, err := l.DB.Exec("DELETE FROM "+l.tableName+" WHERE email=$1", l.Email)
	if err != nil {
		return fmt.Errorf("couldn't delete the email from the %s", l.tableName)
	}

	return nil
}

// Seek looks for the specified email in the list.
func (l *List) Seek(email string) error {
	_, err := l.DB.Exec("SELECT * FROM "+l.tableName+" WHERE email=$1", l.Email)
	if err != nil {
		return errors.New("email not found")
	}

	return nil
}
