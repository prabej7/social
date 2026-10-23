
package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	mu      sync.Mutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (l *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	counts, exists := l.clients[ip]

	if exists && counts >= l.limit {
		return false, l.window
	}

	l.clients[ip] = counts + 1

	if !exists {
		go l.resetCount(ip)
	}

	return true, 0
}

func (l *FixedWindowLimiter) resetCount(ip string) {
	time.Sleep(l.window)

	l.mu.Lock()
	delete(l.clients, ip)
	l.mu.Unlock()
}
