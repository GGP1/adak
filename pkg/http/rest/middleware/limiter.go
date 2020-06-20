/*
Package middleware includes all the middlewares used in the api
*/
package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor holds the rate limiter for each visitor
// and the last time the visitor was seen
type visitor struct {
	limiter   *rate.Limiter
	lastVisit time.Time
}

// visitors map holds the values of the visitors
var visitors = make(map[string]*visitor)
var mutex sync.Mutex

func init() {
	go cleanupVisitors()
}

// LimitRate limits the number of requests allowed
// per user per second
func LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks if the visitors map exists and creates a new one
// if not, update last visit and return the visitor limiter
func getVisitor(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// Implement a 3 tokens bucket, initially full
		// and refilled at a rate of 1 token per second
		limiter := rate.NewLimiter(1, 3)
		// Save visitor
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update visitor last visit
	v.lastVisit = time.Now()
	return v.limiter
}

// Every minute check the visitors in the map that
// haven't been seen for 15 minutes
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mutex.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastVisit) > 15*time.Minute {
				delete(visitors, ip)
			}
		}
	}
}
