package listing

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage.
type Repository interface {
	GetProducts(db *gorm.DB, product *[]model.Product) error
	GetProductByID(db *gorm.DB, product *model.Product, id string) error

	GetReviews(db *gorm.DB, review *[]model.Review) error
	GetReviewByID(db *gorm.DB, review *model.Review, id string) error

	GetShops(db *gorm.DB, shop *[]model.Shop) error
	GetShopByID(db *gorm.DB, shop *model.Shop, id string) error

	GetUsers(db *gorm.DB, user *[]model.User) error
	GetUserByID(db *gorm.DB, user *model.User, id string) error
}

// Service provides models listing operations.
type Service interface {
	GetProducts(db *gorm.DB, product *[]model.Product) error
	GetProductByID(db *gorm.DB, product *model.Product, id string) error

	GetReviews(db *gorm.DB, review *[]model.Review) error
	GetReviewByID(db *gorm.DB, review *model.Review, id string) error

	GetShops(db *gorm.DB, shop *[]model.Shop) error
	GetShopByID(db *gorm.DB, shop *model.Shop, id string) error

	GetUsers(db *gorm.DB, user *[]model.User) error
	GetUserByID(db *gorm.DB, user *model.User, id string) error
}
