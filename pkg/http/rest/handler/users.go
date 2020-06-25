package handler

import (
	"encoding/json"
	"fmt"
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
			response.Error(w, r, http.StatusInternalServerError, err)
			return
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
			response.Error(w, r, http.StatusNotFound, err)
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
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer r.Body.Close()

		err := adding.AddUser(&user)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
		fmt.Fprintln(w, "Please validate your email")
	}
}

// UpdateUser updates the user with the given id
func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer r.Body.Close()

		err := updating.UpdateUser(&user, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
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
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.Text(w, r, http.StatusOK, "User deleted successfully.")
	}
}
