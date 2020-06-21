package middleware

import (
	"net/http"
)

// IsLoggedIn verifies if the user is logged in
func IsLoggedIn(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/", "/verify", "/users", "/shops", "/products", "/reviews"}
		requestPath := r.URL.Path

		// Iterate through the paths allowed to access without logging in and
		// compare them to the current url, if true, pass to the next middleware
		for _, value := range notAuth {
			if value == requestPath {
				f(w, r)
				return
			}
		}

		cookie, err := r.Cookie("SID")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		if cookie == nil {
			http.Error(w, "Please log in to access", http.StatusUnauthorized)
		}

		f(w, r)
	})
}
