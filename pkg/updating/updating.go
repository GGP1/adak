// Package updating includes database updating operations.
package updating

import (
	"context"
	"time"

	"github.com/GGP1/palo/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repository provides access to the storage.
type Repository interface {
	UpdateProduct(ctx context.Context, p *model.Product, id string) error
	UpdateShop(ctx context.Context, shop *model.Shop, id string) error
	UpdateUser(ctx context.Context, u *model.User, id string) error
}

// Service provides models updating operations.
type Service interface {
	UpdateProduct(ctx context.Context, p *model.Product, id string) error
	UpdateShop(ctx context.Context, shop *model.Shop, id string) error
	UpdateUser(ctx context.Context, u *model.User, id string) error
}

type service struct {
	r  Repository
	DB *sqlx.DB
}

// NewService creates a updating service with the necessary dependencies.
func NewService(r Repository, db *sqlx.DB) Service {
	return &service{r, db}
}

// UpdateProduct updates a product, returns an error.
func (s *service) UpdateProduct(ctx context.Context, p *model.Product, id string) error {
	q := `UPDATE products SET stock=$2, brand=$3, category=$4, type=$5,
	description=$6, weight=$7, discount=$8, taxes=$9, subtotal=$10, total=$11
	WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, q, id, p.Stock, p.Brand, p.Category, p.Type,
		p.Description, p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't update the product")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// UpdateShop updates a shop, returns an error.
func (s *service) UpdateShop(ctx context.Context, shop *model.Shop, id string) error {
	q := `UPDATE shops SET name=$2, country=$3, city=$4, address=$5
	WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, q, id, shop.Name, shop.Location.Country,
		shop.Location.City, shop.Location.Address)
	if err != nil {
		return errors.Wrap(err, "couldn't update the shop")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// UpdateUser returns updates a user, returns an error.
func (s *service) UpdateUser(ctx context.Context, u *model.User, id string) error {
	q := `UPDATE users SET name=$2, email=$3 WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, q, id, u.Username, u.Email)
	if err != nil {
		return errors.Wrap(err, "couldn't update the user")
	}

	select {
	case <-time.After(0 * time.Nanosecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
