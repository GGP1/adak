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
)

// Server holds a http server
type Server struct {
	*http.Server
}

// New returns a new Server
func New(port string, router http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:           ":" + port,
			Handler:        router,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Start runs the server listening for errors
func (srv *Server) Start() error {
	serverErr := make(chan error, 1)

	go func() {
		fmt.Println("Listening on port", srv.Addr)
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

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := srv.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("main: Graceful shutdown did not complete in %v : %v", timeout, err)
		}

		err = srv.Close()
		if err != nil {
			return fmt.Errorf("main: Couldn't stop server gracefully : %v", err)
		}
		return nil
	}
}
