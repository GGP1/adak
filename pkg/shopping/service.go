package shopping

import (
	"errors"
	"strings"
	"sync"

	"github.com/GGP1/palo/pkg/model"
)

var (
	errNotFound = errors.New("no products found")
)

// Cart stores the products that the user chose to buy
type Cart struct {
	Products map[uint]model.Product
	Weight   float32
	Discount float32
	Tax      float32
	Subtotal float32
	Total    float32
	sync.RWMutex
}

// NewCart returns a cart with the default values
func NewCart() *Cart {
	return &Cart{
		Products: make(map[uint]model.Product),
		Weight:   0,
		Discount: 0,
		Tax:      0,
		Total:    0,
	}
}

// Add a product to the cart
func (c *Cart) Add(product *model.Product) {
	c.Lock()
	defer c.Unlock()

	var tax float32
	var discount float32

	tax = (product.Subtotal / 100) * product.Taxes
	discount = (product.Subtotal / 100) * product.Discount

	c.Products[product.ID] = *product
	c.Weight = c.Weight + product.Weight
	c.Discount = c.Discount + discount
	c.Tax = c.Tax + tax
	c.Subtotal = c.Subtotal + product.Subtotal
	c.Total = c.Total + c.Subtotal + c.Tax - c.Discount
}

// Checkout takes all the products and returns the final purchase
func (c *Cart) Checkout() float32 {
	c.Lock()
	defer c.Unlock()

	total := c.Total + c.Tax - c.Discount

	return total
}

// FilterByBrand looks for products with the specified brand
func (c *Cart) FilterByBrand(productBrand string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	for k, v := range c.Products {
		if strings.ToLower(productBrand) == strings.ToLower(v.Brand) {
			var products []model.Product
			products = append(products, c.Products[k])

			return products, nil
		}
	}

	return nil, errNotFound
}

// FilterByCategory looks for products with the specified category
func (c *Cart) FilterByCategory(productCategory string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	for k, v := range c.Products {
		if strings.ToLower(productCategory) == strings.ToLower(v.Category) {
			var products []model.Product
			products = append(products, c.Products[k])

			return products, nil
		}
	}

	return nil, errNotFound
}

// FilterByTotal looks for products within the total price range specified
func (c *Cart) FilterByTotal(minTotal, maxTotal float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	for k, v := range c.Products {
		if v.Total >= minTotal && v.Total <= maxTotal {
			var products []model.Product
			products = append(products, c.Products[k])

			return products, nil
		}
	}

	return nil, errNotFound
}

// FilterByType looks for products with the specified type
func (c *Cart) FilterByType(productType string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	for k, v := range c.Products {
		if strings.ToLower(productType) == strings.ToLower(v.Type) {
			var products []model.Product
			products = append(products, c.Products[k])

			return products, nil
		}
	}

	return nil, errNotFound
}

// FilterByWeight looks for products within the weight range specified
func (c *Cart) FilterByWeight(minWeight, maxWeight float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	for k, v := range c.Products {
		if v.Weight >= minWeight && v.Weight <= maxWeight {
			var products []model.Product
			products = append(products, c.Products[k])

			return products, nil
		}
	}

	return nil, errNotFound
}

// Remove takes away a product from the cart
func (c *Cart) Remove(key uint) {
	c.Lock()
	delete(c.Products, key)
	c.Unlock()
}

// Reset cart products
func (c *Cart) Reset() {
	c.Lock()
	defer c.Unlock()

	// Delete all the key/values from the map
	for key := range c.Products {
		delete(c.Products, key)
	}

	// Set cart variables to 0
	c.Weight = 0
	c.Discount = 0
	c.Tax = 0
	c.Subtotal = 0
	c.Total = 0
}

// ShowItems prints cart items
func (c *Cart) ShowItems() map[uint]model.Product {
	c.RLock()
	defer c.RUnlock()

	return c.Products
}

// Size returns the amount of products in the cart
func (c *Cart) Size() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.Products)
}
