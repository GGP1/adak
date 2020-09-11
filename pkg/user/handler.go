package user

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/shopping/cart"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Handler handles user endpoints.
type Handler struct {
	Service Service
}

// Create creates a new user and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user AddUser
			ctx  = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		confirmationCode := token.GenerateRunes(20)
		errCh := make(chan error)

		go email.SendValidation(ctx, user.Username, user.Email, confirmationCode, errCh)

		select {
		case <-ctx.Done():
			response.Error(w, r, http.StatusInternalServerError, ctx.Err())
		case <-errCh:
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(<-errCh, "failed sending validation email"))
		default:
			if err := h.Service.Create(ctx, &user); err != nil {
				response.Error(w, r, http.StatusBadRequest, err)
				return
			}

			response.HTMLText(w, r, http.StatusCreated, "Your account was successfully created.\nWe've sent you an email to validate your account.")
		}
	}
}

// Delete removes a user.
func (h *Handler) Delete(db *sqlx.DB, s auth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")

		var (
			ctx = r.Context()
		)

		if err := token.CheckPermits(id, uID.Value); err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		user, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
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

		if err := token.CheckPermits(id, uID.Value); err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, user); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.Service.Update(ctx, &user, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "User updated successfully.")
	}
}
