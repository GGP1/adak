// Package auth provides user authentication and authorization support.
package auth

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/token"
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
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
}

type userData struct {
	email    string
	lastSeen time.Time
}

type session struct {
	sync.RWMutex

	DB *sqlx.DB

	// user session
	store map[string]userData
	// time to wait after failing x times (increments every fail)
	delay map[string]time.Time
	// number of tries to log in
	tries map[string][]struct{}
	// last time the user logged in
	cleaned time.Time
	// session length
	length int
}

// NewSession creates a new session with the necessary dependencies.
func NewSession(db *sqlx.DB) Session {
	return &session{
		DB:      db,
		store:   make(map[string]userData),
		cleaned: time.Now(),
		length:  0,
		tries:   make(map[string][]struct{}),
		delay:   make(map[string]time.Time),
	}
}

// AlreadyLoggedIn checks if the user has an active session or not.
func (s *session) AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("SID")
	if err != nil {
		return false
	}

	s.Lock()
	defer s.Unlock()

	user, ok := s.store[cookie.Value]
	if ok {
		user.lastSeen = time.Now()
		s.store[cookie.Value] = user
	}

	cookie.MaxAge = s.length

	return ok
}

// Clean deletes all the sessions that have expired.
func (s *session) Clean() {
	s.Lock()
	defer s.Unlock()
	for key, value := range s.store {
		if time.Now().Sub(value.lastSeen) > (time.Hour * 120) {
			delete(s.store, key)
		}
	}
	s.cleaned = time.Now()
}

// Login authenticates users.
func (s *session) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, email, password string) error {
	s.Lock()
	defer s.Unlock()

	ip := tracking.GetUserIP(r)

	if !s.delay[ip].IsZero() && s.delay[ip].Sub(time.Now()) > 0 {
		return errors.Errorf("please wait %v before trying again", s.delay[ip].Sub(time.Now()))
	}

	var user User
	query := `SELECT id, cart_id, username, email, password, verified_email FROM users WHERE email=$1`

	// Check if the email exists and if it is verified
	if err := s.DB.GetContext(ctx, &user, query, email); err != nil {
		s.loginDelay(ip)
		return errors.New("invalid email or password")
	}

	if !user.VerfiedEmail {
		return errors.New("please verify your email to log in")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.loginDelay(ip)
		return errors.New("invalid email or password")
	}

	for _, admin := range AdminList {
		if admin == user.Email {
			admID := token.RandString(8)
			cookie.Set(w, "AID", admID, "/", s.length)
		}
	}

	// -SID- used to add the user to the session map
	sID := token.RandString(27)
	cookie.Set(w, "SID", sID, "/", s.length)

	s.store[sID] = userData{user.Email, time.Now()}
	delete(s.tries, ip)
	delete(s.delay, ip)

	// -UID- used to deny users from making requests to other accounts
	userID, err := token.GenerateFixedJWT(user.ID)
	if err != nil {
		return errors.Wrap(err, "failed generating a jwt token")
	}
	cookie.Set(w, "UID", userID, "/", s.length)

	// -CID- used to identify which cart belongs to each user
	cookie.Set(w, "CID", user.CartID, "/", s.length)

	return nil
}

// Login authenticates users using OAuth2.
func (s *session) LoginOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) error {
	var user User
	query := `SELECT id, cart_id, username, email, password, verified_email FROM users WHERE email=$1`

	// Check if the email exists and if it is verified
	if err := s.DB.GetContext(ctx, &user, query, email); err != nil {
		return errors.New("please register before logging in")
	}

	for _, admin := range AdminList {
		if admin == user.Email {
			admID := token.RandString(8)
			cookie.Set(w, "AID", admID, "/", s.length)
		}
	}

	// -SID- used to add the user to the session map
	sID := token.RandString(27)
	if err := cookie.Set(w, "SID", sID, "/", s.length); err != nil {
		return err
	}

	s.Lock()
	s.store[sID] = userData{user.Email, time.Now()}
	s.Unlock()

	userID, err := token.GenerateFixedJWT(user.ID)
	if err != nil {
		return errors.Wrap(err, "failed generating a jwt token")
	}

	// -UID- used to deny users from making requests to other accounts
	if err := cookie.Set(w, "UID", userID, "/", s.length); err != nil {
		return err
	}

	// -CID- used to identify which cart belongs to each user
	if err := cookie.Set(w, "CID", user.CartID, "/", s.length); err != nil {
		return err
	}

	return nil
}

// Logout removes the user session and its cookies.
func (s *session) Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie) {
	admin, _ := r.Cookie("AID")
	if admin != nil {
		cookie.Delete(w, "AID")
	}

	cookie.Delete(w, "SID")
	cookie.Delete(w, "UID")
	cookie.Delete(w, "CID")

	s.Lock()
	delete(s.store, c.Value)
	s.Unlock()

	if time.Now().Sub(s.cleaned) > (time.Second * 30) {
		go s.Clean()
	}
}

// loginDelay increments the time that the user will have to wait after failing.
func (s *session) loginDelay(ip string) {
	s.Lock()
	defer s.Unlock()

	s.tries[ip] = append(s.tries[ip], struct{}{})

	delay := (len(s.tries[ip]) * 2)
	s.delay[ip] = time.Now().Add(time.Second * time.Duration(delay))
}
