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
	*sqlx.DB
	user.Service
}

// TODO: use cache

// AdminsOnly requires the user to be an administrator to proceed.
func (a *Auth) AdminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionID, err := cookie.Get(r, "SID")
		if err != nil {
			response.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		id := strings.Split(sessionID.Value, ":")[0]
		us, err := a.Service.GetByID(ctx, id)
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

func (a *Auth) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionID, err := cookie.Get(r, "SID")
		if err != nil {
			response.Error(w, http.StatusUnauthorized, errors.New("please log in to access"))
			return
		}

		// sID = id:username:salt
		sID := strings.Split(sessionID.Value, ":")

		us, err := a.Service.GetByID(ctx, sID[0])
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		if us.Username != sID[1] {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
