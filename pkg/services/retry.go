package services

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// RetryConfig holds retry configuration for network operations
type RetryConfig struct {
	MaxRetries      int           `json:"max_retries"`
	BaseDelay       time.Duration `json:"base_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	BackoffFactor   float64       `json:"backoff_factor"`
	RetryableErrors []string      `json:"retryable_errors"`
	RetryableStatus []int         `json:"retryable_status"`
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		BaseDelay:     100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		RetryableErrors: []string{
			"connection refused",
			"timeout",
			"network unreachable",
			"temporary failure",
			"connection reset",
			"broken pipe",
			"eof",
			"connection timeout",
			"deadline exceeded",
		},
		RetryableStatus: []int{
			http.StatusRequestTimeout,      // 408
			http.StatusTooManyRequests,     // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout,      // 504
		},
	}
}

// NetworkErrorClassifier helps classify different types of network errors
type NetworkErrorClassifier struct {
	logger Logger
}

// NewNetworkErrorClassifier creates a new error classifier
func NewNetworkErrorClassifier(logger Logger) *NetworkErrorClassifier {
	return &NetworkErrorClassifier{logger: logger}
}

// ErrorType represents different types of network errors
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeConnectionRefused
	ErrorTypeTimeout
	ErrorTypeDNSResolution
	ErrorTypeNetworkUnreachable
	ErrorTypeConnectionReset
	ErrorTypeTemporaryFailure
	ErrorTypeHTTPError
)

// ClassifyError classifies the error into specific types
func (nec *NetworkErrorClassifier) ClassifyError(err error) ErrorType {
	if err == nil {
		return ErrorTypeUnknown
	}

	errStr := strings.ToLower(err.Error())

	// Connection refused errors
	if strings.Contains(errStr, "connection refused") {
		return ErrorTypeConnectionRefused
	}

	// Timeout errors
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return ErrorTypeTimeout
	}

	// DNS resolution errors
	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "name resolution") {
		return ErrorTypeDNSResolution
	}

	// Network unreachable errors
	if strings.Contains(errStr, "network unreachable") || strings.Contains(errStr, "no route to host") {
		return ErrorTypeNetworkUnreachable
	}

	// Connection reset errors
	if strings.Contains(errStr, "connection reset") || strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "eof") {
		return ErrorTypeConnectionReset
	}

	// Check for URL parsing errors
	if _, urlErr := err.(*url.Error); urlErr {
		return ErrorTypeTemporaryFailure
	}

	// Check for net errors
	if _, netErr := err.(*net.OpError); netErr {
		return ErrorTypeTemporaryFailure
	}

	return ErrorTypeUnknown
}

// IsRetryable determines if an error is retryable based on its type and configuration
func (nec *NetworkErrorClassifier) IsRetryable(err error, config RetryConfig) bool {
	if err == nil {
		return false
	}

	errorType := nec.ClassifyError(err)
	errStr := strings.ToLower(err.Error())

	// Check if error type is generally retryable
	switch errorType {
	case ErrorTypeConnectionRefused, ErrorTypeTimeout, ErrorTypeNetworkUnreachable,
		 ErrorTypeConnectionReset, ErrorTypeTemporaryFailure:
		return true
	case ErrorTypeDNSResolution:
		// DNS errors are generally not retryable unless they're temporary
		return strings.Contains(errStr, "temporary")
	}

	// Check specific error messages
	for _, retryableErr := range config.RetryableErrors {
		if strings.Contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// IsRetryableStatus determines if an HTTP status code is retryable
func (nec *NetworkErrorClassifier) IsRetryableStatus(statusCode int, config RetryConfig) bool {
	for _, retryableStatus := range config.RetryableStatus {
		if statusCode == retryableStatus {
			return true
		}
	}
	return false
}

// RetryWithBackoff executes a function with exponential backoff retry logic
func (nec *NetworkErrorClassifier) RetryWithBackoff(
	ctx context.Context,
	config RetryConfig,
	operation func() (*models.Response, error),
) (*models.Response, error) {
	var lastErr error
	var resp *models.Response

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff and jitter
			delay := nec.calculateBackoffDelay(attempt, config)

			nec.logger.Debug("Retrying network operation",
				"attempt", attempt,
				"max_retries", config.MaxRetries,
				"delay", delay,
				"last_error", lastErr,
			)

			// Wait before retry or exit if context is cancelled
			select {
			case <-ctx.Done():
				return resp, ctx.Err()
			case <-time.After(delay):
				// Continue with retry
			}
		}

		// Execute the operation
		resp, lastErr = operation()
		if lastErr == nil {
			// Success - check if response status is retryable
			if resp == nil || !nec.IsRetryableStatus(resp.StatusCode, config) {
				return resp, nil
			}

			nec.logger.Debug("Received retryable HTTP status",
				"status", resp.StatusCode,
				"attempt", attempt,
			)

			// Treat retryable status as an error for retry logic
			lastErr = fmt.Errorf("HTTP %d: retryable status", resp.StatusCode)
		} else {
			// Check if this error is retryable
			if !nec.IsRetryable(lastErr, config) {
				return resp, lastErr
			}

			nec.logger.Debug("Network operation failed with retryable error",
				"attempt", attempt,
				"error_type", nec.ClassifyError(lastErr),
				"error", lastErr,
			)
		}
	}

	return resp, lastErr
}

// calculateBackoffDelay calculates exponential backoff delay with jitter
func (nec *NetworkErrorClassifier) calculateBackoffDelay(attempt int, config RetryConfig) time.Duration {
	// Exponential backoff: delay = baseDelay * (backoffFactor ^ (attempt-1))
	delay := float64(config.BaseDelay) * math.Pow(config.BackoffFactor, float64(attempt-1))

	// Add jitter to prevent thundering herd (±25% random variation)
	jitter := delay * 0.25 * (rand.Float64()*2 - 1)
	delay += jitter

	// Ensure delay doesn't exceed maximum
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// ConnectionPoolConfig holds connection pool settings
type ConnectionPoolConfig struct {
	MaxIdleConns        int           `json:"max_idle_conns"`
	MaxIdleConnsPerHost int           `json:"max_idle_conns_per_host"`
	MaxConnsPerHost     int           `json:"max_conns_per_host"`
	IdleConnTimeout     time.Duration `json:"idle_conn_timeout"`
	DialTimeout         time.Duration `json:"dial_timeout"`
	KeepAlive           time.Duration `json:"keep_alive"`
}

// DefaultConnectionPoolConfig returns default connection pool configuration
func DefaultConnectionPoolConfig() ConnectionPoolConfig {
	return ConnectionPoolConfig{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     90 * time.Second,
		DialTimeout:         10 * time.Second,
		KeepAlive:           30 * time.Second,
	}
}

// HTTPClientConfig holds complete HTTP client configuration
type HTTPClientConfig struct {
	Timeout             time.Duration            `json:"timeout"`
	RetryConfig         RetryConfig              `json:"retry_config"`
	ConnectionPoolConfig ConnectionPoolConfig     `json:"connection_pool_config"`
}

// DefaultHTTPClientConfig returns default HTTP client configuration
func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		Timeout:             30 * time.Second,
		RetryConfig:         DefaultRetryConfig(),
		ConnectionPoolConfig: DefaultConnectionPoolConfig(),
	}
}

// NetworkDiagnostics provides network connectivity diagnostics
type NetworkDiagnostics struct {
	logger Logger
}

// NewNetworkDiagnostics creates a new network diagnostics instance
func NewNetworkDiagnostics(logger Logger) *NetworkDiagnostics {
	return &NetworkDiagnostics{logger: logger}
}

// DiagnoseConnectivity performs network connectivity diagnostics
func (nd *NetworkDiagnostics) DiagnoseConnectivity(ctx context.Context, serverURL string) map[string]interface{} {
	diagnostics := make(map[string]interface{})

	start := time.Now()

	// Test basic TCP connectivity
	if host, port, err := nd.parseHostPort(serverURL); err == nil {
		diagnostics["tcp_test"] = nd.testTCPConnectivity(ctx, host, port, 5*time.Second)
	}

	// Test HTTP connectivity
	diagnostics["http_test"] = nd.testHTTPConnectivity(ctx, serverURL, 10*time.Second)

	// Test DNS resolution
	diagnostics["dns_test"] = nd.testDNSResolution(serverURL)

	diagnostics["total_duration"] = time.Since(start)

	return diagnostics
}

// parseHostPort extracts host and port from URL
func (nd *NetworkDiagnostics) parseHostPort(serverURL string) (string, int, error) {
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return "", 0, err
	}

	host := parsed.Hostname()
	port := parsed.Port()
	if port == "" {
		if parsed.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return host, 0, fmt.Errorf("port parsing not implemented")
}

// testTCPConnectivity tests basic TCP connectivity
func (nd *NetworkDiagnostics) testTCPConnectivity(ctx context.Context, host string, port int, timeout time.Duration) map[string]interface{} {
	result := make(map[string]interface{})

	start := time.Now()

	// Create dialer with timeout
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
	duration := time.Since(start)

	result["duration"] = duration
	result["success"] = err == nil

	if err != nil {
		result["error"] = err.Error()
		nd.logger.Debug("TCP connectivity test failed", "host", host, "port", port, "error", err)
	} else {
		conn.Close()
		nd.logger.Debug("TCP connectivity test successful", "host", host, "port", port, "duration", duration)
	}

	return result
}

// testHTTPConnectivity tests HTTP connectivity
func (nd *NetworkDiagnostics) testHTTPConnectivity(ctx context.Context, serverURL string, timeout time.Duration) map[string]interface{} {
	result := make(map[string]interface{})

	start := time.Now()

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		result["error"] = err.Error()
		result["success"] = false
		return result
	}

	resp, err := client.Do(req)
	duration := time.Since(start)

	result["duration"] = duration

	if err != nil {
		result["success"] = false
		result["error"] = err.Error()
		nd.logger.Debug("HTTP connectivity test failed", "url", serverURL, "error", err)
	} else {
		defer resp.Body.Close()
		result["success"] = true
		result["status_code"] = resp.StatusCode
		nd.logger.Debug("HTTP connectivity test successful", "url", serverURL, "status", resp.StatusCode, "duration", duration)
	}

	return result
}

// testDNSResolution tests DNS resolution
func (nd *NetworkDiagnostics) testDNSResolution(serverURL string) map[string]interface{} {
	result := make(map[string]interface{})

	start := time.Now()

	parsed, err := url.Parse(serverURL)
	if err != nil {
		result["error"] = err.Error()
		result["success"] = false
		return result
	}

	host := parsed.Hostname()

	ips, err := net.LookupIP(host)
	duration := time.Since(start)

	result["duration"] = duration
	result["host"] = host

	if err != nil {
		result["success"] = false
		result["error"] = err.Error()
		nd.logger.Debug("DNS resolution test failed", "host", host, "error", err)
	} else {
		result["success"] = true
		var ipStrings []string
		for _, ip := range ips {
			ipStrings = append(ipStrings, ip.String())
		}
		result["ips"] = ipStrings
		nd.logger.Debug("DNS resolution test successful", "host", host, "ips", ipStrings, "duration", duration)
	}

	return result
}