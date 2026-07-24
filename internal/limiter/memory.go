package limiter

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.Mutex
	clients map[string]*client
	cfg     Config
	stopCh  chan struct{}
}

func NewMemoryStore(cfg Config) *MemoryStore {
	ms := &MemoryStore{
		clients: make(map[string]*client),
		cfg:     cfg,
		stopCh:  make(chan struct{}),
	}
	go ms.startCleanup()
	return ms
}

func (ms *MemoryStore) Allow(_ context.Context, key string) (bool, error) {
	ms.mu.Lock()
	c, exists := ms.clients[key]
	if !exists {
		c = newClient(ms.cfg.MaxTokens, ms.cfg.Rate)
		ms.clients[key] = c
	}
	ms.mu.Unlock()
	return c.Allow(), nil
}

func (ms *MemoryStore) startCleanup() {
	ticker := time.NewTicker(ms.cfg.CleanupTTL / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.cleanExpiredClients()
		case <-ms.stopCh:
			return
		}
	}
}
func (ms *MemoryStore) cleanExpiredClients() {
	ms.mu.Lock()
	defer ms.mu.Unlock() // Гарантированно отпустим мьютекс при выходе из функции

	now := time.Now()
	for key, c := range ms.clients {
		if now.Sub(c.lastSeen) > ms.cfg.CleanupTTL {
			c.stop()
			delete(ms.clients, key)
		}
	}
}

func (ms *MemoryStore) Close() error {
	close(ms.stopCh)
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for _, c := range ms.clients {
		c.stop()
	}
	return nil
}
