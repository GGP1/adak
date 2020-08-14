package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/stripe/stripe-go"
)

// Server holds a server configurations.
type Server struct {
	*http.Server
	TimeoutShutdown time.Duration
	Logger
}

// Logger holds the logger level.
type Logger struct {
	Stripe
}

// Stripe holds stripe configurations.
type Stripe struct {
	Level stripe.Level
}

// New returns a new server.
func New(c *Config, router http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:           c.Server.Host + ":" + c.Server.Port,
			Handler:        router,
			ReadTimeout:    c.Server.Timeout.Read * time.Second,
			WriteTimeout:   c.Server.Timeout.Write * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		c.Server.Timeout.Shutdown * time.Second,
		Logger{
			Stripe{
				c.Logger.Stripe.Level,
			},
		},
	}
}

// Start runs the server listening for errors.
func (srv *Server) Start() error {
	stripe.Key = cfg.StripeKey

	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: srv.Logger.Stripe.Level,
	}

	serverErr := make(chan error, 1)

	go func() {
		fmt.Println("Listening on", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Shutdown
	select {
	case err := <-serverErr:
		return fmt.Errorf("error: Listening and serving failed %s", err)

	case <-shutdown:
		log.Println("main: Start shutdown")

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), srv.TimeoutShutdown)
		defer cancel()

		// Asking listener to shutdown and load shed
		err := srv.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("main: Graceful shutdown did not complete in %v : %v", srv.TimeoutShutdown, err)
		}

		err = srv.Close()
		if err != nil {
			return fmt.Errorf("main: Couldn't stop server gracefully : %v", err)
		}
		return nil
	}
}
