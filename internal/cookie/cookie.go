package cookie

import (
	"encoding/hex"
	"net/http"

	"github.com/GGP1/adak/internal/crypt"

	"github.com/pkg/errors"
)

// Delete a cookie.
func Delete(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "0",
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// Get deciphers and returns the cookie.
func Get(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}

	ciphertext, err := hex.DecodeString(cookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "decoding cookie value")
	}

	plaintext, err := crypt.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	cookie.Value = string(plaintext)

	return cookie, nil
}

// GetValue is like Get but returns only the value.
func GetValue(r *http.Request, name string) (string, error) {
	cookie, err := Get(r, name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

// IsSet returns whether the cookie is set or not.
func IsSet(r *http.Request, name string) bool {
	c, _ := r.Cookie(name)

	return c != nil
}

// Set a cookie.
func Set(w http.ResponseWriter, name, value, path string, age int) error {
	ciphertext, err := crypt.Encrypt([]byte(value))
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    hex.EncodeToString(ciphertext),
		Path:     path,
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true, // True means no scripts, http requests only. It does not refer to http(s)
		SameSite: http.SameSiteStrictMode,
		MaxAge:   age,
	})

	return nil
}
