package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Text colors
const (
	// Set 1
	// reset = "\033[0m"
	// red    = "\033[31m"
	// green  = "\033[32m"
	// yellow = "\033[33m"
	// blue   = "\033[34m"
	// purple = "\033[35m"
	// cyan   = "\033[36m"
	// white  = "\033[37m"

	// Set 2
	green  = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	blue   = "\033[97;44m"
	purple = "\033[97;45m"
	cyan   = "\033[97;46m"
	reset  = "\033[0m"
)

// LoggingResponseWriter helps us intercept the response status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Create a new logging response writer
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader intercepts write header input (status code) and store it in our
// loggingResponseWriter struct to use it later
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LogFormatter prints the server requests on the console.
// It uses different colors depending on the request status
// and the method called.
func LogFormatter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		timestamp := time.Now().Format("2006/01/02 15:04:05")
		method := r.Method
		status := lrw.statusCode
		path := r.URL.Path
		latency := time.Since(start)

		statusColor := statusCodeColor(status)
		methodColor := methodColor(method)
		resetColor := resetColor()

		log := fmt.Sprintf("[PALO] %v | %s %3d %s | %s %-7s %s | %#v | %v",
			timestamp,
			statusColor, status, resetColor,
			methodColor, method, resetColor,
			path,
			latency)

		fmt.Println(log)
	})
}

func statusCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return green
	case http.MethodPost:
		return blue
	case http.MethodPut:
		return cyan
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return yellow
	case http.MethodHead:
		return purple
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func resetColor() string {
	return reset
}
