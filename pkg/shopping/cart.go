package shopping

import (
	"errors"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
)

// Cart represents a temporary record of items that the customer
// selected for purchase.
type Cart struct {
	ID string `json:"id"`
	// Counter contains the quantity of products placed in the cart
	Counter  int            `json:"counter"`
	Weight   float64        `json:"weight"`
	Discount float64        `json:"discount"`
	Taxes    float64        `json:"taxes"`
	Subtotal float64        `json:"subtotal"`
	Total    float64        `json:"total"`
	Products []*CartProduct `json:"products"`
}

// CartProduct represents a product that has been appended to the cart.
type CartProduct struct {
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

// NewCart returns a cart with the default values.
func NewCart(id string) *Cart {
	return &Cart{
		ID:       id,
		Counter:  0,
		Weight:   0,
		Discount: 0,
		Taxes:    0,
		Subtotal: 0,
		Total:    0,
		Products: []*CartProduct{},
	}
}

// Add adds a product to the cart.
func Add(db *sqlx.DB, cartID string, p *CartProduct, quantity int) (*CartProduct, error) {
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

	if err := db.Get(&cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return nil, fmt.Errorf("couldn't find the cart: %v", err)
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

	if sum == 0 {
		_, err := db.Exec(pQuery, p.ID, cartID, p.Quantity, p.Brand, p.Category, p.Type, p.Description,
			p.Weight, p.Discount, p.Taxes, p.Subtotal, p.Total)
		if err != nil {
			return nil, fmt.Errorf("couldn't create the product: %v", err)
		}
	}

	if sum != 0 {
		p.Quantity += sum

		_, err := db.Exec("UPDATE cart_products SET quantity=$2 WHERE cart_id=$1", cartID, p.Quantity)
		if err != nil {
			return nil, fmt.Errorf("couldn't update the product: %v", err)
		}
	}

	_, err := db.Exec(cQuery, cartID, cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Subtotal,
		cart.Total)
	if err != nil {
		return nil, fmt.Errorf("couldn't update the cart: %v", err)
	}

	return p, nil
}

// Checkout takes all the products and returns the total price.
func Checkout(db *sqlx.DB, cartID string) (float64, error) {
	var cart Cart

	if err := db.Get(&cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return 0, fmt.Errorf("couldn't find the cart: %v", err)
	}

	total := cart.Total + cart.Taxes - cart.Discount

	return total, nil
}

// DeleteCart takes a cart from the database and permanently deletes it.
func DeleteCart(db *sqlx.DB, cartID string) error {
	_, err := db.Exec("DELETE FROM carts WHERE id=$1", cartID)
	if err != nil {
		return errors.New("couldn't delete the cart")
	}

	return nil
}

// Get returns the user cart.
func Get(db *sqlx.DB, cartID string) (Cart, error) {
	var (
		cart     Cart
		products []*CartProduct
	)

	if err := db.Get(&cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return Cart{}, fmt.Errorf("couldn't find the cart: %v", err)
	}

	if err := db.Select(&products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		return Cart{}, fmt.Errorf("couldn't find cart products: %v", err)
	}

	cart.Products = products

	return cart, nil
}

// Items prints cart products.
func Items(db *sqlx.DB, cartID string) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Select(&products, "SELECT * FROM cart_products WHERE cart_id=$1", cartID); err != nil {
		return nil, fmt.Errorf("couldn't find the cart: %v", err)
	}

	if len(products) == 0 {
		return nil, errors.New("cart is empty")
	}

	return products, nil
}

// Remove takes away the specified quantity of products from the cart.
func Remove(db *sqlx.DB, cartID string, pID string, quantity int) error {
	var (
		cart Cart
		p    CartProduct
	)

	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	if err := db.Get(&cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return fmt.Errorf("couldn't find the cart: %v", err)
	}

	if err := db.Get(&p, "SELECT * FROM cart_products WHERE id = $1 AND cart_id=$2", pID, cartID); err != nil {
		return errors.New("product not found")
	}

	if quantity > p.Quantity {
		return fmt.Errorf("quantity inserted (%d) is higher than the stock of products (%d)", quantity, p.Quantity)
	}

	if quantity == p.Quantity {
		_, err := db.Exec("DELETE FROM cart_products WHERE id=$1 AND cart_id=$2", pID, cartID)
		if err != nil {
			return fmt.Errorf("couldn't delete the product: %v", err)
		}
	}

	if cart.Counter == 1 {
		if err := Reset(db, cartID); err != nil {
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

	_, err := db.Exec(cQuery, cartID, cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Subtotal,
		cart.Total)
	if err != nil {
		return fmt.Errorf("couldn't update the cart: %v", err)
	}

	return nil
}

// Reset sets the cart to its default values.
func Reset(db *sqlx.DB, cartID string) error {
	cQuery := `UPDATE carts SET counter=$2, weight=$3, discount=$4, taxes=$5, 
	subtotal=$6, total=$7 WHERE id=$1`

	_, err := db.Exec("DELETE FROM cart_products WHERE cart_id=$1", cartID)
	if err != nil {
		return fmt.Errorf("couldn't delete the product: %v", err)
	}

	// Set cart values to 0
	_, err = db.Exec(cQuery, cartID, 0, 0, 0, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("couldn't update the cart: %v", err)
	}

	return nil
}

// Size returns the quantity of products in the cart.
func Size(db *sqlx.DB, cartID string) (int, error) {
	var cart Cart

	if err := db.Get(&cart, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return 0, fmt.Errorf("couldn't find the cart: %v", err)
	}

	return cart.Counter, nil
}

// String returns a string with the cart details.
func String(db *sqlx.DB, cartID string) (string, error) {
	var c Cart

	if err := db.Get(&c, "SELECT * FROM carts WHERE id=$1", cartID); err != nil {
		return "", fmt.Errorf("couldn't find the cart: %v", err)
	}

	const details = `The cart has %d products, a weight of %2.fkg, $%2.f of discounts, 
	$%2.f of taxes and a total of $%2.f`

	return fmt.Sprintf(details, c.Counter, c.Weight, c.Discount, c.Taxes, c.Total), nil
}
