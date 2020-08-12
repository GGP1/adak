package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	green  = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	blue   = "\033[97;44m"
	purple = "\033[97;45m"
	cyan   = "\033[97;46m"
	reset  = "\033[0m"
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
		n := rand.Int()
		reqID := strconv.Itoa(n)

		w.Header().Set("X-Request-ID", reqID)

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		timestamp := time.Now().Format("15:04:05 02/01/2006")
		method := r.Method
		status := lrw.statusCode
		path := r.URL.Path
		latency := time.Since(start)

		statusColor := statusCodeColor(status)
		methodColor := methodColor(method)
		resetColor := resetColor()

		log := fmt.Sprintf("[PALO] %v | %s %3d %s | %-10v | %s %-7s %s | %#v",
			timestamp,
			statusColor, status, resetColor,
			latency,
			methodColor, method, resetColor,
			path)

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

func methodColor(reqMethod string) string {
	method := reqMethod

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
