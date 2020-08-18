// Package ordering provides users interfaces for requesting products.
package ordering

import (
	"errors"
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
	if cart.Counter == 0 {
		return nil, errors.New("ordering 0 products is not permitted")
	}

	id := uuid.GenerateRandRunes(24)

	order := Order{
		ID:           id,
		UserID:       userID,
		Currency:     currency,
		Address:      address,
		City:         city,
		State:        state,
		ZipCode:      zipcode,
		Status:       PendingState,
		OrderedAt:    time.Now(),
		DeliveryDate: deliveryDate,
		CartID:       cart.ID,
		Cart: OrderCart{
			Counter:  cart.Counter,
			Weight:   cart.Weight,
			Discount: cart.Discount,
			Taxes:    cart.Taxes,
			Subtotal: cart.Subtotal,
			Total:    cart.Total,
		},
	}

	var wg sync.WaitGroup

	for _, product := range cart.Products {
		wg.Add(1)

		go func(products []OrderProduct, product *shopping.CartProduct) {
			defer wg.Done()

			id := uuid.GenerateRandRunes(20)

			orderP := OrderProduct{
				ID:          id,
				ProductID:   product.ID,
				Quantity:    product.Quantity,
				Brand:       product.Brand,
				Category:    product.Category,
				Type:        product.Type,
				Description: product.Description,
				Weight:      product.Weight,
				Discount:    product.Discount,
				Taxes:       product.Taxes,
				Subtotal:    product.Subtotal,
				Total:       product.Total,
			}

			db.Create(&orderP)

			order.Products = append(order.Products, orderP)
		}(order.Products, product)
	}
	wg.Wait()

	if err := db.Create(&order).Error; err != nil {
		return nil, fmt.Errorf("couldn't create the order: %v", err)
	}

	if err := shopping.Reset(db, cart.ID); err != nil {
		return nil, fmt.Errorf("couldn't reset the cart: %v", err)
	}

	return &order, nil
}

// Delete removes an order.
func Delete(db *gorm.DB, orderID int) error {
	var (
		order        Order
		orderCart    OrderCart
		orderProduct OrderProduct
	)

	if err := db.Delete(&order, orderID).Error; err != nil {
		return fmt.Errorf("couldn't delete the order: %v", err)
	}

	if err := db.Where("order_id=?", orderID).Delete(&orderCart).Error; err != nil {
		return fmt.Errorf("couldn't delete the order cart: %v", err)
	}

	if err := db.Where("order_id=?", orderID).Delete(&orderProduct).Error; err != nil {
		return fmt.Errorf("couldn't delete the order products: %v", err)
	}

	return nil
}

// Get removes an order.
func Get(db *gorm.DB, orders *[]Order) error {
	if err := db.Preload("Cart").Preload("Products").Find(&orders).Error; err != nil {
		return fmt.Errorf("couldn't find the orders: %v", err)
	}

	return nil
}
