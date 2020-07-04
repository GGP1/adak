/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"net/http"

	"github.com/GGP1/palo/internal/response"
	h "github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/http/rest/middleware"
	"github.com/GGP1/palo/pkg/shopping"

	"github.com/gorilla/mux"
)

// NewRouter creates and returns a mux router
func NewRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Create auth session
	repo := new(h.AuthRepository)
	session := h.NewSession(*repo)

	// Create shopping cart
	cart := shopping.NewCart()

	// Auth
	r.HandleFunc("/login", h.Session.Login(session)).Methods("POST")
	r.HandleFunc("/logout", h.Session.Logout(session)).Methods("GET")

	// Products
	r.HandleFunc("/products", h.GetProducts()).Methods("GET")
	r.HandleFunc("/products/{id}", h.GetProductByID()).Methods("GET")
	r.HandleFunc("/products/add", h.AddProduct()).Methods("POST")
	r.HandleFunc("/products/{id}", middleware.IsLoggedIn(h.UpdateProduct())).Methods("PUT")
	r.HandleFunc("/products/{id}", middleware.IsLoggedIn(h.DeleteProduct())).Methods("DELETE")

	// Reviews
	r.HandleFunc("/reviews", h.GetReviews()).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.GetReviewByID()).Methods("GET")
	r.HandleFunc("/reviews/add", middleware.IsLoggedIn(h.AddReview())).Methods("POST")
	r.HandleFunc("/reviews/{id}", middleware.IsLoggedIn(h.DeleteReview())).Methods("DELETE")

	// Shopping
	r.HandleFunc("/shopping", h.CartShowItems(cart)).Methods("GET")
	r.HandleFunc("/shopping/add", h.CartAdd(cart)).Methods("POST")
	r.HandleFunc("/shopping/checkout", h.CartCheckout(cart)).Methods("GET")
	r.HandleFunc("/shopping/{id}", h.CartRemove(cart)).Methods("DELETE")
	r.HandleFunc("/shopping/reset", h.CartReset(cart)).Methods("GET")
	r.HandleFunc("/shopping/size", h.CartSize(cart)).Methods("GET")

	// Shops
	r.HandleFunc("/shops", h.GetShops()).Methods("GET")
	r.HandleFunc("/shops/{id}", h.GetShopByID()).Methods("GET")
	r.HandleFunc("/shops/add", h.AddShop()).Methods("POST")
	r.HandleFunc("/shops/{id}", middleware.IsLoggedIn(h.UpdateShop())).Methods("PUT")
	r.HandleFunc("/shops/{id}", middleware.IsLoggedIn(h.DeleteShop())).Methods("DELETE")

	// Users
	r.HandleFunc("/users", h.GetUsers()).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetUserByID()).Methods("GET")
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
		response.HTMLText(w, r, http.StatusOK, "Welcome to the Palo home page")
	}
}

// Email verification page
func verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.HTMLText(w, r, http.StatusOK, "You have successfully confirmed your email!")
	}
}
