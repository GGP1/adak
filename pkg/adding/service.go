package adding

import (
	"palo/pkg/models"
	stg "palo/pkg/storage"
)

// Adder provides user and product adding operations.
type Adder interface {
	AddUser(models.User) error
	AddProduct(models.Product) error
	AddReview(models.Review) error
}

// AddUser returns a new user and appends it to the database
func AddUser(user *models.User) (err error) {
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

// AddReview returns a product and appends it to the database
func AddReview(review *models.Review) (err error) {
	if err = stg.DB.Create(review).Error; err != nil {
		return err
	}
	return nil
}
