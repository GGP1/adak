package auth

import (
	"time"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"
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
