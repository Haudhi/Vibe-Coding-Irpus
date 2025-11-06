package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// DefaultConfig returns default Redis configuration
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}
}

// RedisCache implements cache interface using Redis
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(config *Config, logger *zap.Logger) (*RedisCache, error) {
	if config == nil {
		config = DefaultConfig()
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis successfully",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.Int("db", config.DB),
	)

	return &RedisCache{
		client: rdb,
		logger: logger,
	}, nil
}

// Set stores a value in the cache
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("Failed to marshal cache value",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	if err := r.client.Set(ctx, key, jsonValue, expiration).Err(); err != nil {
		r.logger.Error("Failed to set cache value",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	r.logger.Debug("Cache value set",
		zap.String("key", key),
		zap.Duration("expiration", expiration),
	)

	return nil
}

// Get retrieves a value from the cache
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("Cache miss", zap.String("key", key))
			return ErrCacheMiss
		}
		r.logger.Error("Failed to get cache value",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get cache value: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		r.logger.Error("Failed to unmarshal cache value",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	r.logger.Debug("Cache hit", zap.String("key", key))
	return nil
}

// Delete removes a value from the cache
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("Failed to delete cache value",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete cache value: %w", err)
	}

	r.logger.Debug("Cache value deleted", zap.String("key", key))
	return nil
}

// DeletePattern removes values matching a pattern
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.Error("Failed to get keys for pattern",
			zap.String("pattern", pattern),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get keys for pattern: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		r.logger.Error("Failed to delete cache values by pattern",
			zap.String("pattern", pattern),
			zap.Int("key_count", len(keys)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete cache values by pattern: %w", err)
	}

	r.logger.Debug("Cache values deleted by pattern",
		zap.String("pattern", pattern),
		zap.Int("deleted_count", len(keys)),
	)

	return nil
}

// Exists checks if a key exists in the cache
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to check cache key existence",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check cache key existence: %w", err)
	}

	exists := result > 0
	r.logger.Debug("Cache key existence checked",
		zap.String("key", key),
		zap.Bool("exists", exists),
	)

	return exists, nil
}

// SetTTL updates the expiration time for a key
func (r *RedisCache) SetTTL(ctx context.Context, key string, expiration time.Duration) error {
	if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
		r.logger.Error("Failed to set cache TTL",
			zap.String("key", key),
			zap.Duration("expiration", expiration),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set cache TTL: %w", err)
	}

	r.logger.Debug("Cache TTL updated",
		zap.String("key", key),
		zap.Duration("expiration", expiration),
	)

	return nil
}

// GetTTL returns the remaining time to live for a key
func (r *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to get cache TTL",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to get cache TTL: %w", err)
	}

	r.logger.Debug("Cache TTL retrieved",
		zap.String("key", key),
		zap.Duration("ttl", ttl),
	)

	return ttl, nil
}

// Increment increments a numeric value in the cache
func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to increment cache value",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to increment cache value: %w", err)
	}

	r.logger.Debug("Cache value incremented",
		zap.String("key", key),
		zap.Int64("value", result),
	)

	return result, nil
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	r.logger.Info("Closing Redis connection")
	return r.client.Close()
}

// Health checks the Redis connection health
func (r *RedisCache) Health(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	return nil
}

// GetStats returns Redis statistics
func (r *RedisCache) GetStats(ctx context.Context) (*RedisStats, error) {
	info, err := r.client.Info(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %w", err)
	}

	// Parse info to extract key statistics
	stats := &RedisStats{
		Info: info,
		// You can parse specific fields from info as needed
		// For now, just storing the raw info
	}

	return stats, nil
}

// RedisStats holds Redis statistics
type RedisStats struct {
	Info string
}

// Cache errors
var (
	ErrCacheMiss = fmt.Errorf("cache miss")
)

// Cache interface defines cache operations
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetTTL(ctx context.Context, key string, expiration time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	Increment(ctx context.Context, key string) (int64, error)
	Close() error
	Health(ctx context.Context) error
}

// Ensure RedisCache implements Cache interface
var _ Cache = (*RedisCache)(nil)