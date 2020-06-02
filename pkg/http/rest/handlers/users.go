package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/utils/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/models"
	"github.com/GGP1/palo/pkg/updating"

	"github.com/gorilla/mux"
)

// GetUsers lists all the users
func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user []models.User

		err := listing.GetUsers(&user)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.Respond(w, r, http.StatusOK, user)
	}
}

// GetOneUser lists one user based on the id
func GetOneUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAUser(&user, id)

		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.Respond(w, r, http.StatusOK, user)
	}
}

// AddUser creates a new user and saves it
func AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var buf bytes.Buffer
		var err error

		err = json.NewEncoder(&buf).Encode(&user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Review not found")
		}

		err = adding.AddUser(&user)
		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.Respond(w, r, http.StatusOK, user)
	}
}

// UpdateUser updates a user
func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		param := mux.Vars(r)
		id := param["id"]

		err := updating.UpdateUser(&user, id)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.Respond(w, r, http.StatusOK, user)
	}
}

// DeleteUser deletes a user
func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteUser(&user, id)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.Respond(w, r, http.StatusOK, user)
	}
}
