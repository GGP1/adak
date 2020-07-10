/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"net/http"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"

	// h -> handler
	h "github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/updating"

	// m -> middleware
	m "github.com/GGP1/palo/pkg/http/rest/middleware"
	"github.com/GGP1/palo/pkg/shopping"

	"github.com/gorilla/mux"
)

// NewRouter creates and returns a mux router
func NewRouter(a adding.Service, d deleting.Service, l listing.Service, u updating.Service) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Create users mail lists
	pendingList := email.NewPendingList()
	validatedList := email.NewValidatedList()

	// Create auth session
	repo := *new(h.AuthRepository)
	session := h.NewSession(repo)

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
	r.HandleFunc("/products", h.GetProducts(l)).Methods("GET")
	r.HandleFunc("/products/{id}", h.GetProductByID(l)).Methods("GET")
	r.HandleFunc("/products/add", h.AddProduct(a)).Methods("POST")
	r.HandleFunc("/products/{id}", m.RequireLogin(h.UpdateProduct(u))).Methods("PUT")
	r.HandleFunc("/products/{id}", m.RequireLogin(h.DeleteProduct(d))).Methods("DELETE")

	// ==========
	// 	Reviews
	// ==========
	r.HandleFunc("/reviews", h.GetReviews(l)).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.GetReviewByID(l)).Methods("GET")
	r.HandleFunc("/reviews/add", m.RequireLogin(h.AddReview(a))).Methods("POST")
	r.HandleFunc("/reviews/{id}", m.RequireLogin(h.DeleteReview(d))).Methods("DELETE")

	// ==========
	// 	Shopping
	// ==========
	r.HandleFunc("/shopping", m.RequireLogin(h.CartGet(cart))).Methods("GET")
	r.HandleFunc("/shopping/{id}", m.RequireLogin(h.CartRemove(cart))).Methods("DELETE")
	r.HandleFunc("/shopping/add/{amount}", m.RequireLogin(h.CartAdd(cart))).Methods("POST")
	r.HandleFunc("/shopping/checkout", m.RequireLogin(h.CartCheckout(cart))).Methods("GET")
	r.HandleFunc("/shopping/filter/brand/{brand}", m.RequireLogin(h.CartFilterByBrand(cart))).Methods("GET")
	r.HandleFunc("/shopping/filter/category/{category}", m.RequireLogin(h.CartFilterByCategory(cart))).Methods("GET")
	r.HandleFunc("/shopping/filter/total/{min}/{max}", m.RequireLogin(h.CartFilterByTotal(cart))).Methods("GET")
	r.HandleFunc("/shopping/filter/type/{type}", m.RequireLogin(h.CartFilterByType(cart))).Methods("GET")
	r.HandleFunc("/shopping/filter/weight/{min}/{max}", m.RequireLogin(h.CartFilterByWeight(cart))).Methods("GET")
	r.HandleFunc("/shopping/items", m.RequireLogin(h.CartItems(cart))).Methods("GET")
	r.HandleFunc("/shopping/reset", m.RequireLogin(h.CartReset(cart))).Methods("GET")
	r.HandleFunc("/shopping/size", m.RequireLogin(h.CartSize(cart))).Methods("GET")

	// ==========
	// 	Shops
	// ==========
	r.HandleFunc("/shops", h.GetShops(l)).Methods("GET")
	r.HandleFunc("/shops/{id}", h.GetShopByID(l)).Methods("GET")
	r.HandleFunc("/shops/add", h.AddShop(a)).Methods("POST")
	r.HandleFunc("/shops/{id}", m.RequireLogin(h.UpdateShop(u))).Methods("PUT")
	r.HandleFunc("/shops/{id}", m.RequireLogin(h.DeleteShop(d))).Methods("DELETE")

	// ==========
	// 	Users
	// ==========
	r.HandleFunc("/users", h.GetUsers(l)).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetUserByID(l)).Methods("GET")
	r.HandleFunc("/users/add", h.AddUser(a, pendingList)).Methods("POST")
	r.HandleFunc("/users/{id}", m.RequireLogin(h.UpdateUser(u))).Methods("PUT")
	r.HandleFunc("/users/{id}", m.RequireLogin(h.DeleteUser(d, pendingList, validatedList))).Methods("DELETE")

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	http.Handle("/", r)

	return r
}
