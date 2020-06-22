/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"io"
	"net/http"

	h "github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/http/rest/middleware"

	"github.com/gorilla/mux"
)

// NewRouter returns mux router
func NewRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Create auth session
	repo := new(h.AuthRepository)
	session := h.NewSession(*repo)

	// Auth
	r.HandleFunc("/login", h.Session.Login(session)).Methods("POST")
	r.HandleFunc("/logout", h.Session.Logout(session)).Methods("GET")

	// Products
	r.HandleFunc("/products", h.GetProducts()).Methods("GET")
	r.HandleFunc("/products/{id}", h.GetOneProduct()).Methods("GET")
	r.HandleFunc("/products/add", h.AddProduct()).Methods("POST")
	r.HandleFunc("/products/{id}", middleware.IsLoggedIn(h.UpdateProduct())).Methods("PUT")
	r.HandleFunc("/products/{id}", middleware.IsLoggedIn(h.DeleteProduct())).Methods("DELETE")

	// Reviews
	r.HandleFunc("/reviews", h.GetReviews()).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.GetOneReview()).Methods("GET")
	r.HandleFunc("/reviews/add", middleware.IsLoggedIn(h.AddReview())).Methods("POST")
	r.HandleFunc("/reviews/{id}", middleware.IsLoggedIn(h.DeleteReview())).Methods("DELETE")

	// Shops
	r.HandleFunc("/shops", h.GetShops()).Methods("GET")
	r.HandleFunc("/shops/{id}", h.GetOneShop()).Methods("GET")
	r.HandleFunc("/shops/add", h.AddShop()).Methods("POST")
	r.HandleFunc("/shops/{id}", middleware.IsLoggedIn(h.UpdateShop())).Methods("PUT")
	r.HandleFunc("/shops/{id}", middleware.IsLoggedIn(h.DeleteShop())).Methods("DELETE")

	// Users
	r.HandleFunc("/users", h.GetUsers()).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetOneUser()).Methods("GET")
	r.HandleFunc("/users/add", h.AddUser()).Methods("POST")
	r.HandleFunc("/users/{id}", middleware.IsLoggedIn(h.UpdateUser())).Methods("PUT")
	r.HandleFunc("/users/{id}", middleware.IsLoggedIn(h.DeleteUser())).Methods("DELETE")

	// Home
	r.HandleFunc("/", Home()).Methods("GET")

	// Email verification
	r.HandleFunc("/verify", verify()).Methods("GET")

	// Middlewares
	r.Use(middleware.AllowCrossOrigin)
	r.Use(middleware.LimitRate)
	r.Use(middleware.LogFormatter)

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
