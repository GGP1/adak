package shop

import (
	"context"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	Create(ctx context.Context, shop *Shop) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Shop, error)
	GetByID(ctx context.Context, id string) (Shop, error)
	Search(ctx context.Context, search string) ([]Shop, error)
	Update(ctx context.Context, shop *Shop, id string) error
}

// Service provides shop operations.
type Service interface {
	Create(ctx context.Context, shop *Shop) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Shop, error)
	GetByID(ctx context.Context, id string) (Shop, error)
	Search(ctx context.Context, search string) ([]Shop, error)
	Update(ctx context.Context, shop *Shop, id string) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a deleting service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// Create creates a shop.
func (s *service) Create(ctx context.Context, shop *Shop) error {
	sQuery := `INSERT INTO shops
	(id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4)`

	lQuery := `INSERT INTO locations
	(shop_id, country, state, zip_code, city, address)
	VALUES ($1, $2, $3, $4, $5, $6)`

	id := token.RandString(30)
	shop.CreatedAt = time.Now()

	_, err := s.DB.ExecContext(ctx, sQuery, id, shop.Name, shop.CreatedAt, shop.UpdatedAt)
	if err != nil {
		logger.Log.Errorf("failed creating a shop: %v", err)
		return errors.Wrap(err, "couldn't create the shop")
	}

	_, err = s.DB.ExecContext(ctx, lQuery, id, shop.Location.Country, shop.Location.State,
		shop.Location.ZipCode, shop.Location.City, shop.Location.Address)
	if err != nil {
		logger.Log.Errorf("failed creating the location: %v", err)
		return errors.Wrap(err, "couldn't create the location")
	}

	return nil
}

// Delete permanently deletes a shop from the database.
func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM shops WHERE id=$1", id)
	if err != nil {
		logger.Log.Errorf("failed deleting shop: %v", err)
		return errors.Wrap(err, "couldn't delete the shop")
	}

	return nil
}

// Get returns a list with all the shops stored in the database.
func (s *service) Get(ctx context.Context) ([]Shop, error) {
	var shops []Shop

	if err := s.DB.SelectContext(ctx, &shops, "SELECT * FROM shops"); err != nil {
		logger.Log.Errorf("failed listing shops: %v", err)
		return nil, errors.Wrap(err, "couldn't find the shops")
	}

	list, err := getRelationships(ctx, s.DB, shops)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByID retrieves the shop requested from the database.
func (s *service) GetByID(ctx context.Context, id string) (Shop, error) {
	var (
		shop     Shop
		location Location
		reviews  []review.Review
		products []product.Product
	)

	if err := s.DB.GetContext(ctx, &shop, "SELECT * FROM shops WHERE id=$1", id); err != nil {
		logger.Log.Errorf("failed listing shops: %v", err)
		return Shop{}, errors.Wrap(err, "couldn't find the shop")
	}

	if err := s.DB.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", id); err != nil {
		logger.Log.Errorf("failed listing shop's locations: %v", err)
		return Shop{}, errors.Wrap(err, "couldn't find the location")
	}

	if err := s.DB.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE shop_id=$1", id); err != nil {
		logger.Log.Errorf("failed listing shop's reviews: %v", err)
		return Shop{}, errors.Wrap(err, "couldn't find the reviews")
	}

	if err := s.DB.SelectContext(ctx, &products, "SELECT * FROM products WHERE shop_id=$1", id); err != nil {
		logger.Log.Errorf("failed listing shop's products: %v", err)
		return Shop{}, errors.Wrap(err, "couldn't find the products")
	}

	shop.Location = location
	shop.Reviews = reviews
	shop.Products = products

	return shop, nil
}

// Search looks for the shops that contain the value specified. (Only text fields)
func (s *service) Search(ctx context.Context, search string) ([]Shop, error) {
	if strings.ContainsAny(search, ";-\\|@#~€¬<>_()[]}{¡^'") {
		return nil, errors.New("invalid search")
	}

	shops := []Shop{}
	q := `SELECT * FROM shops WHERE
	to_tsvector(id || ' ' || name) @@ plainto_tsquery($1)`

	if err := s.DB.SelectContext(ctx, &shops, q, search); err != nil {
		logger.Log.Errorf("failed listing shops: %v", err)
		return nil, errors.Wrap(err, "couldn't find shops")
	}

	list, err := getRelationships(ctx, s.DB, shops)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Update updates shop fields.
func (s *service) Update(ctx context.Context, shop *Shop, id string) error {
	q := `UPDATE shops SET name=$2, country=$3, city=$4, address=$5
	WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, q, id, shop.Name, shop.Location.Country,
		shop.Location.City, shop.Location.Address)
	if err != nil {
		logger.Log.Errorf("failed updating shop: %v", err)
		return errors.Wrap(err, "couldn't update the shop")
	}

	return nil
}

func getRelationships(ctx context.Context, db *sqlx.DB, shops []Shop) ([]Shop, error) {
	var list []Shop

	ch, errCh := make(chan Shop), make(chan error, 1)

	for _, shop := range shops {
		go func(shop Shop) {
			var (
				location Location
				reviews  []review.Review
				products []product.Product
			)

			if err := db.GetContext(ctx, &location, "SELECT * FROM locations WHERE shop_id=$1", shop.ID); err != nil {
				logger.Log.Errorf("failed listing shop's locations: %v", err)
				errCh <- errors.Wrap(err, "couldn't find the location")
			}

			if err := db.Select(&reviews, "SELECT * FROM reviews WHERE shop_id=$1", shop.ID); err != nil {
				logger.Log.Errorf("failed listing shop's reviews: %v", err)
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			if err := db.Select(&products, "SELECT * FROM products WHERE shop_id=$1", shop.ID); err != nil {
				logger.Log.Errorf("failed listing shop's products: %v", err)
				errCh <- errors.Wrap(err, "couldn't find the products")
			}

			shop.Location = location
			shop.Products = products
			shop.Reviews = reviews

			ch <- shop
		}(shop)
	}

	for range shops {
		select {
		case shop := <-ch:
			list = append(list, shop)
		case err := <-errCh:
			return nil, err
		}
	}
	return list, nil
}
