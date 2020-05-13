package listing

import (
	"github.com/GGP1/palo/pkg/model"
	stg "github.com/GGP1/palo/pkg/storage"
)

// Service provides user and product listing operations.
type Service interface {
	GetUsers([]model.User) error
	GetAUser(model.User, string) error
	GetProducts([]model.Product) error
	GetAProduct(model.Product, string) error
	GetReviews([]model.Review) error
	GetAReview(model.Review, string) error
}

// GetUsers returns all the users in the database
func GetUsers(user *[]model.User) (err error) {
	if err = stg.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

// GetAUser returns a single user
func GetAUser(user *model.User, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

// GetProducts returns all the products in the database
func GetProducts(product *[]model.Product) (err error) {
	if err = stg.DB.Find(product).Error; err != nil {
		return err
	}
	return nil
}

// GetAProduct returns a single product
func GetAProduct(product *model.Product, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(product).Error; err != nil {
		return err
	}
	return nil
}

// GetReviews returns all the products in the database
func GetReviews(review *[]model.Review) (err error) {
	if err = stg.DB.Find(review).Error; err != nil {
		return err
	}
	return nil
}

// GetAReview returns a single product
func GetAReview(review *model.Review, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(review).Error; err != nil {
		return err
	}
	return nil
}
