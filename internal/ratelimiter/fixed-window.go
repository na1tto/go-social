package rateLimiter

import (
	"sync"
	"time"
)

type clientData struct {
	count   int
	resetAt time.Time
}

type FixedWindowLimiter struct {
	sync.Mutex
	clients map[string]clientData
	limit   int
	window  time.Duration
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		clients: make(map[string]clientData),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	now := time.Now()

	rl.Lock()
	defer rl.Unlock()

	client, exists := rl.clients[ip]

	if !exists || now.After(client.resetAt) {
		rl.clients[ip] = clientData{
			count:   1,
			resetAt: now.Add(rl.window),
		}

		return true, 0
	}

	if client.count >= rl.limit {
		retryAfter := time.Until(client.resetAt)
		return false, retryAfter
	}

	client.count++
	rl.clients[ip] = client

	return true, 0
}
