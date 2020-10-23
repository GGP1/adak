package cart

import (
	"context"
	"fmt"
	"math"
	"net"

	"google.golang.org/grpc"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Shopping implements the shopping service.
type Shopping struct {
	db   *sqlx.DB
	Cart *Cart
}

// NewService returns a new shopping server.
func NewService(db *sqlx.DB) *Shopping {
	return &Shopping{
		db:   db,
		Cart: &Cart{},
	}
}

// Run starts the server.
func (s *Shopping) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errors.Wrapf(err, "shopping: failed listening on port %d", port)
	}

	srv := grpc.NewServer()
	RegisterShoppingServer(srv, s)

	return srv.Serve(lis)
}

// New returns a new cart with the id provided.
func (s *Shopping) New(ctx context.Context, req *NewRequest) (*NewResponse, error) {
	return &NewResponse{
		Cart: &Cart{
			ID: req.ID,
		},
	}, nil
}

// Add adds a product to the cart.
func (s *Shopping) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	var sum int64

	pQuery := `INSERT INTO cart_products
	(id, cart_id, quantity, brand, category, type, description, weight, 
	discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	if err := s.db.GetContext(ctx, &s.Cart, "SELECT * FROM carts WHERE id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	AddProduct(s.Cart, req.Product, req.Quantity)

	// Check how many products are in the cart with the product id provided
	// If sum == 0 (there is no product with the same id and cart_id), create the product
	s.db.QueryRow("SELECT SUM(quantity) FROM cart_products WHERE id=$1 AND cart_id=$2", req.Product.ID, req.CartID).Scan(&sum)
	if sum == 0 {
		_, err := s.db.ExecContext(ctx, pQuery, req.Product.ID, req.CartID, req.Product.Quantity, req.Product.Brand, req.Product.Category, req.Product.Type,
			req.Product.Description, req.Product.Weight, req.Product.Discount, req.Product.Taxes, req.Product.Subtotal, req.Product.Total)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create the product")
		}
	}
	// If sum != 0 (product already exists), update the quantity
	if sum != 0 {
		req.Product.Quantity += sum

		_, err := s.db.ExecContext(ctx, "UPDATE cart_products SET quantity=$2 WHERE cart_id=$1", req.CartID, req.Product.Quantity)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update the product")
		}
	}

	_, err := s.db.ExecContext(ctx, cQuery, req.CartID, s.Cart.Counter, s.Cart.Weight, s.Cart.Discount, s.Cart.Taxes, s.Cart.Subtotal,
		s.Cart.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the cart")
	}

	return &AddResponse{Product: req.Product}, nil
}

// Checkout returns the cart total.
func (s *Shopping) Checkout(ctx context.Context, req *CheckoutRequest) (*CheckoutResponse, error) {
	if err := s.db.GetContext(ctx, &s.Cart, "SELECT * FROM carts WHERE id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	return &CheckoutResponse{Total: s.Cart.Total}, nil
}

// Delete permanently deletes a cart from the database.
func (s *Shopping) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	_, err := s.db.ExecContext(ctx, "DELETE FROM carts WHERE id=$1", req.CartID)
	if err != nil {
		return nil, errors.New("couldn't delete the cart")
	}

	return nil, nil
}

// Get returns the user cart.
func (s *Shopping) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	var products []*Product

	if err := s.db.GetContext(ctx, &s.Cart, "SELECT * FROM carts WHERE id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	if err := s.db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart products")
	}

	s.Cart.Products = products

	return &GetResponse{Cart: s.Cart}, nil
}

// Products returns the cart products.
func (s *Shopping) Products(ctx context.Context, req *ProductsRequest) (*ProductsResponse, error) {
	var products []*Product

	if err := s.db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	if len(products) == 0 {
		return nil, errors.New("cart is empty")
	}

	return &ProductsResponse{Products: products}, nil
}

// Remove takes away the specified quantity of products from the cart.
func (s *Shopping) Remove(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	var p *Product

	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	if err := s.db.GetContext(ctx, &s.Cart, "SELECT * FROM carts WHERE id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	if err := s.db.GetContext(ctx, &p, "SELECT * FROM cart_products WHERE id = $1 AND cart_id=$2", req.ProductID, req.CartID); err != nil {
		return nil, errors.New("couldn't find the product")
	}

	if req.Quantity > p.Quantity {
		return nil, fmt.Errorf("quantity inserted (%d) is higher than the stock of products (%d)", req.Quantity, p.Quantity)
	}

	if req.Quantity == p.Quantity {
		_, err := s.db.ExecContext(ctx, "DELETE FROM cart_products WHERE id=$1 AND cart_id=$2", req.ProductID, req.CartID)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't delete the product")
		}
	}

	if s.Cart.Counter == 1 {
		_, err := s.Reset(ctx, &ResetRequest{CartID: req.CartID})
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	RemoveProduct(s.Cart, p, req.Quantity)

	_, err := s.db.ExecContext(ctx, cQuery, req.CartID, s.Cart.Counter, s.Cart.Weight, s.Cart.Discount, s.Cart.Taxes, s.Cart.Subtotal,
		s.Cart.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the cart")
	}

	return &RemoveResponse{}, nil
}

// Reset sets cart values to default.
func (s *Shopping) Reset(ctx context.Context, req *ResetRequest) (*ResetResponse, error) {
	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	_, err := s.db.ExecContext(ctx, "DELETE FROM cart_products WHERE cart_id=$1", req.CartID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete cart products")
	}

	_, err = s.db.ExecContext(ctx, cQuery, req.CartID, 0, 0, 0, 0, 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't reset the cart")
	}

	return &ResetResponse{}, nil
}

// Size returns the quantity of products inside the cart.
func (s *Shopping) Size(ctx context.Context, req *SizeRequest) (*SizeResponse, error) {
	if err := s.db.GetContext(ctx, &s.Cart, "SELECT * FROM carts WHERE id=$1", req.CartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	return &SizeResponse{Counter: s.Cart.Counter}, nil
}

// String returns a string with the cart details.
func String(ctx context.Context, db *sqlx.DB, cartID string) (string, error) {
	var c Cart

	if err := db.GetContext(ctx, &c, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return "", errors.Wrap(err, "couldn't find the cart")
	}

	const details = `The cart has %d products, a weight of %2.fkg, $%2.f of discounts, 
	$%2.f of taxes and a total of $%2.f`

	return fmt.Sprintf(details, c.Counter, c.Weight, c.Discount, c.Taxes, c.Total), nil
}

// AddProduct executes the mathematical process that takes adding a product to the cart.
func AddProduct(c *Cart, p *Product, quantity int64) {
	// percentages -> numeric values
	taxes := ((p.Subtotal / 100) * p.Taxes)
	discount := ((p.Subtotal / 100) * p.Discount)

	p.CartID = c.ID
	p.Total = p.Total + p.Subtotal + taxes - discount

	// math.Ceil(x*100)/100 is used to round floats
	for i := 0; i < int(quantity); i++ {
		c.Counter++
		p.Quantity++
		c.Weight += math.Ceil(p.Weight*100) / 100
		c.Discount += math.Ceil(discount*100) / 100
		c.Taxes += math.Ceil(taxes*100) / 100
		c.Subtotal += math.Ceil(p.Subtotal*100) / 100
		c.Total = c.Total + p.Subtotal + taxes - discount
	}
}

// RemoveProduct executes the mathematical process that takes removing a product from the cart.
func RemoveProduct(c *Cart, p *Product, quantity int64) {
	taxes := (p.Subtotal / 100) * p.Taxes
	discount := (p.Subtotal / 100) * p.Discount

	// math.Ceil(x*100)/100 is used to round float numbers
	for i := 0; i < int(quantity); i++ {
		c.Counter--
		p.Quantity--
		c.Weight -= math.Ceil(p.Weight*100) / 100
		c.Discount -= math.Ceil(discount*100) / 100
		c.Taxes -= math.Ceil(taxes*100) / 100
		c.Subtotal -= math.Ceil(p.Subtotal*100) / 100
		c.Total = c.Total - p.Subtotal - taxes + discount
	}
}
