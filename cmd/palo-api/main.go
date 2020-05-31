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
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	storage.Connect()
	defer storage.DB.Close()

	// Router setup
	r := rest.SetupRouter()

	// Server setup
	server := &http.Server{
		Addr:           ":4000",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

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
		log.Fatalf("error: Listening and serving %s", err)

	case <-shutdown:
		log.Println("main: Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = server.Close()
		}

		if err != nil {
			log.Fatalf("main : Could not stop server gracefully : %v", err)
		}
	}
}
