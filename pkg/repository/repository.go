package repository

import (
	"net/http"

	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
)

// Repo contains all services methods.
type Repo interface {
	// Creating
	CreateProduct(db *gorm.DB, product *model.Product) error
	CreateReview(db *gorm.DB, review *model.Review) error
	CreateShop(db *gorm.DB, shop *model.Shop) error
	CreateUser(db *gorm.DB, user *model.User) error

	// Deleting
	DeleteProduct(db *gorm.DB, id string) error
	DeleteReview(db *gorm.DB, id string) error
	DeleteShop(db *gorm.DB, id string) error
	DeleteUser(db *gorm.DB, id string) error

	// Listing
	GetProducts(db *gorm.DB, product *[]model.Product) error
	GetProductByID(db *gorm.DB, product *model.Product, id string) error
	GetReviews(db *gorm.DB, review *[]model.Review) error
	GetReviewByID(db *gorm.DB, review *model.Review, id string) error
	GetShops(db *gorm.DB, shop *[]model.Shop) error
	GetShopByID(db *gorm.DB, shop *model.Shop, id string) error
	GetUsers(db *gorm.DB, user *[]model.User) error
	GetUserByID(db *gorm.DB, user *model.User, id string) error

	// Searching
	SearchProducts(db *gorm.DB, products *[]model.Product, search string) error
	SearchShops(db *gorm.DB, shops *[]model.Shop, search string) error
	SearchUsers(db *gorm.DB, users *[]model.User, search string) error

	// Updating
	UpdateProduct(db *gorm.DB, product *model.Product, id string) error
	UpdateShop(db *gorm.DB, shop *model.Shop, id string) error
	UpdateUser(db *gorm.DB, user *model.User, id string) error

	// Auth session
	AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool
	Clean()
	EmailChange(id, newEmail, token string, validatedList email.Service) error
	Login(w http.ResponseWriter, email, password string) error
	Logout(w http.ResponseWriter, r *http.Request, c *http.Cookie)
	PasswordChange(id, oldPass, newPass string) error

	// Email
	Add(email, token string) error
	Read() (emailList map[string]string, err error)
	Remove(key string) error
	Seek(email string) error
}
