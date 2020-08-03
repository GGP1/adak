// Package handler contains the handlers used by the router
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"
	"github.com/go-chi/chi"

	"github.com/pkg/errors"
)

// Login takes a user credentials and authenticates it.
func Login(s auth.Session, validatedList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.AlreadyLoggedIn(w, r) {
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
			response.Error(w, r, http.StatusUnauthorized, fmt.Errorf("please verify your email before logging in: %w", err))
			return
		}

		// Authenticate user
		err = s.Login(w, user.Email, user.Password)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You logged in!")
	}
}

// Logout logs the user out from the session and removes cookies.
func Logout(s auth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("SID")

		if c == nil {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("error: you cannot log out without a session"))
			return
		}

		// Logout user from the session and delete cookies
		s.Logout(w, r, c)

		response.HTMLText(w, r, http.StatusOK, "You are now logged out.")
	}
}

// ValidateEmail saves the user email into the validated list.
// Once in the validated list, the user is able to log in.
func ValidateEmail(pendingList, validatedList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var validated bool
		token := chi.URLParam(r, "token")

		pList, err := pendingList.Read()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		for k, v := range pList {
			if v == token {
				err := validatedList.Add(k, v)
				if err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
					return
				}
				validated = true
			}
		}

		if !validated {
			response.Error(w, r, http.StatusInternalServerError, errors.New("error: email validation failed"))
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully validated your email!")
	}
}
