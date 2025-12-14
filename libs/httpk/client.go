package httpk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client with retry and circuit breaker support
type Client struct {
	httpClient     *http.Client
	retry          *RetryConfig
	circuitBreaker *CircuitBreaker
	logger         Logger
}

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	Ctx     context.Context
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

// Config holds client configuration
type Config struct {
	Timeout        time.Duration
	Retry          *RetryConfig
	CircuitBreaker *CircuitBreakerConfig
	Logger         Logger
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout: 30 * time.Second,
		Retry:   DefaultRetryConfig(),
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 5,
			SuccessThreshold: 2,
			Timeout:          60 * time.Second,
		},
		Logger: nil, // Will use default logger if not set
	}
}

// NewClient creates a new HTTP client
func NewClient(cfg *Config) *Client {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		retry:  cfg.Retry,
		logger: cfg.Logger,
	}

	// Initialize circuit breaker if config is provided
	if cfg.CircuitBreaker != nil {
		client.circuitBreaker = NewCircuitBreaker(*cfg.CircuitBreaker, client.logger)
	}

	return client
}

// Do executes an HTTP request with retry and circuit breaker
func (c *Client) Do(req *Request) (*Response, error) {
	// Use circuit breaker if enabled
	if c.circuitBreaker != nil {
		return c.circuitBreaker.Execute(func() (*Response, error) {
			return c.doWithRetry(req)
		})
	}

	return c.doWithRetry(req)
}

// doWithRetry executes request with retry logic
func (c *Client) doWithRetry(req *Request) (*Response, error) {
	var lastErr error
	var resp *Response

	retryCount := 0
	if c.retry != nil {
		retryCount = c.retry.MaxRetries
	}

	for attempt := 0; attempt <= retryCount; attempt++ {
		if attempt > 0 && c.logger != nil {
			c.logger.Infof("Retrying request (attempt %d/%d)", attempt, retryCount)
		}

		resp, lastErr = c.doRequest(req)

		// If there's a network/request error, don't retry (fail fast)
		if lastErr != nil {
			break
		}

		// Success - status code is OK and not in retry list
		if !c.shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		// Don't retry on last attempt
		if attempt == retryCount {
			break
		}

		// Calculate backoff and sleep (only for status code retries)
		if c.retry != nil {
			backoff := c.retry.Backoff(attempt)
			if c.logger != nil {
				c.logger.Debugf("Backing off for %v before retry (status %d)", backoff, resp.StatusCode)
			}
			time.Sleep(backoff)
		}
	}

	// Return error or last response
	if lastErr != nil {
		return nil, fmt.Errorf("request failed: %w", lastErr)
	}

	// If we get here, we exhausted retries on a bad status code
	return resp, nil
}

// doRequest executes a single HTTP request
func (c *Client) doRequest(req *Request) (*Response, error) {
	start := time.Now()

	// Prepare request body
	var bodyReader io.Reader
	if req.Body != nil {
		jsonBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	// Create HTTP request
	ctx := req.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set default content type if not set
	if req.Body != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Log request
	if c.logger != nil {
		c.logger.Debugf("HTTP %s %s", req.Method, req.URL)
	}

	// Execute request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		if c.logger != nil {
			c.logger.Errorf("Request failed: %v", err)
		}
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	duration := time.Since(start)

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       respBody,
		Duration:   duration,
	}

	// Log response
	if c.logger != nil {
		c.logger.Debugf("HTTP %d %s (took %v)", resp.StatusCode, req.URL, duration)
	}

	return resp, nil
}

// shouldRetry determines if request should be retried based on status code
func (c *Client) shouldRetry(statusCode int) bool {
	if c.retry == nil {
		return false
	}

	for _, code := range c.retry.RetryStatuses {
		if statusCode == code {
			return true
		}
	}

	return false
}

// GET performs a GET request
func (c *Client) GET(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
		Ctx:     ctx,
	})
}

// POST performs a POST request
func (c *Client) POST(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodPost,
		URL:     url,
		Headers: headers,
		Body:    body,
		Ctx:     ctx,
	})
}

// PUT performs a PUT request
func (c *Client) PUT(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodPut,
		URL:     url,
		Headers: headers,
		Body:    body,
		Ctx:     ctx,
	})
}

// PATCH performs a PATCH request
func (c *Client) PATCH(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodPatch,
		URL:     url,
		Headers: headers,
		Body:    body,
		Ctx:     ctx,
	})
}

// DELETE performs a DELETE request
func (c *Client) DELETE(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
		Ctx:     ctx,
	})
}

// DecodeJSON decodes response body to target interface
func (r *Response) DecodeJSON(target interface{}) error {
	if err := json.Unmarshal(r.Body, target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return nil
}
