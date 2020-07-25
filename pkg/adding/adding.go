/*
Package adding includes database adding operations
*/
package adding

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
)

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// AddProduct takes a new product and appends it to the database
func (s *service) AddProduct(db *gorm.DB, product *model.Product) error {
	err := product.Validate()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := db.Create(product).Error; err != nil {
		return fmt.Errorf("couldn't create the product")
	}

	return nil
}

// AddReview takes a new review and appends it to the database
func (s *service) AddReview(db *gorm.DB, review *model.Review) error {
	if err := db.Create(review).Error; err != nil {
		return fmt.Errorf("couldn't create the review")
	}

	return nil
}

// AddShop takes a new shop and appends it to the database
func (s *service) AddShop(db *gorm.DB, shop *model.Shop) error {
	err := shop.Validate()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := db.Create(shop).Error; err != nil {
		return fmt.Errorf("couldn't create the shop")
	}

	return nil
}

// AddUser takes a new user, hashes its password, sends
// a verification email and appends it to the database
func (s *service) AddUser(db *gorm.DB, user *model.User) error {
	err := user.Validate("")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	rowsAffected := db.Where("email = ?", user.Email).First(&user).RowsAffected
	if rowsAffected != 0 {
		return fmt.Errorf("email is already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create a cart for each user
	id := uuid.New()
	user.CartID = id.String()

	cart := shopping.NewCart(user.CartID)

	if err := db.Create(cart).Error; err != nil {
		return fmt.Errorf("couldn't create the cart")
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("couldn't create the user")
	}

	return nil
}
