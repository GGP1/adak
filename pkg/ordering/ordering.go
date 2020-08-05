// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"math/rand"
	"time"

	"github.com/GGP1/palo/pkg/shopping"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Order represents the user purchase request.
type Order struct {
	ID           int            `json:"order_id"`
	UserID       string         `json:"user_id"`
	OrderedAt    time.Time      `json:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date"`
	CartID       string         `json:"cart_id"`
	Cart         OrderCart      `json:"cart" gorm:"foreignkey:OrderID"`
	Products     []OrderProduct `json:"products" gorm:"foreignkey:OrderID"`
}

// OrderCart represents the cart ordered by the user.
type OrderCart struct {
	OrderID  int     `json:"order_id"`
	Counter  int     `json:"counter"`
	Weight   float32 `json:"weight"`
	Discount float32 `json:"discount"`
	Taxes    float32 `json:"taxes"`
	Subtotal float32 `json:"subtotal"`
	Total    float32 `json:"total"`
}

// OrderProduct represents the a product place into the cart ordered by the user.
type OrderProduct struct {
	OrderID     int     `json:"order_id"`
	ID          int     `json:"id"`
	Quantity    int     `json:"quantity"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Weight      float32 `json:"weight"`
	Discount    float32 `json:"discount"`
	Taxes       float32 `json:"taxes"`
	Subtotal    float32 `json:"subtotal"`
	Total       float32 `json:"total"`
}

// Delete removes an order.
func Delete(db *gorm.DB, orderID int) error {
	var order Order
	var orderCart OrderCart
	var orderProduct OrderProduct

	err := db.Delete(&order, orderID).Error
	if err != nil {
		return errors.Wrap(err, "couldn't delete the order")
	}

	err = db.Where("order_id=?", orderID).Delete(&orderCart).Error
	if err != nil {
		return errors.Wrap(err, "couldn't delete the order cart")
	}

	err = db.Where("order_id=?", orderID).Delete(&orderProduct).Error
	if err != nil {
		return errors.Wrap(err, "couldn't delete the order products")
	}

	return nil
}

// Get removes an order.
func Get(db *gorm.DB, orders *[]Order) error {
	err := db.Preload("Cart").Preload("Products").Find(&orders).Error
	if err != nil {
		return errors.Wrap(err, "couldn't find the orders")
	}

	return nil
}

// New creates an order.
func New(db *gorm.DB, userID string, cart shopping.Cart, deliveryDate time.Time) (*Order, error) {
	var order Order

	if cart.Counter == 0 {
		return nil, errors.New("ordering 0 products is not permitted")
	}

	id := rand.Intn(2147483647)

	// Set fields
	order.ID = id
	order.UserID = userID
	order.OrderedAt = time.Now()
	order.DeliveryDate = deliveryDate

	order.CartID = cart.ID
	order.Cart.Counter = cart.Counter
	order.Cart.Weight = cart.Weight
	order.Cart.Discount = cart.Discount
	order.Cart.Taxes = cart.Taxes
	order.Cart.Subtotal = cart.Subtotal
	order.Cart.Total = cart.Total

	for _, product := range cart.Products {
		var orderP OrderProduct

		orderP.ID = product.ID
		orderP.Quantity = product.Quantity
		orderP.Brand = product.Brand
		orderP.Category = product.Category
		orderP.Type = product.Type
		orderP.Description = product.Description
		orderP.Discount = product.Discount
		orderP.Weight = product.Weight
		orderP.Taxes = product.Taxes
		orderP.Subtotal = product.Subtotal
		orderP.Total = product.Total

		err := db.Create(&orderP).Error
		if err != nil {
			return nil, err
		}

		order.Products = append(order.Products, orderP)
	}

	err := db.Create(&order).Error
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order")
	}

	err = shopping.Reset(db, cart.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't reset the cart")
	}
	return &order, nil
}
