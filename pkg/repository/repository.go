package repository

import (
	"net/http"

	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// MonoRepo contains the functions of all the services
type MonoRepo interface {
	// Adding
	AddProduct(db *gorm.DB, product *model.Product) error
	AddReview(db *gorm.DB, review *model.Review) error
	AddShop(db *gorm.DB, shop *model.Shop) error
	AddUser(db *gorm.DB, user *model.User) error

	// Deleting
	DeleteProduct(db *gorm.DB, product *model.Product, id string) error
	DeleteReview(db *gorm.DB, review *model.Review, id string) error
	DeleteShop(db *gorm.DB, shop *model.Shop, id string) error
	DeleteUser(db *gorm.DB, user *model.User, id string) error

	// Listing
	GetProducts(db *gorm.DB, product *[]model.Product) error
	GetProductByID(db *gorm.DB, product *model.Product, id string) error
	GetReviews(db *gorm.DB, review *[]model.Review) error
	GetReviewByID(db *gorm.DB, review *model.Review, id string) error
	GetShops(db *gorm.DB, shop *[]model.Shop) error
	GetShopByID(db *gorm.DB, shop *model.Shop, id string) error
	GetUsers(db *gorm.DB, user *[]model.User) error
	GetUserByID(db *gorm.DB, user *model.User, id string) error

	// Updating
	UpdateProduct(db *gorm.DB, product *model.Product, id string) error
	UpdateShop(db *gorm.DB, shop *model.Shop, id string) error
	UpdateUser(db *gorm.DB, user *model.User, id string) error

	// Auth session
	Login(db *gorm.DB, validatedList email.Service) http.HandlerFunc
	Logout() http.HandlerFunc

	// Email
	Add(email, token string) error
	Read() (emailList map[string]string, err error)
	Remove(key string) error
	Seek(email string) error
}
