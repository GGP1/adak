// Package handler contains the handlers used by the router
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

type changeEmail struct {
	Email string `json:"email"`
}

type changePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// EmailChange takes the new email and send an email confirmation.
func EmailChange(db *gorm.DB, validatedList email.Emailer, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var new changeEmail
		var user model.User

		err := json.NewDecoder(r.Body).Decode(&new)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer r.Body.Close()

		if err := user.ValidateEmail(user.Email); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		exists := validatedList.Exists(new.Email)
		if exists {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("email is already taken"))
		}

		uID, _ := r.Cookie("UID")
		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = l.GetUserByID(db, &user, userID)
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

		err := s.EmailChange(id, email, token, validatedList)
		if err != nil {
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

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}
		defer r.Body.Close()

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

// PasswordChange updates the user password.
func PasswordChange(s auth.Session, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var changePass changePassword

		uID, _ := r.Cookie("UID")

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		err = json.NewDecoder(r.Body).Decode(&changePass)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer r.Body.Close()

		err = s.PasswordChange(userID, changePass.OldPassword, changePass.NewPassword)
		if err != nil {
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

		for k, v := range pList {
			if v == token {
				err := validatedList.Add(k, v)
				if err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
					return
				}

				err = pendingList.Remove(k)
				if err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
					return
				}

				validated = true
			}
		}

		if !validated {
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("error: email validation failed"))
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully validated your email!")
	}
}
