package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"github.com/openfoundry/oms/internal/domain/repository"
	"go.uber.org/zap"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(config RedisConfig) (*RedisCache, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	addr, password, db, ttl, logger := config.Addr, config.Password, config.DB, config.TTL, config.Logger
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		logger: logger,
		ttl:    ttl,
	}, nil
}

// Set stores a value in the cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if ttl == 0 {
		ttl = c.ttl
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		c.logger.Error("Failed to set cache value", 
			zap.String("key", key),
			zap.Duration("ttl", ttl),
			zap.Error(err))
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	return nil
}

// Get retrieves a value from the cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return repository.ErrCacheMiss
		}
		c.logger.Error("Failed to get cache value", 
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to get cache value: %w", err)
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Error("Failed to delete cache value", 
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to delete cache value: %w", err)
	}

	return nil
}

// Invalidate removes multiple keys from the cache
func (c *RedisCache) Invalidate(ctx context.Context, pattern string) error {
	// Use SCAN to find all matching keys
	var cursor uint64
	var keys []string

	for {
		var batch []string
		var err error
		batch, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			c.logger.Error("Failed to scan keys", 
				zap.String("pattern", pattern),
				zap.Error(err))
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		// Delete in batches to avoid overloading Redis
		const batchSize = 1000
		for i := 0; i < len(keys); i += batchSize {
			end := i + batchSize
			if end > len(keys) {
				end = len(keys)
			}

			if err := c.client.Del(ctx, keys[i:end]...).Err(); err != nil {
				c.logger.Error("Failed to delete keys batch", 
					zap.Int("batch_start", i),
					zap.Int("batch_end", end),
					zap.Error(err))
				return fmt.Errorf("failed to delete keys batch: %w", err)
			}
		}
	}

	c.logger.Info("Invalidated cache keys", 
		zap.String("pattern", pattern),
		zap.Int("count", len(keys)))

	return nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// ObjectTypeCache provides caching for object types
type ObjectTypeCache struct {
	cache  *RedisCache
	prefix string
}

// NewObjectTypeCache creates a new object type cache
func NewObjectTypeCache(cache *RedisCache) *ObjectTypeCache {
	return &ObjectTypeCache{
		cache:  cache,
		prefix: "object_type:",
	}
}

// GetByID retrieves an object type by ID from cache
func (c *ObjectTypeCache) GetByID(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error) {
	key := fmt.Sprintf("%sid:%s", c.prefix, id.String())
	var objectType entity.ObjectType
	err := c.cache.Get(ctx, key, &objectType)
	if err != nil {
		return nil, err
	}
	return &objectType, nil
}

// GetByName retrieves an object type by name from cache
func (c *ObjectTypeCache) GetByName(ctx context.Context, name string) (*entity.ObjectType, error) {
	key := fmt.Sprintf("%sname:%s", c.prefix, name)
	var objectType entity.ObjectType
	err := c.cache.Get(ctx, key, &objectType)
	if err != nil {
		return nil, err
	}
	return &objectType, nil
}

// Set stores an object type in cache with multiple keys
func (c *ObjectTypeCache) Set(ctx context.Context, objectType *entity.ObjectType) error {
	// Cache by ID
	idKey := fmt.Sprintf("%sid:%s", c.prefix, objectType.ID.String())
	if err := c.cache.Set(ctx, idKey, objectType, 0); err != nil {
		return err
	}

	// Cache by name
	nameKey := fmt.Sprintf("%sname:%s", c.prefix, objectType.Name)
	if err := c.cache.Set(ctx, nameKey, objectType, 0); err != nil {
		// Rollback ID cache on error
		_ = c.cache.Delete(ctx, idKey)
		return err
	}

	return nil
}

// Delete removes an object type from cache
func (c *ObjectTypeCache) Delete(ctx context.Context, objectType *entity.ObjectType) error {
	// Delete by ID
	idKey := fmt.Sprintf("%sid:%s", c.prefix, objectType.ID.String())
	if err := c.cache.Delete(ctx, idKey); err != nil {
		return err
	}

	// Delete by name
	nameKey := fmt.Sprintf("%sname:%s", c.prefix, objectType.Name)
	if err := c.cache.Delete(ctx, nameKey); err != nil {
		return err
	}

	// Invalidate any list caches
	return c.cache.Invalidate(ctx, c.prefix+"list:*")
}

// InvalidateAll removes all object types from cache
func (c *ObjectTypeCache) InvalidateAll(ctx context.Context) error {
	return c.cache.Invalidate(ctx, c.prefix+"*")
}

// LinkTypeCache provides caching for link types
type LinkTypeCache struct {
	cache  *RedisCache
	prefix string
}

// NewLinkTypeCache creates a new link type cache
func NewLinkTypeCache(cache *RedisCache) *LinkTypeCache {
	return &LinkTypeCache{
		cache:  cache,
		prefix: "link_type:",
	}
}

// GetByID retrieves a link type by ID from cache
func (c *LinkTypeCache) GetByID(ctx context.Context, id uuid.UUID) (*entity.LinkType, error) {
	key := fmt.Sprintf("%sid:%s", c.prefix, id.String())
	var linkType entity.LinkType
	err := c.cache.Get(ctx, key, &linkType)
	if err != nil {
		return nil, err
	}
	return &linkType, nil
}

// GetByName retrieves a link type by name from cache
func (c *LinkTypeCache) GetByName(ctx context.Context, name string) (*entity.LinkType, error) {
	key := fmt.Sprintf("%sname:%s", c.prefix, name)
	var linkType entity.LinkType
	err := c.cache.Get(ctx, key, &linkType)
	if err != nil {
		return nil, err
	}
	return &linkType, nil
}

// Set stores a link type in cache
func (c *LinkTypeCache) Set(ctx context.Context, linkType *entity.LinkType) error {
	// Cache by ID
	idKey := fmt.Sprintf("%sid:%s", c.prefix, linkType.ID.String())
	if err := c.cache.Set(ctx, idKey, linkType, 0); err != nil {
		return err
	}

	// Cache by name
	nameKey := fmt.Sprintf("%sname:%s", c.prefix, linkType.Name)
	if err := c.cache.Set(ctx, nameKey, linkType, 0); err != nil {
		// Rollback ID cache on error
		_ = c.cache.Delete(ctx, idKey)
		return err
	}

	return nil
}

// Delete removes a link type from cache
func (c *LinkTypeCache) Delete(ctx context.Context, linkType *entity.LinkType) error {
	// Delete by ID
	idKey := fmt.Sprintf("%sid:%s", c.prefix, linkType.ID.String())
	if err := c.cache.Delete(ctx, idKey); err != nil {
		return err
	}

	// Delete by name
	nameKey := fmt.Sprintf("%sname:%s", c.prefix, linkType.Name)
	if err := c.cache.Delete(ctx, nameKey); err != nil {
		return err
	}

	// Invalidate any list caches
	return c.cache.Invalidate(ctx, c.prefix+"list:*")
}

// InvalidateAll removes all link types from cache
func (c *LinkTypeCache) InvalidateAll(ctx context.Context) error {
	return c.cache.Invalidate(ctx, c.prefix+"*")
}