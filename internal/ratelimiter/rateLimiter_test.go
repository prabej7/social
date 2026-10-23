package ratelimiter

import (
	"testing"
	"time"
)

func TestFixedWindowLimiter(t *testing.T) {
	limiter := NewFixedWindowLimiter(10, time.Minute)
	ip := "192.168.1.100"

	// First 10 requests should be allowed.
	for i := 1; i <= 10; i++ {
		allowed, _ := limiter.Allow(ip)
		if !allowed {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	// 11th request should be rejected.
	allowed, retryAfter := limiter.Allow(ip)
	if allowed {
		t.Fatal("request 11 should have been rejected")
	}

	t.Logf("Rate limit correctly enforced; retry after: %v", retryAfter)
}
