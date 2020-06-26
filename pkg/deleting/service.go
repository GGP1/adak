/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Service provides models deleting operations.
type Service interface {
	DeleteProduct(*model.Product, string) error
	DeleteReview(*model.Review, string) error
	DeleteShop(*model.Shop, string) error
	DeleteUser(*model.User, string) error
}

// DeleteProduct takes a product from the database and permanently deletes it
func DeleteProduct(product *model.Product, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(product, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the product")
	}
	return nil
}

// DeleteReview takes a review from the database and permanently deletes it
func DeleteReview(review *model.Review, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(review, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the review")
	}
	return nil
}

// DeleteShop takes a shop from the database and permanently deletes it
func DeleteShop(shop *model.Shop, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(shop, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the shop")
	}
	return nil
}

// DeleteUser takes a user from the database and permanently deletes it
func DeleteUser(user *model.User, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(user, id).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the user")
	}
	return nil
}
