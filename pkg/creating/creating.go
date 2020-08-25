/*
Package creating includes database creating operations
*/
package creating

import (
	"context"
	"time"

	"github.com/GGP1/palo/internal/random"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Repository provides access to the storage.
type Repository interface {
	CreateProduct(ctx context.Context, p *model.Product) error
	CreateReview(ctx context.Context, r *model.Review, userID string) error
	CreateShop(ctx context.Context, shop *model.Shop) error
	CreateUser(ctx context.Context, user *model.User) error
}

// Service provides models adding operations.
type Service interface {
	CreateProduct(ctx context.Context, p *model.Product) error
	CreateReview(ctx context.Context, r *model.Review, userID string) error
	CreateShop(ctx context.Context, shop *model.Shop) error
	CreateUser(ctx context.Context, user *model.User) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// CreateProduct creates a product.
func (s *service) CreateProduct(ctx context.Context, p *model.Product) error {
	q := `INSERT INTO products 
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

	_, err := s.DB.ExecContext(ctx, q, id, p.ShopID, p.Stock, p.Brand, p.Category, p.Type, p.Description,
		p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the product")
	}

	return nil
}

// CreateReview creates a review.
func (s *service) CreateReview(ctx context.Context, r *model.Review, userID string) error {
	q := `INSERT INTO reviews
	(id, stars, comment, user_id, product_id, shop_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	err := r.Validate()
	if err != nil {
		return err
	}

	id := random.GenerateRunes(30)
	r.CreatedAt = time.Now()

	_, err = s.DB.ExecContext(ctx, q, id, r.Stars, r.Comment, userID, r.ProductID, r.ShopID, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the review")
	}

	return nil
}

// CreateShop creates a shop.
func (s *service) CreateShop(ctx context.Context, shop *model.Shop) error {
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

	_, err := s.DB.ExecContext(ctx, sQuery, id, shop.Name, shop.CreatedAt, shop.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the shop")
	}

	_, err = s.DB.ExecContext(ctx, lQuery, id, shop.Location.Country, shop.Location.State,
		shop.Location.ZipCode, shop.Location.City, shop.Location.Address)
	if err != nil {
		return errors.Wrap(err, "couldn't create the shop")
	}

	return nil
}

// CreateUser creates a user.
func (s *service) CreateUser(ctx context.Context, user *model.User) error {
	cartQuery := `INSERT INTO carts
	(id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	userQuery := `INSERT INTO users
	(id, cart_id, username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if err := user.Validate(""); err != nil {
		return err
	}

	err := s.DB.GetContext(ctx, &user, "SELECT email FROM users WHERE email=$1", user.Email)
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

	_, err = s.DB.ExecContext(ctx, cartQuery, cart.ID, cart.Counter, cart.Weight,
		cart.Discount, cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't create the cart")
	}

	userID := random.GenerateRunes(30)
	user.CreatedAt = time.Now()

	_, err = s.DB.ExecContext(ctx, userQuery, userID, cart.ID, user.Username, user.Email,
		user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "couldn't create the user")
	}

	return nil
}
