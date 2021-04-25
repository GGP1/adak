// Package auth provides user authentication and authorization support.
package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/pkg/tracking"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Clean()
	Login(ctx context.Context, w http.ResponseWriter, r *http.Request, email, password string) error
	LoginOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) error
	Logout(w http.ResponseWriter, r *http.Request)
}

type login struct {
	// time to wait until next attempt
	delay time.Time
	// cumulative number of attempts
	attempts int64
}

type session struct {
	sync.RWMutex

	DB *sqlx.DB

	dev bool
	// user sessions store consisting of map[sessionID]lastSeen
	store map[string]time.Time
	// frustrated login attempts
	loginFails map[string]login
	// last time the session was cleaned
	cleaned time.Time
	// session length in seconds
	length int
}

// NewSession creates a new session with the necessary dependencies.
func NewSession(db *sqlx.DB, dev bool) Session {
	return &session{
		DB:         db,
		dev:        dev,
		store:      make(map[string]time.Time),
		cleaned:    time.Now(),
		length:     0,
		loginFails: make(map[string]login),
	}
}

// AlreadyLoggedIn checks if the user has an active session or not.
func (s *session) AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := cookie.Get(r, "SID")
	if err != nil {
		return false
	}

	s.Lock()
	_, ok := s.store[cookie.Value]
	if ok {
		s.store[cookie.Value] = time.Now()
	}
	s.Unlock()

	// Refresh cookie max age
	cookie.MaxAge = s.length

	return ok
}

// Clean deletes all the sessions that have expired.
// TODO: Use a cron to run it.
func (s *session) Clean() {
	s.Lock()
	for key, value := range s.store {
		if time.Now().Sub(value) > (time.Hour * 168) {
			delete(s.store, key)
		}
	}
	s.cleaned = time.Now()
	s.Unlock()
}

// Login authenticates users.
func (s *session) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, email, password string) error {
	ip := tracking.GetUserIP(r)

	s.RLock()
	delayTime := s.loginFails[ip].delay
	if !delayTime.IsZero() && delayTime.Sub(time.Now()) > 0 {
		return errors.Errorf("please wait %v before trying again", delayTime.Sub(time.Now()))
	}
	s.RUnlock()

	var user User
	query := `SELECT id, cart_id, username, email, password, verified_email FROM users WHERE email=$1`

	// Check if the email exists and if it is verified
	if err := s.DB.GetContext(ctx, &user, query, email); err != nil {
		s.addDelay(ip)
		return errors.New("invalid email or password")
	}

	if !user.VerfiedEmail && !s.dev {
		return errors.New("please verify your email to log in")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.addDelay(ip)
		return errors.New("invalid email or password")
	}

	sID, err := sessionID(user.ID, user.Username)
	if err != nil {
		return err
	}

	s.Lock()
	s.store[sID] = time.Now()
	delete(s.loginFails, ip)
	s.Unlock()

	// -SID- session id
	if err := cookie.Set(w, "SID", sID, "/", s.length); err != nil {
		return err
	}
	// -UID- user id, used to deny users from making requests to other accounts
	if err := cookie.Set(w, "UID", user.ID, "/", s.length); err != nil {
		return err
	}
	// -CID- cart id, used to identify which cart belongs to each user
	return cookie.Set(w, "CID", user.CartID, "/", s.length)
}

// Login authenticates users using OAuth2.
func (s *session) LoginOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) error {
	var user User
	query := `SELECT id, cart_id, username, email, password, verified_email FROM users WHERE email=$1`

	// Check if the email exists and if it is verified
	if err := s.DB.GetContext(ctx, &user, query, email); err != nil {
		return errors.New("please register before logging in")
	}

	sID, err := sessionID(user.ID, user.Username)
	if err != nil {
		return err
	}

	s.Lock()
	s.store[sID] = time.Now()
	s.Unlock()

	// -SID- session id
	if err := cookie.Set(w, "SID", sID, "/", s.length); err != nil {
		return err
	}
	// -UID- user id, used to deny users from making requests to other accounts
	if err := cookie.Set(w, "UID", user.ID, "/", s.length); err != nil {
		return err
	}
	// -CID- cart id, used to identify which cart belongs to each user
	return cookie.Set(w, "CID", user.CartID, "/", s.length)
}

// Logout removes the user session and its cookies.
func (s *session) Logout(w http.ResponseWriter, r *http.Request) {
	// The cookie is already validated
	sessionID, _ := cookie.GetValue(r, "SID")

	cookie.Delete(w, "SID")
	cookie.Delete(w, "UID")
	cookie.Delete(w, "CID")

	s.Lock()
	delete(s.store, sessionID)
	s.Unlock()

	if time.Now().Sub(s.cleaned) > (time.Minute * 30) {
		go s.Clean()
	}
}

// addDelay increments the time that the user will have to wait after failing.
func (s *session) addDelay(ip string) {
	s.Lock()
	login := s.loginFails[ip]
	login.attempts++
	login.delay = time.Now().Add(time.Second * time.Duration(login.attempts*2))
	s.Unlock()
}

func sessionID(userID, username string) (string, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", errors.Wrap(err, "generating salt")
	}

	sessionID := fmt.Sprintf("%s:%s:%s", userID, username, string(salt))
	return sessionID, nil
}
