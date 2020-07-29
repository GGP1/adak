/*
Package middleware includes all the middlewares used in the api
*/
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor holds the rate limiter for each visitor and the
// last time the visitor was seen.
type visitor struct {
	// Control how frequent the events are allowed to happen
	limiter *rate.Limiter
	// The latest time a visitor have been seen
	lastSeen time.Time
}

// visitors map holds the values of the visitors
var visitors = make(map[string]*visitor)
var mutex sync.RWMutex

func init() {
	go cleanupVisitors()
}

// LimitRate limits the number of requests allowed
// per user per second.
func LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip, 1, 3)
		// Control how frequent events may happen
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks if the visitors map exists and creates a new one
// if not, update last visit and return the visitor limiter.
func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mutex.RLock()
	defer mutex.RUnlock()

	v, exists := visitors[ip]
	if !exists {
		// Implement a "token bucket" of size b, initially
		// full and refilled at a rate of r token per second
		limiter := rate.NewLimiter(r, b)
		// Save visitor
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update visitor last event
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute look for the visitors in the map that
// haven't been seen for 10 minutes.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mutex.RLock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(visitors, ip)
			}
		}
		mutex.RUnlock()
	}
}
