package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/postgres"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestPostgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	c, err := config.New()
	assert.NoError(t, err)

	db, err := postgres.Connect(ctx, &c.Database)
	assert.NoError(t, err)
	defer db.Close()

	assert.NoError(t, db.Ping())
}

func TestPostgresErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := postgres.Connect(ctx, &config.DatabaseConfig{})
	assert.Error(t, err)
}
