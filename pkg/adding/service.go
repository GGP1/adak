/*
Package adding includes database adding operations
*/
package adding

import (
	"github.com/GGP1/palo/internal/utils/email"
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Service provides models adding operations.
type Service interface {
	AddUser(*model.User) error
	AddProduct(*model.Product) error
	AddReview(*model.Review) error
	AddShop(*model.Shop) error
}

// AddUser returns a new user and appends it to the database
func AddUser(user *model.User, db *gorm.DB) error {
	// Hash password
	hash, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create user
	if err := db.Create(user).Error; err != nil {
		return err
	}

	// Send confirmation email to the user
	email.Confirmation(user)

	return nil
}

// AddProduct returns a product and appends it to the database
func AddProduct(product *model.Product, db *gorm.DB) error {
	if err := db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

// AddReview returns a review and appends it to the database
func AddReview(review *model.Review, db *gorm.DB) error {
	if err := db.Create(review).Error; err != nil {
		return err
	}
	return nil
}

// AddShop returns a shop and appends it to the database
func AddShop(shop *model.Shop, db *gorm.DB) error {
	if err := db.Create(shop).Error; err != nil {
		return err
	}
	return nil
}
