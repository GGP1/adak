package middleware

import (
	"net/http"
)

// RequireLogin verifies if the user is logged in.
func RequireLogin(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("SID")
		if err != nil {
			http.Error(w, "Please log in to access.", http.StatusUnauthorized)
			return
		}

		f(w, r)
	})
}
