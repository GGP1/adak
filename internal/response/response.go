package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GGP1/adak/internal/bufferpool"

	"github.com/bradfitz/gomemcache/memcache"
)

type msgResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type errResponse struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

// EncodedJSON writes a response from a buffer with json encoded content.
//
// The status is predefined as 200 (OK).
func EncodedJSON(w http.ResponseWriter, buf []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(buf); err != nil {
		Error(w, http.StatusInternalServerError, err)
	}
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
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if _, err := io.Copy(w, buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// JSONAndCache works just like JSON but saves the encoding of v to the cache before writing the response.
//
// The status should always be 200 (OK).
func JSONAndCache(mc *memcache.Client, w http.ResponseWriter, key string, v interface{}) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	if err := mc.Set(&memcache.Item{Key: key, Value: buf.Bytes()}); err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(buf.Bytes()); err != nil {
		Error(w, http.StatusInternalServerError, err)
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
