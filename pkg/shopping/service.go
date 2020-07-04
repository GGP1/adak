package shopping

import (
	"sync"

	"github.com/GGP1/palo/pkg/model"
)

// Cart stores the products that the user chose to buy
type Cart struct {
	Products map[uint]model.Product
	Price    float32
	Weight   float32
	Discount float32
	Tax      float32
	Total    float32
	mutex    sync.RWMutex
}

// NewCart returns a cart with the default values
func NewCart() *Cart {
	return &Cart{
		Products: make(map[uint]model.Product),
		Price:    0,
		Weight:   0,
		Discount: 0,
		Tax:      0,
		Total:    0,
	}
}

// Add a product to the cart
func (c *Cart) Add(product *model.Product) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var tax float32
	var discount float32

	tax = (product.Subtotal / 100) * product.Taxes
	discount = (product.Subtotal / 100) * product.Discount

	c.Products[product.ID] = *product
	c.Price = c.Price + product.Subtotal
	c.Weight = c.Weight + product.Weight
	c.Discount = c.Discount + discount
	c.Tax = c.Tax + tax
	c.Total = c.Total + product.Total
}

// Checkout takes all the products and returns the final purchase
func (c *Cart) Checkout() float32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	total := c.Total + c.Tax - c.Discount

	return total
}

// Remove takes away a product from the cart
func (c *Cart) Remove(key uint) {
	c.mutex.Lock()
	delete(c.Products, key)
	c.mutex.Unlock()
}

// Reset cart products
func (c *Cart) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Delete all the key/values from the map
	for key := range c.Products {
		delete(c.Products, key)
	}

	// Set cart variables to 0
	c.Price = 0
	c.Weight = 0
	c.Discount = 0
	c.Tax = 0
	c.Total = 0
}

// ShowItems prints cart items
func (c *Cart) ShowItems() *Cart {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c
}

// Size returns the amount of products in the cart
func (c *Cart) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.Products)
}
