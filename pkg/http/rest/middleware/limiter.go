package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/pkg/tracking"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/pkg/errors"
)

// RateLimiter uses a leaky bucket algorithm for limiting the requests to the API from the same host.
type RateLimiter struct {
	limiter *redis_rate.Limiter
	rate    int
}

// NewRateLimiter returns a rate limiter with the configuration values passed.
func NewRateLimiter(config config.RateLimiter, rdb *redis.Client) *RateLimiter {
	rl := &RateLimiter{
		limiter: redis_rate.NewLimiter(rdb),
		rate:    config.Rate,
	}

	return rl
}

// Limit make sure no one abuses the API by using token bucket algorithm.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := tracking.GetUserIP(r)
		if ip == "" {
			// Try hard to avoid this as an attacker able to hide ips from
			// headers will be able to perform DDOS.
			next.ServeHTTP(w, r)
			return
		}

		res, err := rl.limiter.Allow(r.Context(), ip, redis_rate.PerMinute(rl.rate))
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			w.Header().Add("RateLimit-Limit", strconv.Itoa(res.Limit.Burst))
			w.Header().Add("RateLimit-Reset", strconv.Itoa(int(res.ResetAfter/time.Second)))
			w.Header().Add("Retry-After", strconv.Itoa(int(res.RetryAfter/time.Second)))
			response.Error(w, http.StatusTooManyRequests, errors.New("Too Many Requests"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
