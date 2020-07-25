package shopping

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	errProductNotFound = errors.New("no products found")
)

// FilterByBrand looks for products with the specified brand
func FilterByBrand(db *gorm.DB, cartID, brand string) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND brand = ?", cartID, brand).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByCategory looks for products with the specified category
func FilterByCategory(db *gorm.DB, cartID, category string) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND category = ?", cartID, category).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByDiscount looks for products within the percentage of discount range specified
func FilterByDiscount(db *gorm.DB, cartID string, min, max float32) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND discount >= ? AND discount <= ?", cartID, min, max).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterBySubtotal looks for products within the subtotal price range specified
func FilterBySubtotal(db *gorm.DB, cartID string, min, max float32) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND subtotal >= ? AND subtotal <= ?", cartID, min, max).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByTaxes looks for products within the percentage of taxes range specified
func FilterByTaxes(db *gorm.DB, cartID string, min, max float32) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND taxes >= ? AND taxes <= ?", cartID, min, max).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByTotal looks for products within the total price range specified
func FilterByTotal(db *gorm.DB, cartID string, min, max float32) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND total >= ? AND total <= ?", cartID, min, max).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByType looks for products with the specified type
func FilterByType(db *gorm.DB, cartID, pType string) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND type = ?", cartID, pType).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByWeight looks for products within the weight range specified
func FilterByWeight(db *gorm.DB, cartID string, min, max float32) ([]CartProduct, error) {
	var products []CartProduct

	if err := db.Where("cart_id = ? AND weight >= ? AND weight <= ?", cartID, min, max).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "couldn't find the products")
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}
