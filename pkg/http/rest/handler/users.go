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

	"github.com/gorilla/mux"
)

// GetUsers lists all the users
func GetUsers(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user []model.User

		err := l.GetUsers(&user)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// GetUserByID lists the user with the id requested
func GetUserByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]

		err := l.GetUserByID(&user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// AddUser creates a new user and saves it
func AddUser(a adding.Service, pendingList *email.PendingList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := a.AddUser(&user)
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

// UpdateUser updates the user with the given id
func UpdateUser(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := u.UpdateUser(&user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// DeleteUser deletes a user
func DeleteUser(d deleting.Service, pendingList *email.PendingList, validatedList *email.ValidatedList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		id := mux.Vars(r)["id"]

		err := d.DeleteUser(&user, id)
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
