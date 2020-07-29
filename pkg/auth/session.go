package auth

import "net/http"

// Repository provides access to the auth storage.
type Repository interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Login(w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	SessionClean()
}

// Session provides auth operations.
type Session interface {
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Login(w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	SessionClean()
}
