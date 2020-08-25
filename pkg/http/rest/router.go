package rest

import (
	"net/http"

	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/email"
	"github.com/GGP1/palo/pkg/product"
	"github.com/GGP1/palo/pkg/review"
	"github.com/GGP1/palo/pkg/shop"
	"github.com/GGP1/palo/pkg/shopping/cart"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/GGP1/palo/pkg/shopping/payment/stripe"
	"github.com/GGP1/palo/pkg/tracking"
	"github.com/GGP1/palo/pkg/user"
	"github.com/GGP1/palo/pkg/user/account"

	// m -> middleware
	m "github.com/GGP1/palo/pkg/http/rest/middleware"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// NewRouter initializes services, creates and returns a mux router
func NewRouter(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()

	// Service repositories
	aRepo := *new(account.Repository)
	pRepo := *new(product.Repository)
	rRepo := *new(review.Repository)
	sRepo := *new(shop.Repository)
	uRepo := *new(user.Repository)

	// Services
	a := account.NewService(aRepo, db)
	p := product.NewService(pRepo, db)
	rev := review.NewService(rRepo, db)
	s := shop.NewService(sRepo, db)
	u := user.NewService(uRepo, db)

	// -- Auth session --
	session := auth.NewSession(db)
	// -- Email lists --
	pendingList := email.NewList(db, "pending_list")
	validatedList := email.NewList(db, "validated_list")
	// -- Tracker --
	tracker := tracking.NewTracker(db, "")

	// Middlewares
	r.Use(m.AllowCrossOrigin)
	r.Use(m.LimitRate)
	r.Use(m.LogFormatter)

	// Auth
	r.Post("/login", auth.Login(session, validatedList))
	r.Get("/logout", m.RequireLogin(auth.Logout(session)))

	// Cart
	cart := cart.Handler{DB: db}
	r.Get("/cart", m.RequireLogin(cart.Get()))
	r.Post("/cart/add/{quantity}", m.RequireLogin(cart.Add()))
	r.Get("/cart/brand/{brand}", m.RequireLogin(cart.FilterByBrand()))
	r.Get("/cart/category/{category}", m.RequireLogin(cart.FilterByCategory()))
	r.Get("/cart/discount/{min}/{max}", m.RequireLogin(cart.FilterByDiscount()))
	r.Get("/cart/checkout", m.RequireLogin(cart.Checkout()))
	r.Get("/cart/products", m.RequireLogin(cart.Products()))
	r.Delete("/cart/remove/{id}/{quantity}", m.RequireLogin(cart.Remove()))
	r.Get("/cart/reset", m.RequireLogin(cart.Reset()))
	r.Get("/cart/size", m.RequireLogin(cart.Size()))
	r.Get("/cart/taxes/{min}/{max}", m.RequireLogin(cart.FilterByTaxes()))
	r.Get("/cart/total/{min}/{max}", m.RequireLogin(cart.FilterByTotal()))
	r.Get("/cart/type/{type}", m.RequireLogin(cart.FilterByType()))
	r.Get("/cart/weight/{min}/{max}", m.RequireLogin(cart.FilterByWeight()))

	// Home
	r.Get("/", Home(tracker))

	// Ordering
	r.Get("/orders", m.AdminsOnly(ordering.GetOrder(db)))
	r.Delete("/order/{id}", m.AdminsOnly(ordering.DeleteOrder(db)))
	r.Post("/order/new", m.RequireLogin(ordering.NewOrder(db)))

	// Product
	r.Post("/products/create", m.AdminsOnly(product.Create(p)))
	r.Delete("/products/{id}", m.AdminsOnly(product.Delete(p)))
	r.Get("/products", product.Get(p))
	r.Get("/products/{id}", product.GetByID(p))
	r.Get("/products/search/{query}", product.Search(p))
	r.Put("/products/{id}", m.AdminsOnly(product.Update(p)))

	// Review
	r.Post("/reviews/create", m.RequireLogin(review.Create(rev)))
	r.Delete("/reviews/{id}", m.AdminsOnly(review.Delete(rev)))
	r.Get("/reviews", review.Get(rev))
	r.Get("/reviews/{id}", review.GetByID(rev))

	// Shop
	r.Post("/shops/create", m.AdminsOnly(shop.Create(s)))
	r.Delete("/shops/{id}", m.AdminsOnly(shop.Delete(s)))
	r.Get("/shops", shop.Get(s))
	r.Get("/shops/{id}", shop.GetByID(s))
	r.Get("/shops/search/{query}", shop.Search(s))
	r.Put("/shops/{id}", m.AdminsOnly(shop.Update(s)))

	// Stripe
	stripe := stripe.Handler{}
	r.Get("/stripe/balance", m.AdminsOnly(stripe.GetBalance()))
	r.Get("/stripe/event/{event}", m.AdminsOnly(stripe.GetEvent()))
	r.Get("/stripe/transactions/{txID}", m.AdminsOnly(stripe.GetTxBalance()))
	r.Get("/stripe/events", m.AdminsOnly(stripe.ListEvents()))
	r.Get("/stripe/transactions", m.AdminsOnly(stripe.ListTxs()))

	// Tracking
	r.Get("/tracker", m.AdminsOnly(tracking.GetHits(tracker)))
	r.Delete("/tracker/{id}", m.AdminsOnly(tracking.DeleteHit(tracker)))
	r.Get("/tracker/search/{search}", m.AdminsOnly(tracking.SearchHit(tracker)))
	r.Get("/tracker/{field}/{value}", m.AdminsOnly(tracking.SearchHitByField(tracker)))

	// User
	r.Post("/users/create", user.Create(u, pendingList))
	r.Delete("/users/{id}", m.RequireLogin(user.Delete(db, u, session, pendingList, validatedList)))
	r.Get("/users", user.Get(u))
	r.Get("/users/{id}", user.GetByID(u))
	r.Get("/users/search/{query}", user.Search(u))
	r.Put("/users/{id}", m.RequireLogin(user.Update(u)))
	r.Get("/users/{id}/qrcode", user.QRCode(u))
	// Account
	r.Post("/settings/email", m.RequireLogin(account.ChangeEmail(db, pendingList)))
	r.Post("/settings/password", m.RequireLogin(account.ChangePassword(a)))
	r.Get("/verification/{token}", account.SendEmailValidation(pendingList, validatedList))
	r.Get("/verification/{token}/{email}/{id}", account.ValidateEmailChange(a, validatedList))

	http.Handle("/", r)

	return r
}
