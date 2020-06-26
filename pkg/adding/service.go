/*
Package adding includes database adding operations
*/
package adding

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Service provides models adding operations.
type Service interface {
	AddProduct(*model.Product) error
	AddReview(*model.Review) error
	AddShop(*model.Shop) error
	AddUser(*model.User) error
}

// AddProduct takes a new product and appends it to the database
func AddProduct(product *model.Product) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Create(product).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the product")
	}

	return nil
}

// AddReview takes a new review and appends it to the database
func AddReview(review *model.Review) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Create(review).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the review")
	}

	return nil
}

// AddShop takes a new shop and appends it to the database
func AddShop(shop *model.Shop) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Create(shop).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the shop")
	}

	return nil
}

// AddUser takes a new user, hashes its password, sends
// a verification email and appends it to the database
func AddUser(user *model.User) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = user.Validate("login")
	if err != nil {
		return err
	}

	// Hash password
	hash, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create user
	if err := db.Create(user).Error; err != nil {
		return errors.Wrap(err, "couldn't create the user")
	}

	// Send confirmation email to the user
	email.Confirmation(user)

	return nil
}
