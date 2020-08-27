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
	accountService := account.NewService(aRepo, db)
	productService := product.NewService(pRepo, db)
	reviewService := review.NewService(rRepo, db)
	shopService := shop.NewService(sRepo, db)
	userService := user.NewService(uRepo, db)

	// -- Auth session --
	session := auth.NewSession(db)
	// -- Email lists --
	pendingList := email.NewService(db, "pending_list")
	validatedList := email.NewService(db, "validated_list")
	// -- Tracking --
	trackingService := tracking.NewService(db, "")

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
	r.Get("/", Home(trackingService))

	// Ordering
	order := ordering.Handler{DB: db}
	r.Get("/orders", m.AdminsOnly(order.Get()))
	r.Delete("/order/{id}", m.AdminsOnly(order.Delete()))
	r.Post("/order/new", m.RequireLogin(order.New()))

	// Product
	product := product.Handler{Service: productService}
	r.Post("/products/create", m.AdminsOnly(product.Create()))
	r.Delete("/products/{id}", m.AdminsOnly(product.Delete()))
	r.Get("/products", product.Get())
	r.Get("/products/{id}", product.GetByID())
	r.Get("/products/search/{query}", product.Search())
	r.Put("/products/{id}", m.AdminsOnly(product.Update()))

	// Review
	review := review.Handler{Service: reviewService}
	r.Post("/reviews/create", m.RequireLogin(review.Create()))
	r.Delete("/reviews/{id}", m.AdminsOnly(review.Delete()))
	r.Get("/reviews", review.Get())
	r.Get("/reviews/{id}", review.GetByID())

	// Shop
	shop := shop.Handler{Service: shopService}
	r.Post("/shops/create", m.AdminsOnly(shop.Create()))
	r.Delete("/shops/{id}", m.AdminsOnly(shop.Delete()))
	r.Get("/shops", shop.Get())
	r.Get("/shops/{id}", shop.GetByID())
	r.Get("/shops/search/{query}", shop.Search())
	r.Put("/shops/{id}", m.AdminsOnly(shop.Update()))

	// Stripe
	stripe := stripe.Handler{}
	r.Get("/stripe/balance", m.AdminsOnly(stripe.GetBalance()))
	r.Get("/stripe/event/{event}", m.AdminsOnly(stripe.GetEvent()))
	r.Get("/stripe/transactions/{txID}", m.AdminsOnly(stripe.GetTxBalance()))
	r.Get("/stripe/events", m.AdminsOnly(stripe.ListEvents()))
	r.Get("/stripe/transactions", m.AdminsOnly(stripe.ListTxs()))

	// Tracking
	tracker := tracking.Handler{TrackerSv: trackingService}
	r.Get("/tracker", m.AdminsOnly(tracker.GetHits()))
	r.Delete("/tracker/{id}", m.AdminsOnly(tracker.DeleteHit()))
	r.Get("/tracker/search/{query}", m.AdminsOnly(tracker.SearchHit()))
	r.Get("/tracker/{field}/{value}", m.AdminsOnly(tracker.SearchHitByField()))

	// User
	user := user.Handler{Service: userService}
	r.Post("/users/create", user.Create(pendingList))
	r.Delete("/users/{id}", m.RequireLogin(user.Delete(db, session, pendingList, validatedList)))
	r.Get("/users", user.Get())
	r.Get("/users/{id}", user.GetByID())
	r.Get("/users/search/{query}", user.Search())
	r.Put("/users/{id}", m.RequireLogin(user.Update()))
	r.Get("/users/{id}/qrcode", user.QRCode())
	// Account
	account := account.Handler{Service: accountService}
	r.Post("/settings/email", m.RequireLogin(account.SendChangeConfirmation(userService, pendingList)))
	r.Post("/settings/password", m.RequireLogin(account.ChangePassword()))
	r.Get("/verification/{token}", account.SendEmailValidation(pendingList, validatedList))
	r.Get("/verification/{token}/{email}/{id}", account.ChangeEmail(validatedList))

	http.Handle("/", r)

	return r
}
