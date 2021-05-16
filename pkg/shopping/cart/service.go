package cart

import (
	"context"

	"github.com/GGP1/adak/pkg/product"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
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
	Add(ctx context.Context, cartProduct Product) error
	Checkout(ctx context.Context, cartID string) (int64, error)
	Create(ctx context.Context, cartID string) error
	Delete(ctx context.Context, cartID string) error
	FilterBy(ctx context.Context, cartID, field, args string) ([]product.Product, error)
	Get(ctx context.Context, cartID string) (Cart, error)
	CartProduct(ctx context.Context, cartID, productID string) (Product, error)
	CartProducts(ctx context.Context, cartID string) ([]Product, error)
	Remove(ctx context.Context, cartID string, pID string, quantity int64) error
	Reset(ctx context.Context, cartID string) error
	Size(ctx context.Context, cartID string) (int64, error)
}

type service struct {
	db      *sqlx.DB
	mc      *memcache.Client
	metrics metrics
}

// NewService returns a new cart service.
func NewService(db *sqlx.DB, mc *memcache.Client) Service {
	return &service{db, mc, initMetrics()}
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
func (s *service) Add(ctx context.Context, cartProduct Product) error {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Add"}).Inc()

	var p product.Product
	if err := s.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=$1", cartProduct.ID); err != nil {
		return errors.Wrap(err, "couldn't find product")
	}

	if err := s.createOrUpdateProduct(ctx, cartProduct); err != nil {
		return err
	}

	q := `UPDATE carts SET 
	counter=counter+$2, weight=weight+$3, 
	discount=discount+$4, taxes=taxes+$5, 
	subtotal=subtotal+$6, total=total+$7 
	WHERE id=$1`
	_, err := s.db.ExecContext(ctx, q, cartProduct.CartID, cartProduct.Quantity,
		p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
	if err != nil {
		return errors.Wrap(err, "updating cart")
	}

	if err := s.mc.Delete(cartProduct.CartID.String); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting cart from cache")
	}

	return nil
}

// Checkout returns the cart total.
func (s *service) Checkout(ctx context.Context, cartID string) (int64, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Checkout"}).Inc()

	var cart Cart
	if err := s.db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	total := cart.Total.Int64 + cart.Taxes.Int64 - cart.Discount.Int64
	return total, nil
}

// Create a cart.
func (s *service) Create(ctx context.Context, cartID string) error {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Create"}).Inc()

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
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Delete"}).Inc()

	if _, err := s.db.ExecContext(ctx, "DELETE FROM carts WHERE id=$1", cartID); err != nil {
		return errors.New("deleting cart from postgres")
	}

	if err := s.mc.Delete(cartID); err != nil && err != memcache.ErrCacheMiss {
		return errors.Wrap(err, "deleting cart from cache")
	}

	return nil
}

// FilterBy filters the cart products field by the given args.
func (s *service) FilterBy(ctx context.Context, cartID, field, args string) ([]product.Product, error) {
	condition, ok := fields[field]
	if !ok {
		return nil, errors.Errorf("%q, invalid field", field)
	}

	s.metrics.methodCalls.With(prometheus.Labels{"method": "FilterBy " + field}).Inc()

	query := "SELECT * FROM products WHERE cart_id=$1 AND " + condition
	var products []product.Product
	if err := s.db.SelectContext(ctx, &products, query, cartID, args); err != nil {
		return nil, errors.Wrap(err, "no products found")
	}

	if len(products) == 0 {
		return nil, errors.New("no products found")
	}

	return products, nil
}

// Get returns the user cart.
func (s *service) Get(ctx context.Context, cartID string) (Cart, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Get"}).Inc()

	q := `SELECT c.*, cp.*
	FROM carts AS c
	LEFT JOIN cart_products AS cp ON c.id=cp.cart_id
	WHERE c.id=$1`

	rows, err := s.db.QueryContext(ctx, q, cartID)
	if err != nil {
		return Cart{}, errors.Wrap(err, "fetching cart")
	}
	defer rows.Close()

	var cart Cart
	for rows.Next() {
		p := Product{}
		err := rows.Scan(
			&cart.ID, &cart.Counter, &cart.Weight, &cart.Discount,
			&cart.Taxes, &cart.Subtotal, &cart.Total,
			&p.ID, &p.CartID, &p.Quantity,
		)
		if err != nil {
			return Cart{}, errors.Wrap(err, "couldn't scan cart")
		}

		cart.Products = append(cart.Products, p)
	}

	return cart, nil
}

// CartProduct returns a cart product.
func (s *service) CartProduct(ctx context.Context, cartID, productID string) (Product, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Product"}).Inc()

	var product Product
	q := "SELECT * FROM cart_products WHERE id=$1 AND cart_id=$2"
	if err := s.db.GetContext(ctx, &product, q, productID, cartID); err != nil {
		return Product{}, errors.Wrap(err, "couldn't find cart product")
	}

	return product, nil
}

// Products returns the cart products.
func (s *service) CartProducts(ctx context.Context, cartID string) ([]Product, error) {
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Products"}).Inc()

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
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Remove"}).Inc()

	cartProduct, err := s.CartProduct(ctx, cartID, pID)
	if err != nil {
		return err
	}

	if quantity > cartProduct.Quantity.Int64 {
		return errors.Errorf("quantity to remove (%d) is higher than the stock of products (%v)",
			quantity, cartProduct.Quantity)
	}

	if quantity == cartProduct.Quantity.Int64 {
		_, err := s.db.ExecContext(ctx, "DELETE FROM cart_products WHERE id=$1 AND cart_id=$2", pID, cartID)
		if err != nil {
			return errors.Wrap(err, "couldn't delete the product")
		}
	}

	var product product.Product
	if err := s.db.GetContext(ctx, &product, "SELECT * FROM products WHERE id = $1", pID); err != nil {
		return errors.New("couldn't find the product")
	}

	q := `UPDATE carts SET 
	counter=counter-$2, weight=weight-$3, discount=discount-$4, 
	taxes=taxes-$5, subtotal=subtotal-$6, total=total-$7 
	WHERE id=$1`
	_, err = s.db.ExecContext(ctx, q, cartID, quantity, product.Weight,
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
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Reset"}).Inc()

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
	s.metrics.methodCalls.With(prometheus.Labels{"method": "Size"}).Inc()

	var size int64
	row := s.db.QueryRowContext(ctx, "SELECT counter FROM carts WHERE id=$1", cartID)
	if err := row.Scan(&size); err != nil {
		return 0, errors.Wrap(err, "couldn't scan cart size")
	}
	return size, nil
}

func (s *service) createOrUpdateProduct(ctx context.Context, cartProduct Product) error {
	productsQ := `INSERT INTO cart_products
	(id, cart_id, quantity)
	VALUES ($1, $2, $3)
	ON CONFLICT (id) DO UPDATE SET 
	quantity=EXCLUDED.quantity+$3`
	_, err := s.db.ExecContext(ctx, productsQ, cartProduct.ID, cartProduct.CartID, cartProduct.Quantity)
	if err != nil {
		return errors.Wrap(err, "couldn't create the product")
	}

	return nil
}
