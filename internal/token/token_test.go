package token_test

import (
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/GGP1/adak/internal/crypt"
	"github.com/GGP1/adak/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestGenerateString(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := token.RandString(10)
		u := token.RandString(10)

		assert.NotEqual(t, s, u)
	}
}

func TestCheckPermits(t *testing.T) {
	id := "checkPermitsTest"
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	ciphertext, err := crypt.Encrypt([]byte(id))
	assert.NoError(t, err)

	r.AddCookie(&http.Cookie{
		Name:  "UID",
		Value: hex.EncodeToString(ciphertext),
		Path:  "/",
	})

	t.Run("Success", func(t *testing.T) {
		err = token.CheckPermits(r, id)
		assert.NoError(t, err, "Failed checking permits")
	})

	t.Run("Permission denied", func(t *testing.T) {
		err := token.CheckPermits(r, id+"fail")
		assert.Error(t, err)
	})

	t.Run("ID too long", func(t *testing.T) {
		err := token.CheckPermits(r, "9 }NkbKPLja;As[0<|d4nMG!5l3>x$+Qp")
		assert.Error(t, err)
	})
}
