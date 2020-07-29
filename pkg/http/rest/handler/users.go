package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// Add creates a new user and saves it.
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
			return
		}

		err = pendingList.Add(user.Email, token)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = email.SendValidation(user, token)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed sending validation email: %w", err))
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Your account was successfully created.\nPlease validate your email to start using Palo.")
	}
}

// Delete removes a user.
func (us *Users) Delete(d deleting.Service, s auth.Session, pendingList, validatedList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]
		cookie, _ := r.Cookie("UID")
		// Generate a fixed token of the id and compare it with the cookie
		// to check if it's the same user
		userID, err := auth.GenerateFixedJWT(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if userID != cookie.Value {
			response.Error(w, r, http.StatusUnauthorized, fmt.Errorf("not allowed to delete others user"))
			return
		}

		// Remove user from email lists
		err = pendingList.Remove(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = validatedList.Remove(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = d.DeleteUser(us.DB, &user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// If the user is logged in, log him out
		s.Logout(w, r, cookie)

		response.HTMLText(w, r, http.StatusOK, "User deleted successfully.")
	}
}

// Find lists all the users
func (us *Users) Find(l listing.Service) http.HandlerFunc {
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

// FindByID lists the user with the id requested.
func (us *Users) FindByID(l listing.Service) http.HandlerFunc {
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

// Update updates the user with the given id.
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
