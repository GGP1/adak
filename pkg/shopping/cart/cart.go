package cart

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Cart represents a temporary record of items that the customer
// selected for purchase.
type Cart struct {
	ID       string     `json:"id"`
	Counter  int        `json:"counter"` // Counter contains the quantity of products placed in the cart
	Weight   float64    `json:"weight"`
	Discount float64    `json:"discount"`
	Taxes    float64    `json:"taxes"`
	Subtotal float64    `json:"subtotal"`
	Total    float64    `json:"total"`
	Products []*Product `json:"products"`
}

// Product represents a product that has been appended to the cart.
type Product struct {
	ID          string  `json:"id"`
	CartID      string  `json:"cart_id" db:"cart_id"`
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
		Products: []*Product{},
	}
}

// Add adds a product to the cart.
func Add(ctx context.Context, db *sqlx.DB, cartID string, p *Product, quantity int) (*Product, error) {
	var (
		cart Cart
		sum  int
	)

	pQuery := `INSERT INTO cart_products
	(id, cart_id, quantity, brand, category, type, description, weight, 
	discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	taxes := ((p.Subtotal / 100) * p.Taxes)
	discount := ((p.Subtotal / 100) * p.Discount)

	p.CartID = cartID
	p.Total = p.Total + p.Subtotal + taxes - discount

	// math.Ceil(x*100)/100 is used to round float numbers
	for i := 0; i < quantity; i++ {
		cart.Counter++
		p.Quantity++
		cart.Weight += math.Ceil(p.Weight*100) / 100
		cart.Discount += math.Ceil(discount*100) / 100
		cart.Taxes += math.Ceil(taxes*100) / 100
		cart.Subtotal += math.Ceil(p.Subtotal*100) / 100
		cart.Total = cart.Total + p.Subtotal + taxes - discount
	}

	db.QueryRow("SELECT SUM(quantity) FROM cart_products WHERE id=$1 AND cart_id=$2", p.ID, cartID).Scan(&sum)
	// If sum == 0 (does not exist a product with the same id and cart_id), create the product.
	if sum == 0 {
		_, err := db.ExecContext(ctx, pQuery, p.ID, cartID, p.Quantity, p.Brand, p.Category, p.Type, p.Description,
			p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create the product")
		}
	}
	// If sum != 0 (product already exists), update the quantity.
	if sum != 0 {
		p.Quantity += sum

		_, err := db.ExecContext(ctx, "UPDATE cart_products SET quantity=$2 WHERE cart_id=$1", cartID, p.Quantity)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update the product")
		}
	}

	_, err := db.ExecContext(ctx, cQuery, cartID, cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Subtotal,
		cart.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the cart")
	}

	return p, nil
}

// Checkout returns the cart total.
func Checkout(ctx context.Context, db *sqlx.DB, cartID string) (float64, error) {
	var cart Cart

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	total := cart.Total + cart.Taxes - cart.Discount

	return total, nil
}

// Delete permanently deletes a cart from the database.
func Delete(ctx context.Context, db *sqlx.DB, cartID string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM carts WHERE id=$1", cartID)
	if err != nil {
		return errors.New("couldn't delete the cart")
	}

	return nil
}

// Get returns the user cart.
func Get(ctx context.Context, db *sqlx.DB, cartID string) (Cart, error) {
	var (
		cart     Cart
		products []*Product
	)

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return Cart{}, errors.Wrap(err, "couldn't find the cart")
	}

	if err := db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		return Cart{}, errors.Wrap(err, "couldn't find cart products")
	}

	cart.Products = products

	return cart, nil
}

// Products returns the cart products.
func Products(ctx context.Context, db *sqlx.DB, cartID string) ([]Product, error) {
	var products []Product

	if err := db.SelectContext(ctx, &products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
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
		return errors.Wrap(err, "couldn't find the cart")
	}

	if err := db.GetContext(ctx, &p, "SELECT * FROM cart_products WHERE id = $1 AND cart_id=$2", pID, cartID); err != nil {
		return errors.New("product not found")
	}

	if quantity > p.Quantity {
		return fmt.Errorf("quantity inserted (%d) is higher than the stock of products (%d)", quantity, p.Quantity)
	}

	if quantity == p.Quantity {
		_, err := db.ExecContext(ctx, "DELETE FROM cart_products WHERE id=$1 AND cart_id=$2", pID, cartID)
		if err != nil {
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

	// math.Ceil(x*100)/100 is used to round float numbers
	for i := 0; i < quantity; i++ {
		cart.Counter--
		p.Quantity--
		cart.Weight -= math.Ceil(p.Weight*100) / 100
		cart.Discount -= math.Ceil(discount*100) / 100
		cart.Taxes -= math.Ceil(taxes*100) / 100
		cart.Subtotal -= math.Ceil(p.Subtotal*100) / 100
		cart.Total = cart.Total - p.Subtotal - taxes + discount
	}

	_, err := db.ExecContext(ctx, cQuery, cartID, cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Subtotal,
		cart.Total)
	if err != nil {
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
		return errors.Wrap(err, "couldn't delete cart products")
	}

	_, err = db.ExecContext(ctx, cQuery, cartID, 0, 0, 0, 0, 0, 0)
	if err != nil {
		return errors.Wrap(err, "couldn't reset the cart")
	}

	return nil
}

// Size returns the quantity of products inside the cart.
func Size(ctx context.Context, db *sqlx.DB, cartID string) (int, error) {
	var cart Cart

	if err := db.GetContext(ctx, &cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	return cart.Counter, nil
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
