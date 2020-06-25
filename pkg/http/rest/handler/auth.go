/*
Package handler contains the methods used by the router
*/
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/model"

	"github.com/google/uuid"
)

// AuthRepository provides access to the auth storage
type AuthRepository interface {
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
	alreadyLoggedIn(http.ResponseWriter, *http.Request) bool
}

// Session provides auth operations
type Session interface {
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
	alreadyLoggedIn(http.ResponseWriter, *http.Request) bool
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
func (s *session) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := model.User{}

		// Check if the user is already logged in or not
		if s.alreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Decode request body
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}
		defer r.Body.Close()

		// Validate it has no empty values
		err = user.Validate("login")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		// Authenticate user
		_, err = auth.SignIn(user.Email, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Invalid email or password")
			return
		}

		// Create uuid and cookie
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
		// Map store - cookieValue: userInfo
		s.store[cookie.Value] = userInfo{user.Email, time.Now()}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "You are now logged in.")
	}
}

// Logout removes the authentication cookie
func (s *session) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("SID")

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

		// Clean up session
		if time.Now().Sub(s.clean) > (time.Second * 30) {
			go s.sessionClean()
		}

		response.Text(w, r, http.StatusOK, "You are now logged out.")
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
	http.SetCookie(w, cookie)
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
