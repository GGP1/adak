package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// LoggingResponseWriter helps us intercept the response status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Create a new logging response writer.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader intercepts write header input (status code) and store it in our
// loggingResponseWriter struct to use it later.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LogFormatter prints the server requests on the console.
func LogFormatter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Each request must have a id with it
		n := rand.Int63()
		reqID := strconv.Itoa(int(n))

		w.Header().Set("X-Request-ID", reqID)

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		timestamp := time.Now().Format("15:04:05 02/01/2006")
		method := r.Method
		status := lrw.statusCode
		path := r.URL.Path
		latency := time.Since(start)

		log := fmt.Sprintf("[PALO] %v | %3d | %-10v | %-7s | %#v",
			timestamp,
			status,
			latency,
			method,
			path)

		fmt.Println(log)
	})
}
