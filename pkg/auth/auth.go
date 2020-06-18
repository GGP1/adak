/*
Package auth is used for authenticating users functions
*/
package auth

import (
	"github.com/GGP1/palo/internal/utils/database"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
)

// SignIn authenticates users and returns a jwt token
func SignIn(email, password string) (string, error) {
	user := model.User{}

	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	err = db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}

	err = security.ComparePasswords(user.Password, password)
	if err != nil {
		return "", err
	}

	token, err := GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
