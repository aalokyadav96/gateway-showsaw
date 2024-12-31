package main

import (
	"time"

	"github.com/sony/gobreaker"
)

var circuitBreakers = make(map[string]*gobreaker.CircuitBreaker)

func initCircuitBreaker(serviceName string) *gobreaker.CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	})
	circuitBreakers[serviceName] = cb
	return cb
}
