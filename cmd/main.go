package main

import (
	"context"
	"log"

	"github.com/GGP1/palo/cmd/server"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	db, close, err := storage.PostgresConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	esClient, err := storage.Elasticsearch(ctx)
	if err != nil {
		log.Fatal(err)
	}

	router := rest.NewRouter(ctx, db, esClient)

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
