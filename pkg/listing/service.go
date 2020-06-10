/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Service provides models listing operations.
type Service interface {
	GetUsers() error
	GetAUser() error
	GetProducts() error
	GetAProduct() error
	GetReviews() error
	GetAReview() error
	GetShops() error
	GetAShop() error
}

// GetUsers returns all the users in the database
func GetUsers(user *[]model.User, db *gorm.DB) error {
	if err := db.Find(user).Error; err != nil {
		return err
	}
	return nil
}

// GetAUser returns a single user
func GetAUser(user *model.User, id string, db *gorm.DB) error {
	if err := db.First(user, id).Error; err != nil {
		return err
	}
	return nil
}

// GetProducts returns all the products in the database
func GetProducts(product *[]model.Product, db *gorm.DB) error {
	if err := db.Find(product).Error; err != nil {
		return err
	}
	return nil
}

// GetAProduct returns a single product
func GetAProduct(product *model.Product, id string, db *gorm.DB) error {
	if err := db.First(product, id).Error; err != nil {
		return err
	}
	return nil
}

// GetReviews returns all the reviews in the database
func GetReviews(review *[]model.Review, db *gorm.DB) error {
	if err := db.Find(review).Error; err != nil {
		return err
	}
	return nil
}

// GetAReview returns a single review
func GetAReview(review *model.Review, id string, db *gorm.DB) error {
	if err := db.First(review, id).Error; err != nil {
		return err
	}
	return nil
}

// GetShops returns all the shops in the database
func GetShops(shop *[]model.Shop, db *gorm.DB) error {
	if err := db.Find(shop).Error; err != nil {
		return err
	}
	return nil
}

// GetAShop returns a single shop
func GetAShop(shop *model.Shop, id string, db *gorm.DB) error {
	if err := db.First(shop, id).Error; err != nil {
		return err
	}
	return nil
}
