// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"context"
	"time"

	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/shopping/cart"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Order status
const (
	PendingState  = "pending"
	PaidState     = "paid"
	ShippingState = "shipping"
	ShippedState  = "shipped"
	FailedState   = "failed"
)

// New creates an order.
func New(ctx context.Context, db *sqlx.DB, userID, currency, address, city, country, state, zipcode string, deliveryDate time.Time, c cart.Cart) (*Order, error) {
	orderQuery := `INSERT INTO orders
	(id, user_id, currency, address, city, country, state, zip_code, status, ordered_at, delivery_date, cart_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	orderCQuery := `INSERT INTO order_carts
	(order_id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	orderPQuery := `INSERT INTO order_products
	(product_id, order_id, quantity, brand, category, type, description, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	if c.Counter == 0 {
		return nil, errors.New("ordering 0 products is not permitted")
	}

	id := token.GenerateRunes(30)

	order := Order{
		ID:           id,
		UserID:       userID,
		Currency:     currency,
		Address:      address,
		City:         city,
		State:        state,
		ZipCode:      zipcode,
		Status:       PendingState,
		OrderedAt:    time.Now(),
		DeliveryDate: deliveryDate,
		CartID:       c.ID,
		Cart: OrderCart{
			Counter:  c.Counter,
			Weight:   c.Weight,
			Discount: c.Discount,
			Taxes:    c.Taxes,
			Subtotal: c.Subtotal,
			Total:    c.Total,
		},
	}

	for _, product := range c.Products {
		_, err := db.ExecContext(ctx, orderPQuery, product.ID, id, product.Quantity, product.Brand,
			product.Category, product.Type, product.Description, product.Weight,
			product.Discount, product.Taxes, product.Subtotal, product.Total)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create order products")
		}
	}

	_, err := db.ExecContext(ctx, orderQuery, id, userID, currency, address, city, country,
		state, zipcode, PendingState, time.Now(), deliveryDate, c.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order")
	}

	_, err = db.ExecContext(ctx, orderCQuery, id, c.Counter, c.Weight, c.Discount,
		c.Taxes, c.Subtotal, c.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order cart")
	}

	if err := cart.Reset(ctx, db, c.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't reset the cart")
	}

	return &order, nil
}

// Delete removes an order.
func Delete(ctx context.Context, db *sqlx.DB, orderID string) error {
	if _, err := db.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", orderID); err != nil {
		return errors.Wrap(err, "couldn't delete the order")
	}

	_, err := db.ExecContext(ctx, "DELETE FROM order_carts WHERE order_id=$1", orderID)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the order cart")
	}

	_, err = db.ExecContext(ctx, "DELETE FROM order_products WHERE order_id=$1", orderID)
	if err != nil {
		return errors.Wrap(err, "couldn't delete the order products")
	}

	return nil
}

// Get retrieves all the orders.
func Get(ctx context.Context, db *sqlx.DB) ([]Order, error) {
	var (
		orders []Order
		list   []Order
	)

	if err := db.SelectContext(ctx, &orders, "SELECT * FROM orders"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	for _, order := range orders {
		var (
			cart     OrderCart
			products []OrderProduct
		)

		if err := db.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
			return nil, errors.Wrap(err, "couldn't find the order cart")
		}

		if err := db.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
			return nil, errors.Wrap(err, "couldn't find order products")
		}

		order.Cart = cart
		order.Products = products

		list = append(list, order)
	}

	return list, nil
}

// GetByID returns the order with the id provided.
func GetByID(ctx context.Context, db *sqlx.DB, orderID string) (Order, error) {
	var (
		order    Order
		cart     OrderCart
		products []OrderProduct
	)

	if err := db.GetContext(ctx, "SELECT * FROM orders WHERE id=$1", orderID); err != nil {
		return Order{}, errors.Wrap(err, "order not found")
	}

	if err := db.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
		return Order{}, errors.Wrap(err, "order cart not found")
	}

	if err := db.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
		return Order{}, errors.Wrap(err, "order products not found")
	}

	order.Cart = cart
	order.Products = products

	return order, nil
}

// GetByUserID retrieves orders depending on the user requested.
func GetByUserID(ctx context.Context, db *sqlx.DB, userID string) ([]Order, error) {
	var (
		orders []Order
		list   []Order
	)

	if err := db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", userID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	for _, order := range orders {
		var (
			cart     OrderCart
			products []OrderProduct
		)

		if err := db.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
			return nil, errors.Wrap(err, "couldn't find the order cart")
		}

		if err := db.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
			return nil, errors.Wrap(err, "couldn't find order products")
		}

		order.Cart = cart
		order.Products = products

		list = append(list, order)
	}

	return list, nil
}

// UpdateStatus updates the order status.
func UpdateStatus(ctx context.Context, db *sqlx.DB, oID, oStatus string) error {
	_, err := db.ExecContext(ctx, "UPDATE orders SET status=$2 WHERE id=$1", oID, oStatus)
	if err != nil {
		return errors.Wrap(err, "couldn't update the order status")
	}

	return nil
}
