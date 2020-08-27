package user

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/email"
	"github.com/GGP1/palo/pkg/shopping/cart"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Handler handles user endpoints.
type Handler struct {
	Service Service
}

// Create creates a new user and saves it.
func (h *Handler) Create(pendingList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user User
			ctx  = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		token, err := token.GenerateJWT(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(err, "could not generate the jwt token"))
			return
		}

		if err := pendingList.Add(ctx, user.Email, token); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		errCh := make(chan error)

		go email.SendValidation(ctx, user.Username, user.Email, token, errCh)

		select {
		case <-ctx.Done():
			response.Error(w, r, http.StatusInternalServerError, ctx.Err())
		case <-errCh:
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(<-errCh, "failed sending validation email"))
		default:
			if err = h.Service.Create(ctx, &user); err != nil {
				response.Error(w, r, http.StatusBadRequest, err)
			}

			response.HTMLText(w, r, http.StatusCreated, "Your account was successfully created.\nPlease validate your email to start using Palo.")
		}
	}
}

// Delete removes a user.
func (h *Handler) Delete(db *sqlx.DB, s auth.Session, pendingList, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")

		var (
			ctx = r.Context()
		)

		// Check if it's the same user
		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if userID != id {
			response.Error(w, r, http.StatusUnauthorized, errors.New("not allowed to delete others user"))
			return
		}

		user, err := h.Service.GetByID(ctx, userID)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		if err := validatedList.Remove(ctx, user.Email); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := cart.Delete(ctx, db, user.CartID); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := h.Service.Delete(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Logout(w, r, uID)

		response.HTMLText(w, r, http.StatusOK, "User deleted successfully.")
	}
}

// Get lists all the users.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		users, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, users)
	}
}

// GetByID lists the user with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		user, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, user)
	}
}

// QRCode shows the user id in a qrcode format.
func (h *Handler) QRCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var (
			user ListUser
			ctx  = r.Context()
		)

		user, err := h.Service.GetByID(ctx, id)
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

// Search looks for the products with the given value.
func (h *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")

		ctx := r.Context()

		users, err := h.Service.Search(ctx, query)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, users)
	}
}

// Update updates the user with the given id.
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")

		var (
			user UpdateUser
			ctx  = r.Context()
		)

		// Check if it's the same user
		userID, err := token.ParseFixedJWT(uID.Value)
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

		if err := h.Service.Update(ctx, &user, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "User updated successfully.")
	}
}
