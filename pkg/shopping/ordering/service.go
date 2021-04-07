// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/shopping/cart"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Order status
const (
	PendingState  = "pending"
	PaidState     = "paid"
	ShippingState = "shipping"
	ShippedState  = "shipped"
	FailedState   = "failed"
)

// Ordering implements the ordering service
type Ordering struct {
	db *sqlx.DB

	Order          *Order
	shoppingClient cart.ShoppingClient
}

// NewService returns a new ordering server.
func NewService(db *sqlx.DB, shoppingConn *grpc.ClientConn) *Ordering {
	return &Ordering{
		db:             db,
		shoppingClient: cart.NewShoppingClient(shoppingConn),
	}
}

// Run starts the server.
func (o *Ordering) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "ordering: failed connecting to the server on port %d", port)
	}

	srv := grpc.NewServer()
	RegisterOrderingServer(srv, o)

	return srv.Serve(lis)
}

// New creates an order.
func (o *Ordering) New(ctx context.Context, req *NewRequest) (*NewResponse, error) {
	orderQuery := `INSERT INTO orders
	(id, user_id, currency, address, city, country, state, zip_code, status, ordered_at, delivery_date, cart_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	orderCQuery := `INSERT INTO order_carts
	(order_id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	orderPQuery := `INSERT INTO order_products
	(product_id, order_id, quantity, brand, category, type, description, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	if req.Cart.Counter == 0 {
		return nil, errors.New("ordering 0 products is not permitted")
	}

	id := token.GenerateRunes(30)

	order := Order{
		ID:           id,
		UserID:       req.UserID,
		Currency:     req.Currency,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		ZipCode:      req.ZipCode,
		Status:       PendingState,
		OrderedAt:    timestamppb.Now(),
		DeliveryDate: req.DeliveryDate,
		CartID:       req.Cart.ID,
		Cart: &OrderCart{
			Counter:  req.Cart.Counter,
			Weight:   req.Cart.Weight,
			Discount: req.Cart.Discount,
			Taxes:    req.Cart.Taxes,
			Subtotal: req.Cart.Subtotal,
			Total:    req.Cart.Total,
		},
	}

	// Save order products
	for _, p := range req.Cart.Products {
		_, err := o.db.ExecContext(ctx, orderPQuery, p.ID, id, p.Quantity, p.Brand, p.Category,
			p.Type, p.Description, p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create order products")
		}
	}

	// Save order
	_, err := o.db.ExecContext(ctx, orderQuery, id, req.UserID, req.Currency, req.Address, req.City, req.Country,
		req.State, req.ZipCode, PendingState, time.Now(), req.DeliveryDate, req.Cart.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order")
	}

	// Save order cart
	_, err = o.db.ExecContext(ctx, orderCQuery, id, req.Cart.Counter, req.Cart.Weight, req.Cart.Discount,
		req.Cart.Taxes, req.Cart.Subtotal, req.Cart.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order cart")
	}

	_, err = o.shoppingClient.Reset(ctx, &cart.ResetRequest{CartID: req.Cart.ID})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't reset the cart")
	}

	return &NewResponse{Order: &order}, nil
}

// Delete removes an order.
func (o *Ordering) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	if _, err := o.db.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", req.OrderID); err != nil {
		return nil, errors.Wrap(err, "couldn't delete the order")
	}

	_, err := o.db.ExecContext(ctx, "DELETE FROM order_carts WHERE order_id=$1", req.OrderID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete the order cart")
	}

	_, err = o.db.ExecContext(ctx, "DELETE FROM order_products WHERE order_id=$1", req.OrderID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete the order products")
	}

	return nil, nil
}

// Get retrieves all the orders.
func (o *Ordering) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	var orders []*Order

	if err := o.db.SelectContext(ctx, &orders, "SELECT * FROM orders"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	list, err := o.getRelationships(ctx, orders)
	if err != nil {
		return nil, err
	}

	return &GetResponse{Orders: list}, nil
}

// GetByID returns the order with the id provided.
func (o *Ordering) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	var (
		order    Order
		cart     *OrderCart
		products []*OrderProduct
	)

	if err := o.db.GetContext(ctx, "SELECT * FROM orders WHERE id=$1", req.OrderID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the order")
	}

	if err := o.db.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the order cart")
	}

	if err := o.db.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the order products")
	}

	order.Cart = cart
	order.Products = products

	return &GetByIDResponse{Order: &order}, nil
}

// GetByUserID retrieves orders depending on the user requested.
func (o *Ordering) GetByUserID(ctx context.Context, req *GetByUserIDRequest) (*GetByUserIDResponse, error) {
	var orders []*Order

	if err := o.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", req.UserID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the orders")
	}

	list, err := o.getRelationships(ctx, orders)
	if err != nil {
		return nil, err
	}

	return &GetByUserIDResponse{Orders: list}, nil
}

// UpdateStatus updates the order status.
func (o *Ordering) UpdateStatus(ctx context.Context, req *UpdateStatusRequest) (*UpdateStatusResponse, error) {
	_, err := o.db.ExecContext(ctx, "UPDATE orders SET status=$2 WHERE id=$1", req.OrderID, req.OrderStatus)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the order status")
	}

	return nil, nil
}

func (o *Ordering) getRelationships(ctx context.Context, orders []*Order) ([]*Order, error) {
	var list []*Order

	ch, errCh := make(chan *Order), make(chan error, 1)

	for _, order := range orders {
		go func(order *Order) {
			var (
				cart     *OrderCart
				products []*OrderProduct
			)

			// Remove expired leaving a gap of 1 week to compute the shipping order
			if order.DeliveryDate.Seconds > time.Now().Add(time.Hour*168).Unix() {
				_, err := o.Delete(ctx, &DeleteRequest{OrderID: order.ID})
				if err != nil {
					errCh <- err
				}
				return
			}

			if err := o.db.GetContext(ctx, &cart, "SELECT * FROM order_carts WHERE order_id=$1", order.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the order cart")
			}

			if err := o.db.SelectContext(ctx, &products, "SELECT * FROM order_products WHERE order_id=$1", order.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the order products")
			}

			order.Cart = cart
			order.Products = products

			ch <- order
		}(order)
	}

	for i := 0; i < len(orders); i++ {
		select {
		case order := <-ch:
			list = append(list, order)
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil
}
