package cookie

import (
	"encoding/hex"
	"net/http"

	"github.com/GGP1/adak/internal/crypt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var secretKey = []byte(viper.GetString("token.secretkey"))

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

// Get deciphers and returns the cookie value.
func Get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", errors.Wrapf(err, "failed retrieving %s cookie", name)
	}

	ciphertext, err := hex.DecodeString(cookie.Value)
	if err != nil {
		return "", errors.Wrap(err, "decoding cookie value")
	}

	plaintext, err := crypt.Decrypt(secretKey, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Set a cookie.
func Set(w http.ResponseWriter, name, value, path string, age int) error {
	ciphertext, err := crypt.Encrypt(secretKey, []byte(value))
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    hex.EncodeToString(ciphertext),
		Path:     path,
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   age,
	})

	return nil
}
