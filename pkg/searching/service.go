package searching

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repository provides access to the storage.
type Repository interface {
	SearchProducts(db *gorm.DB, products *[]model.Product, search string) error
	SearchShops(db *gorm.DB, users *[]model.Shop, search string) error
	SearchUsers(db *gorm.DB, users *[]model.User, search string) error
}

// Service provides models searching operations.
type Service interface {
	SearchProducts(db *gorm.DB, products *[]model.Product, search string) error
	SearchShops(db *gorm.DB, users *[]model.Shop, search string) error
	SearchUsers(db *gorm.DB, users *[]model.User, search string) error
}
