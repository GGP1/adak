package ordering

import (
	"context"
	"time"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/shopping/cart"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Service contains order functionalities.
type Service interface {
	New(ctx context.Context, userID string, cartID string, oParams OrderParams, deliveryDate time.Time, cartService cart.Service) (*Order, error)
	Delete(ctx context.Context, orderID string) error
	Get(ctx context.Context) ([]Order, error)
	GetByID(ctx context.Context, orderID string) (Order, error)
	GetByUserID(ctx context.Context, userID string) ([]Order, error)
	GetCartByID(ctx context.Context, orderID string) (OrderCart, error)
	GetProductsByID(ctx context.Context, orderID string) ([]OrderProduct, error)
	UpdateStatus(ctx context.Context, orderID string, status status) error
}

type service struct {
	DB *sqlx.DB
}

// NewService returns a new ordering service.
func NewService(db *sqlx.DB) Service {
	return &service{DB: db}
}

// New creates an order.
func (s *service) New(ctx context.Context, userID, cartID string, oParams OrderParams, deliveryDate time.Time, cartService cart.Service) (*Order, error) {
	orderQuery := `INSERT INTO orders
	(id, user_id, currency, address, city, country, state, zip_code, status, ordered_at, delivery_date, cart_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	orderCartsQuery := `INSERT INTO order_carts
	(order_id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	orderProdQuery := `INSERT INTO order_products
	(product_id, order_id, quantity, brand, category, type, description, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	cart, err := cartService.Get(ctx, cartID)
	if err != nil {
		return nil, err
	}

	if cart.Counter == 0 {
		return nil, errors.New("ordering 0 products is not permitted")
	}

	id := token.RandString(30)

	order := Order{
		ID:           id,
		UserID:       userID,
		Currency:     oParams.Currency,
		Address:      oParams.Address,
		City:         oParams.City,
		State:        oParams.State,
		ZipCode:      oParams.ZipCode,
		Status:       pending,
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

	// Save order products
	for _, product := range cart.Products {
		_, err := s.DB.ExecContext(ctx, orderProdQuery, product.ID, id, product.Quantity, product.Brand,
			product.Category, product.Type, product.Description, product.Weight,
			product.Discount, product.Taxes, product.Subtotal, product.Total)
		if err != nil {
			logger.Log.Errorf("failed creating order's products: %v", err)
			return nil, errors.Wrap(err, "couldn't create order products")
		}
	}

	// Save order
	_, err = s.DB.ExecContext(ctx, orderQuery, id, userID, oParams.Currency, oParams.Address, oParams.City, oParams.Country,
		oParams.State, oParams.ZipCode, pending, time.Now(), deliveryDate, cart.ID)
	if err != nil {
		logger.Log.Errorf("failed creating order: %v", err)
		return nil, errors.Wrap(err, "couldn't create the order")
	}

	// Save order cart
	_, err = s.DB.ExecContext(ctx, orderCartsQuery, id, cart.Counter, cart.Weight, cart.Discount,
		cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		logger.Log.Errorf("failed creating order's cart: %v", err)
		return nil, errors.Wrap(err, "couldn't create the order cart")
	}

	if err := cartService.Reset(ctx, cart.ID); err != nil {
		logger.Log.Errorf("failed resetting cart after an order: %v", err)
		return nil, errors.Wrap(err, "couldn't reset the cart")
	}

	return &order, nil
}

// Delete removes an order.
func (s *service) Delete(ctx context.Context, orderID string) error {
	if _, err := s.DB.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", orderID); err != nil {
		logger.Log.Errorf("failed deleting order: %v", err)
		return errors.Wrap(err, "couldn't delete the order")
	}

	_, err := s.DB.ExecContext(ctx, "DELETE FROM order_carts WHERE order_id=$1", orderID)
	if err != nil {
		logger.Log.Errorf("failed deleting order's cart: %v", err)
		return errors.Wrap(err, "couldn't delete the order cart")
	}

	_, err = s.DB.ExecContext(ctx, "DELETE FROM order_products WHERE order_id=$1", orderID)
	if err != nil {
		logger.Log.Errorf("failed deleting order's products: %v", err)
		return errors.Wrap(err, "couldn't delete the order products")
	}

	return nil
}

// Get retrieves all the orders.
func (s *service) Get(ctx context.Context) ([]Order, error) {
	var orders []Order

	if err := s.DB.SelectContext(ctx, &orders, "SELECT * FROM orders"); err != nil {
		logger.Log.Errorf("failed listing orders: %v", err)
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	list, err := s.getRelationships(ctx, orders)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByID returns the order with the id provided.
func (s *service) GetByID(ctx context.Context, orderID string) (Order, error) {
	var order Order

	if err := s.DB.GetContext(ctx, "SELECT * FROM orders WHERE id=$1", orderID); err != nil {
		return Order{}, errors.Wrap(err, "couldn't find the order")
	}

	cart, err := s.GetCartByID(ctx, order.ID)
	if err != nil {
		return Order{}, err
	}

	products, err := s.GetProductsByID(ctx, order.ID)
	if err != nil {
		return Order{}, err
	}

	order.Cart = cart
	order.Products = products

	return order, nil
}

// GetByUserID retrieves orders depending on the user requested.
func (s *service) GetByUserID(ctx context.Context, userID string) ([]Order, error) {
	var orders []Order

	if err := s.DB.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", userID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	list, err := s.getRelationships(ctx, orders)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetCartByID returns the cart with the order id provided.
func (s *service) GetCartByID(ctx context.Context, orderID string) (OrderCart, error) {
	var cart OrderCart

	if err := s.DB.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", orderID); err != nil {
		logger.Log.Errorf("failed order's cart: %v", err)
		return OrderCart{}, errors.Wrap(err, "couldn't find the order cart")
	}

	return cart, nil
}

// GetProductsByID returns the products with the order id provided.
func (s *service) GetProductsByID(ctx context.Context, orderID string) ([]OrderProduct, error) {
	var products []OrderProduct

	if err := s.DB.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", orderID); err != nil {
		logger.Log.Errorf("failed order's products: %v", err)
		return nil, errors.Wrap(err, "couldn't find the order products")
	}

	return products, nil
}

// UpdateStatus updates the order status.
func (s *service) UpdateStatus(ctx context.Context, orderID string, status status) error {
	_, err := s.DB.ExecContext(ctx, "UPDATE orders SET status=$2 WHERE id=$1", orderID, status)
	if err != nil {
		logger.Log.Errorf("failed updating order's status: %v", err)
		return errors.Wrap(err, "couldn't update the order status")
	}

	return nil
}

func (s *service) getRelationships(ctx context.Context, orders []Order) ([]Order, error) {
	ch, errCh := make(chan Order), make(chan error, 1)

	for _, order := range orders {
		go func(order Order) {
			// Remove expired leaving a gap of 1 week to compute the shipping order
			if order.DeliveryDate.Sub(time.Now().Add(time.Hour*168)) < 0 {
				if err := s.Delete(ctx, order.ID); err != nil {
					errCh <- err
				}
				return
			}

			cart, err := s.GetCartByID(ctx, order.ID)
			if err != nil {
				errCh <- err
			}

			products, err := s.GetProductsByID(ctx, order.ID)
			if err != nil {
				errCh <- err
			}

			order.Cart = cart
			order.Products = products

			ch <- order
		}(order)
	}

	list := make([]Order, len(orders))
	for i := range orders {
		select {
		case order := <-ch:
			list[i] = order
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil
}
