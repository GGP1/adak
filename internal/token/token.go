package token

import (
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/GGP1/adak/internal/cookie"

	"github.com/pkg/errors"
)

var pool = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// RandString returns a random string.
func RandString(length int) string {
	b := make([]rune, length)

	for i := range b {
		// Don't handle error as len(pool) is always greater than 0
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		b[i] = pool[n.Int64()]
	}

	return string(b)
}

// CheckPermits cheks if the user is trying to perform and action on his own
// account (return nil) or not (return error).
func CheckPermits(r *http.Request, paramID string) error {
	userID, err := cookie.Get(r, "UID")
	if err != nil {
		return err
	}

	if userID.Value != paramID {
		return errors.New("it is not allowed to perform this action on third party accounts")
	}

	return nil
}
