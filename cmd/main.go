package main

import (
	"log"
	"os"

	"github.com/GGP1/palo/cmd/server"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	db, close, err := storage.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	router := rest.NewRouter(db)

	godotenv.Load("../.env")
	port := os.Getenv("SRV_PORT")

	srv := server.New(port, router)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}
