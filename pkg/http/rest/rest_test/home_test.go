package rest_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/palo/internal/config"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/storage"
	"github.com/GGP1/palo/pkg/tracking"
)

func TestHome(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		t.Fatalf("Creating config failed: %v", err)
	}

	db, close, err := storage.PostgresConnect(ctx, &conf.Database)
	if err != nil {
		t.Fatalf("Database failed connecting: %v", err)
	}
	defer close()

	trackingService := tracking.NewService(db, "")

	req := httptest.NewRequest("GET", "http://localhost:4000", nil)
	rec := httptest.NewRecorder()

	rest.Home(trackingService).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Should be status %d: got %v", http.StatusOK, res.StatusCode)
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Should read response body: %v", err)
	}
}
