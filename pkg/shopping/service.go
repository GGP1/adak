package shopping

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/GGP1/palo/pkg/model"
)

var (
	errNotFound = errors.New("no products found")
)

// Cart stores the products that the user chose to buy
type Cart struct {
	sync.RWMutex
	Error error

	Products map[uint]model.Product
	Weight   float32
	Discount float32
	Taxes    float32
	Subtotal float32
	Total    float32
}

// NewCart returns a cart with the default values
func NewCart() *Cart {
	return &Cart{
		Error:    nil,
		Products: make(map[uint]model.Product),
		Weight:   0,
		Discount: 0,
		Taxes:    0,
		Total:    0,
	}
}

// Add a product to the cart
func (c *Cart) Add(product *model.Product) {
	c.Lock()
	defer c.Unlock()

	taxes := (product.Subtotal / 100) * product.Taxes
	discount := (product.Subtotal / 100) * product.Discount

	c.Products[product.ID] = *product
	c.Weight = c.Weight + product.Weight
	c.Discount = c.Discount + discount
	c.Taxes = c.Taxes + taxes
	c.Subtotal = c.Subtotal + product.Subtotal
	c.Total = c.Total + product.Subtotal + taxes - discount
}

// Checkout takes all the products and returns the final purchase
func (c *Cart) Checkout() float32 {
	c.Lock()
	defer c.Unlock()

	total := c.Total + c.Taxes - c.Discount

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

// Items prints cart items
func (c *Cart) Items() map[uint]model.Product {
	c.RLock()
	defer c.RUnlock()

	return c.Products
}

// Remove takes away a product from the cart
func (c *Cart) Remove(key uint) {
	c.Lock()
	defer c.Unlock()

	product := c.Products[key]

	tax := (product.Subtotal / 100) * product.Taxes
	discount := (product.Subtotal / 100) * product.Discount

	c.Weight = c.Weight - product.Weight
	c.Discount = c.Discount - discount
	c.Taxes = c.Taxes - tax
	c.Total = c.Total - product.Subtotal - tax + discount

	delete(c.Products, key)
}

// Reset cart products
func (c *Cart) Reset() {
	c.Lock()
	defer c.Unlock()

	c.Products = map[uint]model.Product{}
	c.Weight = 0
	c.Discount = 0
	c.Taxes = 0
	c.Subtotal = 0
	c.Total = 0
}

// Size returns the amount of products in the cart
func (c *Cart) Size() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.Products)
}

// String returns a string with the cart details
func (c *Cart) String() string {
	return fmt.Sprintf("The cart has a weight of %2.f kg, $%2.f of discounts, $%2.f of taxes and a total of $%2.f", c.Weight, c.Discount, c.Taxes, c.Total)
}
