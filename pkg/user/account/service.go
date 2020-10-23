package account

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/GGP1/palo/pkg/user"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

// Accounts implements the accounts interface.
type Accounts struct {
	db *sqlx.DB

	userClient user.UsersClient
}

// NewService returns a new accounts server.
func NewService(db *sqlx.DB, userConn *grpc.ClientConn) *Accounts {
	return &Accounts{
		db:         db,
		userClient: user.NewUsersClient(userConn),
	}
}

// Run starts the server.
func (a *Accounts) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errors.Wrapf(err, "products: failed listening on port %d", port)
	}

	srv := grpc.NewServer()
	RegisterAccountsServer(srv, a)

	return srv.Serve(lis)
}

// ChangeEmail changes the user email.
func (a *Accounts) ChangeEmail(ctx context.Context, req *ChangeEmailRequest) (*ChangeEmailResponse, error) {
	u, err := a.userClient.GetByID(ctx, &user.GetByIDRequest{ID: req.ID})
	if err != nil {
		return nil, err
	}

	if u.User.CreatedAt.Seconds > time.Now().Add(72*time.Hour).Unix() {
		return nil, errors.New("accounts must be 3 days old to change email")
	}

	_, err = a.userClient.Update(ctx, &user.UpdateRequest{
		ID: req.ID,
		User: &user.UpdateUser{
			Email: req.NewEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ChangePassword changes the user password.
func (a *Accounts) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	u, err := a.userClient.GetByID(ctx, &user.GetByIDRequest{ID: req.ID})
	if err != nil {
		return nil, err
	}

	if u.User.CreatedAt.Seconds > time.Now().Add(72*time.Hour).Unix() {
		return nil, errors.New("accounts must be 3 days old to change password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.User.Password), []byte(req.OldPass)); err != nil {
		return nil, errors.Wrap(err, "invalid old password")
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't generate the password hash")
	}
	u.User.Password = string(newPassHash)

	_, err = a.userClient.Update(ctx, &user.UpdateRequest{
		ID: req.ID,
		User: &user.UpdateUser{
			Password: u.User.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateEmail sets the time when the user validated its email and the token he received.
func (a *Accounts) ValidateEmail(ctx context.Context, req *ValidateEmailRequest) (*ValidateEmailResponse, error) {
	u, err := a.userClient.GetByEmail(ctx, &user.GetByEmailRequest{Email: req.Email})
	if err != nil {
		return nil, err
	}

	_, err = a.userClient.Update(ctx, &user.UpdateRequest{
		ID: u.User.ID,
		User: &user.UpdateUser{
			VerifiedEmail:    req.VerifiedEmail,
			ConfirmationCode: req.ConfirmationCode,
		},
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
