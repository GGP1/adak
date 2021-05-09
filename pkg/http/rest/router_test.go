package rest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/http/rest"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	mux := rest.NewRouter(config.Config{}, nil, nil, nil)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	assert.NoError(t, err)

	buf, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "Welcome to the Adak home page\n", string(buf))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	h := res.Header
	assert.Equal(t, h.Get("Access-Control-Allow-Credentials"), "true")
	assert.Equal(t, h.Get("Access-Control-Allow-Headers"),
		"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, accept, origin, Cache-Control, X-Requested-With")
	assert.Equal(t, h.Get("Access-Control-Allow-Methods"), "POST, OPTIONS, GET, PUT, DELETE, HEAD")
	assert.Equal(t, h.Get("Access-Control-Expose-Headers"), "UID, SID, CID")
	assert.Equal(t, h.Get("Content-Security-Policy"), "default-src 'self';")
	assert.Equal(t, h.Get("Content-Type"), "text/html; charset=UTF-8")
	assert.Equal(t, h.Get("Feature-Policy"), "microphone 'none'; camera 'none'")
	assert.Equal(t, h.Get("Referrer-Policy"), "no-referrer")
	assert.Equal(t, h.Get("Strict-Transport-Security"), "max-age=31536000; includeSubDomains; preload")
	assert.Equal(t, h.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, h.Get("X-Frame-Options"), "SAMEORIGIN")
	assert.Equal(t, h.Get("X-Permitted-Cross-Domain-Policies"), "none")
	assert.Equal(t, h.Get("X-Xss-Protection"), "1; mode=block")
}
