package middleware

import "net/http"

// Secure adds security headers to the http connection.
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-XSS-Protection: stops a page from loading when it detects XSS attacks
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		// HTTP Strict Transport Security:
		// lets a web site tell browsers that it should only be accessed using HTTPS, instead of using HTTP
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		// X-Frame-Options:
		// indicate whether or not a browser should be allowed to render a page in a <frame>, <iframe>, <embed> or <object>
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		// X-Content-Type-Options:
		// is a marker used by the server to indicate that the MIME types advertised in the Content-Type headers
		// should not be changed and be followed
		w.Header().Add("X-Content-Type-Options", "nosniff")
		// Content Security Policy: allows web site administrators to control resources the user agent is allowed to load for a given page
		w.Header().Add("Content-Security-Policy", "default-src 'self';")
		// X-Permitted-Cross-Domain-Policies: allow other systems to access the domain
		w.Header().Add("X-Permitted-Cross-Domain-Policies", "none")
		// Referrer-Policy: sets the parameter for amount of information sent along with Referer Header while making a request
		w.Header().Add("Referrer-Policy", "no-referrer")
		// Feature-Policy: provides a mechanism to allow and deny the use of browser features in its own frame,
		// and in content within any <iframe> elements in the document
		w.Header().Add("Feature-Policy", "microphone 'none'; camera 'none'")

		next.ServeHTTP(w, r)
	})
}
