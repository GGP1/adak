package main

import (
	"log"

	"github.com/GGP1/palo/cmd/server"
	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	db, close, err := storage.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	router := rest.NewRouter(db)

	srv := server.New(cfg.SrvPort, router)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}
