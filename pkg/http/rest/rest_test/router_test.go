package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/http/rest"
	"github.com/GGP1/adak/pkg/postgres"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

// Fix
func TestRouter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.New()
	assert.NoError(t, err)

	db, err := postgres.Connect(ctx, &conf.Database)
	if err != nil {
		t.Fatal("Database failed connecting")
	}
	defer db.Close()

	cache, err := lru.New(conf.Cache.Size)
	if err != nil {
		t.Fatalf("couldn't create the cache: %v", err)
	}

	srv := httptest.NewServer(rest.NewRouter(db, cache))
	defer srv.Close()

	res, err := http.Get("http://localhost:4000/")
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
