// Package handler contains the handlers used by the router
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"

	"github.com/go-chi/chi"
)

type changeEmail struct {
	Email string `json:"email"`
}

type changePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// EmailChange takes the new email and sends an email confirmation.
func EmailChange(validatedList email.Emailer, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			new changeEmail
		)

		if err := json.NewDecoder(r.Body).Decode(&new); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := model.ValidateEmail(new.Email); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		exists := validatedList.Exists(new.Email)
		if exists {
			response.Error(w, r, http.StatusBadRequest, errors.New("email is already taken"))
		}

		uID, _ := r.Cookie("UID")
		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		user, err := l.GetUserByID(userID)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		token, err := auth.GenerateJWT(user)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("could not generate the jwt token: %w", err))
			return
		}

		errCh := make(chan error)

		go email.SendChangeConfirmation(user, token, new.Email, errCh)

		select {
		case <-errCh:
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed sending confirmation email: %w", <-errCh))
			return
		default:
			response.HTMLText(w, r, http.StatusOK, "We sent you an email to confirm that it is you.")
		}
	}
}

// EmailChangeConfirmation changes the user email to the specified one.
func EmailChangeConfirmation(s auth.Session, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		email := chi.URLParam(r, "email")
		id := chi.URLParam(r, "id")

		if err := s.EmailChange(id, email, token, validatedList); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully changed your email!")
	}
}

// Login takes a user credentials and authenticates it.
func Login(s auth.Session, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.AlreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := user.Validate("login"); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Check if the email is validated
		if err := validatedList.Seek(user.Email); err != nil {
			response.Error(w, r, http.StatusUnauthorized, fmt.Errorf("please verify your email before logging in: %w", err))
			return
		}

		if err := s.Login(w, user.Email, user.Password); err != nil {
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
			response.Error(w, r, http.StatusBadRequest, errors.New("error: you cannot log out without a session"))
			return
		}

		// Logout user from the session and delete cookies
		s.Logout(w, r, c)

		response.HTMLText(w, r, http.StatusOK, "You are now logged out.")
	}
}

// PasswordChange updates the user password.
func PasswordChange(s auth.Session, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var changePass changePassword

		uID, _ := r.Cookie("UID")

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if err := json.NewDecoder(r.Body).Decode(&changePass); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := s.PasswordChange(userID, changePass.OldPassword, changePass.NewPassword); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Password changed successfully.")
	}
}

// ValidateEmail saves the user email into the validated list.
// Once in the validated list, the user is able to log in.
func ValidateEmail(pendingList, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		var validated bool

		pList, err := pendingList.Read()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		for _, v := range pList {
			if v.Token == token {
				if err := validatedList.Add(v.Email, v.Token); err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
					return
				}

				if err := pendingList.Remove(v.Email); err != nil {
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
