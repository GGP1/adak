// Package middleware provides http services.
package middleware

import (
	"errors"
	"net/http"

	"github.com/GGP1/adak/internal/response"
)

// AdminsOnly checks if the user is an admin or not.
func AdminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aID, err := r.Cookie("AID")
		if err != nil {
			response.Error(w, http.StatusUnauthorized, errors.New("401 unauthorized"))
			return
		}

		if aID.Value == "" {
			response.Error(w, http.StatusUnauthorized, errors.New("401 unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
