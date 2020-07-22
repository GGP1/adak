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

	Products map[uint]*model.Product
	Counter  int
	Weight   float32
	Discount float32
	Taxes    float32
	Subtotal float32
	Total    float32
}

// NewCart returns a cart with the default values
func NewCart() *Cart {
	return &Cart{
		Products: make(map[uint]*model.Product),
		Counter:  0,
		Weight:   0,
		Discount: 0,
		Taxes:    0,
		Total:    0,
	}
}

// Add a product to the cart
func (c *Cart) Add(product *model.Product, quantity int) error {
	c.Lock()
	defer c.Unlock()

	if c.Products[product.ID] != nil {
		existingProduct := c.Products[product.ID]
		if product.Brand != existingProduct.Brand || product.Type != existingProduct.Type {
			return errors.New("the product id is already in use")
		}
		existingProduct.Quantity += quantity
	} else {
		c.Products[product.ID] = product
	}

	taxes := ((product.Subtotal / 100) * product.Taxes)
	discount := ((product.Subtotal / 100) * product.Discount)

	for i := 0; i < quantity; i++ {
		c.Counter++
		product.Quantity++
		c.Weight += product.Weight
		c.Discount += discount
		c.Taxes += taxes
		c.Subtotal += product.Subtotal
		c.Total = c.Total + product.Subtotal + taxes - discount
	}

	return nil
}

// Checkout takes all the products and returns the total price
func (c *Cart) Checkout() float32 {
	c.RLock()
	defer c.RUnlock()

	total := c.Total + c.Taxes - c.Discount

	return total
}

// FilterByBrand looks for products with the specified brand
func (c *Cart) FilterByBrand(brand string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if strings.ToLower(brand) == strings.ToLower(v.Brand) {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByCategory looks for products with the specified category
func (c *Cart) FilterByCategory(category string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if strings.ToLower(category) == strings.ToLower(v.Category) {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByDiscount looks for products within the percentage of discount range specified
func (c *Cart) FilterByDiscount(min, max float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if v.Subtotal >= min && v.Discount <= max {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterBySubtotal looks for products within the subtotal price range specified
func (c *Cart) FilterBySubtotal(min, max float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if v.Subtotal >= min && v.Subtotal <= max {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByTaxes looks for products within the percentage of taxes range specified
func (c *Cart) FilterByTaxes(min, max float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if v.Taxes >= min && v.Taxes <= max {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByTotal looks for products within the total price range specified
func (c *Cart) FilterByTotal(min, max float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if v.Total >= min && v.Total <= max {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByType looks for products with the specified type
func (c *Cart) FilterByType(pType string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if strings.ToLower(pType) == strings.ToLower(v.Type) {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByWeight looks for products within the weight range specified
func (c *Cart) FilterByWeight(min, max float32) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for k, v := range c.Products {
		wg.Add(1)
		go func(k uint, v *model.Product) {
			defer wg.Done()
			if v.Weight >= min && v.Weight <= max {
				products = append(products, *c.Products[k])
			}
		}(k, v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// Items prints cart items
func (c *Cart) Items() map[uint]*model.Product {
	c.RLock()
	defer c.RUnlock()

	return c.Products
}

// Remove takes away the specified quantity of products from the cart
func (c *Cart) Remove(key uint, quantity int) error {
	c.Lock()
	defer c.Unlock()

	if c.Products[key] == nil {
		return errors.New("product not found")
	}

	if c.Counter == 1 {
		c.Reset()
		return nil
	}

	product := c.Products[key]
	taxes := (product.Subtotal / 100) * product.Taxes
	discount := (product.Subtotal / 100) * product.Discount

	if quantity > product.Quantity {
		return errors.New("quantity ")
	}

	if quantity == product.Quantity {
		delete(c.Products, key)
	}

	for i := 0; i < quantity; i++ {
		c.Counter--
		product.Quantity--
		c.Weight -= product.Weight
		c.Discount -= discount
		c.Taxes -= taxes
		c.Subtotal -= product.Subtotal
		c.Total = c.Total - product.Subtotal - taxes + discount
	}

	return nil
}

// Reset cart products
func (c *Cart) Reset() {
	c.Lock()
	defer c.Unlock()

	c.Products = map[uint]*model.Product{}
	c.Counter = 0
	c.Weight = 0
	c.Discount = 0
	c.Taxes = 0
	c.Subtotal = 0
	c.Total = 0
}

// Size returns the quantity of products in the cart
func (c *Cart) Size() int {
	c.RLock()
	defer c.RUnlock()

	return c.Counter
}

// String returns a string with the cart details
func (c *Cart) String() string {
	return fmt.Sprintf("The cart has %d products, a weight of %2.fkg, $%2.f of discounts, $%2.f of taxes and a total of $%2.f", c.Counter, c.Weight, c.Discount, c.Taxes, c.Total)
}
