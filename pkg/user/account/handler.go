package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/email"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/pkg/user"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// Handler handles account endpoints.
type Handler struct {
	accountService Service
	userService    user.Service
	emailer        email.Emailer
}

type changeEmail struct {
	Email string `json:"email"`
}

// NewHandler returns a new account handler.
func NewHandler(accountS Service, userS user.Service, emailer email.Emailer) Handler {
	return Handler{
		accountService: accountS,
		userService:    userS,
		emailer:        emailer,
	}
}

// ChangeEmail changes the user email to the specified one.
func (h *Handler) ChangeEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := chi.URLParam(r, "token")
		email := chi.URLParam(r, "email")

		id, err := params.URLID(ctx)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		if err := h.accountService.ChangeEmail(ctx, id, email, token); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("email changed to %q", email))
	}
}

type changePassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"email,required"`
}

// ChangePassword updates the user password.
func (h *Handler) ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var changePass changePassword
		ctx := r.Context()

		userID, err := cookie.GetValue(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&changePass); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.accountService.ChangePassword(ctx, userID, changePass.OldPassword, changePass.NewPassword); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, "successfully changed password")
	}
}

// SendChangeConfirmation takes the new email and sends an email confirmation.
func (h *Handler) SendChangeConfirmation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var new changeEmail
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&new); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if _, err := h.userService.GetByEmail(ctx, new.Email); err == nil {
			response.Error(w, http.StatusBadRequest, errors.New("email is already taken"))
			return
		}

		userID, err := cookie.GetValue(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		user, err := h.userService.GetByID(ctx, userID)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		token := uuid.NewString()
		if err := h.emailer.SendChangeConfirmation(user.ID, user.Username, user.Email, token, new.Email); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, "verification email sent")
	}
}

// SendEmailValidation saves the user email into the validated list.
// Once in the validated list, the user is able to log in.
func (h *Handler) SendEmailValidation(u user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		token := chi.URLParam(r, "token")
		ctx := r.Context()

		usr, err := u.GetByEmail(ctx, email)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		if err := h.accountService.ValidateUserEmail(ctx, usr.ID, token, true); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("validated %q", email))
	}
}
