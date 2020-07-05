/*
Package auth provides authentication and authorization support.
*/
package auth

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// SignIn authenticates users and returns a jwt token
func SignIn(email, password string) error {
	user := model.User{}

	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}

	err = db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return err
	}

	err = security.ComparePasswords(user.Password, password)
	if err != nil {
		return err
	}

	return nil
}
