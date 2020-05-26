/*
Package adding provides the methods needed to create objects and append them to the database
*/
package adding

import (
	"github.com/GGP1/palo/pkg/auth/security"
	"github.com/GGP1/palo/pkg/model"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides user and product adding operations.
type Service interface {
	AddUser(*model.User) error
	AddProduct(*model.Product) error
	AddReview(*model.Review) error
	AddShop(*model.Shop) error
}

// AddUser returns a new user and appends it to the database
func AddUser(user *model.User) (err error) {
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
func AddProduct(product *model.Product) (err error) {
	if err = stg.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}

// AddReview returns a review and appends it to the database
func AddReview(review *model.Review) (err error) {
	if err = stg.DB.Create(review).Error; err != nil {
		return err
	}
	return nil
}

// AddShop returns a shop and appends it to the database
func AddShop(shop *model.Shop) (err error) {
	if err = stg.DB.Create(shop).Error; err != nil {
		return err
	}
	return nil
}
