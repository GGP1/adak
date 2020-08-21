// Package searching provides a service for searching specific information
// in the database related to the core api models.
package searching

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to the storage.
type Repository interface {
	SearchProducts(db *sqlx.DB, search string) ([]model.Product, error)
	SearchShops(db *sqlx.DB, search string) ([]model.Shop, error)
	SearchUsers(db *sqlx.DB, search string) ([]model.User, error)
}

// Service provides models searching operations.
type Service interface {
	SearchProducts(db *sqlx.DB, search string) ([]model.Product, error)
	SearchShops(db *sqlx.DB, search string) ([]model.Shop, error)
	SearchUsers(db *sqlx.DB, search string) ([]model.User, error)
}

type service struct {
	r Repository
}

// NewService creates a searching service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// SearchProducts looks for the products that contain the value specified. (Only text fields)
func (s *service) SearchProducts(db *sqlx.DB, search string) ([]model.Product, error) {
	var (
		products []model.Product
		result   []model.Product
	)

	query := `SELECT * FROM products WHERE
	to_tsvector(id || ' ' || shop_id || ' ' || brand || ' ' || type || ' ' || category || ' ' || description)
	@@ to_tsquery($1)`

	if err := db.Select(&products, query, search); err != nil {
		return nil, fmt.Errorf("couldn't find products: %v", err)
	}

	for _, product := range products {
		var reviews []model.Review

		if err := db.Select(&reviews, "SELECT * FROM reviews WHERE product_id=$1", product.ID); err != nil {
			return nil, fmt.Errorf("error fetching reviews: %v", err)
		}

		product.Reviews = reviews

		result = append(result, product)
	}

	return result, nil
}

// SearchShops looks for the shops that contain the value specified. (Only text fields)
func (s *service) SearchShops(db *sqlx.DB, search string) ([]model.Shop, error) {
	var (
		shops  []model.Shop
		result []model.Shop
	)

	query := `SELECT * FROM shops WHERE
	to_tsvector(id || ' ' || name) @@ to_tsquery($1)`

	if err := db.Select(&shops, query, search); err != nil {
		return nil, fmt.Errorf("couldn't find shops: %v", err)
	}

	for _, shop := range shops {
		var (
			location model.Location
			reviews  []model.Review
			products []model.Product
		)

		if err := db.Get(&location, "SELECT * FROM locations WHERE shop_id=$1", shop.ID); err != nil {
			return nil, fmt.Errorf("location not found: %v", err)
		}

		if err := db.Select(&reviews, "SELECT * FROM reviews WHERE shop_id=$1", shop.ID); err != nil {
			return nil, fmt.Errorf("reviews not found: %v", err)
		}

		if err := db.Select(&products, "SELECT * FROM products WHERE shop_id=$1", shop.ID); err != nil {
			return nil, fmt.Errorf("products not found: %v", err)
		}

		shop.Location = location
		shop.Reviews = reviews
		shop.Products = products

		result = append(result, shop)
	}

	return result, nil
}

// SearchUsers looks for the users that contain the value specified. (Only text fields)
func (s *service) SearchUsers(db *sqlx.DB, search string) ([]model.User, error) {
	var (
		users  []model.User
		result []model.User
	)

	query := `SELECT * FROM users WHERE
	to_tsvector(id || ' ' || username || ' ' || email) 
	@@ to_tsquery($1)`

	if err := db.Select(&users, query, search); err != nil {
		return nil, fmt.Errorf("couldn't find users: %v", err)
	}

	for _, user := range users {
		var reviews []model.Review

		orders, err := ordering.GetByUserID(db, user.ID)
		if err != nil {
			return nil, err
		}

		if err := db.Select(&reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
			return nil, fmt.Errorf("error fetching reviews: %v", err)
		}

		user.Orders = orders
		user.Reviews = reviews

		result = append(result, user)
	}

	return result, nil
}
