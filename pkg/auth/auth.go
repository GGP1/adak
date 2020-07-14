/*
Package auth provides authentication and authorization support.
*/
package auth

import (
	"strconv"

	"github.com/GGP1/palo/pkg/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
)

// SignIn authenticates users and returns a jwt token
func SignIn(db *gorm.DB, email, password string) (string, error) {
	user := model.User{}

	err := db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	// Convert user id to string and generate a jwt token
	id := strconv.Itoa(int(user.ID))
	userID, err := GenerateFixedJWT(id)
	if err != nil {
		return "", err
	}

	return userID, nil
}
