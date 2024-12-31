package main

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

var (
	//limiter    = rate.NewLimiter(1, 3) // 1 request per second with a burst of 3
	limiters   = make(map[string]*rate.Limiter)
	limitersMu sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()

	// Check if a limiter already exists for this IP
	if limiter, exists := limiters[ip]; exists {
		return limiter
	}

	// Create a new limiter for this IP
	limiter := rate.NewLimiter(1, 3)
	limiters[ip] = limiter
	return limiter
}

func rateLimit(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ip := r.RemoteAddr // Get the user's IP address
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next(w, r, ps) // Call the next handler
	}
}
