package rest

import (
	"net/http"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/pkg/user/account"

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
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		_, err := s.accountClient.ChangeEmail(ctx, &account.ChangeEmailRequest{ID: id, NewEmail: email, Token: token})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "You have successfully changed your email!")
	}
}
