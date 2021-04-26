package postgres_test

import (
	"context"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/postgres"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	env := []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres", "listen_addresses = '*'"}
	pool, resource := test.NewResource(t, "postgres", "13.2-alpine", env)

	err := pool.Retry(func() error {
		db, err := postgres.Connect(context.TODO(), config.Postgres{
			Username: "postgres",
			Host:     "localhost",
			Port:     resource.GetPort("5432/tcp"),
			Name:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		})
		assert.NoError(t, err)

		return db.Close()
	})
	assert.NoError(t, err)
}
