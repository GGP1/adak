/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"net/http"

	"github.com/GGP1/palo/internal/email"
	// h -> handler
	h "github.com/GGP1/palo/pkg/http/rest/handler"
	// m -> middleware
	m "github.com/GGP1/palo/pkg/http/rest/middleware"
	"github.com/GGP1/palo/pkg/shopping"

	"github.com/gorilla/mux"
)

// NewRouter creates and returns a mux router
func NewRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Create users mail lists
	pendingList := email.NewPendingList()
	validatedList := email.NewValidatedList()

	// Create auth session
	repo := new(h.AuthRepository)
	session := h.NewSession(*repo)

	// Create shopping cart
	cart := shopping.NewCart()

	// ==========
	// 	Auth
	// ==========
	r.HandleFunc("/login", h.Session.Login(session, validatedList)).Methods("POST")
	r.HandleFunc("/logout", h.Session.Logout(session)).Methods("GET")

	// ==========
	// 	Email
	// ==========
	r.HandleFunc("/email/{token}", h.ValidateEmail(pendingList, validatedList)).Methods("GET")

	// ==========
	// 	Home
	// ==========
	r.HandleFunc("/", h.Home()).Methods("GET")

	// ==========
	// 	Products
	// ==========
	r.HandleFunc("/products", h.GetProducts()).Methods("GET")
	r.HandleFunc("/products/{id}", h.GetProductByID()).Methods("GET")
	r.HandleFunc("/products/add", h.AddProduct()).Methods("POST")
	r.HandleFunc("/products/{id}", m.RequireLogin(h.UpdateProduct())).Methods("PUT")
	r.HandleFunc("/products/{id}", m.RequireLogin(h.DeleteProduct())).Methods("DELETE")

	// ==========
	// 	Reviews
	// ==========
	r.HandleFunc("/reviews", h.GetReviews()).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.GetReviewByID()).Methods("GET")
	r.HandleFunc("/reviews/add", m.RequireLogin(h.AddReview())).Methods("POST")
	r.HandleFunc("/reviews/{id}", m.RequireLogin(h.DeleteReview())).Methods("DELETE")

	// ==========
	// 	Shopping
	// ==========
	r.HandleFunc("/shopping", h.CartGet(cart)).Methods("GET")
	r.HandleFunc("/shopping/items", h.CartItems(cart)).Methods("GET")
	r.HandleFunc("/shopping/add", m.RequireLogin(h.CartAdd(cart))).Methods("POST")
	r.HandleFunc("/shopping/checkout", m.RequireLogin(h.CartCheckout(cart))).Methods("GET")
	r.HandleFunc("/shopping/filter/brand/{brand}", h.CartFilterByBrand(cart)).Methods("GET")
	r.HandleFunc("/shopping/filter/category/{category}", h.CartFilterByCategory(cart)).Methods("GET")
	r.HandleFunc("/shopping/filter/total/{min}/{max}", h.CartFilterByTotal(cart)).Methods("GET")
	r.HandleFunc("/shopping/filter/type/{type}", h.CartFilterByType(cart)).Methods("GET")
	r.HandleFunc("/shopping/filter/weight/{min}/{max}", h.CartFilterByWeight(cart)).Methods("GET")
	r.HandleFunc("/shopping/remove/{id}", m.RequireLogin(h.CartRemove(cart))).Methods("DELETE")
	r.HandleFunc("/shopping/reset", m.RequireLogin(h.CartReset(cart))).Methods("GET")
	r.HandleFunc("/shopping/size", m.RequireLogin(h.CartSize(cart))).Methods("GET")

	// ==========
	// 	Shops
	// ==========
	r.HandleFunc("/shops", h.GetShops()).Methods("GET")
	r.HandleFunc("/shops/{id}", h.GetShopByID()).Methods("GET")
	r.HandleFunc("/shops/add", h.AddShop()).Methods("POST")
	r.HandleFunc("/shops/{id}", m.RequireLogin(h.UpdateShop())).Methods("PUT")
	r.HandleFunc("/shops/{id}", m.RequireLogin(h.DeleteShop())).Methods("DELETE")

	// ==========
	// 	Users
	// ==========
	r.HandleFunc("/users", h.GetUsers()).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetUserByID()).Methods("GET")
	r.HandleFunc("/users/add", h.AddUser(pendingList)).Methods("POST")
	r.HandleFunc("/users/{id}", m.RequireLogin(h.UpdateUser())).Methods("PUT")
	r.HandleFunc("/users/{id}", m.RequireLogin(h.DeleteUser(pendingList, validatedList))).Methods("DELETE")

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	http.Handle("/", r)

	return r
}
