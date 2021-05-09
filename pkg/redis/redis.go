package redis

import (
	"context"
	"net"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/logger"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Connect establishes a connection with the redis client.
func Connect(ctx context.Context, config config.Redis) (*redis.Client, error) {
	addr := net.JoinHostPort(config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     addr,
		Password: config.Password,
		DB:       0,
	}).WithContext(ctx)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping error")
	}

	logger.Infof("Connected to redis on %s", addr)
	return rdb, nil
}
