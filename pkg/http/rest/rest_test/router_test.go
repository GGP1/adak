package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/http/rest"
	"github.com/GGP1/adak/pkg/memcached"
	"github.com/GGP1/adak/pkg/postgres"

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
	assert.NoError(t, err)
	defer db.Close()

	mc, err := memcached.Connect(conf.Memcached)
	assert.NoError(t, err)

	srv := httptest.NewServer(rest.NewRouter(db, mc))
	defer srv.Close()

	res, err := http.Get("http://localhost:4000/")
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
