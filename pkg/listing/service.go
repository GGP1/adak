package listing

import (
	"palo/pkg/models"
	stg "palo/pkg/storage"
)

// Lister provides user and product listing operations.
type Lister interface {
	GetUsers([]models.User) error
	GetAUser(models.User, string) error
	GetProducts([]models.Product) error
	GetAProduct(models.Product, string) error
	GetReviews([]models.Review, string) error
	GetAReview(models.Review, string) error
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
	if err = stg.DB.Where("id = ?", id).First(user).Error; err != nil {
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
	if err = stg.DB.Where("id = ?", id).First(product).Error; err != nil {
		return err
	}
	return nil
}

// GetReviews returns all the products in the database
func GetReviews(review *[]models.Review) (err error) {
	if err = stg.DB.Find(review).Error; err != nil {
		return err
	}
	return nil
}

// GetAReview returns a single product
func GetAReview(review *models.Review, id string) (err error) {
	if err = stg.DB.Where("id = ?", id).First(review).Error; err != nil {
		return err
	}
	return nil
}
