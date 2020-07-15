/*
Package handler contains the methods used by the router
*/
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"

	"github.com/google/uuid"
)

// AuthRepository provides access to the auth storage
type AuthRepository interface {
	Login(db *gorm.DB, validatedList email.Service) http.HandlerFunc
	Logout() http.HandlerFunc
}

// Session provides auth operations
type Session interface {
	Login(db *gorm.DB, validatedList email.Service) http.HandlerFunc
	Logout() http.HandlerFunc
}

type userInfo struct {
	email    string
	lastSeen time.Time
}

type session struct {
	store      map[string]userInfo
	clean      time.Time
	length     int
	repository AuthRepository
}

// NewSession creates a new session with the necessary dependencies
func NewSession(r AuthRepository) Session {
	return &session{
		store:      make(map[string]userInfo),
		clean:      time.Now(),
		length:     0,
		repository: r,
	}
}

// Login takes a user and authenticates it
func (s *session) Login(db *gorm.DB, validatedList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.alreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		user := model.User{}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}
		defer r.Body.Close()

		// Validate it has no empty values
		err = user.Validate("login")
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Check if the email is validated
		err = validatedList.Seek(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, errors.New("Please verify your email before logging in"))
			return
		}

		// Authenticate user and get the jwt token of its id
		userID, err := auth.SignIn(db, user.Email, user.Password)
		if err != nil {
			response.HTMLText(w, r, http.StatusUnauthorized, "error: Invalid email or password")
			return
		}

		// Set a token of the user id in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "UID",
			Value:    userID,
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
			MaxAge:   s.length,
		})

		sID := uuid.New()
		cookie := &http.Cookie{
			Name:     "SID",
			Value:    sID.String(),
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
			MaxAge:   s.length,
		}
		http.SetCookie(w, cookie)

		s.store[cookie.Value] = userInfo{user.Email, time.Now()}

		response.HTMLText(w, r, http.StatusOK, "You logged in!")
	}
}

// Logout removes the authentication cookie
func (s *session) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("SID")

		if c == nil {
			response.Error(w, r, http.StatusBadRequest, errors.New("You cannot log out without a session"))
			return
		}

		// Delete map key equal to the cookieValue
		delete(s.store, c.Value)

		cookie := &http.Cookie{
			Name:     "SID",
			Value:    "0",
			Expires:  time.Unix(1414414788, 1414414788000),
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
			MaxAge:   -1,
		}

		http.SetCookie(w, cookie)

		http.SetCookie(w, &http.Cookie{
			Name:     "UID",
			Value:    "0",
			Expires:  time.Unix(1414414788, 1414414788000),
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
			MaxAge:   -1,
		})

		// Clean up session
		if time.Now().Sub(s.clean) > (time.Second * 30) {
			go s.sessionClean()
		}

		response.HTMLText(w, r, http.StatusOK, "You are now logged out.")
	}
}

// alreadyLoggedIn checks if the user have previously logged in
func (s *session) alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
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

// sessionClean deletes all the sessions that have expired
func (s *session) sessionClean() {
	for key, value := range s.store {
		if time.Now().Sub(value.lastSeen) > (time.Hour * 120) {
			delete(s.store, key)
		}
	}
	s.clean = time.Now()
}
