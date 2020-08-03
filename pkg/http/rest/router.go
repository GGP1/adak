/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"net/http"

	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/creating"
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

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

// NewRouter initializes services, creates and returns a mux router
func NewRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	// Create repository
	repo := *new(repository.Repo)

	// Services
	c := creating.NewService(repo)
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
	tracker := *tracking.NewTracker(db, "")

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	// Auth
	r.Post("/login", h.Login(session, validatedList))
	r.Get("/logout", m.RequireLogin(h.Logout(session)))
	r.Get("/verification/{token}", h.ValidateEmail(pendingList, validatedList))

	// Home
	r.Get("/", h.Home(tracker))

	// Products
	products := h.Products{DB: db}
	r.Get("/products", products.Get(l))
	r.Get("/products/{id}", products.GetByID(l))
	r.Post("/products/create", products.Create(c))
	r.Put("/products/{id}", m.AdminsOnly(products.Update(u)))
	r.Delete("/products/{id}", m.AdminsOnly(products.Delete(d)))
	r.Get("/products/search/{search}", products.Search(s))

	// Reviews
	reviews := h.Reviews{DB: db}
	r.Get("/reviews", reviews.Get(l))
	r.Get("/reviews/{id}", reviews.GetByID(l))
	r.Post("/reviews/create", m.RequireLogin(reviews.Create(c)))
	r.Delete("/reviews/{id}", m.AdminsOnly(reviews.Delete(d)))

	// Shopping
	r.Get("/cart", m.RequireLogin(h.CartGet(db)))
	r.Post("/cart/create/{quantity}", m.RequireLogin(h.CartAdd(db)))
	r.Get("/cart/brand/{brand}", m.RequireLogin(h.CartFilterByBrand(db)))
	r.Get("/cart/category/{category}", m.RequireLogin(h.CartFilterByCategory(db)))
	r.Get("/cart/discount/{min}/{max}", m.RequireLogin(h.CartFilterByDiscount(db)))
	r.Get("/cart/checkout", m.RequireLogin(h.CartCheckout(db)))
	r.Get("/cart/items", m.RequireLogin(h.CartItems(db)))
	r.Delete("/cart/remove/{id}/{quantity}", m.RequireLogin(h.CartRemove(db)))
	r.Get("/cart/reset", m.RequireLogin(h.CartReset(db)))
	r.Get("/cart/size", m.RequireLogin(h.CartSize(db)))
	r.Get("/cart/taxes/{min}/{max}", m.RequireLogin(h.CartFilterByTaxes(db)))
	r.Get("/cart/total/{min}/{max}", m.RequireLogin(h.CartFilterByTotal(db)))
	r.Get("/cart/type/{type}", m.RequireLogin(h.CartFilterByType(db)))
	r.Get("/cart/weight/{min}/{max}", m.RequireLogin(h.CartFilterByWeight(db)))

	// Shops
	shops := h.Shops{DB: db}
	r.Get("/shops", shops.Get(l))
	r.Get("/shops/{id}", shops.GetByID(l))
	r.Post("/shops/create", shops.Create(c))
	r.Put("/shops/{id}", m.AdminsOnly(shops.Update(u)))
	r.Delete("/shops/{id}", m.AdminsOnly(shops.Delete(d)))
	r.Get("/shops/search/{search}", shops.Search(s))

	// Tracking
	r.Get("/tracker", m.AdminsOnly(h.GetHits(tracker)))
	r.Delete("/tracker/{id}", m.AdminsOnly(h.DeleteHit(tracker)))
	r.Get("/tracker/search/{search}", m.AdminsOnly(h.SearchHit(tracker)))
	r.Get("/tracker/{field}/{value}", m.AdminsOnly(h.SearchHitByField(tracker)))

	// Users
	users := h.Users{DB: db}
	r.Get("/users", users.Get(l))
	r.Get("/users/{id}", users.GetByID(l))
	r.Post("/users/create", users.Create(c, pendingList))
	r.Put("/users/{id}", m.RequireLogin(users.Update(u)))
	r.Delete("/users/{id}", m.RequireLogin(users.Delete(d, session, pendingList, validatedList)))
	r.Get("/users/search/{search}", users.Search(s))

	http.Handle("/", r)

	return r
}
