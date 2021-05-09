package ordering

import (
	"context"
	"time"

	"github.com/GGP1/adak/pkg/shopping/cart"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4/zero"
)

// Service contains order functionalities.
type Service interface {
	New(ctx context.Context, id, userID string, cartID string, oParams OrderParams, cartService cart.Service) (*Order, error)
	Delete(ctx context.Context, orderID string) error
	Get(ctx context.Context) ([]Order, error)
	GetByID(ctx context.Context, orderID string) (Order, error)
	GetByUserID(ctx context.Context, userID string) ([]Order, error)
	GetCartByID(ctx context.Context, orderID string) (OrderCart, error)
	GetProductsByID(ctx context.Context, orderID string) ([]OrderProduct, error)
	UpdateStatus(ctx context.Context, orderID string, status zero.Int) error
}

type service struct {
	db *sqlx.DB
}

// NewService returns a new ordering service.
func NewService(db *sqlx.DB) Service {
	return &service{db}
}

// New creates an order.
func (s *service) New(ctx context.Context, id, userID, cartID string,
	oParams OrderParams, cartService cart.Service) (*Order, error) {
	cart, err := cartService.Get(ctx, cartID)
	if err != nil {
		return nil, err
	}

	if cart.Counter.Int64 == 0 {
		return nil, errors.New("ordering zero products is not permitted")
	}

	// Format delivery date
	deliveryDate := time.Date(oParams.Date.Year, time.Month(oParams.Date.Month), oParams.Date.Day,
		oParams.Date.Hour, oParams.Date.Minutes, 0, 0, time.Local)
	if deliveryDate.Before(time.Now()) {
		return nil, errors.New("past dates are not valid")
	}

	orderQ := `INSERT INTO orders
	(id, user_id, currency, address, city, country, state, zip_code, 
	status, ordered_at, delivery_date, cart_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = s.db.ExecContext(ctx, orderQ, id, userID, oParams.Currency,
		oParams.Address, oParams.City, oParams.Country, oParams.State, oParams.ZipCode,
		zero.IntFrom(int64(Pending)), zero.TimeFrom(time.Now()),
		zero.TimeFrom(deliveryDate), cart.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order")
	}

	if err := s.saveOrderCart(ctx, id, cart); err != nil {
		return nil, err
	}

	if err := s.saveOrderProducts(ctx, id, cart.Products); err != nil {
		return nil, err
	}

	order := Order{
		ID:           zero.StringFrom(id),
		UserID:       zero.StringFrom(userID),
		Currency:     zero.StringFrom(oParams.Currency),
		Address:      zero.StringFrom(oParams.Address),
		City:         zero.StringFrom(oParams.City),
		State:        zero.StringFrom(oParams.State),
		ZipCode:      zero.StringFrom(oParams.ZipCode),
		Status:       zero.IntFrom(int64(Pending)),
		OrderedAt:    zero.TimeFrom(time.Now()),
		DeliveryDate: zero.TimeFrom(deliveryDate),
		CartID:       zero.StringFrom(cart.ID),
		Cart: OrderCart{
			OrderID:  zero.StringFrom(id),
			Counter:  cart.Counter,
			Weight:   cart.Weight,
			Discount: cart.Discount,
			Taxes:    cart.Taxes,
			Subtotal: cart.Subtotal,
			Total:    cart.Total,
		},
	}

	return &order, nil
}

// Delete removes an order.
func (s *service) Delete(ctx context.Context, orderID string) error {
	if _, err := s.db.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", orderID); err != nil {
		return errors.Wrap(err, "couldn't delete the order")
	}

	return nil
}

// Get retrieves all the orders.
func (s *service) Get(ctx context.Context) ([]Order, error) {
	var orders []Order

	if err := s.db.SelectContext(ctx, &orders, "SELECT * FROM orders"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	return orders, nil
}

// GetByID returns the order with the id provided.
func (s *service) GetByID(ctx context.Context, orderID string) (Order, error) {
	q := `SELECT o.*, c.*, p.*
	FROM orders AS o
	LEFT JOIN order_carts AS c ON o.id=c.order_id
	LEFT JOIN order_products AS p ON o.id=p.order_id
	WHERE o.id=$1`

	rows, err := s.db.QueryContext(ctx, q, orderID)
	if err != nil {
		return Order{}, errors.Wrap(err, "fetching orders")
	}
	defer rows.Close()

	var order Order
	for rows.Next() {
		c := OrderCart{}
		p := OrderProduct{}
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Currency, &order.Address, &order.City,
			&order.State, &order.ZipCode, &order.Country, &order.Status, &order.OrderedAt,
			&order.DeliveryDate, &order.CartID,
			&c.OrderID, &c.Counter, &c.Weight, &c.Discount, &c.Taxes, &c.Subtotal, &c.Total,
			&p.ProductID, &p.OrderID, &p.Quantity, &p.Brand, &p.Category, &p.Type, &p.Description,
			&p.Weight, &p.Discount, &p.Taxes, &p.Subtotal, &p.Total,
		)
		if err != nil {
			return Order{}, errors.Wrap(err, "couldn't scan order")
		}

		order.Cart = c
		order.Products = append(order.Products, p)
	}

	return order, nil
}

// GetByUserID retrieves orders depending on the user requested.
func (s *service) GetByUserID(ctx context.Context, userID string) ([]Order, error) {
	q := `SELECT o.*, c.*, p.*
	FROM orders AS o
	LEFT JOIN order_carts AS c ON o.id=c.order_id
	LEFT JOIN order_products AS p ON o.id=p.order_id
	WHERE o.user_id=$1`

	rows, err := s.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching orders")
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		o := Order{}
		c := OrderCart{}
		p := OrderProduct{}
		err := rows.Scan(
			&o.ID, &o.UserID, &o.Currency, &o.Address, &o.City,
			&o.State, &o.ZipCode, &o.Country, &o.Status, &o.OrderedAt,
			&o.DeliveryDate, &o.CartID,
			&c.OrderID, &c.Counter, &c.Weight, &c.Discount, &c.Taxes, &c.Subtotal, &c.Total,
			&p.ProductID, &p.OrderID, &p.Quantity, &p.Brand, &p.Category, &p.Type,
			&p.Description, &p.Weight, &p.Discount, &p.Taxes, &p.Subtotal, &p.Total,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't scan order")
		}

		o.Cart = c
		o.Products = append(o.Products, p)
		orders = append(orders, o)
	}

	return orders, nil
}

// GetCartByID returns the cart with the order id provided.
func (s *service) GetCartByID(ctx context.Context, orderID string) (OrderCart, error) {
	var cart OrderCart

	if err := s.db.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", orderID); err != nil {
		return OrderCart{}, errors.Wrap(err, "couldn't find the order cart")
	}

	return cart, nil
}

// GetProductsByID returns the products with the order id provided.
func (s *service) GetProductsByID(ctx context.Context, orderID string) ([]OrderProduct, error) {
	var products []OrderProduct

	if err := s.db.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", orderID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the order products")
	}

	return products, nil
}

// UpdateStatus updates the order status.
func (s *service) UpdateStatus(ctx context.Context, orderID string, status zero.Int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE orders SET status=$2 WHERE id=$1", orderID, status)
	if err != nil {
		return errors.Wrap(err, "couldn't update the order status")
	}

	return nil
}

// saveOrderCart saves the current user cart to the database.
func (s *service) saveOrderCart(ctx context.Context, id string, cart *cart.Cart) error {
	cartQ := `INSERT INTO order_carts
	(order_id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, cartQ, id, cart.Counter, cart.Weight,
		cart.Discount, cart.Taxes, cart.Subtotal, cart.Total)
	if err != nil {
		return errors.Wrap(err, "couldn't save the order cart")
	}

	return nil
}

// saveOrderProducts saves cart products to the database using batch insert.
func (s *service) saveOrderProducts(ctx context.Context, id string, cartProducts []cart.Product) error {
	products := make([]OrderProduct, len(cartProducts))
	for i, p := range cartProducts {
		products[i] = OrderProduct{
			ProductID:   p.ID,
			OrderID:     zero.StringFrom(id),
			Quantity:    p.Quantity,
			Brand:       p.Brand,
			Category:    p.Category,
			Description: p.Description,
			Discount:    p.Discount,
			Taxes:       p.Taxes,
			Type:        p.Type,
			Subtotal:    p.Subtotal,
			Total:       p.Total,
		}
	}

	productsQ := `INSERT INTO order_products
	(order_id, product_id, quantity, brand, category, type, description, weight, 
	discount, taxes, subtotal, total)
	VALUES 
	(:order_id, :product_id, :quantity, :brand, :category, :type, :description, 
	:weight, :discount, :taxes, :subtotal, :total)`
	if _, err := s.db.NamedExecContext(ctx, productsQ, products); err != nil {
		return errors.Wrap(err, "couldn't save order products")
	}

	return nil
}
