package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/utils/response"
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

		err := listing.GetUsers(&user)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// GetOneUser lists one user based on the id
func GetOneUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAUser(&user, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "User not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// AddUser creates a new user and saves it
func AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.AddUser(&user)
		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// UpdateUser updates a user
func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := updating.UpdateUser(&user, id)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// DeleteUser deletes a user
func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteUser(&user, id)
		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "User deleted")
	}
}
