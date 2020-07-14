package auth

import (
	"os"
	"time"

	"github.com/GGP1/palo/internal/env"
	"github.com/GGP1/palo/pkg/model"
	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT creates a new jwt token - changes over time -
func GenerateJWT(user model.User) (string, error) {
	env.Load()

	key := []byte(os.Getenv("SECRET_KEY"))

	claim := jwt.MapClaims{
		"user": user.Email,
		"exp":  time.Now().Add(time.Minute * 27).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(key)
}

// GenerateFixedJWT creates a jwt token that does not vary
func GenerateFixedJWT(id string) (string, error) {
	env.Load()

	key := []byte(os.Getenv("SECRET_KEY"))

	claim := jwt.MapClaims{
		"id": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(key)
}
