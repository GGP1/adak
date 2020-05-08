/*
Package listing: lists products/users taking them from the database
*/
package listing

// Product model
type Product struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Brand    string `json:"brand"`
	Weight   int    `json:"weight"`
	Cost     int    `json:"cost"`
	Review   Review `json:"review"`
}
