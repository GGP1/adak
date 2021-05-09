package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/GGP1/adak/cmd/server"
	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/http/rest"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	c := config.Config{
		Server: config.Server{
			Host: "localhost",
			Port: "61111",
		},
	}
	srv := server.New(c, rest.NewRouter(c, nil, nil, nil))
	ctx := context.Background()

	go func() {
		// Wait for the server to start
		time.Sleep(50 * time.Millisecond)

		res, err := http.Get("http://localhost:61111/")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		err = srv.Shutdown(ctx)
		assert.NoError(t, err)
		err = srv.Close()
		assert.NoError(t, err)
	}()

	err := srv.Start(ctx)
	assert.Error(t, err) // http: Server closed
}
