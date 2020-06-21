package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	run()
}

// func run for easier testing
func run() {
	// Create a database connection, automigrate and
	// check tables existence
	_, close, err := storage.Database()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	// Create auth session
	repo := new(handler.Repository)
	session := handler.NewSession(*repo)

	// New router
	r := rest.NewRouter(session)

	// Server setup
	server := &http.Server{
		Addr:           ":4000",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = startServer(server)
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(server *http.Server) error {
	serverErrors := make(chan error, 1)
	// Start server listening for errors
	go func() {
		fmt.Println("Listening on port", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: Listening and serving failed %s", err)

	case <-shutdown:
		log.Println("main: Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("main: Graceful shutdown did not complete in %v : %v", timeout, err)
		}

		err = server.Close()
		if err != nil {
			return fmt.Errorf("main: Couldn't stop server gracefully : %v", err)
		}
		return nil
	}
}
