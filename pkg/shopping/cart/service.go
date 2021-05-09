package cart

import (
	"context"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4/zero"
)

var fields = map[string]string{
	"brand":    "brand=$2",
	"category": "category=$2",
	"discount": "discount >= $2 AND discount <= $3",
	"subtotal": "subtotal >= $2 AND subtotal <= $3",
	"taxes":    "taxes >= $2 AND taxes <= $3",
	"total":    "total >= $2 AND total <= $3",
	"type":     "type=$2",
}

// Service contains order functionalities.
type Service interface {
	Add(ctx context.Context, cartID string, p *Product) error
	Checkout(ctx context.Context, cartID string) (int64, error)
	Create(ctx context.Context, cartID string) error
	Delete(ctx context.Context, cartID string) error
	FilterBy(ctx context.Context, cartID, field, args string) ([]Product, error)
	Get(ctx context.Context, cartID string) (*Cart, error)
	Products(ctx context.Context, cartID string) ([]Product, error)
	Remove(ctx context.Context, cartID string, pID string, quantity int64) error
	Reset(ctx context.Context, cartID string) error
	Size(ctx context.Context, cartID string) (int64, error)
}

type service struct {
	db *sqlx.DB
	mc *memcache.Client
}

// NewService returns a new cart service.
func NewService(db *sqlx.DB, mc *memcache.Client) Service {
	return &service{db, mc}
}

// New returns a cart with the default values.
func New(id string) *Cart {
	return &Cart{
		ID:       id,
		Counter:  zero.IntFrom(0),
		Weight:   zero.IntFrom(0),
		Discount: zero.IntFrom(0),
		Taxes:    zero.IntFrom(0),
		Subtotal: zero.IntFrom(0),
		Total:    zero.IntFrom(0),
		Products: []Product{},
	}
}

// Add adds a product to the cart.
func (s *service) Add(ctx context.Context, cartID string, p *Product) error {
	p.CartID = zero.StringFrom(cartID)
	if p.Total.Int64 == 0 {
		p.Total = zero.IntFrom(p.Subtotal.Int64 + p.Taxes.Int64 - p.Discount.Int64)
	}

	if err := s.createOrUpdateProduct(ctx, cartID, p); err != nil {
		return err
	}

	q := `UPDATE carts SET 
	counter=counter+$2, weight=weight+$3, 
	discount=discount+$4, taxes=taxes+$5, 
	subtotal=subtotal+$6, total=total+$7 
	WHERE id=$1`
	_, err := s.db.ExecContext(ctx, q, cartID, p.Quantity, p.Weight, p.Discount,
		p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return errors.Wrap(err, "updating cart")
	}

	if err := s.mc.Delete(cartID); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting cart from cache")
	}

	return nil
}

// Checkout returns the cart total.
func (s *service) Checkout(ctx context.Context, cartID string) (int64, error) {
	var cart Cart

	if err := s.db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	total := cart.Total.Int64 + cart.Taxes.Int64 - cart.Discount.Int64

	return total, nil
}

// Create a cart.
func (s *service) Create(ctx context.Context, cartID string) error {
	cartQuery := `INSERT INTO carts
	(id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, cartQuery, cartID, 0, 0, 0, 0, 0, 0)
	if err != nil {
		return errors.Wrap(err, "couldn't create the cart")
	}

	return nil
}

// Delete permanently deletes a cart from the database.
func (s *service) Delete(ctx context.Context, cartID string) error {
	if _, err := s.db.ExecContext(ctx, "DELETE FROM carts WHERE id=$1", cartID); err != nil {
		return errors.New("deleting cart from postgres")
	}

	if err := s.mc.Delete(cartID); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting cart from cache")
	}

	return nil
}

// FilterBy filters the cart products field by the given args.
func (s *service) FilterBy(ctx context.Context, cartID, field, args string) ([]Product, error) {
	var products []Product

	condition, ok := fields[field]
	if !ok {
		return nil, errors.Errorf("%q, invalid field", field)
	}

	query := "SELECT * FROM cart_products WHERE cart_id=$1 AND " + condition
	if err := s.db.SelectContext(ctx, &products, query, cartID, args); err != nil {
		return nil, errors.Wrap(err, "no products found")
	}

	if len(products) == 0 {
		return nil, errors.New("no products found")
	}

	return products, nil
}

// Get returns the user cart.
func (s *service) Get(ctx context.Context, cartID string) (*Cart, error) {
	q := `SELECT c.*, cp.*
	FROM carts AS c
	LEFT JOIN cart_products AS cp ON c.id=cp.cart_id
	WHERE c.id=$1`

	rows, err := s.db.QueryContext(ctx, q, cartID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching cart")
	}
	defer rows.Close()

	var cart Cart
	for rows.Next() {
		p := Product{}
		err := rows.Scan(
			&cart.ID, &cart.Counter, &cart.Weight, &cart.Discount,
			&cart.Taxes, &cart.Subtotal, &cart.Total,
			&p.ID, &p.CartID, &p.Quantity,
			&p.Brand, &p.Category, &p.Type,
			&p.Description, &p.Weight, &p.Discount,
			&p.Taxes, &p.Subtotal, &p.Total,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't scan cart")
		}

		cart.Products = append(cart.Products, p)
	}

	return &cart, nil
}

// Products returns the cart products.
func (s *service) Products(ctx context.Context, cartID string) ([]Product, error) {
	var products []Product

	if err := s.db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart products")
	}

	if len(products) == 0 {
		return nil, errors.New("cart is empty")
	}

	return products, nil
}

// Remove takes away the specified quantity of products from the cart.
func (s *service) Remove(ctx context.Context, cartID string, pID string, quantity int64) error {
	var product Product

	productQ := "SELECT * FROM cart_products WHERE id = $1 AND cart_id=$2"
	if err := s.db.GetContext(ctx, &product, productQ, pID, cartID); err != nil {
		return errors.New("couldn't find the product")
	}

	if quantity > product.Quantity.Int64 {
		return errors.Errorf("quantity to remove (%d) is higher than the stock of products (%v)",
			quantity, product.Quantity)
	}

	if quantity == product.Quantity.Int64 {
		_, err := s.db.ExecContext(ctx, "DELETE FROM cart_products WHERE id=$1 AND cart_id=$2", pID, cartID)
		if err != nil {
			return errors.Wrap(err, "couldn't delete the product")
		}
	}

	cartQ := `UPDATE carts SET 
	counter=counter-$2, weight=weight-$3, discount=discount-$4, 
	taxes=taxes-$5, subtotal=subtotal-$6, total=total-$7 
	WHERE id=$1`
	_, err := s.db.ExecContext(ctx, cartQ, cartID, quantity, product.Weight,
		product.Discount, product.Taxes, product.Subtotal, product.Total)
	if err != nil {
		return errors.Wrap(err, "updating cart")
	}

	if err := s.mc.Delete(cartID); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting cart from cache")
	}

	return nil
}

// Reset sets cart values to default.
func (s *service) Reset(ctx context.Context, cartID string) error {
	del := "DELETE FROM cart_products WHERE cart_id=$1"
	if _, err := s.db.ExecContext(ctx, del, cartID); err != nil {
		return errors.Wrap(err, "couldn't delete cart products")
	}

	upt := `UPDATE carts SET 
	counter=$2, weight=$3, discount=$4, taxes=$5, subtotal=$6, total=$7 
	WHERE id=$1`
	if _, err := s.db.ExecContext(ctx, upt, cartID, 0, 0, 0, 0, 0, 0); err != nil {
		return errors.Wrap(err, "updating cart")
	}

	if err := s.mc.Delete(cartID); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting cart from cache")
	}

	return nil
}

// Size returns the quantity of products inside the cart.
func (s *service) Size(ctx context.Context, cartID string) (int64, error) {
	var size int64
	row := s.db.QueryRowContext(ctx, "SELECT counter FROM carts WHERE id=$1", cartID)
	if err := row.Scan(&size); err != nil {
		return 0, errors.Wrap(err, "couldn't scan cart size")
	}
	return size, nil
}

func (s *service) createOrUpdateProduct(ctx context.Context, cartID string, p *Product) error {
	productsQ := `INSERT INTO cart_products
	(id, cart_id, quantity, brand, category, type, description, weight, 
	discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	ON CONFLICT (id) DO UPDATE SET 
	quantity=EXCLUDED.quantity+$3`
	_, err := s.db.ExecContext(ctx, productsQ, p.ID, cartID, p.Quantity, p.Brand, p.Category,
		p.Type, p.Description, p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't create the product")
	}

	return nil
}
