/*
Package auth is used for authenticating users functions
*/
package auth

import (
	"github.com/GGP1/palo/internal/utils/cfg"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// SignIn authenticates users and returns a jwt token
func SignIn(email, password string) (string, error) {
	user := model.User{}

	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return "err", nil
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
