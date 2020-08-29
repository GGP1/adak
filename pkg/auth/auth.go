// Package auth provides authentication and authorization support.
package auth

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/GGP1/palo/internal/token"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Clean()
	Login(ctx context.Context, w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
}

type userInfo struct {
	email    string
	lastSeen time.Time
}

type session struct {
	sync.RWMutex

	DB     *sqlx.DB
	store  map[string]userInfo
	clean  time.Time
	length int
}

// NewSession creates a new session with the necessary dependencies.
func NewSession(db *sqlx.DB) Session {
	return &session{
		DB:     db,
		store:  make(map[string]userInfo),
		clean:  time.Now(),
		length: 0,
	}
}

// AlreadyLoggedIn checks if the user has an active session or not.
func (session *session) AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("SID")
	if err != nil {
		return false
	}

	session.Lock()
	defer session.Unlock()

	user, ok := session.store[cookie.Value]
	if ok {
		user.lastSeen = time.Now()
		session.store[cookie.Value] = user
	}

	cookie.MaxAge = session.length

	return ok
}

// Clean deletes all the sessions that have expired.
func (session *session) Clean() {
	for key, value := range session.store {
		if time.Now().Sub(value.lastSeen) > (time.Hour * 240) {
			delete(session.store, key)
		}
	}
	session.clean = time.Now()
}

// Login authenticates users.
func (session *session) Login(ctx context.Context, w http.ResponseWriter, email, password string) error {
	var user User

	q := `SELECT id, cart_id, username, email, password, email_verified_at FROM users WHERE email=$1`

	// Check if the email exists and if it is verified
	if err := session.DB.GetContext(ctx, &user, q, email); err != nil {
		return errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("invalid password")
	}

	for _, admin := range adminList {
		if admin == user.Email {
			admID := token.GenerateRunes(8)
			setCookie(w, "AID", admID, "/", session.length)
		}
	}

	// -SID- used to add the user to the session map
	sID := token.GenerateRunes(27)
	setCookie(w, "SID", sID, "/", session.length)

	session.Lock()
	session.store[sID] = userInfo{user.Email, time.Now()}
	session.Unlock()

	// -UID- used to deny users from making requests to other accounts
	userID, err := token.GenerateFixedJWT(user.ID)
	if err != nil {
		return errors.Wrap(err, "failed generating a jwt token")
	}
	setCookie(w, "UID", userID, "/", session.length)

	// -CID- used to identify wich cart belongs to each user
	setCookie(w, "CID", user.CartID, "/", session.length)

	return nil
}

// Logout removes the user session and its cookies.
func (session *session) Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie) {
	admin, _ := r.Cookie("AID")
	if admin != nil {
		deleteCookie(w, "AID")
	}

	deleteCookie(w, "SID")
	deleteCookie(w, "UID")
	deleteCookie(w, "CID")

	session.Lock()
	delete(session.store, c.Value)
	session.Unlock()

	if time.Now().Sub(session.clean) > (time.Second * 30) {
		go session.Clean()
	}
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "0",
		Expires:  time.Unix(1414414788, 1414414788000),
		Path:     "/",
		Domain:   "127.0.0.1",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func setCookie(w http.ResponseWriter, name, value, path string, lenght int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   "127.0.0.1",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   lenght,
	})
}
