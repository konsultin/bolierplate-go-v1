package httpk

import (
	"math"
	"math/rand"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retry attempts
	InitialBackoff time.Duration // Initial backoff duration
	MaxBackoff     time.Duration // Maximum backoff duration
	Multiplier     float64       // Backoff multiplier
	Jitter         bool          // Add randomization to backoff
	RetryStatuses  []int         // HTTP status codes that should trigger retry
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		Multiplier:     2.0,
		Jitter:         true,
		RetryStatuses:  []int{429, 500, 502, 503, 504},
	}
}

// Backoff calculates the backoff duration for the given attempt
func (r *RetryConfig) Backoff(attempt int) time.Duration {
	// Calculate exponential backoff
	backoff := float64(r.InitialBackoff) * math.Pow(r.Multiplier, float64(attempt))

	// Cap at max backoff
	if backoff > float64(r.MaxBackoff) {
		backoff = float64(r.MaxBackoff)
	}

	// Add jitter if enabled (randomize between 0.5x and 1.5x of calculated backoff)
	if r.Jitter {
		jitter := 0.5 + rand.Float64() // Random value between 0.5 and 1.5
		backoff *= jitter
	}

	return time.Duration(backoff)
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(maxRetries int) func(*RetryConfig) {
	return func(cfg *RetryConfig) {
		cfg.MaxRetries = maxRetries
	}
}

// WithInitialBackoff sets the initial backoff duration
func WithInitialBackoff(duration time.Duration) func(*RetryConfig) {
	return func(cfg *RetryConfig) {
		cfg.InitialBackoff = duration
	}
}

// WithMaxBackoff sets the maximum backoff duration
func WithMaxBackoff(duration time.Duration) func(*RetryConfig) {
	return func(cfg *RetryConfig) {
		cfg.MaxBackoff = duration
	}
}

// WithMultiplier sets the backoff multiplier
func WithMultiplier(multiplier float64) func(*RetryConfig) {
	return func(cfg *RetryConfig) {
		cfg.Multiplier = multiplier
	}
}

// WithJitter enables/disables jitter
func WithJitter(enabled bool) func(*RetryConfig) {
	return func(cfg *RetryConfig) {
		cfg.Jitter = enabled
	}
}

// WithRetryStatuses sets the HTTP status codes that should trigger retry
func WithRetryStatuses(statuses ...int) func(*RetryConfig) {
	return func(cfg *RetryConfig) {
		cfg.RetryStatuses = statuses
	}
}

// NewRetryConfig creates a new retry configuration with options
func NewRetryConfig(opts ...func(*RetryConfig)) *RetryConfig {
	cfg := DefaultRetryConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
