package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/pkg/user"

	"github.com/jmoiron/sqlx"
)

// Auth contains the elements needed to authorize users.
type Auth struct {
	DB          *sqlx.DB
	UserService user.Service
}

// AdminsOnly requires the user to be an administrator to proceed.
func (a *Auth) AdminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionID, err := cookie.GetValue(r, "SID")
		if err != nil {
			response.Error(w, http.StatusForbidden, errors.New("Unauthorized"))
			return
		}

		id := strings.Split(sessionID, ":")[0]

		us, err := a.UserService.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		if !us.IsAdmin {
			// Return 404 instead of 401 to not give additional information
			response.Error(w, http.StatusNotFound, errors.New("Not Found"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireLogin makes sure the user is logged in before forwarding the request,
// it returns an error otherwise.
func (a *Auth) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionID, err := cookie.GetValue(r, "SID")
		if err != nil {
			response.Error(w, http.StatusForbidden, errors.New("please log in to access"))
			return
		}

		// sID = id:username:salt
		sID := strings.Split(sessionID, ":")
		id := sID[0]
		username := sID[1]

		// TODO: Use redis to save sessions

		us, err := a.UserService.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		if us.Username != username {
			response.Error(w, http.StatusNotFound, errors.New("Not Found"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
