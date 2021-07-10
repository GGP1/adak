package postgres

import (
	"bytes"

	"github.com/GGP1/adak/internal/params"
)

const (
	// Orders table
	Orders table = "orders"
	// Products table
	Products table = "products"
	// Reviews table
	Reviews table = "reviews"
	// Shops table
	Shops table = "shops"
	// Users table
	Users table = "users"
)

type table string

// AddPagination adds pagination to a query, returns it and the arguments to be used.
func AddPagination(query string, params params.Query) (string, []interface{}) {
	buf := bytes.NewBufferString(query)

	args := []interface{}{params.Limit}
	if params.Cursor.Used {
		buf.WriteString(" WHERE created_at < $2 OR (created_at = $2 AND id < $3)")
		args = append(args, params.Cursor.CreatedAt, params.Cursor.ID) // Respect query args order
	}
	buf.WriteString(" ORDER BY created_at DESC, id DESC LIMIT $1")

	return buf.String(), args
}
