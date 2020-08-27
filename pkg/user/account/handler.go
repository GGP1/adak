package account

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/email"
	"github.com/GGP1/palo/pkg/user"

	"github.com/pkg/errors"

	"github.com/go-chi/chi"
)

// Handler handles account endpoints.
type Handler struct {
	Service Service
}

type changeEmail struct {
	Email string `json:"email"`
}

// ChangeEmail changes the user email to the specified one.
func (h *Handler) ChangeEmail(validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		email := chi.URLParam(r, "email")
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := h.Service.ChangeEmail(ctx, id, email, token, validatedList); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully changed your email!")
	}
}

type changePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword updates the user password.
func (h *Handler) ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, _ := r.Cookie("UID")

		var (
			changePass changePassword
			ctx        = r.Context()
		)

		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if err := json.NewDecoder(r.Body).Decode(&changePass); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.Service.ChangePassword(ctx, userID, changePass.OldPassword, changePass.NewPassword); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Password changed successfully.")
	}
}

// SendChangeConfirmation takes the new email and sends an email confirmation.
func (h *Handler) SendChangeConfirmation(u user.Service, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			new changeEmail
			ctx = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&new); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := email.Validate(new.Email); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		exists := validatedList.Exists(ctx, new.Email)
		if exists {
			response.Error(w, r, http.StatusBadRequest, errors.New("email is already taken"))
		}

		uID, _ := r.Cookie("UID")
		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		user, err := u.GetByID(ctx, userID)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		token, err := token.GenerateJWT(user.Email)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(err, "could not generate the jwt token"))
			return
		}

		errCh := make(chan error)

		go email.SendChangeConfirmation(user.ID, user.Username, user.Email, token, new.Email, errCh)

		select {
		case <-errCh:
			response.Error(w, r, http.StatusInternalServerError, errors.Wrap(<-errCh, "failed sending confirmation email"))
			return
		default:
			response.HTMLText(w, r, http.StatusOK, "We sent you an email to confirm that it is you.")
		}
	}
}

// SendEmailValidation saves the user email into the validated list.
// Once in the validated list, the user is able to log in.
func (h *Handler) SendEmailValidation(pendingList, validatedList email.Emailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		var (
			validated bool
			ctx       = r.Context()
		)

		pList, err := pendingList.Read(ctx)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		for _, v := range pList {
			if v.Token == token {
				if err := validatedList.Add(ctx, v.Email, v.Token); err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
					return
				}

				if err := pendingList.Remove(ctx, v.Email); err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
					return
				}

				validated = true
			}
		}

		if !validated {
			response.Error(w, r, http.StatusInternalServerError, errors.New("email validation failed"))
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully validated your email!")
	}
}
