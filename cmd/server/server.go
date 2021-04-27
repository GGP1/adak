package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/stripe/stripe-go"
)

// Server holds a server configurations.
type Server struct {
	*http.Server
	Stripe

	ShutdownTimeout time.Duration
}

// Stripe holds stripe configurations.
type Stripe struct {
	SecretKey string
	Level     stripe.Level
}

// New returns a new server.
func New(c *config.Config, router http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:           c.Server.Host + ":" + c.Server.Port,
			Handler:        router,
			ReadTimeout:    c.Server.Timeout.Read * time.Second,
			WriteTimeout:   c.Server.Timeout.Write * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		Stripe{
			SecretKey: c.Stripe.SecretKey,
			Level:     c.Stripe.Logger.Level,
		},
		c.Server.Timeout.Shutdown * time.Second,
	}
}

// Start runs the server listening for errors.
func (srv *Server) Start() error {
	stripe.Key = srv.Stripe.SecretKey
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: srv.Stripe.Level,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info().Msgf("Listening on %s", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return errors.Wrap(err, "Listen and serve failed")

	case <-interrupt:
		log.Info().Msg("Start shutdown")

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), srv.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed
		if err := srv.Shutdown(ctx); err != nil {
			return errors.Wrapf(err, "Graceful shutdown did not complete in %v", srv.ShutdownTimeout)
		}

		if err := srv.Close(); err != nil {
			return errors.Wrap(err, "Couldn't stop server gracefully")
		}
		return nil
	}
}
