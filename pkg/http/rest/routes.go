/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"io"
	"net/http"

	h "github.com/GGP1/palo/pkg/http/rest/handler"

	"github.com/gorilla/mux"
)

// SetupRouter returns gin Engine
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Auth
	r.HandleFunc("/login", h.Login()).Methods("POST")
	r.HandleFunc("/logout", h.Logout()).Methods("GET")

	// Products
	r.HandleFunc("/products", h.GetProducts()).Methods("GET")
	r.HandleFunc("/products/{id}", h.GetOneProduct()).Methods("GET")
	r.HandleFunc("/products/{id}", h.AddProduct()).Methods("POST")
	r.HandleFunc("/products/{id}", h.UpdateProduct()).Methods("PUT")
	r.HandleFunc("/products/{id}", h.DeleteProduct()).Methods("DELETE")

	// Reviews
	r.HandleFunc("/reviews", h.GetReviews()).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.GetOneReview()).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.AddReview()).Methods("POST")
	r.HandleFunc("/reviews/{id}", h.DeleteReview()).Methods("DELETE")

	// Shops
	r.HandleFunc("/shops", h.GetShops()).Methods("GET")
	r.HandleFunc("/shops/{id}", h.GetOneShop()).Methods("GET")
	r.HandleFunc("/shops/{id}", h.AddShop()).Methods("POST")
	r.HandleFunc("/shops/{id}", h.UpdateShop()).Methods("PUT")
	r.HandleFunc("/shops/{id}", h.DeleteShop()).Methods("DELETE")

	// Users
	r.HandleFunc("/users", h.GetUsers()).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetOneUser()).Methods("GET")
	r.HandleFunc("/users/{id}", h.AddUser()).Methods("POST")
	r.HandleFunc("/users/{id}", h.UpdateUser()).Methods("PUT")
	r.HandleFunc("/users/{id}", h.DeleteUser()).Methods("DELETE")

	// Home
	r.HandleFunc("/", home()).Methods("GET")

	http.Handle("/", r)

	return r
}

func home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Welcome to the Palo home page")
	}
}
