package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestAPI(t *testing.T) {
	t.Run("database", database)
	// Can't test because the firewall security alert popps everytime the server runs, commented until fixed
	// t.Run("server", server)
}

func database(t *testing.T) {
	t.Log("Given the need to test database connection.")
	{
		t.Logf("\tTest 0:\tWhen checking the database connection")
		{
			db, err := storage.Connect()
			if err != nil {
				t.Fatalf("\t%s\tShould be able to connect to the database : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to connect to the database.", succeed)

			defer db.Close()
		}
	}
}

func server(t *testing.T) {
	sv := &http.Server{
		Addr:           ":4000",
		Handler:        nil,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	t.Log("Given the need to test server listening.")
	{
		t.Logf("\tTest 0:\t When checking server listening on port %s", sv.Addr)
		{
			err := sv.ListenAndServe()
			if err != nil {
				t.Fatalf("\t%s\tShould be able to listen and serve: %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to listen and serve.", succeed)
		}
	}
}
