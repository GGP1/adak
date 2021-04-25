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
	"github.com/stripe/stripe-go"
)

// Server holds a server configurations.
type Server struct {
	*http.Server
	TLS
	Stripe

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
		&http.Server{
			Addr:           c.Server.Host + ":" + c.Server.Port,
			Handler:        router,
			ReadTimeout:    c.Server.Timeout.Read * time.Second,
			WriteTimeout:   c.Server.Timeout.Write * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		TLS{
			KeyFile:  c.Server.TLS.KeyFile,
			CertFile: c.Server.TLS.CertFile,
		},
		Stripe{
			SecretKey: c.Stripe.SecretKey,
			Level:     c.Stripe.Logger.Level,
		},
		c.Server.Timeout.Shutdown * time.Second,
	}
}

// Start runs the server listening for errors.
func (srv *Server) Start(ctx context.Context) error {
	stripe.Key = srv.Stripe.SecretKey
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: srv.Stripe.Level,
	}

	serverErr := make(chan error, 1)

	go func() {
		if srv.CertFile != "" && srv.KeyFile != "" {
			logger.Log.Infof("Listening on https://%s", srv.Addr)
			serverErr <- srv.ListenAndServeTLS(srv.CertFile, srv.KeyFile)
			return
		}

		logger.Log.Infof("Listening on http://%s", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Shutdown
	select {
	case err := <-serverErr:
		return errors.Wrap(err, "Listening and serve failed")

	case <-shutdown:
		logger.Log.Info("Start shutdown")

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

		logger.Log.Info("Server shutdown gracefully")
		return nil
	}
}
