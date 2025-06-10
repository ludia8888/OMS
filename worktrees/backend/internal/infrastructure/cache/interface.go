package cache

import (
	"context"
	"time"
)

// CacheService defines the interface for caching
type CacheService interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error
	
	// Set stores a value in cache with TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	
	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error
	
	// InvalidatePattern removes all keys matching a pattern
	InvalidatePattern(ctx context.Context, pattern string) error
	
	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)
	
	// Close closes the cache connection
	Close() error
}