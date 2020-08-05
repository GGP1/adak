package auth

import (
	"fmt"
	"time"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT creates a new jwt token - changes over time -.
func GenerateJWT(user model.User) (string, error) {
	key := []byte(cfg.SecretKey)

	claim := jwt.MapClaims{
		"user": user.Email,
		"exp":  time.Now().Add(time.Minute * 27).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(key)
}

// GenerateFixedJWT creates a jwt token that does not vary.
func GenerateFixedJWT(id string) (string, error) {
	key := []byte(cfg.SecretKey)

	claim := jwt.MapClaims{
		"id": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(key)
}

// ParseFixedJWT takes the claims from the token and returns the id value.
func ParseFixedJWT(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}

		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed parsing the token")
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims["id"], nil
}
