package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

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
func (us *Users) Add(a adding.Service, pendingList *email.PendingList) http.HandlerFunc {
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
			fmt.Printf("%v", err)
		}

		// Add user mail and token to the email pending confirmation list
		pendingList.Add(user, token)

		// Send validation email
		email.SendValidation(user, token)

		response.JSON(w, r, http.StatusOK, user)
		fmt.Fprintln(w, "Please validate your email")
	}
}

// Update updates the user with the given id
func (us *Users) Update(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := u.UpdateUser(us.DB, &user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// Delete deletes a user
func (us *Users) Delete(d deleting.Service, pendingList *email.PendingList, validatedList *email.ValidatedList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]

		err := d.DeleteUser(us.DB, &user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Remove user from email lists
		pendingList.Remove(user.Email)
		validatedList.Remove(user.Email)

		// If the user is logged in, log him out
		c, _ := r.Cookie("SID")
		if c != nil {
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
		}

		response.HTMLText(w, r, http.StatusOK, "User deleted successfully.")
	}
}
