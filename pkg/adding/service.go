package adding

import (
	stg "palo/pkg/storage"
)

// Service provides user and product adding operations.
type Service interface {
	AddUser(User) error
	AddProduct(Product) error
}

// AddUser returns a new user and appends it to the database
func AddUser(user *User) (err error) {
	if err = stg.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// AddProduct returns a product and appends it to the database
func AddProduct(product *Product) (err error) {
	if err = stg.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}
