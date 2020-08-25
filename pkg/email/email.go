package email

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jmoiron/sqlx"
)

// Emailer provides email operations.
type Emailer interface {
	Add(ctx context.Context, email, token string) error
	Exists(ctx context.Context, email string) bool
	Read(ctx context.Context) ([]List, error)
	Remove(ctx context.Context, email string) error
	Seek(ctx context.Context, email string) error
}

// List represents a list of emails and provides the methods to make
// queries, each list may contain a unique table.
type List struct {
	DB *sqlx.DB
	// tableName is used to distinguish email tables
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
func (l *List) Add(ctx context.Context, email, token string) error {
	q := `INSERT INTO ` + l.tableName + `
	(email, token)
	VALUES ($1, $2)`

	l.Email = email
	l.Token = token

	_, err := l.DB.ExecContext(ctx, q, l.Email, l.Token)
	if err != nil {
		return fmt.Errorf("couldn't create the %s list", l.tableName)
	}

	return nil
}

// Exists checks if the email is already stored in the database.
func (l *List) Exists(ctx context.Context, email string) bool {
	row := l.DB.QueryRowContext(ctx, "SELECT * FROM "+l.tableName+" WHERE email=$1", l.Email)
	if row != nil {
		return true
	}

	return false
}

// Read returns a map with the email list or an error.
func (l *List) Read(ctx context.Context) ([]List, error) {
	var list []List

	if err := l.DB.SelectContext(ctx, &list, "SELECT * FROM "+l.tableName+""); err != nil {
		return nil, errors.New("list not found")
	}

	return list, nil
}

// Remove deletes an email from the list.
func (l *List) Remove(ctx context.Context, email string) error {
	_, err := l.DB.ExecContext(ctx, "DELETE FROM "+l.tableName+" WHERE email=$1", l.Email)
	if err != nil {
		return fmt.Errorf("couldn't delete the email from the %s", l.tableName)
	}

	return nil
}

// Seek looks for the specified email in the list.
func (l *List) Seek(ctx context.Context, email string) error {
	_, err := l.DB.ExecContext(ctx, "SELECT * FROM "+l.tableName+" WHERE email=$1", l.Email)
	if err != nil {
		return errors.New("email not found")
	}

	return nil
}

// Validate checks if the email is valid.
func Validate(email string) error {
	emailRegexp, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		return err
	}

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}

	return nil
}
