// Package test contains testing helpers.
package test

import (
	"fmt"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

// NewResource returns a new pool, a docker container and handles its purge.
func NewResource(t testing.TB, repository, tags string, env []string) (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	resource, err := pool.Run(repository, tags, env)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, pool.Purge(resource), "Couldn't free resources")
	})

	return pool, resource
}

// StartMemcached initializes a docker container with memcached running in it.
func StartMemcached(t testing.TB) *memcache.Client {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	resource, err := pool.Run("memcached", "1.6.9-alpine", nil)
	assert.NoError(t, err)

	var mc *memcache.Client
	err = pool.Retry(func() error {
		mc = memcache.New(fmt.Sprintf("localhost:%s", resource.GetPort("11211/tcp")))
		return mc.Ping()
	})
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, pool.Purge(resource), "Couldn't free resources")
	})

	return mc
}

// StartPostgres initializes a docker container with postgres running in it.
func StartPostgres(t testing.TB) *sqlx.DB {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	// The database name will be taken from the user name
	env := []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres", "listen_addresses = '*'"}
	resource, err := pool.Run("postgres", "13.2-alpine", env)
	assert.NoError(t, err)

	var db *sqlx.DB
	err = pool.Retry(func() error {
		url := fmt.Sprintf("host=localhost port=%s user=postgres password=postgres dbname=postgres sslmode=disable",
			resource.GetPort("5432/tcp"))
		db, err = sqlx.Connect("postgres", url)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, db.Close(), "Couldn't close the database")
		assert.NoError(t, pool.Purge(resource), "Couldn't free resources")
	})

	return db
}
