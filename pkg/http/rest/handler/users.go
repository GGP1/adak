package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// Users handles users routes
type Users struct {
	DB *gorm.DB
}

// GetAll lists all the users
func (us *Users) GetAll(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user []model.User

		err := l.GetUsers(us.DB, &user)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// GetByID lists the user with the id requested
func (us *Users) GetByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]

		err := l.GetUserByID(us.DB, &user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// Add creates a new user and saves it
func (us *Users) Add(a adding.Service, pendingList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := a.AddUser(us.DB, &user)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := auth.GenerateJWT(user)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("could not generate a jwt token: %w", err))
		}

		// Add user mail and token to the email pending confirmation list
		pendingList.Add(user.Email, token)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		// Send validation email
		err = email.SendValidation(user, token)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed sending validation email: %w", err))
		}

		fmt.Fprintln(w, "You account was successfully created. Please validate your email to start using Palo.")
	}
}

// Update updates the user with the given id
func (us *Users) Update(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]
		c, _ := r.Cookie("UID")
		// Generate a fixed token of the id and compare it with the cookie
		// to check if it's the same user
		userID, err := auth.GenerateFixedJWT(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if userID != c.Value {
			response.Error(w, r, http.StatusUnauthorized, fmt.Errorf("not allowed to update others user"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err = u.UpdateUser(us.DB, &user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// Delete deletes a user
func (us *Users) Delete(d deleting.Service, pendingList, validatedList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]
		c, _ := r.Cookie("UID")
		// Generate a fixed token of the id and compare it with the cookie
		// to check if it's the same user
		userID, err := auth.GenerateFixedJWT(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if userID != c.Value {
			response.Error(w, r, http.StatusUnauthorized, fmt.Errorf("not allowed to delete others user"))
			return
		}

		err = d.DeleteUser(us.DB, &user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Remove user from email lists
		pendingList.Remove(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		validatedList.Remove(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		// If the user is logged in, log him out
		c2, _ := r.Cookie("SID")
		if c2 != nil {
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
		}

		response.HTMLText(w, r, http.StatusOK, "User deleted successfully.")
	}
}
