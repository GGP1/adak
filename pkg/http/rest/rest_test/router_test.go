package rest_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/palo/internal/config"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"
)

// Checkboxes
const (
	succeed = "\u2713"
	failed  = "\u2717"
)

// Fix
func TestRouting(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, close, err := storage.PostgresConnect(ctx, &conf.Database)
	if err != nil {
		t.Fatal("Database failed connecting")
	}
	defer close()

	srv := httptest.NewServer(rest.NewRouter(db))
	defer srv.Close()

	t.Log("Given the need to test the router.")
	{
		t.Logf("\tTest 0: When checking GET request.")
		{
			res, err := http.Get("http://localhost:4000/")
			if err != nil {
				t.Errorf("\t%s\tShould return a response: %v", failed, err)
			}
			t.Logf("\t%s\tShould return a response.", succeed)

			if res.StatusCode != http.StatusOK {
				t.Errorf("\t%s\tShould be status OK: got %v", failed, res.StatusCode)
			}
			t.Logf("\t%s\tShould be status OK.", succeed)
		}
	}

}
