// Package test contains testing helpers.
package test

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/GGP1/adak/internal/crypt"
	"github.com/GGP1/adak/pkg/postgres"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	// Used to open a connection with postgres database
	_ "github.com/lib/pq"
)

// AddCookie encrypts and adds a cookie to the request passed.
func AddCookie(t testing.TB, r *http.Request, name, value string) {
	t.Helper()
	if viper.Get("token.secretKey") == "" {
		viper.Set("token.secretKey", "1")
	}

	c, err := crypt.Encrypt([]byte(value))
	assert.NoError(t, err)

	r.AddCookie(&http.Cookie{
		Name:  name,
		Value: hex.EncodeToString(c),
	})
}

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

// StartMemcached starts a memcached container and makes the cleanup.
func StartMemcached(t testing.TB) *memcache.Client {
	pool, resource, mc, err := RunMemcached()
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, pool.Purge(resource), "Couldn't free resources")
	})

	return mc
}

// RunMemcached initializes a docker container with memcached running in it.
func RunMemcached() (*dockertest.Pool, *dockertest.Resource, *memcache.Client, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, err
	}

	resource, err := pool.Run("memcached", "1.6.9-alpine", nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var mc *memcache.Client
	err = pool.Retry(func() error {
		mc = memcache.New(fmt.Sprintf("localhost:%s", resource.GetPort("11211/tcp")))
		return mc.Ping()
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return pool, resource, mc, nil
}

// StartPostgres starts a postgres container and makes the cleanup.
func StartPostgres(t testing.TB) *sqlx.DB {
	pool, resource, db, err := RunPostgres()
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, db.Close(), "Couldn't close the connection with postgres")
		assert.NoError(t, pool.Purge(resource), "Couldn't free resources")
	})

	return db
}

// RunPostgres initializes a docker container with postgres running in it.
func RunPostgres() (*dockertest.Pool, *dockertest.Resource, *sqlx.DB, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, err
	}
	// The database name will be taken from the user name
	env := []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres", "listen_addresses = '*'"}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.2-alpine",
		Env:        env,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, nil, nil, err
	}

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
	if err != nil {
		return nil, nil, nil, err
	}

	if err := postgres.CreateTables(context.Background(), db); err != nil {
		return nil, nil, nil, err
	}

	return pool, resource, db, nil
}

// RunRedis initializes a docker container with redis running in it.
func RunRedis() (*dockertest.Pool, *dockertest.Resource, *redis.Client, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, err
	}

	resource, err := pool.Run("redis", "6.2.1-alpine", nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var rdb *redis.Client
	err = pool.Retry(func() error {
		rdb = redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    net.JoinHostPort("localhost", resource.GetPort("6379/tcp")),
		})
		return rdb.Ping(rdb.Context()).Err()
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return pool, resource, rdb, nil
}

// StartRedis starts a redis container and makes the cleanup.
func StartRedis(t testing.TB) *redis.Client {
	pool, resource, rdb, err := RunRedis()
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, rdb.Close(), "Couldn't close connection with redis")
		assert.NoError(t, pool.Purge(resource), "Couldn't free resources")
	})

	return rdb
}
