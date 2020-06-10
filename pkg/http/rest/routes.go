/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"io"
	"net/http"

	h "github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/http/rest/middleware"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// NewRouter returns mux router
func NewRouter(db *gorm.DB) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Handlers
	auth := h.AuthHandler{DB: db}
	products := h.ProductHandler{DB: db}
	users := h.UserHandler{DB: db}
	shops := h.ShopHandler{DB: db}
	reviews := h.ReviewHandler{DB: db}

	// Auth
	r.HandleFunc("/login", auth.Login()).Methods("POST")
	r.HandleFunc("/logout", auth.Logout()).Methods("GET")

	// Products
	r.HandleFunc("/products", products.Get()).Methods("GET")
	r.HandleFunc("/products/{id}", products.GetOne()).Methods("GET")
	r.HandleFunc("/products/add", products.Add()).Methods("POST")
	r.HandleFunc("/products/{id}", products.Update()).Methods("PUT")
	r.HandleFunc("/products/{id}", products.Delete()).Methods("DELETE")

	// Reviews
	r.HandleFunc("/reviews", reviews.Get()).Methods("GET")
	r.HandleFunc("/reviews/{id}", reviews.GetOne()).Methods("GET")
	r.HandleFunc("/reviews/add", reviews.Add()).Methods("POST")
	r.HandleFunc("/reviews/{id}", reviews.Delete()).Methods("DELETE")

	// Shops
	r.HandleFunc("/shops", shops.Get()).Methods("GET")
	r.HandleFunc("/shops/{id}", shops.GetOne()).Methods("GET")
	r.HandleFunc("/shops/add", shops.Add()).Methods("POST")
	r.HandleFunc("/shops/{id}", shops.Update()).Methods("PUT")
	r.HandleFunc("/shops/{id}", shops.Delete()).Methods("DELETE")

	// Users
	r.HandleFunc("/users", users.Get()).Methods("GET")
	r.HandleFunc("/users/{id}", users.GetOne()).Methods("GET")
	r.HandleFunc("/users/add", users.Add()).Methods("POST")
	r.HandleFunc("/users/{id}", users.Update()).Methods("PUT")
	r.HandleFunc("/users/{id}", users.Delete()).Methods("DELETE")

	// Home
	r.HandleFunc("/", home()).Methods("GET")

	// Email validation
	r.HandleFunc("/verify", verify()).Methods("GET")

	// Middlewares
	r.Use(middleware.AllowCrossOrigin)

	http.Handle("/", r)

	return r
}

// Home page
func home() http.HandlerFunc {
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
