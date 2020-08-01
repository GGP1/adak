/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"context"
	"net/http"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/repository"
	"github.com/GGP1/palo/pkg/searching"
	"github.com/GGP1/palo/pkg/tracking"

	// h -> handler
	h "github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/updating"

	// m -> middleware
	m "github.com/GGP1/palo/pkg/http/rest/middleware"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	elastic "github.com/olivere/elastic/v7"
)

// NewRouter initializes services, creates and returns a mux router
func NewRouter(ctx context.Context, db *gorm.DB, esClient *elastic.Client) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Create repository
	repo := *new(repository.Repo)

	// Services
	a := adding.NewService(repo)
	d := deleting.NewService(repo)
	l := listing.NewService(repo)
	u := updating.NewService(repo)
	s := searching.NewService(repo)
	// -- Auth session --
	session := auth.NewSession(db, repo)
	// -- Email lists --
	pendingList := email.NewList(db, "pending_list", repo)
	validatedList := email.NewList(db, "validated_list", repo)
	// -- Tracker --
	tracker := *tracking.NewTracker(esClient, "salty")

	// Auth
	r.HandleFunc("/login", h.Login(session, validatedList)).Methods("POST")
	r.HandleFunc("/logout", m.RequireLogin(h.Logout(session))).Methods("GET")
	r.HandleFunc("/verification/{token}", h.ValidateEmail(pendingList, validatedList)).Methods("GET")

	// Home
	r.HandleFunc("/", h.Home(ctx, tracker)).Methods("GET")

	// Products
	products := h.Products{DB: db}
	r.HandleFunc("/products", products.Find(l)).Methods("GET")
	r.HandleFunc("/products/{id}", products.FindByID(l)).Methods("GET")
	r.HandleFunc("/products/add", products.Add(a)).Methods("POST")
	r.HandleFunc("/products/{id}", m.AdminsOnly(products.Update(u))).Methods("PUT")
	r.HandleFunc("/products/{id}", m.AdminsOnly(products.Delete(d))).Methods("DELETE")
	r.HandleFunc("/products/search/{search}", products.Search(s)).Methods("GET")

	// Reviews
	reviews := h.Reviews{DB: db}
	r.HandleFunc("/reviews", reviews.Find(l)).Methods("GET")
	r.HandleFunc("/reviews/{id}", reviews.FindByID(l)).Methods("GET")
	r.HandleFunc("/reviews/add", m.RequireLogin(reviews.Add(a))).Methods("POST")
	r.HandleFunc("/reviews/{id}", m.AdminsOnly(reviews.Delete(d))).Methods("DELETE")

	// Shopping
	r.HandleFunc("/cart", m.RequireLogin(h.CartGet(db))).Methods("GET")
	r.HandleFunc("/cart/add/{quantity}", m.RequireLogin(h.CartAdd(db))).Methods("POST")
	r.HandleFunc("/cart/brand/{brand}", m.RequireLogin(h.CartFilterByBrand(db))).Methods("GET")
	r.HandleFunc("/cart/category/{category}", m.RequireLogin(h.CartFilterByCategory(db))).Methods("GET")
	r.HandleFunc("/cart/discount/{min}/{max}", m.RequireLogin(h.CartFilterByDiscount(db))).Methods("GET")
	r.HandleFunc("/cart/checkout", m.RequireLogin(h.CartCheckout(db))).Methods("GET")
	r.HandleFunc("/cart/items", m.RequireLogin(h.CartItems(db))).Methods("GET")
	r.HandleFunc("/cart/remove/{id}/{quantity}", m.RequireLogin(h.CartRemove(db))).Methods("DELETE")
	r.HandleFunc("/cart/reset", m.RequireLogin(h.CartReset(db))).Methods("GET")
	r.HandleFunc("/cart/size", m.RequireLogin(h.CartSize(db))).Methods("GET")
	r.HandleFunc("/cart/taxes/{min}/{max}", m.RequireLogin(h.CartFilterByTaxes(db))).Methods("GET")
	r.HandleFunc("/cart/total/{min}/{max}", m.RequireLogin(h.CartFilterByTotal(db))).Methods("GET")
	r.HandleFunc("/cart/type/{type}", m.RequireLogin(h.CartFilterByType(db))).Methods("GET")
	r.HandleFunc("/cart/weight/{min}/{max}", m.RequireLogin(h.CartFilterByWeight(db))).Methods("GET")

	//	Shops
	shops := h.Shops{DB: db}
	r.HandleFunc("/shops", shops.Find(l)).Methods("GET")
	r.HandleFunc("/shops/{id}", shops.FindByID(l)).Methods("GET")
	r.HandleFunc("/shops/add", shops.Add(a)).Methods("POST")
	r.HandleFunc("/shops/{id}", m.AdminsOnly(shops.Update(u))).Methods("PUT")
	r.HandleFunc("/shops/{id}", m.AdminsOnly(shops.Delete(d))).Methods("DELETE")
	r.HandleFunc("/shops/search/{search}", shops.Search(s)).Methods("GET")

	// Tracking
	r.HandleFunc("/tracker", m.AdminsOnly(h.GetHits(ctx, tracker))).Methods("GET")
	r.HandleFunc("/tracker/{id}", m.AdminsOnly(h.DeleteHit(ctx, tracker))).Methods("DELETE")
	r.HandleFunc("/tracker/search/{search}", m.AdminsOnly(h.SearchHit(ctx, tracker))).Methods("GET")
	r.HandleFunc("/tracker/{field}/{value}", m.AdminsOnly(h.SearchHitByField(ctx, tracker))).Methods("GET")

	// Users
	users := h.Users{DB: db}
	r.HandleFunc("/users", users.Find(l)).Methods("GET")
	r.HandleFunc("/users/{id}", users.FindByID(l)).Methods("GET")
	r.HandleFunc("/users/add", users.Add(a, pendingList)).Methods("POST")
	r.HandleFunc("/users/{id}", m.RequireLogin(users.Update(u))).Methods("PUT")
	r.HandleFunc("/users/{id}", m.RequireLogin(users.Delete(d, session, pendingList, validatedList))).Methods("DELETE")
	r.HandleFunc("/users/search/{search}", users.Search(s)).Methods("GET")

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	http.Handle("/", r)

	return r
}
