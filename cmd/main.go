package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/GGP1/palo/cmd/server"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, close, err := storage.PostgresConnect(ctx)
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

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
