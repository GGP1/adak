package searching

import (
	"github.com/GGP1/palo/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type service struct {
	r Repository
}

// NewService creates a searching service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// SearchProducts looks for the products that contain the value specified. (Only text fields)
func (s *service) SearchProducts(db *gorm.DB, products *[]model.Product, search string) error {
	query := `SELECT * FROM products 
	WHERE to_tsvector(brand || ' ' || type || ' ' || category || ' ' || description) 
	@@ to_tsquery(?)`

	err := db.Raw(query, search).Scan(&products).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find products")
	}

	return nil
}

// SearchShops looks for the shops that contain the value specified. (Only text fields)
func (s *service) SearchShops(db *gorm.DB, users *[]model.Shop, search string) error {
	query := `SELECT * FROM shops 
	WHERE to_tsvector(name) 
	@@ to_tsquery(?)`

	err := db.Raw(query, search).Scan(&users).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find shops")
	}

	return nil
}

// SearchUsers looks for the users that contain the value specified. (Only text fields)
func (s *service) SearchUsers(db *gorm.DB, users *[]model.User, search string) error {
	query := `SELECT * FROM users 
	WHERE to_tsvector(firstname || ' ' || lastname || ' ' || email) 
	@@ to_tsquery(?)`

	err := db.Raw(query, search).Scan(&users).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find users")
	}

	return nil
}
