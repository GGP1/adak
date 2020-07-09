package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/storage"
	"github.com/GGP1/palo/pkg/updating"

	_ "github.com/lib/pq"
)

func main() {
	// Create a database connection, automigrate and
	// check tables existence
	_, close, err := storage.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	// Repos
	addingRepo := *new(adding.Repository)
	deletingRepo := *new(deleting.Repository)
	listingRepo := *new(listing.Repository)
	updatingRepo := *new(updating.Repository)

	// Services
	adder := adding.NewService(addingRepo)
	deleter := deleting.NewService(deletingRepo)
	lister := listing.NewService(listingRepo)
	updater := updating.NewService(updatingRepo)

	// New router
	r := rest.NewRouter(adder, deleter, lister, updater)

	// Server setup
	server := &http.Server{
		Addr:           ":4000",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = StartServer(server)
	if err != nil {
		log.Fatal(err)
	}
}
