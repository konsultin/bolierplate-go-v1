package redis

import (
	"time"
)

// Set stores a key-value pair with optional expiration.
// Pass 0 for expiration to keep the key indefinitely.
func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(c.ctx, key, value, expiration).Err()
}

// SetNX sets a key only if it doesn't exist (atomic). Returns true if set, false if key exists.
func (c *Client) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(c.ctx, key, value, expiration).Result()
}

// SetEX sets a key with mandatory expiration.
func (c *Client) SetEX(key string, value interface{}, expiration time.Duration) error {
	return c.rdb.SetEx(c.ctx, key, value, expiration).Err()
}

// Del deletes one or more keys. Returns the number of keys deleted.
func (c *Client) Del(keys ...string) (int64, error) {
	return c.rdb.Del(c.ctx, keys...).Result()
}

// Expire sets expiration on a key. Returns false if key doesn't exist.
func (c *Client) Expire(key string, expiration time.Duration) (bool, error) {
	return c.rdb.Expire(c.ctx, key, expiration).Result()
}

// Incr increments a key's integer value by 1. Returns new value.
func (c *Client) Incr(key string) (int64, error) {
	return c.rdb.Incr(c.ctx, key).Result()
}

// IncrBy increments a key's integer value by n. Returns new value.
func (c *Client) IncrBy(key string, n int64) (int64, error) {
	return c.rdb.IncrBy(c.ctx, key, n).Result()
}

// Decr decrements a key's integer value by 1. Returns new value.
func (c *Client) Decr(key string) (int64, error) {
	return c.rdb.Decr(c.ctx, key).Result()
}

// DecrBy decrements a key's integer value by n. Returns new value.
func (c *Client) DecrBy(key string, n int64) (int64, error) {
	return c.rdb.DecrBy(c.ctx, key, n).Result()
}
