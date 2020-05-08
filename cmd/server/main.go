package main

import (
	"log"

	"palo/pkg/http/rest"
	db "palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	db.Connect()
	defer db.DB.Close()

	// Setup router
	r := rest.SetupRouter()

	// Run server
	log.Fatal(r.Run(":4000"))
}
