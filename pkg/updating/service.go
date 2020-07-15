package updating

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage
type Repository interface {
	UpdateProduct(db *gorm.DB, product *model.Product, id string) error
	UpdateShop(db *gorm.DB, shop *model.Shop, id string) error
	UpdateUser(db *gorm.DB, user *model.User, id string) error
}

// Service provides models updating operations.
type Service interface {
	UpdateProduct(db *gorm.DB, product *model.Product, id string) error
	UpdateShop(db *gorm.DB, shop *model.Shop, id string) error
	UpdateUser(db *gorm.DB, user *model.User, id string) error
}
