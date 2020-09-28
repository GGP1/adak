package token

import (
	"math/rand"

	"github.com/pkg/errors"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GenerateRunes generates a random string.
func GenerateRunes(length int) string {
	b := make([]rune, length)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
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
