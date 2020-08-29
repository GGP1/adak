package middleware

import "net/http"

// Secure adds security headers to the http connection.
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-XSS-Protection
		w.Header().Add("X-XSS-Protection", "1; mode=block")

		// HTTP Strict Transport Security
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Frame-Options
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")

		// X-Content-Type-Options
		w.Header().Add("X-Content-Type-Options", "nosniff")

		// Content Security Policy
		w.Header().Add("Content-Security-Policy", "default-src 'self';")

		// X-Permitted-Cross-Domain-Policies
		w.Header().Add("X-Permitted-Cross-Domain-Policies", "none")

		// Referrer-Policy
		w.Header().Add("Referrer-Policy", "no-referrer")

		// Feature-Policy
		w.Header().Add("Feature-Policy", "microphone 'none'; camera 'none'")

		next.ServeHTTP(w, r)
	})
}
