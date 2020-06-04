/*
Package adding includes database adding operations
*/
package adding

import (
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/models"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides models adding operations.
type Service interface {
	AddUser(*models.User) error
	AddProduct(*models.Product) error
	AddReview(*models.Review) error
	AddShop(*models.Shop) error
}

// AddUser returns a new user and appends it to the database
func AddUser(user *models.User) (err error) {
	// Hash password
	hash, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create user
	if err = stg.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// AddProduct returns a product and appends it to the database
func AddProduct(product *models.Product) (err error) {
	if err = stg.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}

// AddReview returns a review and appends it to the database
func AddReview(review *models.Review) (err error) {
	if err = stg.DB.Create(review).Error; err != nil {
		return err
	}
	return nil
}

// AddShop returns a shop and appends it to the database
func AddShop(shop *models.Shop) (err error) {
	if err = stg.DB.Create(shop).Error; err != nil {
		return err
	}
	return nil
}
