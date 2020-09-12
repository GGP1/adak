package token

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var secretKey = []byte(viper.GetString("token.secretkey"))

// GenerateJWT creates a new jwt token - changes over time -.
func GenerateJWT(email string) (string, error) {
	key := secretKey

	claim := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 27).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(key)
}

// GenerateFixedJWT creates a jwt token that does not vary.
func GenerateFixedJWT(id string) (string, error) {
	key := secretKey

	claim := jwt.MapClaims{
		"id": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(key)
}

// ParseFixedJWT takes the claims from the token and returns the id value.
// This function is used to take the user id value from the UID cookie.
func ParseFixedJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}

		return secretKey, nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failed parsing the token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("jwt: claim type assertion failed")
	}

	return claims["id"].(string), nil
}
