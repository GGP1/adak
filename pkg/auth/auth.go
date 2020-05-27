/*
Package auth is used for authenticating users functions
*/
package auth

import (
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/models"
	"github.com/GGP1/palo/pkg/storage"
)

// SignIn logs a user in
func SignIn(email, password string) (string, error) {
	var err error
	user := models.User{}

	err = storage.DB.Where("email = ?", email).Take(&user).Error
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
