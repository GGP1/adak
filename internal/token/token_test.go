package token_test

import (
	"testing"

	"github.com/GGP1/adak/internal/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	t.Run("None signing method is not allowed", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
			"test": "none method",
			"must": "fail",
		})

		tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err, "SignedString()")

		_, err = jwt.DecodeSegment(tokenStr)
		assert.Error(t, err, "Expected invalid signing method error")
	})
}

func TestGenerateString(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := token.RandString(10)
		u := token.RandString(10)

		assert.NotEqual(t, s, u)
	}
}

func TestCheckPermits(t *testing.T) {
	id := "checkPermitsTest"

	jwt, err := token.GenerateFixedJWT(id)
	assert.NoError(t, err, "Failed generating fixed JWT")

	err = token.CheckPermits(id, jwt)
	assert.NoError(t, err, "Failed checking permits")

	jwt += "fail"
	err = token.CheckPermits(id, jwt)
	assert.Error(t, err, "Expected an error and got nil")
}
