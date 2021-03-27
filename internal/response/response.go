package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"github.com/GGP1/adak/internal/logger"

	"github.com/pkg/errors"
)

type msgResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type errResponse struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

// Error is the function used to send error resposes.
func Error(w http.ResponseWriter, status int, err error) {
	res := errResponse{
		Status: status,
		Err:    err.Error(),
	}

	JSON(w, status, res)
}

// HTMLText is the function used to send html text resposes.
func HTMLText(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(status)

	if _, err := fmt.Fprintln(w, text); err != nil {
		Error(w, http.StatusInternalServerError, err)
	}
}

// JSON is the function used to send JSON responses.
func JSON(w http.ResponseWriter, status int, v interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		logger.Log.Fatalf("failed writing json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	if _, err := io.Copy(w, &buf); err != nil {
		logger.Log.Fatalf("failed writing to response writer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// JSONText is the function used to send JSON formatted text responses.
func JSONText(w http.ResponseWriter, status int, message string) {
	res := msgResponse{
		Message: message,
		Code:    status,
	}

	JSON(w, status, res)
}

// PNG is used to respond with a png image.
func PNG(w http.ResponseWriter, status int, img image.Image) {
	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		Error(w, http.StatusInternalServerError, errors.New("couldn't encode the PNG image"))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(status)

	if _, err := io.Copy(w, &buf); err != nil {
		Error(w, http.StatusInternalServerError, errors.Wrap(err, "couldn't write to response writer"))
	}
}
