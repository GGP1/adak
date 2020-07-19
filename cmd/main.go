package main

import (
	"log"

	"github.com/GGP1/palo/cmd/server"
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

	configPath, err := server.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	srvConfig, err := server.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(srvConfig, router)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}
