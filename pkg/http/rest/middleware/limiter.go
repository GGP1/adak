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
	mu       sync.Mutex
	visitors = make(map[string]*visitor)
)

// visitor holds the rate limiter for each visitor and the last time the visitor was seen.
type visitor struct {
	// Control how frequent the events are allowed to happen
	limiter *rate.Limiter
	// The latest time a visitor have been seen
	lastSeen time.Time
}

func init() {
	go cleanupVisitors()
}

// LimitRate limits the number of requests allowed per user per second.
func LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := tracking.GetUserIP(r)
		limiter := getVisitor(ip, 0.5, 3)

		if limiter.Allow() == false {
			response.Error(w, http.StatusTooManyRequests, errors.New(http.StatusText(429)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks if the visitors map exists and creates a new one if not, update last visit and return the visitor limiter.
func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// Implement a "token bucket" of size b, initially full and refilled at a rate of r token per second
		limiter := rate.NewLimiter(r, b)
		// Save visitor
		visitors[ip] = &visitor{limiter, time.Now()}

		return limiter
	}

	// Update visitor last event
	v.lastSeen = time.Now()

	return v.limiter
}

// Every 30 minutes look for the visitors in the map that haven't been seen for 24 hours.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute * 30)

		mu.Lock()
		for ip, v := range visitors {
			go func(v *visitor, ip string) {
				if time.Since(v.lastSeen) > 24*time.Hour {
					delete(visitors, ip)
				}
			}(v, ip)
		}
		mu.Unlock()
	}
}
