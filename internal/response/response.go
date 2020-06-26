package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTMLText is the protocol function for html text resposes
func HTMLText(w http.ResponseWriter, r *http.Request, status int, text string) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	w.WriteHeader(status)

	fmt.Fprintln(w, text)
}

// JSON is the protocol function for JSON responses
func JSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)

	if _, err := io.Copy(w, &buf); err != nil {
		fmt.Println("Respond: ", err)
	}
}

// Error is the protocol function for error resposes
func Error(w http.ResponseWriter, r *http.Request, status int, err error) {
	// Set content type, statusCode and write the error
	http.Error(w, err.Error(), status)
}
