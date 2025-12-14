package httpk

import (
	"fmt"
)

// Error types for httpk
const (
	ErrCodeNetworkError       = "NETWORK_ERROR"
	ErrCodeTimeout            = "TIMEOUT"
	ErrCodeCircuitBreakerOpen = "CIRCUIT_BREAKER_OPEN"
	ErrCodeInvalidResponse    = "INVALID_RESPONSE"
	ErrCodeMaxRetriesExceeded = "MAX_RETRIES_EXCEEDED"
)

// HTTPError represents an HTTP error with additional context
type HTTPError struct {
	Code       string
	Message    string
	StatusCode int
	URL        string
	Method     string
	RetryCount int
	Err        error
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v (HTTP %d %s %s, retries: %d)",
			e.Code, e.Message, e.Err, e.StatusCode, e.Method, e.URL, e.RetryCount)
	}
	return fmt.Sprintf("[%s] %s (HTTP %d %s %s, retries: %d)",
		e.Code, e.Message, e.StatusCode, e.Method, e.URL, e.RetryCount)
}

// Unwrap returns the underlying error
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewNetworkError creates a network error
func NewNetworkError(method, url string, err error) *HTTPError {
	return &HTTPError{
		Code:    ErrCodeNetworkError,
		Message: "Network request failed",
		Method:  method,
		URL:     url,
		Err:     err,
	}
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(method, url string) *HTTPError {
	return &HTTPError{
		Code:    ErrCodeTimeout,
		Message: "Request timeout",
		Method:  method,
		URL:     url,
	}
}

// NewCircuitBreakerError creates a circuit breaker error
func NewCircuitBreakerError(method, url string) *HTTPError {
	return &HTTPError{
		Code:    ErrCodeCircuitBreakerOpen,
		Message: "Circuit breaker is open",
		Method:  method,
		URL:     url,
	}
}

// NewInvalidResponseError creates an invalid response error
func NewInvalidResponseError(method, url string, statusCode int, err error) *HTTPError {
	return &HTTPError{
		Code:       ErrCodeInvalidResponse,
		Message:    "Invalid response received",
		StatusCode: statusCode,
		Method:     method,
		URL:        url,
		Err:        err,
	}
}

// NewMaxRetriesError creates a max retries exceeded error
func NewMaxRetriesError(method, url string, retryCount int, err error) *HTTPError {
	return &HTTPError{
		Code:       ErrCodeMaxRetriesExceeded,
		Message:    "Maximum retry attempts exceeded",
		Method:     method,
		URL:        url,
		RetryCount: retryCount,
		Err:        err,
	}
}
