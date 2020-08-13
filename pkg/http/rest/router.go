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

	// Repositories
	cRepo := *new(creating.Repository)
	dRepo := *new(deleting.Repository)
	lRepo := *new(listing.Repository)
	sRepo := *new(searching.Repository)
	uRepo := *new(updating.Repository)

	// Services
	c := creating.NewService(cRepo)
	d := deleting.NewService(dRepo)
	l := listing.NewService(lRepo)
	s := searching.NewService(sRepo)
	u := updating.NewService(uRepo)

	// -- Auth session --
	session := auth.NewSession(db)
	// -- Email lists --
	pendingList := email.NewList(db, "pending_list")
	validatedList := email.NewList(db, "validated_list")
	// -- Tracker --
	tracker := tracking.NewTracker(db, "")

	// Create handlers
	products := h.Products{DB: db}
	reviews := h.Reviews{DB: db}
	shops := h.Shops{DB: db}
	users := h.Users{DB: db}

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	// Auth
	r.Post("/login", h.Login(session, validatedList))
	r.Get("/logout", m.RequireLogin(h.Logout(session)))
	r.Post("/settings/email", m.RequireLogin(h.EmailChange(db, pendingList, l)))
	r.Post("/settings/password", m.RequireLogin(h.PasswordChange(session, l)))
	r.Get("/verification/{token}", h.ValidateEmail(pendingList, validatedList))
	r.Get("/verification/{token}/{email}/{id}", h.EmailChangeConfirmation(session, validatedList))

	// Creating
	r.Post("/products/create", m.AdminsOnly(products.Create(c)))
	r.Post("/reviews/create", m.RequireLogin(reviews.Create(c)))
	r.Post("/shops/create", m.AdminsOnly(shops.Create(c)))
	r.Post("/users/create", users.Create(c, pendingList))

	// Deleting
	r.Delete("/products/{id}", m.AdminsOnly(products.Delete(d)))
	r.Delete("/reviews/{id}", m.AdminsOnly(reviews.Delete(d)))
	r.Delete("/shops/{id}", m.AdminsOnly(shops.Delete(d)))
	r.Delete("/users/{id}", m.RequireLogin(users.Delete(d, session, pendingList, validatedList)))

	// Home
	r.Get("/", h.Home(tracker))

	// Listing
	r.Get("/products", products.Get(l))
	r.Get("/products/{id}", products.GetByID(l))
	r.Get("/reviews", reviews.Get(l))
	r.Get("/reviews/{id}", reviews.GetByID(l))
	r.Get("/shops", shops.Get(l))
	r.Get("/shops/{id}", shops.GetByID(l))
	r.Get("/users", users.Get(l))
	r.Get("/users/{id}", users.GetByID(l))

	// Ordering
	r.Get("/orders", m.AdminsOnly(h.GetOrder(db)))
	r.Delete("/order/{id}", m.AdminsOnly(h.DeleteOrder(db)))
	r.Post("/order/new", m.RequireLogin(h.NewOrder(db)))

	// Payment
	r.Post("/payment", h.CreatePayment())

	// Searching
	r.Get("/products/search/{search}", products.Search(s))
	r.Get("/shops/search/{search}", shops.Search(s))
	r.Get("/users/search/{search}", users.Search(s))

	// Shopping
	r.Get("/cart", m.RequireLogin(h.CartGet(db)))
	r.Post("/cart/add/{quantity}", m.RequireLogin(h.CartAdd(db)))
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

	// Tracking
	r.Get("/tracker", m.AdminsOnly(h.GetHits(tracker)))
	r.Delete("/tracker/{id}", m.AdminsOnly(h.DeleteHit(tracker)))
	r.Get("/tracker/search/{search}", m.AdminsOnly(h.SearchHit(tracker)))
	r.Get("/tracker/{field}/{value}", m.AdminsOnly(h.SearchHitByField(tracker)))

	// Updating
	r.Put("/products/{id}", m.AdminsOnly(products.Update(u)))
	r.Put("/shops/{id}", m.AdminsOnly(shops.Update(u)))
	r.Put("/users/{id}", m.RequireLogin(users.Update(u)))

	// Users
	r.Get("/users/{id}/qrcode", users.QRCode(l))

	http.Handle("/", r)

	return r
}
