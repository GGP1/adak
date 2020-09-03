package auth

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"

	"github.com/pkg/errors"
)

// Login takes a user credentials and authenticates it.
func Login(s Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.AlreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var (
			user User
			ctx  = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := user.Validate(); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.Login(ctx, w, user.Email, user.Password); err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You logged in!")
	}
}

// Logout logs the user out from the session and removes cookies.
func Logout(s Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("SID")

		if c == nil {
			response.Error(w, r, http.StatusBadRequest, errors.New("error: you cannot log out without a session"))
			return
		}

		// Logout user from the session and delete cookies
		s.Logout(w, r, c)

		response.HTMLText(w, r, http.StatusOK, "You are now logged out.")
	}
}
