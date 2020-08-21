// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"errors"
	"fmt"
	"time"

	"github.com/GGP1/palo/internal/random"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/jmoiron/sqlx"
)

const (
	// PendingState is the pending state of an Order
	PendingState = "pending"
	// PaidState is the paid state of an Order
	PaidState = "paid"
	// ShippingState is the shipping state of an order
	ShippingState = "shipping"
	// ShippedState is the shipped state of an Order
	ShippedState = "shipped"
	// FailedState is the failed state of an Order
	FailedState = "failed"
)

// Order represents the user purchase request.
type Order struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id" db:"user_id"`
	Currency     string         `json:"currency"`
	Address      string         `json:"address"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	ZipCode      string         `json:"zip_code" db:"zip_code"`
	Country      string         `json:"country"`
	Status       string         `json:"status"`
	OrderedAt    time.Time      `json:"ordered_at" db:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date" db:"delivery_date"`
	CartID       string         `json:"cart_id" db:"cart_id"`
	Cart         OrderCart      `json:"cart"`
	Products     []OrderProduct `json:"products"`
}

// OrderCart represents the cart ordered by the user.
type OrderCart struct {
	OrderID  string  `json:"order_id" db:"order_id"`
	Counter  int     `json:"counter"`
	Weight   float64 `json:"weight"`
	Discount float64 `json:"discount"`
	Taxes    float64 `json:"taxes"`
	Subtotal float64 `json:"subtotal"`
	Total    float64 `json:"total"`
}

// OrderProduct represents the a product place into the cart ordered by the user.
type OrderProduct struct {
	ProductID   string  `json:"product_id" db:"product_id"`
	OrderID     string  `json:"order_id" db:"order_id"`
	Quantity    int     `json:"quantity"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Discount    float64 `json:"discount"`
	Taxes       float64 `json:"taxes"`
	Subtotal    float64 `json:"subtotal"`
	Total       float64 `json:"total"`
}

// NewOrder creates an order.
func NewOrder(db *sqlx.DB, userID, currency, address, city, country, state, zipcode string, deliveryDate time.Time, cart shopping.Cart) (*Order, error) {
	orderQuery := `INSERT INTO orders
	(id, user_id, currency, address, city, country, state, zip_code, status, ordered_at, delivery_date, cart_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	orderCQuery := `INSERT INTO order_carts
	(order_id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	orderPQuery := `INSERT INTO order_products
	(product_id, order_id, quantity, brand, category, type, description, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	if cart.Counter == 0 {
		return nil, errors.New("ordering 0 products is not permitted")
	}

	id := random.GenerateRunes(30)

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
		CartID:       cart.ID,
		Cart: OrderCart{
			Counter:  cart.Counter,
			Weight:   cart.Weight,
			Discount: cart.Discount,
			Taxes:    cart.Taxes,
			Subtotal: cart.Subtotal,
			Total:    cart.Total,
		},
	}

	for _, product := range cart.Products {
		_, err := db.Exec(orderPQuery, product.ID, id, product.Quantity, product.Brand,
			product.Category, product.Type, product.Description, product.Weight,
			product.Discount, product.Taxes, product.Subtotal, product.Total)
		if err != nil {
			return nil, fmt.Errorf("couldn't create order products: %v", err)
		}
	}

	_, err := db.Exec(orderQuery, id, userID, currency, address, city, country,
		state, zipcode, PendingState, time.Now(), deliveryDate, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the order: %v", err)
	}

	_, err = db.Exec(orderCQuery, id, cart.Counter, cart.Weight, cart.Discount,
		cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the order cart: %v", err)
	}

	if err := shopping.Reset(db, cart.ID); err != nil {
		return nil, fmt.Errorf("couldn't reset the cart: %v", err)
	}

	return &order, nil
}

// Delete removes an order.
func Delete(db *sqlx.DB, orderID string) error {
	if err := db.MustExec("DELETE FROM orders WHERE id=$1", orderID); err != nil {
		return fmt.Errorf("couldn't delete the order: %v", err)
	}

	_, err := db.Exec("DELETE FROM order_carts WHERE order_id=$1", orderID)
	if err != nil {
		return fmt.Errorf("couldn't delete the order cart: %v", err)
	}

	_, err = db.Exec("DELETE FROM order_products WHERE order_id=$1", orderID)
	if err != nil {
		return fmt.Errorf("couldn't delete the order products: %v", err)
	}

	return nil
}

// Get retrieves all the orders.
func Get(db *sqlx.DB) ([]Order, error) {
	var (
		orders []Order
		result []Order
	)

	if err := db.Select(&orders, "SELECT * FROM orders"); err != nil {
		return nil, fmt.Errorf("couldn't find the orders: %v", err)
	}

	for _, order := range orders {
		var (
			cart     OrderCart
			products []OrderProduct
		)

		if err := db.Get(&cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
			return nil, fmt.Errorf("couldn't find the order cart: %v", err)
		}

		if err := db.Select(&products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
			return nil, fmt.Errorf("couldn't find order products: %v", err)
		}

		order.Cart = cart
		order.Products = products

		result = append(result, order)
	}

	return result, nil
}

// GetByUserID retrieves orders depending on their id.
func GetByUserID(db *sqlx.DB, userID string) ([]Order, error) {
	var (
		orders []Order
		result []Order
	)

	if err := db.Select(&orders, "SELECT * FROM orders WHERE user_id=$1", userID); err != nil {
		return nil, fmt.Errorf("couldn't find the orders: %v", err)
	}

	for _, order := range orders {
		var (
			cart     OrderCart
			products []OrderProduct
		)

		if err := db.Get(&cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
			return nil, fmt.Errorf("couldn't find the order cart: %v", err)
		}

		if err := db.Select(&products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
			return nil, fmt.Errorf("couldn't find order products: %v", err)
		}

		order.Cart = cart
		order.Products = products

		result = append(result, order)
	}

	return result, nil
}
