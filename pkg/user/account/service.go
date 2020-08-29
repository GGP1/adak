package account

import (
	"context"
	"time"

	"github.com/GGP1/palo/pkg/user"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Repository provides access to the storage.
type Repository interface {
	ChangeEmail(ctx context.Context, id, newEmail, token string) error
	ChangePassword(ctx context.Context, id, oldPass, newPass string) error
	ValidateUserEmail(ctx context.Context, id, confirmationCode string, verifiedAt time.Time) error
}

// Service provides user account operations.
type Service interface {
	ChangeEmail(ctx context.Context, id, newEmail, token string) error
	ChangePassword(ctx context.Context, id, oldPass, newPass string) error
	ValidateUserEmail(ctx context.Context, id, confirmationCode string, verifiedAt time.Time) error
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
func (s *service) ChangeEmail(ctx context.Context, id, newEmail, token string) error {
	var user user.User

	if err := s.DB.SelectContext(ctx, &user, "SELECT * FROM users WHERE id=?", id); err != nil {
		return errors.Wrap(err, "invalid email")
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

	_, err = s.DB.ExecContext(ctx, "UPDATE users SET password=$2 WHERE id=$1", user.ID, user.Password)
	if err != nil {
		return errors.Wrap(err, "couldn't change the password")
	}

	return nil
}

// ValidateUserEmail sets the time when the user validated its email and the token he received.
func (s *service) ValidateUserEmail(ctx context.Context, id, confirmationCode string, verifiedAt time.Time) error {
	q := "UPDATE users SET email_verified_at=$2, confirmation_code=$3 WHERE id=$1"

	_, err := s.DB.ExecContext(ctx, q, id, verifiedAt, confirmationCode)
	if err != nil {
		return errors.Wrap(err, "couldn't validate the user")
	}

	return nil
}
