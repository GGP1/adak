package rest

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/sanitize"
	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/user"
	"github.com/GGP1/palo/pkg/user/account"

	"github.com/pkg/errors"

	"github.com/go-chi/chi"
)

type changeEmail struct {
	Email string `json:"email"`
}

// AccountChangeEmail changes the user email to the specified one.
func (s *Frontend) AccountChangeEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		email := chi.URLParam(r, "email")
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := sanitize.Normalize(&email); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err := s.accountClient.ChangeEmail(ctx, &account.ChangeEmailRequest{ID: id, NewEmail: email, Token: token})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully changed your email!")
	}
}

type changePassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"email,required"`
}

// AccountChangePassword updates the user password.
func (s *Frontend) AccountChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var changePass changePassword
		uID, _ := r.Cookie("UID")
		ctx := r.Context()

		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if err := json.NewDecoder(r.Body).Decode(&changePass); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := sanitize.Normalize(&changePass.NewPassword); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err = s.accountClient.ChangePassword(ctx, &account.ChangePasswordRequest{
			ID:      userID,
			OldPass: changePass.OldPassword,
			NewPass: changePass.NewPassword,
		})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Password changed successfully.")
	}
}

// AccountSendChangeConfirmation takes the new email and sends an email confirmation.
func (s *Frontend) AccountSendChangeConfirmation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var new changeEmail
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&new); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := sanitize.Normalize(&new.Email); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err := s.userClient.GetByEmail(ctx, &user.GetByEmailRequest{Email: new.Email})
		if err == nil {
			response.Error(w, r, http.StatusBadRequest, errors.New("email is already taken"))
			return
		}

		uID, _ := r.Cookie("UID")
		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		u, err := s.userClient.GetByID(ctx, &user.GetByIDRequest{ID: userID})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		token, err := token.GenerateJWT(u.User.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(err, "could not generate the jwt token"))
			return
		}

		errCh := make(chan error, 1)
		go email.SendChangeConfirmation(u.User.ID, u.User.Username, u.User.Email, token, new.Email, errCh)

		select {
		case err := <-errCh:
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(err, "failed sending confirmation email"))
			return
		default:
			response.HTMLText(w, r, http.StatusOK, "We sent you an email to confirm that it is you.")
		}
	}
}

// AccountSendEmailValidation validates the user account.
func (s *Frontend) AccountSendEmailValidation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		token := chi.URLParam(r, "token")
		var ctx = r.Context()

		_, err := s.accountClient.ValidateEmail(ctx, &account.ValidateEmailRequest{
			Email:            email,
			ConfirmationCode: token,
			VerifiedEmail:    true,
		})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully validated your email!")
	}
}
