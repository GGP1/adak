package deleting

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage.
type Repository interface {
	DeleteProduct(db *gorm.DB, product *model.Product, id string) error
	DeleteReview(db *gorm.DB, review *model.Review, id string) error
	DeleteShop(db *gorm.DB, shop *model.Shop, id string) error
	DeleteUser(db *gorm.DB, user *model.User, id string) error
}

// Service provides models deleting operations.
type Service interface {
	DeleteProduct(db *gorm.DB, product *model.Product, id string) error
	DeleteReview(db *gorm.DB, review *model.Review, id string) error
	DeleteShop(db *gorm.DB, shop *model.Shop, id string) error
	DeleteUser(db *gorm.DB, user *model.User, id string) error
}
