/*
Package adding includes database adding operations
*/
package adding

import (
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides access to the storage
type Repository interface {
	AddProduct(*model.Product) error
	AddReview(*model.Review) error
	AddShop(*model.Shop) error
	AddUser(*model.User) error
}

// Service provides models adding operations.
type Service interface {
	AddProduct(*model.Product) error
	AddReview(*model.Review) error
	AddShop(*model.Shop) error
	AddUser(*model.User) error
}

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// AddProduct takes a new product and appends it to the database
func (s *service) AddProduct(product *model.Product) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Create(product).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the product")
	}

	return nil
}

// AddReview takes a new review and appends it to the database
func (s *service) AddReview(review *model.Review) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Create(review).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the review")
	}

	return nil
}

// AddShop takes a new shop and appends it to the database
func (s *service) AddShop(shop *model.Shop) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Create(shop).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the shop")
	}

	return nil
}

// AddUser takes a new user, hashes its password, sends
// a verification email and appends it to the database
func (s *service) AddUser(user *model.User) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = user.Validate("login")
	if err != nil {
		return err
	}

	existingUser := db.Where("email = ?", user.Email).First(&user).Value
	if existingUser == user.Email {
		return errors.New("error: the email is already used")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	if err := db.Create(user).Error; err != nil {
		return errors.New("error: the email is already used")
	}

	return nil
}
