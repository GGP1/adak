// Package middleware provides http services.
package middleware

import (
	"net/http"
)

// AdminsOnly checks if the user is an admin or not.
func AdminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aID, err := r.Cookie("AID")
		if err != nil {
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		if aID.Value == "" {
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
