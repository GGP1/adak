package shopping

import (
	"strings"
	"sync"

	"github.com/GGP1/palo/pkg/model"
)

// FilterByBrand looks for products with the specified brand
func (c *Cart) FilterByBrand(brand string) ([]model.Product, error) {
	c.RLock()
	defer c.RUnlock()

	var wg sync.WaitGroup
	var products []model.Product

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if strings.ToLower(brand) == strings.ToLower(v.Brand) {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if strings.ToLower(category) == strings.ToLower(v.Category) {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if v.Subtotal >= min && v.Discount <= max {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if v.Subtotal >= min && v.Subtotal <= max {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if v.Taxes >= min && v.Taxes <= max {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if v.Total >= min && v.Total <= max {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if strings.ToLower(pType) == strings.ToLower(v.Type) {
				products = append(products, *v)
			}
		}(v)
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

	for _, v := range c.Products {
		wg.Add(1)
		go func(v *model.Product) {
			defer wg.Done()
			if v.Weight >= min && v.Weight <= max {
				products = append(products, *v)
			}
		}(v)
	}
	wg.Wait()

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}
