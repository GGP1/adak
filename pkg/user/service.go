package user

import (
	"context"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Service provides user operations.
type Service interface {
	Create(ctx context.Context, user *AddUser) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]ListUser, error)
	GetByEmail(ctx context.Context, email string) (ListUser, error)
	GetByID(ctx context.Context, id string) (ListUser, error)
	GetByUsername(ctx context.Context, username string) (ListUser, error)
	Search(ctx context.Context, search string) ([]ListUser, error)
	Update(ctx context.Context, u *UpdateUser, id string) error
}

type service struct {
	db *sqlx.DB
}

// NewService returns a new user service.
func NewService(db *sqlx.DB) Service {
	return &service{db}
}

// Create a user.
func (s *service) Create(ctx context.Context, user *AddUser) error {
	cartQuery := `INSERT INTO carts
	(id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	userQuery := `INSERT INTO users
	(id, cart_id, username, email, password, verified_email, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	var count int
	_ = s.db.GetContext(ctx, &count, "SELECT COUNT(id) FROM users WHERE email=$1", user.Email)
	if count > 0 {
		return errors.New("email is already taken")
	}
	_ = s.db.GetContext(ctx, &count, "SELECT COUNT(id) FROM users WHERE username=$1", user.Username)
	if count > 0 {
		return errors.New("username is already taken")
	}

	// Setting a value other than default blocks forever
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Errorf("failed generating user's password hash: %v", err)
		return errors.Wrap(err, "failed generating the password hash")
	}
	user.Password = string(hash)

	// Create a cart for each user
	cartID := token.RandString(30)
	user.CartID = cartID

	cart := cart.New(user.CartID)

	_, err = s.db.ExecContext(ctx, cartQuery, cart.ID, cart.Counter, cart.Weight,
		cart.Discount, cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		logger.Log.Errorf("failed creating user's cart: %v", err)
		return errors.Wrap(err, "couldn't create the cart")
	}

	userID := token.RandString(30)
	user.CreatedAt = time.Now()

	// Ideally a map should be used but some configuration file types do not support them.
	user.IsAdmin = false
	for _, admin := range viper.GetStringSlice("admin.emails") {
		if admin == user.Email {
			user.IsAdmin = true
			break
		}
	}

	_, err = s.db.ExecContext(ctx, userQuery, userID, cart.ID, user.Username, user.Email,
		user.Password, false, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logger.Log.Errorf("failed creating user: %v", err)
		return errors.Wrap(err, "couldn't create the user")
	}

	return nil
}

// Delete permanently deletes a user from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		logger.Log.Errorf("failed deleting user: %v", err)
		return errors.Wrap(err, "couldn't delete the user")
	}

	return nil
}

// Get returns a list with all the users stored in the database.
func (s *service) Get(ctx context.Context) ([]ListUser, error) {
	var users []ListUser

	if err := s.db.SelectContext(ctx, &users, "SELECT id, cart_id, username, email FROM users"); err != nil {
		logger.Log.Errorf("failed listing users: %v", err)
		return nil, errors.Wrap(err, "couldn't find the users")
	}

	list, err := s.getRelationships(ctx, users)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByEmail retrieves the user requested from the database.
func (s *service) GetByEmail(ctx context.Context, email string) (ListUser, error) {
	var user ListUser

	if err := s.db.GetContext(ctx, &user, "SELECT id, cart_id, username, email, is_admin FROM users WHERE email=$1", email); err != nil {
		return ListUser{}, errors.Wrap(err, "couldn't find the user")
	}

	usr, err := s.getRelationship(ctx, user)
	if err != nil {
		return ListUser{}, err
	}

	return usr, nil
}

// GetByID retrieves the user with the id requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (ListUser, error) {
	var user ListUser

	if err := s.db.GetContext(ctx, &user, "SELECT id, cart_id, username, email, is_admin FROM users WHERE id=$1", id); err != nil {
		return ListUser{}, errors.Wrap(err, "couldn't find the user")
	}

	usr, err := s.getRelationship(ctx, user)
	if err != nil {
		return ListUser{}, err
	}

	return usr, nil
}

// GetByUsername retrieves the user with the username requested from the database.
func (s *service) GetByUsername(ctx context.Context, username string) (ListUser, error) {
	var user ListUser

	if err := s.db.GetContext(ctx, &user, "SELECT id FROM users WHERE username=$1", username); err != nil {
		return ListUser{}, errors.Wrap(err, "couldn't find the user")
	}

	usr, err := s.getRelationship(ctx, user)
	if err != nil {
		return ListUser{}, err
	}

	return usr, nil
}

// Search looks for the users that contain the value specified (only text fields).
func (s *service) Search(ctx context.Context, search string) ([]ListUser, error) {
	if strings.ContainsAny(search, ";-\\|@#~€¬<>_()[]}{¡^'") {
		return nil, errors.New("invalid search")
	}

	users := []ListUser{}
	q := `SELECT id, cart_id, username, email FROM users WHERE
	to_tsvector(id || ' ' || username || ' ' || email) 
	@@ plainto_tsquery($1)`

	if err := s.db.SelectContext(ctx, &users, q, search); err != nil {
		logger.Log.Errorf("failed searching users: %v", err)
		return nil, errors.Wrap(err, "couldn't find the users")
	}

	list, err := s.getRelationships(ctx, users)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Update sets new values for an already existing user.
func (s *service) Update(ctx context.Context, u *UpdateUser, id string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE users SET username=$2 WHERE id=$1", id, u.Username)
	if err != nil {
		logger.Log.Errorf("failed updating user: %v", err)
		return errors.Wrap(err, "couldn't update the user")
	}

	return nil
}

func (s *service) getRelationship(ctx context.Context, user ListUser) (ListUser, error) {
	var (
		reviews []review.Review
		orders  []ordering.Order
	)

	if err := s.db.Select(&reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
		logger.Log.Errorf("failed listing user's reviews: %v", err)
		return ListUser{}, errors.Wrap(err, "couldn't find the reviews")
	}

	if err := s.db.Select(&orders, "SELECT * FROM orders WHERE user_id=$1", user.ID); err != nil {
		logger.Log.Errorf("failed listing user's orders: %v", err)
		return ListUser{}, errors.Wrap(err, "couldn't find the orders")
	}

	user.Orders = orders
	user.Reviews = reviews

	return user, nil
}

func (s *service) getRelationships(ctx context.Context, users []ListUser) ([]ListUser, error) {
	ch, errCh := make(chan ListUser), make(chan error, 1)

	for _, user := range users {
		go func(user ListUser) {
			usr, err := s.getRelationship(ctx, user)
			if err != nil {
				errCh <- err
				return
			}

			ch <- usr
		}(user)
	}

	list := make([]ListUser, len(users))
	for i := range users {
		select {
		case user := <-ch:
			list[i] = user
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil
}
