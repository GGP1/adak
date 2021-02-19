package cart

import (
	"context"
	"fmt"

	"github.com/GGP1/adak/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// New returns a cart with the default values.
func New(id string) *Cart {
	return &Cart{
		ID:       id,
		Counter:  0,
		Weight:   0,
		Discount: 0,
		Taxes:    0,
		Subtotal: 0,
		Total:    0,
		Products: []Product{},
	}
}

// Add adds a product to the cart.
func Add(ctx context.Context, db *sqlx.DB, cartID string, p *Product, quantity int) (*Product, error) {
	var (
		cart Cart
		sum  int
	)

	productsQuery := `INSERT INTO cart_products
	(id, cart_id, quantity, brand, category, type, description, weight, 
	discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	cartQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart: %v", err)
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	// percentages -> numeric values
	taxes := ((p.Subtotal / 100) * p.Taxes)
	discount := ((p.Subtotal / 100) * p.Discount)

	p.CartID = cartID
	p.Total = p.Total + p.Subtotal + taxes - discount

	cart.mu.Lock()
	for i := 0; i < quantity; i++ {
		cart.Counter++
		p.Quantity++
		cart.Weight += p.Weight
		cart.Discount += discount
		cart.Taxes += taxes
		cart.Subtotal += p.Subtotal
		cart.Total = cart.Total + p.Subtotal + taxes - discount
	}
	cart.mu.Unlock()

	db.QueryRow("SELECT SUM(quantity) FROM cart_products WHERE id=$1 AND cart_id=$2", p.ID, cartID).Scan(&sum)
	// Create the product.
	if sum == 0 {
		_, err := db.ExecContext(ctx, productsQuery, p.ID, cartID, p.Quantity, p.Brand, p.Category, p.Type, p.Description,
			p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
		if err != nil {
			logger.Log.Errorf("failed creating cart's products: %v", err)
			return nil, errors.Wrap(err, "couldn't create the product")
		}
	}
	// Update the quantity.
	if sum != 0 {
		p.Quantity += sum

		_, err := db.ExecContext(ctx, "UPDATE cart_products SET quantity=$2 WHERE cart_id=$1", cartID, p.Quantity)
		if err != nil {
			logger.Log.Errorf("failed updating cart's products: %v", err)
			return nil, errors.Wrap(err, "couldn't update the product")
		}
	}

	_, err := db.ExecContext(ctx, cartQuery, cartID, cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Subtotal,
		cart.Total)
	if err != nil {
		logger.Log.Errorf("failed updating the cart: %v", err)
		return nil, errors.Wrap(err, "couldn't update the cart")
	}

	return p, nil
}

// Checkout returns the cart total.
func Checkout(ctx context.Context, db *sqlx.DB, cartID string) (int64, error) {
	var cart Cart

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart: %v", err)
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	total := cart.Total + cart.Taxes - cart.Discount

	return total, nil
}

// Delete permanently deletes a cart from the database.
func Delete(ctx context.Context, db *sqlx.DB, cartID string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM carts WHERE id=$1", cartID)
	if err != nil {
		logger.Log.Errorf("failed deleting cart: %v", err)
		return errors.New("couldn't delete the cart")
	}

	return nil
}

// Get returns the user cart.
func Get(ctx context.Context, db *sqlx.DB, cartID string) (*Cart, error) {
	var (
		cart     Cart
		products []Product
	)

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart: %v", err)
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	if err := db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart's products: %v", err)
		return nil, errors.Wrap(err, "couldn't find the cart products")
	}

	cart.Products = products

	return &cart, nil
}

// Products returns the cart products.
func Products(ctx context.Context, db *sqlx.DB, cartID string) ([]Product, error) {
	var products []Product

	if err := db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart's products: %v", err)
		return nil, errors.Wrap(err, "couldn't find the cart products")
	}

	if len(products) == 0 {
		return nil, errors.New("cart is empty")
	}

	return products, nil
}

// Remove takes away the specified quantity of products from the cart.
func Remove(ctx context.Context, db *sqlx.DB, cartID string, pID string, quantity int) error {
	var (
		cart Cart
		p    Product
	)

	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		logger.Log.Errorf("failed deleting cart: %v", err)
		return errors.Wrap(err, "couldn't find the cart")
	}

	if err := db.GetContext(ctx, &p, "SELECT * FROM cart_products WHERE id = $1 AND cart_id=$2", pID, cartID); err != nil {
		logger.Log.Errorf("failed listing cart's products: %v", err)
		return errors.New("couldn't find the product")
	}

	if quantity > p.Quantity {
		return errors.Errorf("quantity inserted (%d) is higher than the stock of products (%d)", quantity, p.Quantity)
	}

	if quantity == p.Quantity {
		_, err := db.ExecContext(ctx, "DELETE FROM cart_products WHERE id=$1 AND cart_id=$2", pID, cartID)
		if err != nil {
			logger.Log.Errorf("failed deleting cart's products: %v", err)
			return errors.Wrap(err, "couldn't delete the product")
		}
	}

	if cart.Counter == 1 {
		if err := Reset(ctx, db, cartID); err != nil {
			return err
		}
		return nil
	}

	taxes := (p.Subtotal / 100) * p.Taxes
	discount := (p.Subtotal / 100) * p.Discount

	cart.mu.Lock()
	for i := 0; i < quantity; i++ {
		cart.Counter--
		p.Quantity--
		cart.Weight -= p.Weight
		cart.Discount -= discount
		cart.Taxes -= taxes
		cart.Subtotal -= p.Subtotal
		cart.Total = cart.Total - p.Subtotal - taxes + discount
	}
	cart.mu.Unlock()

	_, err := db.ExecContext(ctx, cQuery, cartID, cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Subtotal,
		cart.Total)
	if err != nil {
		logger.Log.Errorf("failed updating cart: %v", err)
		return errors.Wrap(err, "couldn't update the cart")
	}

	return nil
}

// Reset sets cart values to default.
func Reset(ctx context.Context, db *sqlx.DB, cartID string) error {
	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	_, err := db.ExecContext(ctx, "DELETE FROM cart_products WHERE cart_id=$1", cartID)
	if err != nil {
		logger.Log.Errorf("failed deleting cart's products: %v", err)
		return errors.Wrap(err, "couldn't delete cart products")
	}

	_, err = db.ExecContext(ctx, cQuery, cartID, 0, 0, 0, 0, 0, 0)
	if err != nil {
		logger.Log.Errorf("failed resetting cart: %v", err)
		return errors.Wrap(err, "couldn't reset the cart")
	}

	return nil
}

// Size returns the quantity of products inside the cart.
func Size(ctx context.Context, db *sqlx.DB, cartID string) (int, error) {
	var cart Cart

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart: %v", err)
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	return cart.Counter, nil
}

// String returns a string with the cart details.
func String(ctx context.Context, db *sqlx.DB, cartID string) (string, error) {
	var c Cart

	if err := db.GetContext(ctx, &c, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		logger.Log.Errorf("failed listing cart: %v", err)
		return "", errors.Wrap(err, "couldn't find the cart")
	}

	const details = `The cart has %d products, a weight of %2.dkg, $%2.d of discounts, 
	$%2.d of taxes and a total of $%2.d`

	return fmt.Sprintf(details, c.Counter, (c.Weight / 1000), (c.Discount / 100), (c.Taxes / 100), (c.Total / 100)), nil
}
