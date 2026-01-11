package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Get retrieves a value by key. Returns redis.Nil error if key doesn't exist.
func (c *Client) Get(key string) (string, error) {
	return c.rdb.Get(c.ctx, key).Result()
}

// GetBytes retrieves a value as bytes (for binary data).
func (c *Client) GetBytes(key string) ([]byte, error) {
	return c.rdb.Get(c.ctx, key).Bytes()
}

// Exists checks if a key exists.
func (c *Client) Exists(key string) (bool, error) {
	n, err := c.rdb.Exists(c.ctx, key).Result()
	return n > 0, err
}

// TTL gets the remaining time-to-live of a key.
// Returns -1 if key has no expiration, -2 if key doesn't exist.
func (c *Client) TTL(key string) (time.Duration, error) {
	return c.rdb.TTL(c.ctx, key).Result()
}

// MGet retrieves multiple values by keys.
func (c *Client) MGet(keys ...string) ([]interface{}, error) {
	return c.rdb.MGet(c.ctx, keys...).Result()
}

// Keys returns all keys matching a pattern. Use with caution in production.
func (c *Client) Keys(pattern string) ([]string, error) {
	return c.rdb.Keys(c.ctx, pattern).Result()
}

// IsNil checks if error is redis.Nil (key not found).
func IsNil(err error) bool {
	return err == redis.Nil
}
