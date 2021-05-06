// Package validate makes sure to validate structs and fields to avoid making invalid queries to the database.
package validate

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info.
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Struct returns an error if any of the struct's fields is invalid.
func Struct(ctx context.Context, v interface{}) error {
	if err := validate.StructCtx(ctx, v); err != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			return vErr
		}
		return fmt.Errorf("invalid request body: %v", err)
	}

	return nil
}
