/*
Package auth is used for authenticating users functions
*/
package auth

import (
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// SignIn authenticates users and returns a jwt token
func SignIn(email, password string, db *gorm.DB) (string, error) {
	var err error
	user := model.User{}

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
