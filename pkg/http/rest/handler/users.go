package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"

	"github.com/gorilla/mux"
)

// GetUsers lists all the users
func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user []model.User

		err := listing.GetAll(&user)

		if err != nil {
			response.Text(w, r, http.StatusInternalServerError, err)
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// GetOneUser lists the user with the id requested
func GetOneUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetOne(&user, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "User not found")
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// AddUser creates a new user and saves it
func AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Text(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.AddUser(&user)
		if err != nil {
			response.Text(w, r, http.StatusInternalServerError, err)
		}

		response.JSON(w, r, http.StatusOK, user)
		io.WriteString(w, "Please validate your email")
	}
}

// UpdateUser updates the user with the given id
func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Text(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := updating.UpdateUser(&user, id)
		if err != nil {
			response.Text(w, r, http.StatusInternalServerError, err)
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// DeleteUser deletes a user
func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.Delete(&user, id)
		if err != nil {
			response.Text(w, r, http.StatusInternalServerError, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "User deleted")
	}
}
