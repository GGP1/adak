package main

import (
	"context"
	"log"

	"github.com/GGP1/palo/cmd/server"
	"github.com/GGP1/palo/internal/config"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, close, err := storage.PostgresConnect(ctx, &conf.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	router := rest.NewRouter(db)

	srv := server.New(conf, router)

	if err := srv.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
