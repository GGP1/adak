// Package auth provides authentication and authorization support.
package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/GGP1/palo/internal/uuid"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Clean()
	EmailChange(id, newEmail, token string, validatedList email.Emailer) error
	Login(w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	PasswordChange(id, oldPass, newPass string) error
}

type userInfo struct {
	email    string
	lastSeen time.Time
}

type session struct {
	DB *gorm.DB
	sync.RWMutex

	store  map[string]userInfo
	clean  time.Time
	length int
}

// NewSession creates a new session with the necessary dependencies.
func NewSession(db *gorm.DB) Session {
	return &session{
		DB:     db,
		store:  make(map[string]userInfo),
		clean:  time.Now(),
		length: 0,
	}
}

// AlreadyLoggedIn checks if the user have previously logged in or not.
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
		if time.Now().Sub(value.lastSeen) > (time.Hour * 24) {
			delete(session.store, key)
		}
	}
	session.clean = time.Now()
}

// EmailChange changes the user email.
func (session *session) EmailChange(id, newEmail, token string, validatedList email.Emailer) error {
	var user model.User

	err := session.DB.Where("id=?", id).First(&user).Error
	if err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}

	err = validatedList.Remove(user.Email)
	if err != nil {
		return err
	}

	user.Email = newEmail

	err = validatedList.Add(newEmail, token)
	if err != nil {
		return err
	}

	err = session.DB.Save(&user).Error
	if err != nil {
		return fmt.Errorf("couldn't change the email: %v", err)
	}

	return nil
}

// Login authenticates users and returns a jwt token.
// If the user is already logged in, redirects him to the home page.
func (session *session) Login(w http.ResponseWriter, email, password string) error {
	var user model.User

	err := session.DB.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	for _, admin := range adminList {
		if admin == user.Email {
			admID := uuid.GenerateRandRunes(8)
			setCookie(w, "AID", admID, "/", session.length)
		}
	}

	// -SID- used to add the user to the session map
	sID := uuid.GenerateRandRunes(27)
	setCookie(w, "SID", sID, "/", session.length)

	session.Lock()
	session.store[sID] = userInfo{user.Email, time.Now()}
	session.Unlock()

	// -UID- used to deny users from making requests to other accounts
	userID := uuid.GenerateRandRunes(24)
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
func (session *session) PasswordChange(id, oldPass, newPass string) error {
	var user model.User

	err := session.DB.Where("id=?", id).First(&user).Error
	if err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPass))
	if err != nil {
		return fmt.Errorf("invalid old password: %v", err)
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("couldn't generate the password hash: %v", err)
	}
	user.Password = string(newPassHash)

	err = session.DB.Save(&user).Error
	if err != nil {
		return fmt.Errorf("couldn't change the password: %v", err)
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
