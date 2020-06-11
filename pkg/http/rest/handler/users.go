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
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// UserHandler defines all of the handlers related to users. It holds the
// application state needed by the handler methods.
type UserHandler struct {
	DB *gorm.DB
}

// Get lists all the users
func (uh *UserHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user []model.User

		err := listing.GetUsers(&user, uh.DB)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// GetOne lists the user with the id requested
func (uh *UserHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAUser(&user, id, uh.DB)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "User not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// Add creates a new user and saves it
func (uh *UserHandler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.AddUser(&user, uh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.RespondJSON(w, r, http.StatusOK, user)
		io.WriteString(w, "Please validate your email")
	}
}

// Update updates the user with the given id
func (uh *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := updating.UpdateUser(&user, id, uh.DB)

		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		response.RespondJSON(w, r, http.StatusOK, user)
	}
}

// Delete deletes a user
func (uh *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteUser(&user, id, uh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "User deleted")
	}
}
