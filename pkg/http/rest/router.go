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
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter initializes services, creates and returns a mux router
func NewRouter(config config.Config, db *sqlx.DB, mc *memcache.Client, rdb *redis.Client) http.Handler {
	router := chi.NewRouter()

	// Services
	accountService := account.NewService(db)
	cartService := cart.NewService(db, mc)
	orderingService := ordering.NewService(db)
	productService := product.NewService(db, mc)
	reviewService := review.NewService(db, mc)
	shopService := shop.NewService(db, mc)
	userService := user.NewService(db, mc)
	trackingService := tracking.NewService(db)
	session := auth.NewSession(db, rdb, config.Session, config.Development)
	emailer := email.New()

	// Authentication middleware
	mAuth := middleware.Auth{
		DB:          db,
		UserService: userService,
		Session:     session,
	}
	adminsOnly := mAuth.AdminsOnly
	requireLogin := mAuth.RequireLogin
	// Metrics middleware
	metrics := middleware.NewMetrics()

	// Middlewares
	router.Use(middleware.Cors, middleware.Secure, middleware.Recover,
		middleware.LogFormatter, middleware.GZIPCompress, metrics.Scrap)

	// Must be after the other middlewares otherwise they won't have effect when rate limiting
	if config.RateLimiter.Rate > 0 {
		rateLimiter := middleware.NewRateLimiter(config.RateLimiter, rdb)
		router.Use(rateLimiter.Limit)
	}

	// Auth
	router.Post("/login", auth.Login(session))
	router.Get("/login/basic", auth.BasicAuth(session))
	router.With(requireLogin).Get("/logout", auth.Logout(session))
	router.Get("/login/google", auth.LoginGoogle(session))
	router.Get("/login/oauth2/google", auth.OAuth2Google(session))

	// Cart
	cart := cart.NewHandler(cartService, db, mc)
	router.Route("/cart", func(r chi.Router) {
		r.Use(requireLogin)

		r.Get("/", cart.Get())
		r.Post("/add", cart.Add())
		r.Get("/filter/{field}/{args}", cart.FilterBy())
		r.Get("/checkout", cart.Checkout())
		r.Get("/products", cart.Products())
		r.Delete("/remove/{id}/{quantity}", cart.Remove())
		r.Get("/reset", cart.Reset())
		r.Get("/size", cart.Size())
	})

	// Home
	router.Get("/", Home(trackingService))

	// Metrics
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		Registry: prometheus.DefaultRegisterer,
		// The response is already compressed by the gzip middleware, avoid double compression
		DisableCompression: true,
		EnableOpenMetrics:  true,
	}))

	// Ordering
	order := ordering.NewHandler(config.Development, orderingService, cartService, db, mc)
	router.Route("/orders", func(r chi.Router) {
		r.With(adminsOnly).Get("/", order.Get())
		r.With(adminsOnly).Delete("/{id}", order.Delete())
		r.With(adminsOnly).Get("/{id}", order.GetByID())
		r.With(requireLogin).Get("/user/{id}", order.GetByUserID())
		r.With(requireLogin).Post("/new", order.New())
	})

	// Product
	product := product.NewHandler(productService, mc)
	router.Route("/products", func(r chi.Router) {
		r.Get("/", product.Get())
		r.Get("/{id}", product.GetByID())
		r.With(adminsOnly).Put("/{id}", product.Update())
		r.With(adminsOnly).Delete("/{id}", product.Delete())
		r.With(adminsOnly).Post("/create", product.Create())
		r.Get("/search/{query}", product.Search())
	})

	// Review
	review := review.NewHandler(reviewService, mc)
	router.Route("/reviews", func(r chi.Router) {
		r.Get("/", review.Get())
		r.Get("/{id}", review.GetByID())
		r.With(adminsOnly).Delete("/{id}", review.Delete())
		r.With(requireLogin).Post("/create", review.Create())
	})

	// Shop
	shop := shop.NewHandler(shopService, mc)
	router.Route("/shops", func(r chi.Router) {
		r.Get("/", shop.Get())
		r.Get("/{id}", shop.GetByID())
		r.With(adminsOnly).Delete("/{id}", shop.Delete())
		r.With(adminsOnly).Put("/{id}", shop.Update())
		r.With(adminsOnly).Post("/create", shop.Create())
		r.Get("/search/{query}", shop.Search())
	})

	// Stripe
	stripe := stripe.NewHandler()
	router.Route("/stripe", func(r chi.Router) {
		r.Use(adminsOnly)

		r.Get("/balance", stripe.GetBalance())
		r.Get("/event/{event}", stripe.GetEvent())
		r.Get("/transactions/{txID}", stripe.GetTxBalance())
		r.Get("/events", stripe.ListEvents())
		r.Get("/transactions", stripe.ListTxs())
	})

	// Tracking
	tracker := tracking.NewHandler(trackingService)
	router.Route("/tracker", func(r chi.Router) {
		r.Use(adminsOnly)

		r.Get("/", tracker.GetHits())
		r.Delete("/{id}", tracker.DeleteHit())
		r.Get("/search/{query}", tracker.SearchHit())
		r.Get("/{field}/{value}", tracker.SearchHitByField())
	})

	// User
	user := user.NewHandler(config.Development, userService, cartService, emailer, mc)
	router.Route("/users", func(r chi.Router) {
		r.Get("/", user.Get())
		r.Get("/{id}", user.GetByID())
		r.With(requireLogin).Delete("/{id}", user.Delete(session))
		r.With(requireLogin).Put("/{id}", user.Update())
		r.Get("/email/{email}", user.GetByEmail())
		r.Get("/username/{username}", user.GetByUsername())
		r.Post("/create", user.Create())
		r.Get("/search/{query}", user.Search())
	})

	// Account
	account := account.NewHandler(accountService, userService, emailer)
	router.With(requireLogin).Post("/settings/email", account.SendChangeConfirmation())
	router.With(requireLogin).Post("/settings/password", account.ChangePassword())
	router.Get("/verification/{email}/{token}", account.SendEmailValidation(userService))
	router.Get("/verification/{token}/{email}/{id}", account.ChangeEmail())

	http.Handle("/", router)
	return router
}
