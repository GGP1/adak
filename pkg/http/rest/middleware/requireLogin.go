package middleware

import (
	"errors"
	"net/http"

	"github.com/GGP1/adak/internal/response"
)

// RequireLogin verifies if the user is logged in.
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("SID"); err != nil {
			response.Error(w, http.StatusUnauthorized, errors.New("please log in to access"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
