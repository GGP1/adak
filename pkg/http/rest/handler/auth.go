/*
Package handler contains the methods used by the router
*/
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/utils/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/model"

	"github.com/google/uuid"
)

// Login takes a user and authenticates it
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := model.User{}

		// Check if cookie already exists, if not, create it
		_, err := r.Cookie("SID")
		if err != nil {
			id := uuid.New()

			http.SetCookie(w, &http.Cookie{
				Name:     "SID",
				Value:    id.String(),
				Path:     "/",
				Domain:   "localhost",
				Secure:   false,
				HttpOnly: true,
			})
		}

		// Decode request body
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.Respond(w, r, http.StatusUnauthorized, err)
			return
		}

		// Validate it has no empty values
		err = user.Validate("login")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, err.Error())
			return
		}

		// Authenticate user
		token, err := auth.SignIn(user.Email, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Invalid email or password")
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, token)
	}
}

// Logout removes the authentication cookie
func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "SID",
			Value:    "0",
			Expires:  time.Unix(1414414788, 1414414788000),
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "You are now logged out")
	}
}
