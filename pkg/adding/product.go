/*
Package adding: adds products/users to the database
*/
package adding

// Product model
type Product struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Brand    string `json:"brand"`
	Weight   int    `json:"weight"`
	Cost     int    `json:"cost"`
	Review   Review `json:"review"`
}
