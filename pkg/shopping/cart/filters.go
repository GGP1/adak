package cart

import (
	"context"

	"github.com/pkg/errors"
)

var (
	errNotFound = errors.New("no products found")
)

// FilterByBrand looks for products with the specified brand.
func (s *Shopping) FilterByBrand(ctx context.Context, req *FilterTextRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND brand=$2`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Field); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterByCategory looks for products with the specified category.
func (s *Shopping) FilterByCategory(ctx context.Context, req *FilterTextRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND category=$2`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Field); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterByDiscount looks for products within the percentage of discount range specified.
func (s *Shopping) FilterByDiscount(ctx context.Context, req *FilterNumberRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND discount >= $2 AND discount <= $3`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Min, req.Max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterBySubtotal looks for products within the subtotal price range specified.
func (s *Shopping) FilterBySubtotal(ctx context.Context, req *FilterNumberRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND subtotal >= $2 AND subtotal <= $3`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Min, req.Max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterByTaxes looks for products within the percentage of taxes range specified.
func (s *Shopping) FilterByTaxes(ctx context.Context, req *FilterNumberRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND taxes >= $2 AND taxes <= $3`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Min, req.Max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterByTotal looks for products within the total price range specified.
func (s *Shopping) FilterByTotal(ctx context.Context, req *FilterNumberRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND total >= $2 AND total <= $3`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Min, req.Max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterByType looks for products with the specified type.
func (s *Shopping) FilterByType(ctx context.Context, req *FilterTextRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND type=$2`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Field); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}

// FilterByWeight looks for products within the weight range specified.
func (s *Shopping) FilterByWeight(ctx context.Context, req *FilterNumberRequest) (*FilterResponse, error) {
	var products []*Product

	query := `SELECT * FROM cart_products WHERE cart_id=$1 AND AND weight >= $2 AND weight <= $3`

	if err := s.db.SelectContext(ctx, &products, query, req.CartID, req.Min, req.Max); err != nil {
		return nil, errors.Wrap(err, errNotFound.Error())
	}

	if len(products) == 0 {
		return nil, errNotFound
	}

	return &FilterResponse{Products: products}, nil
}
