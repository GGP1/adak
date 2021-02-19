package rest

import (
	"net/http"

	"github.com/GGP1/adak/pkg/auth"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shop"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"
	"github.com/GGP1/adak/pkg/shopping/payment/stripe"
	"github.com/GGP1/adak/pkg/tracking"
	"github.com/GGP1/adak/pkg/user"
	"github.com/GGP1/adak/pkg/user/account"

	// m -> middleware
	m "github.com/GGP1/adak/pkg/http/rest/middleware"

	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"
)

// NewRouter initializes services, creates and returns a mux router
func NewRouter(db *sqlx.DB, cache *lru.Cache) http.Handler {
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
	// -- Tracking --
	trackingService := tracking.NewService(db)

	// Middlewares
	r.Use(m.Cors, m.Secure, m.LimitRate, m.LogFormatter, m.GZIPCompress)

	// Auth
	r.Post("/login", auth.Login(session))
	r.With(m.RequireLogin).Get("/logout", auth.Logout(session))
	r.Get("/login/google", auth.LoginGoogle(session))
	r.Get("/login/oauth2/google", auth.OAuth2Google(session))

	// Cart
	cart := cart.Handler{
		DB:    db,
		Cache: cache,
	}
	r.Route("/cart", func(r chi.Router) {
		r.Use(m.RequireLogin)

		r.Get("/", cart.Get())
		r.Post("/add/{quantity}", cart.Add())
		r.Get("/brand/{brand}", cart.FilterByBrand())
		r.Get("/category/{category}", cart.FilterByCategory())
		r.Get("/discount/{min}/{max}", cart.FilterByDiscount())
		r.Get("/checkout", cart.Checkout())
		r.Get("/products", cart.Products())
		r.Delete("/remove/{id}/{quantity}", cart.Remove())
		r.Get("/reset", cart.Reset())
		r.Get("/size", cart.Size())
		r.Get("/taxes/{min}/{max}", cart.FilterByTaxes())
		r.Get("/total/{min}/{max}", cart.FilterByTotal())
		r.Get("/type/{type}", cart.FilterByType())
		r.Get("/weight/{min}/{max}", cart.FilterByWeight())
	})

	// Home
	r.Get("/", Home(trackingService))

	// Ordering
	order := ordering.Handler{
		DB:    db,
		Cache: cache,
	}
	r.With(m.AdminsOnly).Get("/orders", order.Get())
	r.With(m.AdminsOnly).Delete("/order/{id}", order.Delete())
	r.With(m.AdminsOnly).Get("/order/{id}", order.GetByID())
	r.With(m.RequireLogin).Get("/order/user/{id}", order.GetByUserID())
	r.With(m.RequireLogin).Post("/order/new", order.New())

	// Product
	product := product.Handler{
		Service: productService,
		Cache:   cache,
	}
	r.Route("/products", func(r chi.Router) {
		r.Get("/", product.Get())
		r.Get("/{id}", product.GetByID())
		r.With(m.AdminsOnly).Put("/{id}", product.Update())
		r.With(m.AdminsOnly).Delete("/{id}", product.Delete())
		r.With(m.AdminsOnly).Post("/create", product.Create())
		r.Get("/search/{query}", product.Search())
	})

	// Review
	review := review.Handler{
		Service: reviewService,
		Cache:   cache,
	}
	r.Route("/reviews", func(r chi.Router) {
		r.Get("/", review.Get())
		r.Get("/{id}", review.GetByID())
		r.With(m.AdminsOnly).Delete("/{id}", review.Delete())
		r.With(m.RequireLogin).Post("/create", review.Create())
	})

	// Shop
	shop := shop.Handler{
		Service: shopService,
		Cache:   cache,
	}
	r.Route("/shops", func(r chi.Router) {
		r.Get("/", shop.Get())
		r.Get("/{id}", shop.GetByID())
		r.With(m.AdminsOnly).Delete("/{id}", shop.Delete())
		r.With(m.AdminsOnly).Put("/{id}", shop.Update())
		r.With(m.AdminsOnly).Post("/create", shop.Create())
		r.Get("/search/{query}", shop.Search())
	})

	// Stripe
	stripe := stripe.Handler{}
	r.Route("/stripe", func(r chi.Router) {
		r.Use(m.AdminsOnly)

		r.Get("/balance", stripe.GetBalance())
		r.Get("/event/{event}", stripe.GetEvent())
		r.Get("/transactions/{txID}", stripe.GetTxBalance())
		r.Get("/events", stripe.ListEvents())
		r.Get("/transactions", stripe.ListTxs())
	})

	// Tracking
	tracker := tracking.Handler{Service: trackingService}
	r.Route("/tracker", func(r chi.Router) {
		r.Use(m.AdminsOnly)

		r.Get("/", tracker.GetHits())
		r.Delete("/{id}", tracker.DeleteHit())
		r.Get("/search/{query}", tracker.SearchHit())
		r.Get("/{field}/{value}", tracker.SearchHitByField())
	})

	// User
	user := user.Handler{
		Service: userService,
		Cache:   cache,
	}
	r.Route("/users", func(r chi.Router) {
		r.Get("/", user.Get())
		r.Get("/{id}", user.GetByID())
		r.With(m.RequireLogin).Delete("/{id}", user.Delete(db, session))
		r.With(m.RequireLogin).Put("/{id}", user.Update())
		r.Get("/email/{email}", user.GetByEmail())
		r.Get("/username/{username}", user.GetByUsername())
		r.Get("/{id}/qrcode", user.QRCode())
		r.Post("/create", user.Create())
		r.Get("/search/{query}", user.Search())
	})

	// Account
	account := account.Handler{Service: accountService}
	r.With(m.RequireLogin).Post("/settings/email", account.SendChangeConfirmation(userService))
	r.With(m.RequireLogin).Post("/settings/password", account.ChangePassword())
	r.Get("/verification/{email}/{token}", account.SendEmailValidation(userService))
	r.Get("/verification/{token}/{email}/{id}", account.ChangeEmail())

	http.Handle("/", r)

	return r
}
