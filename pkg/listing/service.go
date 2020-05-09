package listing

import (
	stg "palo/pkg/storage"
)

// Service provides user and product listing operations.
type Service interface {
	GetUsers([]User) error
	GetAUser(User, string) error
	GetProducts([]Product) error
	GetAProduct(Product, string) error
	GetReviews([]Review, string) error
	GetAReview(Review, string) error
}

// GetUsers returns all the users in the database
func GetUsers(user *[]User) (err error) {
	if err = stg.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

// GetAUser returns a single user
func GetAUser(user *User, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

// GetProducts returns all the products in the database
func GetProducts(product *[]Product) (err error) {
	if err = stg.DB.Find(product).Error; err != nil {
		return err
	}
	return nil
}

// GetAProduct returns a single product
func GetAProduct(product *Product, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(product).Error; err != nil {
		return err
	}
	return nil
}

// GetReviews returns all the products in the database
func GetReviews(review *[]Review) (err error) {
	if err = stg.DB.Find(review).Error; err != nil {
		return err
	}
	return nil
}

// GetAReview returns a single product
func GetAReview(review *Review, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(review).Error; err != nil {
		return err
	}
	return nil
}
