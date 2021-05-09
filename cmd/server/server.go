package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/logger"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
)

// Server holds a server configurations.
type Server struct {
	*http.Server
	TLS TLS

	Stripe          Stripe
	TimeoutShutdown time.Duration
}

// TLS contains key and certificate files paths.
type TLS struct {
	KeyFile  string
	CertFile string
}

// Stripe holds stripe configurations.
type Stripe struct {
	SecretKey string
	Level     stripe.Level
}

// New returns a new server.
func New(c config.Config, router http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:           c.Server.Host + ":" + c.Server.Port,
			Handler:        router,
			ReadTimeout:    c.Server.Timeout.Read * time.Second,
			WriteTimeout:   c.Server.Timeout.Write * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		TLS: TLS{
			KeyFile:  c.Server.TLS.KeyFile,
			CertFile: c.Server.TLS.CertFile,
		},
		Stripe: Stripe{
			SecretKey: c.Stripe.SecretKey,
			Level:     c.Stripe.Logger.Level,
		},
		TimeoutShutdown: c.Server.Timeout.Shutdown * time.Second,
	}
}

// Start runs the server listening for errors.
func (srv *Server) Start(ctx context.Context) error {
	logger.Infof("Stripe API version: %s", stripe.APIVersion)
	stripe.Key = srv.Stripe.SecretKey
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: srv.Stripe.Level,
	}

	serverErr := make(chan error, 1)

	go func() {
		if srv.TLS.CertFile != "" && srv.TLS.KeyFile != "" {
			logger.Infof("Listening on https://%s", srv.Addr)
			serverErr <- srv.ListenAndServeTLS(srv.TLS.CertFile, srv.TLS.KeyFile)
			return
		}

		logger.Infof("Listening on http://%s", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return errors.Wrap(err, "Listening and serve failed")

	case <-interrupt:
		logger.Info("Starting shutdown...")

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(ctx, srv.TimeoutShutdown)
		defer cancel()

		// Asking listener to shutdown and load shed
		if err := srv.Shutdown(ctx); err != nil {
			return errors.Wrapf(err, "Graceful shutdown did not complete in %v", srv.TimeoutShutdown)
		}

		if err := srv.Close(); err != nil {
			return errors.Wrap(err, "Couldn't stop server gracefully")
		}

		logger.Info("Server shutdown gracefully")
		return nil
	}
}
