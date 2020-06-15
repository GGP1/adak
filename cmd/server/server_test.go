package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"

	_ "github.com/lib/pq"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestAPI(t *testing.T) {
	t.Run("database", database)
	t.Run("server", server)
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
	db, _ := storage.Connect()
	defer db.Close()

	r := rest.NewRouter(db)

	sv := httptest.NewServer(r)

	t.Log("Given the need to test server listening.")
	{
		t.Logf("\tTest 0:\t When checking server listening on port %s", sv.URL)
		{
			res, err := http.Get(sv.URL)

			if err != nil {
				t.Fatalf("\t%s\tShould be able to listen and serve: %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to listen and serve.", succeed)

			defer res.Body.Close()
		}
	}
}
