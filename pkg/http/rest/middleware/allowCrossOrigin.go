package middleware

import (
	"net/http"
)

// AllowCrossOrigin enables foreign requests.
func AllowCrossOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow localhost on port 3000 to send or receive data from our server
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4001")

		next.ServeHTTP(w, r)
	})
}
