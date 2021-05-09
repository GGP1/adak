package memcached

import (
	"net"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	pool, resource := test.NewResource(t, "memcached", "1.6.9-alpine", nil)

	err := pool.Retry(func() error {
		_, err := Connect(config.Memcached{
			Servers: []string{net.JoinHostPort("localhost", resource.GetPort("11211/tcp"))},
		})
		return err
	})
	assert.NoError(t, err)
}
