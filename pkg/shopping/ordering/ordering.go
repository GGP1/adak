// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"fmt"
	"sync"
	"time"

	"github.com/GGP1/palo/internal/uuid"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/jinzhu/gorm"
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
	ID           string         `json:"order_id"`
	UserID       string         `json:"user_id"`
	Currency     string         `json:"currency"`
	Address      string         `json:"address"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	ZipCode      string         `json:"zip_code"`
	Country      string         `json:"country"`
	Status       string         `json:"status"`
	OrderedAt    time.Time      `json:"ordered_at"`
	DeliveryDate time.Time      `json:"delivery_date"`
	CartID       string         `json:"cart_id"`
	Cart         OrderCart      `json:"cart" gorm:"foreignkey:OrderID"`
	Products     []OrderProduct `json:"products" gorm:"foreignkey:OrderID"`
}

// OrderCart represents the cart ordered by the user.
type OrderCart struct {
	OrderID  string  `json:"order_id"`
	Counter  int     `json:"counter"`
	Weight   float64 `json:"weight"`
	Discount float64 `json:"discount"`
	Taxes    float64 `json:"taxes"`
	Subtotal float64 `json:"subtotal"`
	Total    float64 `json:"total"`
}

// OrderProduct represents the a product place into the cart ordered by the user.
type OrderProduct struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   int     `json:"product_id"`
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

// NewOrder creates an order.
func NewOrder(db *gorm.DB, userID, currency, address, city, country, state, zipcode string, deliveryDate time.Time, cart shopping.Cart) (*Order, error) {
	var order Order

	if cart.Counter == 0 {
		return nil, fmt.Errorf("ordering 0 products is not permitted")
	}

	id := uuid.GenerateRandRunes(24)

	// Set fields
	order.ID = id
	order.UserID = userID
	order.Currency = currency
	order.Address = address
	order.City = city
	order.State = state
	order.ZipCode = zipcode
	order.Country = country
	order.Status = PendingState
	order.OrderedAt = time.Now()
	order.DeliveryDate = deliveryDate

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

			id := uuid.GenerateRandRunes(20)

			orderP.ID = id
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
		return nil, fmt.Errorf("couldn't create the order: %v", err)
	}

	err = shopping.Reset(db, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("couldn't reset the cart: %v", err)
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
		return fmt.Errorf("couldn't delete the order: %v", err)
	}

	err = db.Where("order_id=?", orderID).Delete(&orderCart).Error
	if err != nil {
		return fmt.Errorf("couldn't delete the order cart: %v", err)
	}

	err = db.Where("order_id=?", orderID).Delete(&orderProduct).Error
	if err != nil {
		return fmt.Errorf("couldn't delete the order products: %v", err)
	}

	return nil
}

// Get removes an order.
func Get(db *gorm.DB, orders *[]Order) error {
	err := db.Preload("Cart").Preload("Products").Find(&orders).Error
	if err != nil {
		return fmt.Errorf("couldn't find the orders: %v", err)
	}

	return nil
}
