// Package auth provides authentication and authorization support.
package auth

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/GGP1/palo/internal/random"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Clean()
	EmailChange(ctx context.Context, id, newEmail, token string, validatedList email.Emailer) error
	Login(ctx context.Context, w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	PasswordChange(ctx context.Context, id, oldPass, newPass string) error
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

// EmailChange changes the user email.
func (session *session) EmailChange(ctx context.Context, id, newEmail, token string, validatedList email.Emailer) error {
	var user model.User

	if err := session.DB.SelectContext(ctx, &user, "SELECT * FROM users WHERE id=?", id); err != nil {
		return errors.Wrap(err, "invalid email")
	}

	if err := validatedList.Remove(ctx, user.Email); err != nil {
		return err
	}

	user.Email = newEmail

	if err := validatedList.Add(ctx, newEmail, token); err != nil {
		return err
	}

	_, err := session.DB.ExecContext(ctx, "UPDATE users set email=$2 WHERE id=$1", id, newEmail)
	if err != nil {
		return errors.Wrap(err, "couldn't change the email")
	}

	return nil

}

// Login authenticates users and returns a jwt token.
// If the user is already logged in, redirects him to the home page.
func (session *session) Login(ctx context.Context, w http.ResponseWriter, email, password string) error {
	var user model.User

	if err := session.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return errors.Wrap(err, "invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.Wrap(err, "invalid password")
	}

	for _, admin := range adminList {
		if admin == user.Email {
			admID := random.GenerateRunes(8)
			setCookie(w, "AID", admID, "/", session.length)
		}
	}

	// -SID- used to add the user to the session map
	sID := random.GenerateRunes(27)
	setCookie(w, "SID", sID, "/", session.length)

	session.Lock()
	session.store[sID] = userInfo{user.Email, time.Now()}
	session.Unlock()

	// -UID- used to deny users from making requests to other accounts
	userID, err := GenerateFixedJWT(user.ID)
	if err != nil {
		return errors.Wrap(err, "failed generating a jwt token")
	}
	setCookie(w, "UID", userID, "/", session.length)

	// -CID- used to identify wich cart belongs to each user
	cartID := user.CartID
	setCookie(w, "CID", cartID, "/", session.length)

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

// PasswordChange changes the user password.
func (session *session) PasswordChange(ctx context.Context, id, oldPass, newPass string) error {
	var user model.User

	if err := session.DB.GetContext(ctx, &user, "SELECT password FROM users WHERE id=$1", id); err != nil {
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

	_, err = session.DB.ExecContext(ctx, "UPDATE users SET password=$1", user.Password)
	if err != nil {
		return errors.Wrap(err, "couldn't change the password")
	}

	return nil
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
