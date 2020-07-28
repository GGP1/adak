/*
Package auth provides authentication and authorization support.
*/
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
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
)

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

// AlreadyLoggedIn checks if the user have previously logged in.
func (s *session) AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("SID")
	if err != nil {
		return false
	}

	user, ok := s.store[cookie.Value]
	if ok {
		user.lastSeen = time.Now()
		s.store[cookie.Value] = user
	}

	cookie.MaxAge = s.length

	return ok
}

// Login authenticates users and returns a jwt token.
func (s *session) Login(w http.ResponseWriter, email, password string) error {
	user := model.User{}

	err := s.DB.Where("email = ?", email).Take(&user).Error
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
			setCookie(w, "AID", admID, "/", s.length)
		}
	}

	// SessionID -SID- used to add the user to the session map
	sID := uuid.New()
	setCookie(w, "SID", sID.String(), "/", s.length)

	// Store user in session
	s.Lock()
	s.store[sID.String()] = userInfo{user.Email, time.Now()}
	s.Unlock()

	// Convert user id to string and generate a jwt token
	id := strconv.Itoa(int(user.ID))
	userID, err := GenerateFixedJWT(id)
	if err != nil {
		return fmt.Errorf("failed generating a jwt token")
	}
	// UserID -UID- used to deny users from making requests to other accounts
	setCookie(w, "UID", userID, "/", s.length)

	// CartID -CID- used to identify wich cart belongs to each user
	cartID := user.CartID
	setCookie(w, "CID", cartID, "/", s.length)

	return nil
}

// Logout removes the user session and its cookies.
func (s *session) Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie) {
	admin, _ := r.Cookie("AID")
	if admin != nil {
		deleteCookie(w, "AID")
	}

	// Delete SID, UID and CID cookies
	deleteCookie(w, "SID")
	deleteCookie(w, "UID")
	deleteCookie(w, "CID")

	// Delete user from session
	s.Lock()
	delete(s.store, c.Value)
	s.Unlock()

	// Clean up session
	if time.Now().Sub(s.clean) > (time.Second * 30) {
		go s.SessionClean()
	}
}

// SessionClean deletes all the sessions that have expired.
func (s *session) SessionClean() {
	for key, value := range s.store {
		if time.Now().Sub(value.lastSeen) > (time.Hour * 24) {
			delete(s.store, key)
		}
	}
	s.clean = time.Now()
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
		Path:     "/",
		Domain:   "127.0.0.1",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   lenght,
	})
}
