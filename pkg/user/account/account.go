package account

import (
	"context"

	"github.com/GGP1/palo/pkg/email"
	"github.com/GGP1/palo/pkg/user"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Repository provides access to the storage.
type Repository interface {
	ChangeEmail(ctx context.Context, id, newEmail, token string, validatedList email.Emailer) error
	ChangePassword(ctx context.Context, id, oldPass, newPass string) error
}

// Service provides user account operations.
type Service interface {
	ChangeEmail(ctx context.Context, id, newEmail, token string, validatedList email.Emailer) error
	ChangePassword(ctx context.Context, id, oldPass, newPass string) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a searching service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// Change changes the user email.
func (s *service) ChangeEmail(ctx context.Context, id, newEmail, token string, validatedList email.Emailer) error {
	var user user.User

	if err := s.DB.SelectContext(ctx, &user, "SELECT * FROM users WHERE id=?", id); err != nil {
		return errors.Wrap(err, "invalid email")
	}

	if err := validatedList.Remove(ctx, user.Email); err != nil {
		return err
	}

	user.Email = newEmail

	if err := validatedList.Add(ctx, newEmail, token); err != nil {
		return err
	}

	_, err := s.DB.ExecContext(ctx, "UPDATE users set email=$2 WHERE id=$1", id, newEmail)
	if err != nil {
		return errors.Wrap(err, "couldn't change the email")
	}

	return nil
}

// ChangePassword changes the user password.
func (s *service) ChangePassword(ctx context.Context, id, oldPass, newPass string) error {
	var user user.User

	if err := s.DB.GetContext(ctx, &user, "SELECT password FROM users WHERE id=$1", id); err != nil {
		return errors.Wrap(err, "invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPass)); err != nil {
		return errors.Wrap(err, "invalid old password")
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "couldn't generate the password hash")
	}
	user.Password = string(newPassHash)

	_, err = s.DB.ExecContext(ctx, "UPDATE users SET password=$1", user.Password)
	if err != nil {
		return errors.Wrap(err, "couldn't change the password")
	}

	return nil
}
