/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/storage"
)

// Service provides models listing operations.
type Service interface {
	GetUsers(*[]model.User) error
	GetAUser(*model.User, string) error
	GetProducts(*[]model.Product) error
	GetAProduct(*model.Product, string) error
	GetReviews(*[]model.Review) error
	GetAReview(*model.Review, string) error
	GetShops(*[]model.Shop) error
	GetAShop(*model.Shop, string) error
}

// GetUsers returns all the users in the database
func GetUsers(user *[]model.User) error {
	if err := storage.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

// GetAUser returns a single user
func GetAUser(user *model.User, id string) error {
	if err := storage.DB.First(user, id).Error; err != nil {
		return err
	}
	return nil
}

// GetProducts returns all the products in the database
func GetProducts(product *[]model.Product) error {
	if err := storage.DB.Find(product).Error; err != nil {
		return err
	}
	return nil
}

// GetAProduct returns a single product
func GetAProduct(product *model.Product, id string) error {
	if err := storage.DB.First(product, id).Error; err != nil {
		return err
	}
	return nil
}

// GetReviews returns all the reviews in the database
func GetReviews(review *[]model.Review) error {
	if err := storage.DB.Find(review).Error; err != nil {
		return err
	}
	return nil
}

// GetAReview returns a single review
func GetAReview(review *model.Review, id string) error {
	if err := storage.DB.First(review, id).Error; err != nil {
		return err
	}
	return nil
}

// GetShops returns all the shops in the database
func GetShops(shop *[]model.Shop) error {
	if err := storage.DB.Find(shop).Error; err != nil {
		return err
	}
	return nil
}

// GetAShop returns a single shop
func GetAShop(shop *model.Shop, id string) error {
	if err := storage.DB.First(shop, id).Error; err != nil {
		return err
	}
	return nil
}
