// Package middleware includes all the middlewares used in the api
package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/pkg/tracking"

	"golang.org/x/time/rate"
)

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

// client holds the rate limiter for each client and the last time the client was seen.
type client struct {
	// Control how frequent the events are allowed to happen
	limiter *rate.Limiter
	// The latest time a visitor have been seen
	lastSeen time.Time
}

func init() {
	go cleanupClients()
}

// LimitRate limits the number of requests allowed per user per second.
func LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := tracking.GetUserIP(r)
		limiter := getClient(ip, 0.2, 10)

		if limiter.Allow() == false {
			w.Header().Add("Retry-After", "5 seconds")
			response.Error(w, http.StatusTooManyRequests, errors.New(http.StatusText(429)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks if the clients map exists and creates a new one if not, update last visit and return the visitor limiter.
func getClient(ip string, r rate.Limit, size int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := clients[ip]
	if !exists {
		// Implement a "token bucket" of x size, initially full and refilled at a rate of r token per second
		limiter := rate.NewLimiter(r, size)
		// Save visitor
		clients[ip] = &client{limiter, time.Now()}

		return limiter
	}

	// Update visitor last event
	v.lastSeen = time.Now()

	return v.limiter
}

// Every 30 minutes look for the clients in the map that haven't been seen for 24 hours.
func cleanupClients() {
	for {
		time.Sleep(time.Minute * 30)

		mu.Lock()
		for ip, v := range clients {
			go func(v *client, ip string) {
				if time.Since(v.lastSeen) > 24*time.Hour {
					delete(clients, ip)
				}
			}(v, ip)
		}
		mu.Unlock()
	}
}
