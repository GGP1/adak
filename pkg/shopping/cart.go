package shopping

import (
	"fmt"
	"math"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
	Products []*CartProduct `json:"products" gorm:"foreignkey:CartID"`
}

// CartProduct represents a product that has been appended to the cart.
type CartProduct struct {
	CartID      string  `json:"cart_id"`
	ID          int     `json:"id"`
	Quantity    int     `json:"quantity"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Taxes       float64 `json:"taxes"`
	Discount    float64 `json:"discount"`
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
func Add(db *gorm.DB, cartID string, product *CartProduct, quantity int) (*CartProduct, error) {
	var cart Cart

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	product.CartID = cartID
	taxes := ((product.Subtotal / 100) * product.Taxes)
	discount := ((product.Subtotal / 100) * product.Discount)
	product.Total = product.Total + product.Subtotal + taxes - discount

	// math.Ceil(x*100)/100 is used to round float numbers
	for i := 0; i < quantity; i++ {
		cart.Counter++
		product.Quantity++
		cart.Weight += math.Ceil(product.Weight*100) / 100
		cart.Discount += math.Ceil(discount*100) / 100
		cart.Taxes += math.Ceil(taxes*100) / 100
		cart.Subtotal += math.Ceil(product.Subtotal*100) / 100
		cart.Total = cart.Total + product.Subtotal + taxes - discount
	}

	alreadyExists := db.Where("id=?", product.ID).First(&product).RowsAffected
	if alreadyExists != 0 {
		product.Quantity++
		err := db.Save(&product).Error
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update the product")
		}
	} else {
		err := db.Create(&product).Error
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create the product")
		}
	}

	err := db.Save(&cart).Error
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the cart")
	}

	return product, nil
}

// Checkout takes all the products and returns the total price.
func Checkout(db *gorm.DB, cartID string) (float64, error) {
	var cart Cart

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	total := cart.Total + cart.Taxes - cart.Discount

	return total, nil
}

// DeleteCart takes a cart from the database and permanently deletes it.
func DeleteCart(db *gorm.DB, id string) error {
	var cart Cart

	if err := db.Delete(cart, id).Error; err != nil {
		return fmt.Errorf("couldn't delete the cart")
	}
	return nil
}

// Get returns the user cart.
func Get(db *gorm.DB, cartID string) (Cart, error) {
	var cart Cart

	if err := db.Preload("Products").First(&cart, "id=?", cartID).Error; err != nil {
		return Cart{}, errors.Wrap(err, "couldn't find the cart")
	}

	return cart, nil
}

// Items prints cart products.
func Items(db *gorm.DB, cartID string) ([]CartProduct, error) {
	var cart Cart
	var list []CartProduct

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the cart")
	}

	for _, v := range cart.Products {
		if v != nil {
			list = append(list, *v)
		}
	}

	if len(list) == 0 {
		return nil, errors.New("cart is empty")
	}

	return list, nil
}

// Remove takes away the specified quantity of products from the cart.
func Remove(db *gorm.DB, cartID string, key int, quantity int) error {
	var cart Cart
	var product CartProduct

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return errors.Wrap(err, "couldn't find the cart")
	}

	if err := db.Where("cart_id = ? AND id = ?", cartID, key).Find(&product).Error; err != nil {
		return errors.New("product not found")
	}

	if quantity > product.Quantity {
		return fmt.Errorf("quantity inserted: %d\nis higher than the stock of products %d", quantity, product.Quantity)
	}

	if quantity == product.Quantity {
		err := db.Where("cart_id=?", cartID).Delete(&product, "id=?", key).Error
		if err != nil {
			return errors.Wrap(err, "couldn't delete the product")
		}
	}

	if cart.Counter == 1 {
		err := Reset(db, cartID)
		if err != nil {
			return err
		}
		return nil
	}

	taxes := (product.Subtotal / 100) * product.Taxes
	discount := (product.Subtotal / 100) * product.Discount

	// math.Ceil(x*100)/100 is used to round float numbers
	for i := 0; i < quantity; i++ {
		cart.Counter--
		product.Quantity--
		cart.Weight -= math.Ceil(product.Weight*100) / 100
		cart.Discount -= math.Ceil(discount*100) / 100
		cart.Taxes -= math.Ceil(taxes*100) / 100
		cart.Subtotal -= math.Ceil(product.Subtotal*100) / 100
		cart.Total = cart.Total - product.Subtotal - taxes + discount
	}

	if err := db.Save(&cart).Error; err != nil {
		return errors.Wrap(err, "couldn't update the cart")
	}

	return nil
}

// Reset sets the cart to its default values.
func Reset(db *gorm.DB, cartID string) error {
	var cart Cart
	var product CartProduct

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return errors.Wrap(err, "couldn't find the cart")
	}

	if err := db.Where("cart_id=?", cartID).Delete(&product).Error; err != nil {
		return errors.Wrap(err, "couldn't delete the product")
	}

	cart.Counter = 0
	cart.Weight = 0
	cart.Discount = 0
	cart.Taxes = 0
	cart.Subtotal = 0
	cart.Total = 0

	err := db.Save(&cart).Error
	if err != nil {
		return errors.Wrap(err, "couldn't update the cart")
	}

	return nil
}

// Size returns the quantity of products in the cart.
func Size(db *gorm.DB, cartID string) (int, error) {
	var cart Cart

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return 0, errors.Wrap(err, "couldn't find the cart")
	}

	return cart.Counter, nil
}

// String returns a string with the cart details.
func String(db *gorm.DB, cartID string) (string, error) {
	var cart Cart

	if err := db.Where("id=?", cartID).Find(&cart).Error; err != nil {
		return "", errors.Wrap(err, "couldn't find the cart")
	}

	return fmt.Sprintf(
		"The cart has %d products, a weight of %2.fkg, $%2.f of discounts, $%2.f of taxes and a total of $%2.f",
		cart.Counter, cart.Weight, cart.Discount, cart.Taxes, cart.Total), nil
}
