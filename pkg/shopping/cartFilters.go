package shopping

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	errProductNotFound = fmt.Errorf("no products found")
)

// FilterByBrand looks for products with the specified brand.
func FilterByBrand(db *sqlx.DB, cartID, brand string) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND brand=$2`

	if err := db.Select(&products, query, cartID, brand); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByCategory looks for products with the specified category.
func FilterByCategory(db *sqlx.DB, cartID, category string) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND category=$2`

	if err := db.Select(&products, query, cartID, category); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByDiscount looks for products within the percentage of discount range specified.
func FilterByDiscount(db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND discount >= $2 AND discount <= $3`

	if err := db.Select(&products, query, cartID, min, max); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterBySubtotal looks for products within the subtotal price range specified.
func FilterBySubtotal(db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND subtotal >= $2 AND subtotal <= $3`

	if err := db.Select(&products, query, cartID, min, max); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByTaxes looks for products within the percentage of taxes range specified.
func FilterByTaxes(db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND taxes >= $2 AND taxes <= $3`

	if err := db.Select(&products, query, cartID, min, max); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByTotal looks for products within the total price range specified.
func FilterByTotal(db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND total >= $2 AND total <= $3`

	if err := db.Select(&products, query, cartID, min, max); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByType looks for products with the specified type.
func FilterByType(db *sqlx.DB, cartID, pType string) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND type=$2`

	if err := db.Select(&products, query, cartID, pType); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}

// FilterByWeight looks for products within the weight range specified.
func FilterByWeight(db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND weight >= $2 AND weight <= $3`

	if err := db.Select(&products, query, cartID, min, max); err != nil {
		return nil, fmt.Errorf("couldn't find the products: %v", err)
	}

	if len(products) == 0 {
		return nil, errProductNotFound
	}

	return products, nil
}
