/*
Package creating includes database creating operations
*/
package creating

import (
	"errors"
	"fmt"
	"time"

	"github.com/GGP1/palo/internal/random"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Repository provides access to the storage.
type Repository interface {
	CreateProduct(db *sqlx.DB, product *model.Product) error
	CreateReview(db *sqlx.DB, r *model.Review, userID string) error
	CreateShop(db *sqlx.DB, shop *model.Shop) error
	CreateUser(db *sqlx.DB, user *model.User) error
}

// Service provides models adding operations.
type Service interface {
	CreateProduct(db *sqlx.DB, product *model.Product) error
	CreateReview(db *sqlx.DB, r *model.Review, userID string) error
	CreateShop(db *sqlx.DB, shop *model.Shop) error
	CreateUser(db *sqlx.DB, user *model.User) error
}

type service struct {
	r Repository
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// CreateProduct validates a product and saves it into the database.
func (s *service) CreateProduct(db *sqlx.DB, p *model.Product) error {
	query := `INSERT INTO products 
	(id, shop_id, stock, brand, category, type, description, weight, discount, taxes, subtotal, total, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	if err := p.Validate(); err != nil {
		return err
	}

	id := random.GenerateRunes(35)
	p.CreatedAt = time.Now()

	taxes := ((p.Subtotal / 100) * p.Taxes)
	discount := ((p.Subtotal / 100) * p.Discount)
	p.Total = p.Subtotal + taxes - discount

	_, err := db.Exec(query, id, p.ShopID, p.Stock, p.Brand, p.Category, p.Type, p.Description,
		p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("couldn't create the product: %v", err)
	}

	return nil
}

// CreateReview takes a new review and saves it into the database.
func (s *service) CreateReview(db *sqlx.DB, r *model.Review, userID string) error {
	query := `INSERT INTO reviews
	(id, stars, comment, user_id, product_id, shop_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	err := r.Validate()
	if err != nil {
		return err
	}

	id := random.GenerateRunes(30)
	r.CreatedAt = time.Now()

	_, err = db.Exec(query, id, r.Stars, r.Comment, userID, r.ProductID, r.ShopID, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("couldn't create the review: %v", err)
	}

	return nil
}

// CreateShop validates a shop and saves it into the database.
func (s *service) CreateShop(db *sqlx.DB, shop *model.Shop) error {
	sQuery := `INSERT INTO shops
	(id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4)`

	lQuery := `INSERT INTO locations
	(shop_id, country, state, zip_code, city, address)
	VALUES ($1, $2, $3, $4, $5, $6)`

	if err := shop.Validate(); err != nil {
		return err
	}

	id := random.GenerateRunes(30)
	shop.CreatedAt = time.Now()

	_, err := db.Exec(sQuery, id, shop.Name, shop.CreatedAt, shop.UpdatedAt)
	if err != nil {
		return fmt.Errorf("couldn't create the shop: %v", err)
	}

	_, err = db.Exec(lQuery, id, shop.Location.Country, shop.Location.State,
		shop.Location.ZipCode, shop.Location.City, shop.Location.Address)
	if err != nil {
		return fmt.Errorf("couldn't create the shop: %v", err)
	}

	return nil
}

// CreateUser validates a user, hashes its password, sends
// a verification email and saves it into the database.
func (s *service) CreateUser(db *sqlx.DB, user *model.User) error {
	cartQuery := `INSERT INTO carts
	(id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	userQuery := `INSERT INTO users
	(id, cart_id, username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if err := user.Validate(""); err != nil {
		return err
	}

	err := db.Get(&user, "SELECT email FROM users WHERE email=$1", user.Email)
	if err == nil {
		return errors.New("email is already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Create a cart for each user
	cartID := random.GenerateRunes(30)
	user.CartID = cartID

	cart := shopping.NewCart(user.CartID)

	_, err = db.Exec(cartQuery, cart.ID, cart.Counter, cart.Weight, cart.Discount,
		cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		return fmt.Errorf("couldn't create the cart: %v", err)
	}

	userID := random.GenerateRunes(30)
	user.CreatedAt = time.Now()

	_, err = db.Exec(userQuery, userID, cart.ID, user.Username, user.Email,
		user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("couldn't create the user: %v", err)
	}

	return nil
}
