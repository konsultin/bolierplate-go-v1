package httpk

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	// StateClosed means circuit is closed and requests pass through
	StateClosed CircuitState = iota
	// StateOpen means circuit is open and requests are rejected
	StateOpen
	// StateHalfOpen means circuit is testing if service has recovered
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold int           // Number of failures before opening circuit
	SuccessThreshold int           // Number of successes in half-open before closing
	Timeout          time.Duration // Time to wait before transitioning from open to half-open
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config CircuitBreakerConfig
	logger Logger

	mu              sync.RWMutex
	state           CircuitState
	failures        int
	successes       int
	lastFailureTime time.Time
	lastStateChange time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg CircuitBreakerConfig, logger Logger) *CircuitBreaker {
	return &CircuitBreaker{
		config:          cfg,
		logger:          logger,
		state:           StateClosed,
		failures:        0,
		successes:       0,
		lastStateChange: time.Now(),
	}
}

// Execute runs the given function through the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() (*Response, error)) (*Response, error) {
	// Check circuit state
	if err := cb.beforeRequest(); err != nil {
		return nil, err
	}

	// Execute request
	resp, err := fn()

	// Update circuit state based on result
	cb.afterRequest(err)

	return resp, err
}

// beforeRequest checks if request can proceed
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		// Check if timeout has expired
		if time.Since(cb.lastFailureTime) >= cb.config.Timeout {
			// Transition to half-open
			cb.setState(StateHalfOpen)
			return nil
		}
		return fmt.Errorf("circuit breaker is OPEN")

	case StateHalfOpen, StateClosed:
		return nil

	default:
		return fmt.Errorf("unknown circuit breaker state")
	}
}

// afterRequest updates circuit state after request completion
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

// onFailure handles request failure
func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.FailureThreshold {
			cb.setState(StateOpen)
		}

	case StateHalfOpen:
		// Single failure in half-open state reopens the circuit
		cb.setState(StateOpen)
	}
}

// onSuccess handles request success
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		// Reset failure count on success
		cb.failures = 0

	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.setState(StateClosed)
		}
	}
}

// setState transitions to a new state
func (cb *CircuitBreaker) setState(newState CircuitState) {
	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()

	// Reset counters on state change
	if newState == StateClosed {
		cb.failures = 0
		cb.successes = 0
	} else if newState == StateHalfOpen {
		cb.successes = 0
	}

	// Log state change
	if cb.logger != nil {
		cb.logger.Infof("Circuit breaker state changed: %s -> %s", oldState, newState)
	}
}

// State returns the current circuit breaker state
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Failures returns the current failure count
func (cb *CircuitBreaker) Failures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// Successes returns the current success count (in half-open state)
func (cb *CircuitBreaker) Successes() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.successes
}
