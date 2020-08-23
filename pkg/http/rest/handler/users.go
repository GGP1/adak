package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/creating"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/searching"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
)

// CreateUser creates a new user and saves it.
func CreateUser(c creating.Service, pendingList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user model.User
			ctx  = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		token, err := auth.GenerateJWT(user)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("could not generate the jwt token: %w", err))
			return
		}

		if err := pendingList.Add(ctx, user.Email, token); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		errCh := make(chan error)
		go email.SendValidation(ctx, user, token, errCh)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		select {
		case <-ctx.Done():
			response.Error(w, r, http.StatusInternalServerError, ctx.Err())
		case <-errCh:
			response.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed sending validation email: %w", <-errCh))
		default:
			if err = c.CreateUser(ctx, &user); err != nil {
				response.Error(w, r, http.StatusBadRequest, err)
			}

			response.HTMLText(w, r, http.StatusCreated, "Your account was successfully created.\nPlease validate your email to start using Palo.")
		}
	}
}

// DeleteUser removes a user.
func DeleteUser(d deleting.Service, s auth.Session, db *sqlx.DB, pendingList, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")

		var (
			user model.User
			ctx  = r.Context()
		)

		// Check if it's the same user
		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = db.Get(&user, "SELECT * FROM users WHERE id=$1", userID)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if userID != id {
			response.Error(w, r, http.StatusUnauthorized, errors.New("not allowed to delete others user"))
			return
		}

		if err := validatedList.Remove(ctx, user.Email); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := shopping.DeleteCart(ctx, db, user.CartID); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := d.DeleteUser(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// If the user is logged in, log him out
		s.Logout(w, r, uID)

		response.HTMLText(w, r, http.StatusOK, "User deleted successfully.")
	}
}

// GetUsers lists all the users
func GetUsers(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		users, err := l.GetUsers(ctx)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, users)
	}
}

// GetUserByID lists the user with the id requested.
func GetUserByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		user, err := l.GetUserByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// QRCode shows the user id in a qrcode format.
func QRCode(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var (
			user model.User
			ctx  = r.Context()
		)

		user, err := l.GetUserByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		img, err := user.QRCode()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.PNG(w, r, http.StatusOK, img)
	}
}

// SearchUser looks for the products with the given value.
func SearchUser(s searching.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")

		ctx := r.Context()

		users, err := s.SearchUsers(ctx, query)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, users)
	}
}

// UpdateUser updates the user with the given id.
func UpdateUser(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")

		var (
			user model.User
			ctx  = r.Context()
		)

		// Check if it's the same user
		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if userID != id {
			response.Error(w, r, http.StatusUnauthorized, errors.New("not allowed to update others user"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := u.UpdateUser(ctx, &user, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "User updated successfully.")
	}
}
