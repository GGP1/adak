package rest

import (
	"net/http"

	"github.com/GGP1/adak/cmd/server"
	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/auth"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shop"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"
	"github.com/GGP1/adak/pkg/shopping/payment/stripe"
	"github.com/GGP1/adak/pkg/user"
	"github.com/GGP1/adak/pkg/user/account"

	// m -> middleware
	m "github.com/GGP1/adak/pkg/http/rest/middleware"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

// Frontend implements the frontend service.
type Frontend struct {
	server         *server.Server
	config         *config.Config
	accountClient  account.AccountsClient
	productClient  product.ProductsClient
	reviewClient   review.ReviewsClient
	shopClient     shop.ShopsClient
	userClient     user.UsersClient
	orderingClient ordering.OrderingClient
	sessionClient  auth.SessionClient
	shoppingClient cart.ShoppingClient
}

// NewFrontend returns the frontend server.
func NewFrontend(config *config.Config, accountConn,
	productConn, reviewConn, shopConn, userConn, orderingConn,
	sessionConn, shoppingConn *grpc.ClientConn) *Frontend {
	return &Frontend{
		config:         config,
		accountClient:  account.NewAccountsClient(accountConn),
		productClient:  product.NewProductsClient(productConn),
		reviewClient:   review.NewReviewsClient(reviewConn),
		shopClient:     shop.NewShopsClient(shopConn),
		userClient:     user.NewUsersClient(userConn),
		orderingClient: ordering.NewOrderingClient(orderingConn),
		sessionClient:  auth.NewSessionClient(sessionConn),
		shoppingClient: cart.NewShoppingClient(shopConn),
	}
}

// Run runs the server.
func (s *Frontend) Run(port int) error {
	r := chi.NewRouter()

	// Tracking
	// TODO: write service
	// trackingService := tracking.NewService(db, "")

	// Middlewares
	r.Use(m.Cors, m.Secure, m.LimitRate, m.LogFormatter)

	// Auth
	r.Post("/login", s.Login())
	r.Get("/login/google", s.LoginGoogle())
	r.Get("/login/oauth2/google", s.OAUTH2Google())
	r.With(m.RequireLogin).Get("/logout", s.Logout())

	// Cart
	r.Route("/cart", func(r chi.Router) {
		r.Use(m.RequireLogin)

		r.Get("/", s.CartGet())
		r.Post("/add/{quantity}", s.CartAdd())
		r.Get("/brand/{brand}", s.CartFilterByBrand())
		r.Get("/category/{category}", s.CartFilterByCategory())
		r.Get("/discount/{min}/{max}", s.CartFilterByDiscount())
		r.Get("/checkout", s.CartCheckout())
		r.Get("/products", s.CartProducts())
		r.Delete("/remove/{id}/{quantity}", s.CartRemove())
		r.Get("/reset", s.CartReset())
		r.Get("/size", s.CartSize())
		r.Get("/taxes/{min}/{max}", s.CartFilterByTaxes())
		r.Get("/total/{min}/{max}", s.CartFilterByTotal())
		r.Get("/type/{type}", s.CartFilterByType())
		r.Get("/weight/{min}/{max}", s.CartFilterByWeight())
	})

	// Home
	r.Get("/", Home(nil /*trackingService*/))

	// Ordering
	r.With(m.AdminsOnly).Get("/orders", s.OrderingGet())
	r.With(m.AdminsOnly).Delete("/order/{id}", s.OrderingDelete())
	r.With(m.AdminsOnly).Get("/order/{id}", s.OrderingGetByID())
	r.With(m.RequireLogin).Get("/order/user/{id}", s.OrderingGetByUserID())
	r.With(m.RequireLogin).Post("/order/new", s.OrderingNew())

	// Product
	r.Route("/products", func(r chi.Router) {
		r.Get("/", s.ProductGet())
		r.Get("/{id}", s.ProductGetByID())
		r.With(m.AdminsOnly).Put("/{id}", s.ProductUpdate())
		r.With(m.AdminsOnly).Delete("/{id}", s.ProductDelete())
		r.With(m.AdminsOnly).Post("/create", s.ProductCreate())
		r.Get("/search/{query}", s.ProductSearch())
	})

	// Review
	r.Route("/reviews", func(r chi.Router) {
		r.Get("/", s.ReviewGet())
		r.Get("/{id}", s.ReviewGetByID())
		r.With(m.AdminsOnly).Delete("/{id}", s.ReviewDelete())
		r.With(m.RequireLogin).Post("/create", s.ReviewCreate())
	})

	// Shop
	r.Route("/shops", func(r chi.Router) {
		r.Get("/", s.ShopGet())
		r.Get("/{id}", s.ShopGetByID())
		r.With(m.AdminsOnly).Delete("/{id}", s.ShopDelete())
		r.With(m.AdminsOnly).Put("/{id}", s.ShopUpdate())
		r.With(m.AdminsOnly).Post("/create", s.ShopCreate())
		r.Get("/search/{query}", s.ShopSearch())
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
	// tracker := tracking.Handler{TrackerSv: trackingService}
	// r.Route("/tracker", func(r chi.Router) {
	// 	r.Use(m.AdminsOnly)

	// 	r.Get("/", tracker.GetHits())
	// 	r.Delete("/{id}", tracker.DeleteHit())
	// 	r.Get("/search/{query}", tracker.SearchHit())
	// 	r.Get("/{field}/{value}", tracker.SearchHitByField())
	// })

	// User
	r.Route("/users", func(r chi.Router) {
		r.Get("/", s.UserGet())
		r.Get("/{id}", s.UserGetByID())
		r.With(m.RequireLogin).Delete("/{id}", s.UserDelete())
		r.With(m.RequireLogin).Put("/{id}", s.UserUpdate())
		r.Post("/create", s.UserCreate())
		r.Get("/search/{query}", s.UserSearch())
	})

	// Account
	r.With(m.RequireLogin).Post("/settings/email", s.AccountSendChangeConfirmation())
	r.With(m.RequireLogin).Post("/settings/password", s.AccountChangePassword())
	r.Get("/verification/{email}/{token}", s.AccountSendEmailValidation())
	r.Get("/verification/{token}/{email}/{id}", s.AccountChangeEmail())

	http.Handle("/", r)

	srv := server.New(s.config, r)

	return srv.Start()
}
