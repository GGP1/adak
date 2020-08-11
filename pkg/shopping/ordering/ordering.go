// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"math/rand"
	"sync"
	"time"

	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/shopping/wallet"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
	ID           int            `json:"order_id"`
	UserID       string         `json:"user_id"`
	OrderedAt    time.Time      `json:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date"`
	State        string         `json:"state"`
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
	ID          string  `json:"id"`
	OrderID     int     `json:"order_id"`
	ProductID   int     `json:"product_id"`
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

// NewOrder creates an order.
func NewOrder(db *gorm.DB, userID string, cart shopping.Cart, deliveryDate time.Time) (*Order, error) {
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
	order.State = PendingState

	order.CartID = cart.ID
	order.Cart.Counter = cart.Counter
	order.Cart.Weight = cart.Weight
	order.Cart.Discount = cart.Discount
	order.Cart.Taxes = cart.Taxes
	order.Cart.Subtotal = cart.Subtotal
	order.Cart.Total = cart.Total

	var wg sync.WaitGroup

	for _, product := range cart.Products {
		var orderP OrderProduct
		wg.Add(1)

		go func(orderP OrderProduct, products []OrderProduct, product *shopping.CartProduct) {
			defer wg.Done()

			id := uuid.New()

			orderP.ID = id.String()
			orderP.ProductID = product.ID
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
			db.Create(&orderP)

			order.Products = append(order.Products, orderP)
		}(orderP, order.Products, product)
	}
	wg.Wait()

	err := db.Create(&order).Error
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the order")
	}

	_, err = wallet.SubtractFunds(db, cart.ID, float64(order.Cart.Total))
	if err != nil {
		return nil, err
	}

	err = shopping.Reset(db, cart.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't reset the cart")
	}

	return &order, nil
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
