/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"io"
	"net/http"

	"github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/http/rest/middleware"

	"github.com/gorilla/mux"
)

// NewRouter returns mux router
func NewRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Auth
	r.HandleFunc("/login", handler.Login()).Methods("POST")
	r.HandleFunc("/logout", handler.Logout()).Methods("GET")

	// Products
	r.HandleFunc("/products", handler.GetProducts()).Methods("GET")
	r.HandleFunc("/products/{id}", handler.GetOneProduct()).Methods("GET")
	r.HandleFunc("/products/add", handler.AddProduct()).Methods("POST")
	r.HandleFunc("/products/{id}", handler.UpdateProduct()).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct()).Methods("DELETE")

	// Reviews
	r.HandleFunc("/reviews", handler.GetReviews()).Methods("GET")
	r.HandleFunc("/reviews/{id}", handler.GetOneReview()).Methods("GET")
	r.HandleFunc("/reviews/add", handler.AddReview()).Methods("POST")
	r.HandleFunc("/reviews/{id}", handler.DeleteReview()).Methods("DELETE")

	// Shops
	r.HandleFunc("/shops", handler.GetShops()).Methods("GET")
	r.HandleFunc("/shops/{id}", handler.GetOneShop()).Methods("GET")
	r.HandleFunc("/shops/add", handler.AddShop()).Methods("POST")
	r.HandleFunc("/shops/{id}", handler.UpdateShop()).Methods("PUT")
	r.HandleFunc("/shops/{id}", handler.DeleteShop()).Methods("DELETE")

	// Users
	r.HandleFunc("/users", handler.GetUsers()).Methods("GET")
	r.HandleFunc("/users/{id}", handler.GetOneUser()).Methods("GET")
	r.HandleFunc("/users/add", handler.AddUser()).Methods("POST")
	r.HandleFunc("/users/{id}", handler.UpdateUser()).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser()).Methods("DELETE")

	// Home
	r.HandleFunc("/", Home()).Methods("GET")

	// Email verification
	r.HandleFunc("/verify", verify()).Methods("GET")

	// Middlewares
	r.Use(middleware.AllowCrossOrigin)
	r.Use(middleware.LimitRate)

	http.Handle("/", r)

	return r
}

// Home page
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Welcome to the Palo home page")
	}
}

// Email verification page
func verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "You have successfully confirmed your email!")
	}
}
