package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Client wraps redis client
type Client struct {
	client *redis.Client
	ctx    context.Context
}

// NewClient creates a new Redis client
func NewClient(url string, db int) (*Client, error) {
	var redisClient *redis.Client

	opt, err := redis.ParseURL(url)
	if err != nil {
		// fallback: treat as plain host:port
		redisClient = redis.NewClient(&redis.Options{
			Addr: url,
			DB:   db,
		})
	} else {
		redisClient = redis.NewClient(opt)
	}
	ctx := context.Background()
	// Test connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"db": db,
	}).Info("Successfully connected to Redis")

	return &Client{
		client: redisClient,
		ctx:    ctx,
	}, nil
}

// GetClient returns the underlying redis client
func (c *Client) GetClient() *redis.Client {
	return c.client
}

// GetContext returns the context
func (c *Client) GetContext() context.Context {
	return c.ctx
}

// Set sets a key-value pair with expiration
func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (c *Client) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

// Exists checks if a key exists
func (c *Client) Exists(key string) (bool, error) {
	result, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Delete deletes a key
func (c *Client) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Keys returns all keys matching pattern
func (c *Client) Keys(pattern string) ([]string, error) {
	return c.client.Keys(c.ctx, pattern).Result()
}

// Expire sets expiration on a key
func (c *Client) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(c.ctx, key, expiration).Err()
}

// TTL returns the time to live of a key
func (c *Client) TTL(key string) (time.Duration, error) {
	return c.client.TTL(c.ctx, key).Result()
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.client.Close()
}

// Ping pings the Redis server
func (c *Client) Ping() error {
	return c.client.Ping(c.ctx).Err()
}
