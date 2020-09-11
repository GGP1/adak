// Package middleware provides http services.
package middleware

import (
	"net/http"
)

// AdminsOnly checks if the user is an admin or not.
func AdminsOnly(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adm := r.Header.Get("AID")
		aID, err := r.Cookie("AID")
		if err != nil {
			http.Error(w, "404 page not found", http.StatusUnauthorized)
			return
		}

		if adm == "" && aID.Value == "" {
			http.Error(w, "404 page not found", http.StatusUnauthorized)
			return
		}

		f(w, r)
	})
}
