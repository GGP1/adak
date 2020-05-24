package main

import (
	"fmt"
	"net/http"
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

	// Run server
	fmt.Println("Listening on port", server.Addr)
	server.ListenAndServe()
}
