package middleware

import (
	"net/http"
)

// AdminsOnly checks if the user is an admin or not
func AdminsOnly(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("AID")
		if err != nil {
			http.Error(w, "Restringed access.", http.StatusUnauthorized)
			return
		}

		f(w, r)
	})
}
