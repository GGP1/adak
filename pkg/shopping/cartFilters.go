package shopping

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	errNotFound = errors.New("no products found")
)

// FilterByBrand looks for products with the specified brand.
func FilterByBrand(ctx context.Context, db *sqlx.DB, cartID, brand string) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND brand=$2`

	if err := db.SelectContext(ctx, &products, query, cartID, brand); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByCategory looks for products with the specified category.
func FilterByCategory(ctx context.Context, db *sqlx.DB, cartID, category string) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND category=$2`

	if err := db.SelectContext(ctx, &products, query, cartID, category); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByDiscount looks for products within the percentage of discount range specified.
func FilterByDiscount(ctx context.Context, db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND discount >= $2 AND discount <= $3`

	if err := db.SelectContext(ctx, &products, query, cartID, min, max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterBySubtotal looks for products within the subtotal price range specified.
func FilterBySubtotal(ctx context.Context, db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND subtotal >= $2 AND subtotal <= $3`

	if err := db.SelectContext(ctx, &products, query, cartID, min, max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByTaxes looks for products within the percentage of taxes range specified.
func FilterByTaxes(ctx context.Context, db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND taxes >= $2 AND taxes <= $3`

	if err := db.SelectContext(ctx, &products, query, cartID, min, max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByTotal looks for products within the total price range specified.
func FilterByTotal(ctx context.Context, db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND total >= $2 AND total <= $3`

	if err := db.SelectContext(ctx, &products, query, cartID, min, max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByType looks for products with the specified type.
func FilterByType(ctx context.Context, db *sqlx.DB, cartID, pType string) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND type=$2`

	if err := db.SelectContext(ctx, &products, query, cartID, pType); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}

// FilterByWeight looks for products within the weight range specified.
func FilterByWeight(ctx context.Context, db *sqlx.DB, cartID string, min, max float64) ([]CartProduct, error) {
	var products []CartProduct

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND weight >= $2 AND weight <= $3`

	if err := db.SelectContext(ctx, &products, query, cartID, min, max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return products, nil
}
