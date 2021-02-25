package rest_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/http/rest"
	"github.com/GGP1/adak/pkg/postgres"
	"github.com/GGP1/adak/pkg/tracking"
)

func TestHome(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		t.Fatalf("Creating config failed: %v", err)
	}

	db, err := postgres.Connect(ctx, &conf.Database)
	if err != nil {
		t.Fatalf("Database failed connecting: %v", err)
	}
	defer db.Close()

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
