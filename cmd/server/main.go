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

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/email"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/storage"
	"github.com/GGP1/palo/pkg/updating"

	_ "github.com/lib/pq"
)

func main() {
	// Create a database connection, automigrate and
	// check tables existence
	db, close, err := storage.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	// Repos
	addingRepo := *new(adding.Repository)
	deletingRepo := *new(deleting.Repository)
	listingRepo := *new(listing.Repository)
	updatingRepo := *new(updating.Repository)
	sessionRepo := *new(handler.AuthRepository)
	emailRepo := *new(email.Repository)

	// Services
	adder := adding.NewService(addingRepo)
	deleter := deleting.NewService(deletingRepo)
	lister := listing.NewService(listingRepo)
	updater := updating.NewService(updatingRepo)
	// -- Session--
	session := handler.NewSession(sessionRepo)
	// -- Email lists --
	pendingList := email.NewPendingList(db, emailRepo)
	validatedList := email.NewValidatedList(db, emailRepo)

	// New router
	r := rest.NewRouter(db, adder, deleter, lister, updater, session, pendingList, validatedList)

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
