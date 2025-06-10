package cache

import (
	"errors"
	"time"

	"go.uber.org/zap"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	TTL      time.Duration
	Logger   *zap.Logger
}

// Validate validates the Redis configuration
func (c RedisConfig) Validate() error {
	if c.Addr == "" {
		return errors.New("redis address is required")
	}
	if c.TTL <= 0 {
		return errors.New("TTL must be positive")
	}
	if c.Logger == nil {
		return errors.New("logger is required")
	}
	return nil
}