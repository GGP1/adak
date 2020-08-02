// Package auth provides authentication and authorization support.
package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/GGP1/palo/pkg/model"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Repository provides access to the auth storage.
type Repository interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Login(w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	Clean()
}

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Login(w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	Clean()
}

type userInfo struct {
	email    string
	lastSeen time.Time
}

type session struct {
	DB *gorm.DB
	sync.RWMutex

	store      map[string]userInfo
	clean      time.Time
	length     int
	repository Repository
}

// NewSession creates a new session with the necessary dependencies.
func NewSession(db *gorm.DB, r Repository) Session {
	return &session{
		DB:         db,
		store:      make(map[string]userInfo),
		clean:      time.Now(),
		length:     0,
		repository: r,
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

// Login authenticates users and returns a jwt token.
// If the user is already logged in, redirects him to the home page.
func (session *session) Login(w http.ResponseWriter, email, password string) error {
	var user model.User

	err := session.DB.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return errors.New("invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	num := rand.Int()
	admID := strconv.Itoa(num)

	for _, admin := range adminList {
		if admin == user.Email {
			setCookie(w, "AID", admID, "/", session.length)
		}
	}

	// -SID- used to add the user to the session map
	sID := uuid.New()
	setCookie(w, "SID", sID.String(), "/", session.length)

	session.Lock()
	session.store[sID.String()] = userInfo{user.Email, time.Now()}
	session.Unlock()

	// -UID- used to deny users from making requests to other accounts
	id := strconv.Itoa(int(user.ID))
	userID, err := GenerateFixedJWT(id)
	if err != nil {
		return fmt.Errorf("failed generating a jwt token")
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

// Clean deletes all the sessions that have expired.
func (session *session) Clean() {
	for key, value := range session.store {
		if time.Now().Sub(value.lastSeen) > (time.Hour * 24) {
			delete(session.store, key)
		}
	}
	session.clean = time.Now()
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
