package token

import (
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GenerateRunes generates a random string.
func GenerateRunes(length int) string {
	b := make([]rune, length)

	for i := range b {
		// Don't handle error as it is always greater than 0
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		b[i] = letterRunes[n.Int64()]
	}

	return string(b)
}

// CheckPermits cheks if the user is trying to perform and action on his own
// account or not. If true, return an error.
func CheckPermits(paramID, cookieID string) error {
	userID, err := ParseFixedJWT(cookieID)
	if err != nil {
		return err
	}

	if userID != paramID {
		return errors.New("it is not allowed to perform this action on third party accounts")
	}

	return nil
}
