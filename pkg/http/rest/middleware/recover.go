package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GGP1/adak/internal/response"
)

// Recover recovers any panic (from libraries) and handles the error to prevent the server from shutting down.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				response.Error(w, http.StatusInternalServerError, errors.New(fmt.Sprint(err)))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
