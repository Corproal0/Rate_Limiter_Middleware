package limiter

import (
	"context"
	"time"
)

type Store interface {
	Allow(ctx context.Context, key string) (bool, error)
	Close() error
}

type Config struct {
	Rate       time.Duration
	MaxTokens  int
	CleanupTTL time.Duration
}
