package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/utils/env"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// GenerateJWT creates a new jwt token
func GenerateJWT(user model.User) (string, error) {
	// Load env file
	env.LoadEnv()

	key := []byte(os.Getenv("SECRET_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       user,
		"expiration": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(key)
}

// ExtractToken retrieves the token from headers as a query
func ExtractToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	// Load env variables
	env.LoadEnv()

	key := []byte(os.Getenv("SECRET_KEY"))

	token, err := request.ParseFromRequestWithClaims(
		r,
		request.OAuth2Extractor,
		jwt.MapClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vError := err.(*jwt.ValidationError)
			switch vError.Errors {
			case jwt.ValidationErrorExpired:
				err = errors.New("Your token has expired")
				w.WriteHeader(http.StatusUnauthorized)
				return nil, err
			case jwt.ValidationErrorSignatureInvalid:
				err = errors.New("The signature is invalid")
				w.WriteHeader(http.StatusUnauthorized)
				return nil, err
			default:
				return nil, err
			}
		}
	}

	return token, nil
}
