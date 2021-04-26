package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/http/rest"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

// Fix
func TestRouter(t *testing.T) {
	conf, err := config.New()
	assert.NoError(t, err)

	db := test.StartPostgres(t)
	mc := test.StartMemcached(t)

	srv := httptest.NewServer(rest.NewRouter(conf, db, mc))
	defer srv.Close()

	res, err := http.Get("http://localhost:4000/")
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
