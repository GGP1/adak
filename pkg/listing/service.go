/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Service provides models listing operations.
type Service interface {
	GetProducts(*[]model.Product) error
	GetOneProduct(*model.Product, string) error
	GetReviews(*[]model.Review) error
	GetOneReview(*model.Review, string) error
	GetShops(*[]model.Shop) error
	GetOneShop(*model.Shop, string) error
	GetUsers(*[]model.User) error
	GetOneUser(*model.User, string) error
}

// GetProducts lists all the products stored in the database
func GetProducts(product *[]model.Product) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").Find(product).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetOneProduct lists the product requested from the database
func GetOneProduct(product *model.Product, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").First(product, id).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetReviews lists all the reviews stored in the database
func GetReviews(review *[]model.Review) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Find(review).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetOneReview lists the review requested from the database
func GetOneReview(review *model.Review, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.First(review, id).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetShops lists all the shops stored in the database
func GetShops(shop *[]model.Shop) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Location").Preload("Reviews").Preload("Products").Find(shop).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetOneShop lists the shop requested from the database
func GetOneShop(shop *model.Shop, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Location").Preload("Products").First(shop, id).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetUsers lists all the users stored in the database
func GetUsers(user *[]model.User) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").Find(user).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetOneUser lists the user requested from the database
func GetOneUser(user *model.User, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Preload("Reviews").First(user, id).Error
	if err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}
