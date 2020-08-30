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
)

// Error is the function used to send error resposes.
func Error(w http.ResponseWriter, r *http.Request, status int, err error) {
	e := fmt.Sprintf("status: %d\nerror: %v", status, err)
	// Set content type, statusCode and write the error
	http.Error(w, e, status)
}

// HTMLText is the function used to send html text resposes.
func HTMLText(w http.ResponseWriter, r *http.Request, status int, text string) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	w.WriteHeader(status)

	fmt.Fprintln(w, text)
}

// JSON is the function used to send JSON responses.
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

// PNG is used to respond with a png image.
func PNG(w http.ResponseWriter, r *http.Request, status int, img image.Image) {
	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		http.Error(w, "couldn't encode the PNG image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))

	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, "unable to write image", http.StatusInternalServerError)
	}
}
