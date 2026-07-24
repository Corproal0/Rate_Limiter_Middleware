package limiter

import (
	"sync"
	"time"
)

type client struct {
	tokens   int
	lastSeen time.Time
	mu       sync.Mutex
	stopChan chan struct{}
}

func newClient(maxTokens int, rate time.Duration) *client {
	c := &client{
		tokens:   maxTokens,
		lastSeen: time.Now(),
		stopChan: make(chan struct{}),
	}
	go c.startRefill(rate, maxTokens)
	return c
}

func (c *client) startRefill(rate time.Duration, maxTokens int) {
	ticker := time.NewTicker(rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			if c.tokens < maxTokens {
				c.tokens++
			}
			c.mu.Unlock()
		case <-c.stopChan:
			return
		}
	}
}

func (c *client) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastSeen = time.Now()
	if c.tokens > 0 {
		c.tokens--
		return true
	}
	return false
}

func (c *client) stop() {
	close(c.stopChan)
}
