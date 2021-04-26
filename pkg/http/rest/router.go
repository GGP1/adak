package rest

import (
	"net/http"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/email"
	"github.com/GGP1/adak/pkg/auth"
	"github.com/GGP1/adak/pkg/http/rest/middleware"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shop"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"
	"github.com/GGP1/adak/pkg/shopping/payment/stripe"
	"github.com/GGP1/adak/pkg/tracking"
	"github.com/GGP1/adak/pkg/user"
	"github.com/GGP1/adak/pkg/user/account"
	"github.com/bradfitz/gomemcache/memcache"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// NewRouter initializes services, creates and returns a mux router
func NewRouter(config config.Config, db *sqlx.DB, mc *memcache.Client) http.Handler {
	r := chi.NewRouter()

	// Services
	accountService := account.NewService(db)
	cartService := cart.NewService(db)
	orderingService := ordering.NewService(db)
	productService := product.NewService(db)
	reviewService := review.NewService(db)
	shopService := shop.NewService(db)
	userService := user.NewService(db)
	// -- Auth session --
	session := auth.NewSession(db, config.Development)
	// -- Tracking --
	trackingService := tracking.NewService(db)
	// -- Email --
	emailer := email.New()

	// Authentication middleware
	mAuth := middleware.Auth{
		DB:          db,
		UserService: userService,
	}
	adminsOnly := mAuth.AdminsOnly
	requireLogin := mAuth.RequireLogin

	// Middlewares
	r.Use(middleware.Cors, middleware.Secure, middleware.LimitRate, middleware.Recover,
		middleware.LogFormatter, middleware.GZIPCompress)

	// Auth
	r.Post("/login", auth.Login(session))
	r.With(requireLogin).Get("/logout", auth.Logout(session))
	r.Get("/login/google", auth.LoginGoogle(session))
	r.Get("/login/oauth2/google", auth.OAuth2Google(session))

	// Cart
	cart := cart.Handler{
		Service: cartService,
		DB:      db,
		Cache:   mc,
	}
	r.Route("/cart", func(r chi.Router) {
		r.Use(requireLogin)

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
		Service:     orderingService,
		CartService: cartService,
		DB:          db,
		Cache:       mc,
	}
	r.With(adminsOnly).Get("/orders", order.Get())
	r.With(adminsOnly).Delete("/order/{id}", order.Delete())
	r.With(adminsOnly).Get("/order/{id}", order.GetByID())
	r.With(requireLogin).Get("/order/user/{id}", order.GetByUserID())
	r.With(requireLogin).Post("/order/new", order.New())

	// Product
	product := product.Handler{
		Service: productService,
		Cache:   mc,
	}
	r.Route("/products", func(r chi.Router) {
		r.Get("/", product.Get())
		r.Get("/{id}", product.GetByID())
		r.With(adminsOnly).Put("/{id}", product.Update())
		r.With(adminsOnly).Delete("/{id}", product.Delete())
		r.With(adminsOnly).Post("/create", product.Create())
		r.Get("/search/{query}", product.Search())
	})

	// Review
	review := review.Handler{
		Service: reviewService,
		Cache:   mc,
	}
	r.Route("/reviews", func(r chi.Router) {
		r.Get("/", review.Get())
		r.Get("/{id}", review.GetByID())
		r.With(adminsOnly).Delete("/{id}", review.Delete())
		r.With(requireLogin).Post("/create", review.Create())
	})

	// Shop
	shop := shop.Handler{
		Service: shopService,
		Cache:   mc,
	}
	r.Route("/shops", func(r chi.Router) {
		r.Get("/", shop.Get())
		r.Get("/{id}", shop.GetByID())
		r.With(adminsOnly).Delete("/{id}", shop.Delete())
		r.With(adminsOnly).Put("/{id}", shop.Update())
		r.With(adminsOnly).Post("/create", shop.Create())
		r.Get("/search/{query}", shop.Search())
	})

	// Stripe
	stripe := stripe.Handler{}
	r.Route("/stripe", func(r chi.Router) {
		r.Use(adminsOnly)

		r.Get("/balance", stripe.GetBalance())
		r.Get("/event/{event}", stripe.GetEvent())
		r.Get("/transactions/{txID}", stripe.GetTxBalance())
		r.Get("/events", stripe.ListEvents())
		r.Get("/transactions", stripe.ListTxs())
	})

	// Tracking
	tracker := tracking.Handler{Service: trackingService}
	r.Route("/tracker", func(r chi.Router) {
		r.Use(adminsOnly)

		r.Get("/", tracker.GetHits())
		r.Delete("/{id}", tracker.DeleteHit())
		r.Get("/search/{query}", tracker.SearchHit())
		r.Get("/{field}/{value}", tracker.SearchHitByField())
	})

	// User
	user := user.Handler{
		Development: config.Development,
		Service:     userService,
		CartService: cartService,
		Emailer:     emailer,
		Cache:       mc,
	}
	r.Route("/users", func(r chi.Router) {
		r.Get("/", user.Get())
		r.Get("/{id}", user.GetByID())
		r.With(requireLogin).Delete("/{id}", user.Delete(db, session))
		r.With(requireLogin).Put("/{id}", user.Update())
		r.Get("/email/{email}", user.GetByEmail())
		r.Get("/username/{username}", user.GetByUsername())
		r.Get("/{id}/qrcode", user.QRCode())
		r.Post("/create", user.Create())
		r.Get("/search/{query}", user.Search())
	})

	// Account
	account := account.Handler{
		Service:     accountService,
		UserService: userService,
		Emailer:     emailer,
	}
	r.With(requireLogin).Post("/settings/email", account.SendChangeConfirmation())
	r.With(requireLogin).Post("/settings/password", account.ChangePassword())
	r.Get("/verification/{email}/{token}", account.SendEmailValidation(userService))
	r.Get("/verification/{token}/{email}/{id}", account.ChangeEmail())

	http.Handle("/", r)

	return r
}
