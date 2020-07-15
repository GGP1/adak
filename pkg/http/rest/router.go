/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"net/http"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/repository"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/jinzhu/gorm"

	// h -> handler

	h "github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/updating"

	// m -> middleware

	m "github.com/GGP1/palo/pkg/http/rest/middleware"

	"github.com/gorilla/mux"
)

// NewRouter creates and returns a mux router
func NewRouter(db *gorm.DB) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Create MonoRepo
	repo := *new(repository.MonoRepo)

	// Services
	a := adding.NewService(repo)
	d := deleting.NewService(repo)
	l := listing.NewService(repo)
	u := updating.NewService(repo)
	// -- Auth session --
	session := h.NewSession(repo)
	// -- Email lists --
	pendingList := email.NewList(db, "pending_list", repo)
	validatedList := email.NewList(db, "validated_list", repo)

	// Create cart cart
	cart := shopping.NewCart()

	// ==========
	// 	Auth
	// ==========
	r.HandleFunc("/login", h.Session.Login(session, db, validatedList)).Methods("POST")
	r.HandleFunc("/logout", h.Session.Logout(session)).Methods("GET")

	// ==========
	// 	 Cart
	// ==========
	r.HandleFunc("/cart", h.CartGet(cart)).Methods("GET")
	r.HandleFunc("/cart/{id}", h.CartRemove(cart)).Methods("DELETE")
	r.HandleFunc("/cart/add/{amount}", h.CartAdd(cart)).Methods("POST")
	r.HandleFunc("/cart/checkout", h.CartCheckout(cart)).Methods("GET")
	r.HandleFunc("/cart/filter/brand/{brand}", h.CartFilterByBrand(cart)).Methods("GET")
	r.HandleFunc("/cart/filter/category/{category}", h.CartFilterByCategory(cart)).Methods("GET")
	r.HandleFunc("/cart/filter/total/{min}/{max}", h.CartFilterByTotal(cart)).Methods("GET")
	r.HandleFunc("/cart/filter/type/{type}", h.CartFilterByType(cart)).Methods("GET")
	r.HandleFunc("/cart/filter/weight/{min}/{max}", h.CartFilterByWeight(cart)).Methods("GET")
	r.HandleFunc("/cart/items", h.CartItems(cart)).Methods("GET")
	r.HandleFunc("/cart/reset", h.CartReset(cart)).Methods("GET")
	r.HandleFunc("/cart/size", h.CartSize(cart)).Methods("GET")

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
	products := h.Products{DB: db}
	r.HandleFunc("/products", products.GetAll(l)).Methods("GET")
	r.HandleFunc("/products/{id}", products.GetByID(l)).Methods("GET")
	r.HandleFunc("/products/add", products.Add(a)).Methods("POST")
	r.HandleFunc("/products/{id}", m.RequireLogin(products.Update(u))).Methods("PUT")
	r.HandleFunc("/products/{id}", m.RequireLogin(products.Delete(d))).Methods("DELETE")

	// ==========
	// 	Reviews
	// ==========
	reviews := h.Reviews{DB: db}
	r.HandleFunc("/reviews", reviews.GetAll(l)).Methods("GET")
	r.HandleFunc("/reviews/{id}", reviews.GetByID(l)).Methods("GET")
	r.HandleFunc("/reviews/add", m.RequireLogin(reviews.Add(a))).Methods("POST")
	r.HandleFunc("/reviews/{id}", m.RequireLogin(reviews.Delete(d))).Methods("DELETE")

	// ==========
	// 	Shops
	// ==========
	shops := h.Shops{DB: db}
	r.HandleFunc("/shops", shops.GetAll(l)).Methods("GET")
	r.HandleFunc("/shops/{id}", shops.GetByID(l)).Methods("GET")
	r.HandleFunc("/shops/add", shops.Add(a)).Methods("POST")
	r.HandleFunc("/shops/{id}", m.RequireLogin(shops.Update(u))).Methods("PUT")
	r.HandleFunc("/shops/{id}", m.RequireLogin(shops.Delete(d))).Methods("DELETE")

	// ==========
	// 	Users
	// ==========
	users := h.Users{DB: db}
	r.HandleFunc("/users", users.GetAll(l)).Methods("GET")
	r.HandleFunc("/users/{id}", users.GetByID(l)).Methods("GET")
	r.HandleFunc("/users/add", users.Add(a, pendingList)).Methods("POST")
	r.HandleFunc("/users/{id}", m.RequireLogin(users.Update(u))).Methods("PUT")
	r.HandleFunc("/users/{id}", m.RequireLogin(users.Delete(d, pendingList, validatedList))).Methods("DELETE")

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	http.Handle("/", r)

	return r
}
