package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/pkg/postgres"
	"github.com/GGP1/adak/pkg/review"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4/zero"
)

// Service provides user operations.
type Service interface {
	Create(ctx context.Context, user AddUser) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, params params.Query) ([]ListUser, error)
	GetByEmail(ctx context.Context, email string) (ListUser, error)
	GetByID(ctx context.Context, id string) (ListUser, error)
	GetByUsername(ctx context.Context, username string) (ListUser, error)
	IsAdmin(ctx context.Context, id string) (bool, error)
	Search(ctx context.Context, query string) ([]ListUser, error)
	Update(ctx context.Context, u UpdateUser, id string) error
}

type service struct {
	db      *sqlx.DB
	mc      *memcache.Client
	metrics metrics
}

// NewService returns a new user service.
func NewService(db *sqlx.DB, mc *memcache.Client) Service {
	return &service{db, mc, initMetrics()}
}

// Create a user.
func (s *service) Create(ctx context.Context, user AddUser) error {
	s.metrics.incMethodCalls("Create")

	tx, err := s.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "starting transaction")
	}
	defer tx.Commit()

	var exists bool
	q := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 OR username=$2)"
	_ = tx.GetContext(ctx, &exists, q, user.Email, user.Username)
	if exists {
		return errors.New("email or username is already taken")
	}

	// Setting a value other than default blocks forever
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "failed generating the password hash")
	}
	user.Password = string(hash)

	// TODO: use a map O(1). Either change the configuration or create a map
	// on initalization storing admins' emails in it.
	user.IsAdmin = false
	for _, admin := range viper.GetStringSlice("admins") {
		if admin == user.Email {
			user.IsAdmin = true
			break
		}
	}
	userQuery := `INSERT INTO users
	(id, cart_id, username, email, password, is_admin, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, userQuery, user.ID, user.CartID, user.Username,
		user.Email, user.Password, user.IsAdmin, user.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the user")
	}

	s.metrics.registeredUsers.Inc()
	return nil
}

// Delete permanently deletes a user from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	s.metrics.incMethodCalls("Delete")
	if _, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id); err != nil {
		return errors.Wrap(err, "couldn't delete the user")
	}
	s.metrics.registeredUsers.Dec()

	if err := s.mc.Delete(id); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting user from cache")
	}

	return nil
}

// Get returns a list with all the users stored in the database.
func (s *service) Get(ctx context.Context, params params.Query) ([]ListUser, error) {
	s.metrics.incMethodCalls("Get")

	var users []ListUser
	q, args := postgres.AddPagination("SELECT id, cart_id, username, email, is_admin, created_at, updated_at FROM users", params)
	if err := s.db.SelectContext(ctx, &users, q, args...); err != nil {
		return nil, errors.Wrap(err, "couldn't find the users")
	}

	return users, nil
}

// GetByEmail retrieves the user requested from the database.
func (s *service) GetByEmail(ctx context.Context, email string) (ListUser, error) {
	s.metrics.incMethodCalls("GetByEmail")
	return s.getBy(ctx, "email", email)
}

// GetByID retrieves the user with the id requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (ListUser, error) {
	s.metrics.incMethodCalls("GetByID")
	return s.getBy(ctx, "id", id)
}

// GetByUsername retrieves the user with the username requested from the database.
func (s *service) GetByUsername(ctx context.Context, username string) (ListUser, error) {
	s.metrics.incMethodCalls("GetByUsername")
	return s.getBy(ctx, "username", username)
}

// IsAdmin returns if the user is an admin and an error if the query failed.
func (s *service) IsAdmin(ctx context.Context, id string) (bool, error) {
	s.metrics.incMethodCalls("IsAdmin")
	var isAdmin bool
	row := s.db.QueryRowContext(ctx, "SELECT is_admin FROM users WHERE id=$1", id)
	if err := row.Scan(&isAdmin); err != nil {
		return false, errors.Wrap(err, "couldn't scan user role")
	}

	return isAdmin, nil
}

// Search looks for the users that contain the value specified (only text fields).
func (s *service) Search(ctx context.Context, query string) ([]ListUser, error) {
	s.metrics.incMethodCalls("Search")
	var users []ListUser
	q := `SELECT
	id, cart_id, username, email, is_admin
	FROM users
	WHERE to_tsvector(id || ' ' || username || ' ' || email) 
	@@ plainto_tsquery($1)`

	if err := s.db.SelectContext(ctx, &users, q, query); err != nil {
		return nil, errors.Wrap(err, "couldn't find the users")
	}

	return users, nil
}

// Update sets new values for an already existing user.
func (s *service) Update(ctx context.Context, u UpdateUser, id string) error {
	s.metrics.incMethodCalls("Update")
	q := "UPDATE users SET username=$2, updated_at=$3 WHERE id=$1"
	if _, err := s.db.ExecContext(ctx, q, id, u.Username, zero.TimeFrom(time.Now())); err != nil {
		return errors.Wrap(err, "couldn't update the user")
	}

	if err := s.mc.Delete(id); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "couldn't delete user from cache")
	}

	return nil
}

func (s *service) getBy(ctx context.Context, field, value string) (ListUser, error) {
	// Concatenation preferred over fmt.Sprintf
	q := `SELECT
	u.id, u.cart_id, u.username, u.email, u.is_admin, u.created_at, u.updated_at, r.*
	FROM users AS u
	LEFT JOIN reviews AS r ON u.id = r.user_id
	WHERE u.` + field + `=$1`

	rows, err := s.db.QueryContext(ctx, q, value)
	if err != nil {
		return ListUser{}, errors.Wrap(err, "fetching user")
	}
	defer rows.Close()

	return scan(rows)
}

func scan(rows *sql.Rows) (ListUser, error) {
	var user ListUser

	for rows.Next() {
		r := &review.Review{}
		err := rows.Scan(
			&user.ID, &user.CartID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
			&r.ID, &r.Stars, &r.Comment, &r.UserID, &r.ProductID,
			&r.ShopID, &r.CreatedAt,
		)
		if err != nil {
			return ListUser{}, errors.Wrap(err, "couldn't scan user")
		}

		user.Reviews = append(user.Reviews, *r)
	}

	return user, nil
}
