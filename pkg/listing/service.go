/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/pkg/models"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides user and product listing operations.
type Service interface {
	GetUsers(*[]models.User) error
	GetAUser(*models.User, string) error
	GetProducts(*[]models.Product) error
	GetAProduct(*models.Product, string) error
	GetReviews(*[]models.Review) error
	GetAReview(*models.Review, string) error
	GetShops(*[]models.Shop) error
	GetAShop(*models.Shop, string) error
}

// GetUsers returns all the users in the database
func GetUsers(user *[]models.User) (err error) {
	if err = stg.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

// GetAUser returns a single user
func GetAUser(user *models.User, id string) (err error) {
	if err = stg.DB.First(user, id).Error; err != nil {
		return err
	}
	return nil
}

// GetProducts returns all the products in the database
func GetProducts(product *[]models.Product) (err error) {
	if err = stg.DB.Find(product).Error; err != nil {
		return err
	}
	return nil
}

// GetAProduct returns a single product
func GetAProduct(product *models.Product, id string) (err error) {
	if err = stg.DB.First(product, id).Error; err != nil {
		return err
	}
	return nil
}

// GetReviews returns all the reviews in the database
func GetReviews(review *[]models.Review) (err error) {
	if err = stg.DB.Find(review).Error; err != nil {
		return err
	}
	return nil
}

// GetAReview returns a single review
func GetAReview(review *models.Review, id string) (err error) {
	if err = stg.DB.First(review, id).Error; err != nil {
		return err
	}
	return nil
}

// GetShops returns all the shops in the database
func GetShops(shop *[]models.Shop) (err error) {
	if err = stg.DB.Find(shop).Error; err != nil {
		return err
	}
	return nil
}

// GetAShop returns a single shop
func GetAShop(shop *models.Shop, id string) (err error) {
	if err = stg.DB.First(shop, id).Error; err != nil {
		return err
	}
	return nil
}
