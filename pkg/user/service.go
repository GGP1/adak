package user

import (
	"context"
	"time"

	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/review"
	"github.com/GGP1/palo/pkg/shopping/cart"
	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Repository provides access to the storage.
type Repository interface {
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]ListUser, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (ListUser, error)
	Search(ctx context.Context, search string) ([]ListUser, error)
	Update(ctx context.Context, u *UpdateUser, id string) error
}

// Service provides user operations.
type Service interface {
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]ListUser, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (ListUser, error)
	Search(ctx context.Context, search string) ([]ListUser, error)
	Update(ctx context.Context, u *UpdateUser, id string) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// Create creates a user.
func (s *service) Create(ctx context.Context, user *User) error {
	cartQuery := `INSERT INTO carts
	(id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	userQuery := `INSERT INTO users
	(id, cart_id, username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("email is already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create a cart for each user
	cartID := token.GenerateRunes(30)
	user.CartID = cartID

	cart := cart.New(user.CartID)

	_, err = s.DB.ExecContext(ctx, cartQuery, cart.ID, cart.Counter, cart.Weight,
		cart.Discount, cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't create the cart")
	}

	userID := token.GenerateRunes(30)
	user.CreatedAt = time.Now()

	_, err = s.DB.ExecContext(ctx, userQuery, userID, cart.ID, user.Username, user.Email,
		user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the user")
	}

	return nil
}

// Delete permanently deletes a user from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the user")
	}

	return nil
}

// Get returns a list with all the users stored in the database.
func (s *service) Get(ctx context.Context) ([]ListUser, error) {
	var (
		users []ListUser
		list  []ListUser
	)

	ch := make(chan ListUser)
	errCh := make(chan error)

	if err := s.DB.SelectContext(ctx, &users, "SELECT id, cart_id, username, email FROM users"); err != nil {
		return nil, errors.Wrap(err, "users not found")
	}

	for _, user := range users {
		go func(user ListUser) {
			var reviews []review.Review

			orders, err := ordering.GetByUserID(ctx, s.DB, user.ID)
			if err != nil {
				errCh <- err
			}

			if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			user.Orders = orders
			user.Reviews = reviews

			ch <- user
		}(user)
	}

	select {
	case u := <-ch:
		list = append(list, u)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return list, nil
}

// GetByEmail retrieves the user requested from the database.
func (s *service) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User

	if err := s.DB.GetContext(ctx, &user, "SELECT id, email, username FROM users WHERE email=$1", email); err != nil {
		return User{}, errors.Wrap(err, "user not found")
	}

	return user, nil
}

// GetByID retrieves the user requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (ListUser, error) {
	var (
		user    ListUser
		reviews []review.Review
	)

	if err := s.DB.GetContext(ctx, &user, "SELECT id, cart_id, username, email FROM users WHERE id=$1", id); err != nil {
		return ListUser{}, errors.Wrap(err, "user not found")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", id); err != nil {
		return ListUser{}, errors.Wrap(err, "couldn't find the reviews")
	}

	orders, err := ordering.GetByUserID(ctx, s.DB, id)
	if err != nil {
		return ListUser{}, err
	}

	user.Orders = orders

	return user, nil
}

// Search looks for the users that contain the value specified. (Only text fields)
func (s *service) Search(ctx context.Context, search string) ([]ListUser, error) {
	var (
		users []ListUser
		list  []ListUser
	)

	ch := make(chan ListUser)
	errCh := make(chan error)

	q := `SELECT * FROM users WHERE
	to_tsvector(id || ' ' || username || ' ' || email) 
	@@ to_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &users, q, search); err != nil {
		return nil, errors.Wrap(err, "users not found")
	}

	for _, user := range users {
		go func(user ListUser) {
			var reviews []review.Review

			orders, err := ordering.GetByUserID(ctx, s.DB, user.ID)
			if err != nil {
				errCh <- err
			}

			if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			user.Orders = orders
			user.Reviews = reviews

			ch <- user
		}(user)
	}
	select {
	case u := <-ch:
		list = append(list, u)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return list, nil
}

// Update sets new values for an already existing user.
func (s *service) Update(ctx context.Context, u *UpdateUser, id string) error {
	_, err := s.DB.ExecContext(ctx, "UPDATE users SET username=$2 WHERE id=$1", id, u.Username)
	if err != nil {
		return errors.Wrap(err, "couldn't update the user")
	}

	return nil
}
