package token

import (
	"time"

	"github.com/GGP1/adak/internal/logger"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var secretKey = []byte(viper.GetString("token.secretkey"))

// GenerateJWT creates a new jwt token - changes over time -.
func GenerateJWT(email string) (string, error) {
	claim := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(secretKey)
}

// GenerateFixedJWT creates a jwt token that does not vary.
func GenerateFixedJWT(id string) (string, error) {
	claim := jwt.MapClaims{
		"id": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(secretKey)
}

// GetUserID takes the claims from the token and returns the id value. The tokenString is usually the UID cookie.
func GetUserID(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.Errorf("Invalid token %v", token.Header["alg"])
		}

		return secretKey, nil
	})
	if err != nil {
		logger.Log.Errorf("failed parsing the token: %v", err)
		return "", errors.Wrap(err, "failed parsing the token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Log.Error("invalid jwt claims")
		return "", errors.New("invalid jwt claims")
	}

	return claims["id"].(string), nil
}
