package auth

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/pkg/tracking"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(ctx context.Context, r *http.Request) bool
	Login(ctx context.Context, w http.ResponseWriter, r *http.Request, email, password string) error
	LoginOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) error
	Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type session struct {
	conf    config.Session
	db      *sqlx.DB
	dev     bool
	metrics metrics
	rdb     *redis.Client
}

// NewSession creates a new session with the necessary dependencies.
func NewSession(db *sqlx.DB, rdb *redis.Client, config config.Session, development bool) Session {
	return &session{
		conf:    config,
		db:      db,
		dev:     development,
		metrics: initMetrics(),
		rdb:     rdb,
	}
}

// AlreadyLoggedIn returns if the user is logged in or not.
func (s *session) AlreadyLoggedIn(ctx context.Context, r *http.Request) bool {
	sID, err := cookie.GetValue(r, "SID")
	if err != nil {
		return false
	}

	value, err := s.rdb.Get(ctx, sID).Result()
	if err != nil {
		return false
	}

	return value == sID[len(sID)-16:]
}

// Login attempts to log a user in.
func (s *session) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, email, password string) error {
	// There is no chance of collision with the rate limiter as it uses the prefix "rate:"
	ip := tracking.GetUserIP(r)

	if s.conf.Delay != 0 {
		ttl := s.rdb.TTL(ctx, ip).Val()
		if ttl > 0 {
			return errors.Errorf("please wait %v before trying again", ttl)
		}
	}

	query := "SELECT id, cart_id, username, email, password, verified_email FROM users WHERE email=$1"
	row := s.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.CartID, &user.Username,
		&user.Email, &user.Password, &user.VerifiedEmail)
	if err != nil {
		logger.Debug(err)
		if err := s.addDelay(ctx, ip); err != nil {
			return errors.Wrap(err, "adding delay")
		}
		return errors.New("invalid email or password")
	}

	if !user.VerifiedEmail && !s.dev {
		return errors.New("please verify your email before logging in")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Debug(err)
		if err := s.addDelay(ctx, ip); err != nil {
			return errors.Wrap(err, "adding delay")
		}
		return errors.New("invalid email or password")
	}

	return s.storeSession(ctx, w, user.ID, user.CartID)
}

// LoginOAuth authenticates users using OAuth2.
func (s *session) LoginOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) error {
	query := "SELECT id, cart_id, username, email, password, verified_email FROM users WHERE email=$1"
	row := s.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.CartID, &user.Username,
		&user.Email, &user.Password, &user.VerifiedEmail)
	if err != nil {
		logger.Debug(err)
		return errors.New("invalid email or password")
	}

	if !user.VerifiedEmail && !s.dev {
		return errors.New("please verify your email before logging in")
	}

	return s.storeSession(ctx, w, user.ID, user.CartID)
}

// Logout removes the user session and its cookies.
func (s *session) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// The error is already checked by AlreadyLoggedIn
	sID, _ := cookie.GetValue(r, "SID")
	if err := s.rdb.Del(ctx, sID).Err(); err != nil {
		return errors.Wrap(err, "deleting the session")
	}
	cookie.Delete(w, "SID")
	cookie.Delete(w, "UID")
	cookie.Delete(w, "CID")
	s.metrics.activeSessions.Dec()
	return nil
}

func (s *session) addDelay(ctx context.Context, key string) error {
	if s.conf.Delay == 0 {
		return nil
	}
	// Cannot use pipeline as "v" is needed to set the ttl
	// To make it incremental two values should be stored, attempts and timestamp
	v := s.rdb.Incr(ctx, key).Val()
	if v > s.conf.Attempts {
		return s.rdb.Expire(ctx, key, time.Duration(s.conf.Delay)*time.Minute).Err()
	}
	return nil
}

// storeSession saves the user key and sets the cookies used to authentication.
func (s *session) storeSession(ctx context.Context, w http.ResponseWriter, userID, cartID string) error {
	// The salt that will be used to identify the user's session
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return errors.Wrap(err, "generating salt")
	}

	sID := userID + ":" + string(salt)

	// Store the salt as the value
	if err := s.rdb.Set(ctx, sID, salt, 0).Err(); err != nil {
		return errors.Wrap(err, "saving session")
	}
	// -SID- session id
	if err := cookie.Set(w, "SID", sID, "/", s.conf.Length); err != nil {
		return err
	}
	// -UID- user id, used to deny users from making requests to other accounts
	if err := cookie.Set(w, "UID", userID, "/", s.conf.Length); err != nil {
		return err
	}
	// -CID- cart id, used to identify which cart belongs to each user
	if err := cookie.Set(w, "CID", cartID, "/", s.conf.Length); err != nil {
		return err
	}

	s.metrics.activeSessions.Inc()
	s.metrics.totalSessions.Inc()
	return nil
}
