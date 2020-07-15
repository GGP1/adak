package adding

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage
type Repository interface {
	AddProduct(db *gorm.DB, product *model.Product) error
	AddReview(db *gorm.DB, review *model.Review) error
	AddShop(db *gorm.DB, shop *model.Shop) error
	AddUser(db *gorm.DB, user *model.User) error
}

// Service provides models adding operations.
type Service interface {
	AddProduct(db *gorm.DB, product *model.Product) error
	AddReview(db *gorm.DB, review *model.Review) error
	AddShop(db *gorm.DB, shop *model.Shop) error
	AddUser(db *gorm.DB, user *model.User) error
}
