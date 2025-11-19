package services

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

// Enhanced HTTP client with retry logic and connection pooling
type retryableHTTPService struct {
	*httpService // Embed the original service
	retryConfig  RetryConfig
	classifier   *NetworkErrorClassifier
	diagnostics  *NetworkDiagnostics
}

// NewRetryableHTTPClient creates an enhanced HTTP client with retry logic
func NewRetryableHTTPClient(injector do.Injector) (shared.HTTPClient, error) {
	// Create the base HTTP service
	baseService, err := NewHTTPClient(injector)
	if err != nil {
		return nil, err
	}

	cfg := do.MustInvoke[config.Service](injector)
	logger := do.MustInvoke[shared.Logger](injector)

	// Get configuration
	httpConfig := DefaultHTTPClientConfig()
	if timeout := cfg.GetTimeout(); timeout > 0 {
		httpConfig.Timeout = time.Duration(timeout) * time.Second
	}

	// Create enhanced service
	enhanced := &retryableHTTPService{
		httpService: baseService.(*httpService),
		retryConfig: httpConfig.RetryConfig,
		classifier:  NewNetworkErrorClassifier(logger),
		diagnostics: NewNetworkDiagnostics(logger),
	}

	// Configure the underlying HTTP client with connection pooling
	enhanced.configureHTTPClient(httpConfig.ConnectionPoolConfig)

	logger.Debug("Enhanced HTTP client initialized",
		"max_retries", enhanced.retryConfig.MaxRetries,
		"base_delay", enhanced.retryConfig.BaseDelay,
		"timeout", httpConfig.Timeout,
	)

	return enhanced, nil
}

// configureHTTPClient configures the HTTP client with connection pooling settings
func (e *retryableHTTPService) configureHTTPClient(config ConnectionPoolConfig) {
	// Create custom transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   config.DialTimeout,
			KeepAlive: config.KeepAlive,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
	}

	// Apply transport to HTTP client
	e.httpClient.Transport = transport
}

// Get performs HTTP GET with retry logic
func (e *retryableHTTPService) Get(ctx context.Context, endpoint string) (*models.Response, error) {
	return e.classifier.RetryWithBackoff(ctx, e.retryConfig, func() (*models.Response, error) {
		return e.httpService.Get(ctx, endpoint)
	})
}

// Post performs HTTP POST with retry logic
func (e *retryableHTTPService) Post(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	return e.classifier.RetryWithBackoff(ctx, e.retryConfig, func() (*models.Response, error) {
		return e.httpService.Post(ctx, endpoint, body)
	})
}

// Put performs HTTP PUT with retry logic
func (e *retryableHTTPService) Put(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	return e.classifier.RetryWithBackoff(ctx, e.retryConfig, func() (*models.Response, error) {
		return e.httpService.Put(ctx, endpoint, body)
	})
}

// Delete performs HTTP DELETE with retry logic
func (e *retryableHTTPService) Delete(ctx context.Context, endpoint string) (*models.Response, error) {
	return e.classifier.RetryWithBackoff(ctx, e.retryConfig, func() (*models.Response, error) {
		return e.httpService.Delete(ctx, endpoint)
	})
}

// PostMultipart performs HTTP multipart POST with retry logic
func (e *retryableHTTPService) PostMultipart(ctx context.Context, endpoint string, fields map[string]string, files map[string]io.Reader) (*models.Response, error) {
	// Note: Multipart requests with file uploads are generally not retryable
	// due to stream consumption, so we call the base method directly
	return e.httpService.PostMultipart(ctx, endpoint, fields, files)
}

// Stream performs HTTP streaming with retry logic
func (e *retryableHTTPService) Stream(ctx context.Context, endpoint string, body interface{}) (<-chan []byte, error) {
	// Streaming requests are not retryable due to their stateful nature
	// We call the base method directly
	return e.httpService.Stream(ctx, endpoint, body)
}

// DiagnoseConnectivity performs network diagnostics for the configured API URL
func (e *retryableHTTPService) DiagnoseConnectivity(ctx context.Context) map[string]interface{} {
	apiURL := e.config.GetAPIURL()
	return e.diagnostics.DiagnoseConnectivity(ctx, apiURL)
}

// WithTimeout creates a new HTTP client with custom timeout
func (e *retryableHTTPService) WithTimeout(timeout time.Duration) shared.HTTPClient {
	// Create a copy of the underlying service
	baseCopy := e.httpService.WithTimeout(timeout).(*httpService)

	// Create enhanced wrapper with same retry config
	enhancedCopy := &retryableHTTPService{
		httpService: baseCopy,
		retryConfig: e.retryConfig,
		classifier:  e.classifier,
		diagnostics: e.diagnostics,
	}

	return enhancedCopy
}

// GetRetryConfig returns the current retry configuration
func (e *retryableHTTPService) GetRetryConfig() RetryConfig {
	return e.retryConfig
}

// SetRetryConfig updates the retry configuration
func (e *retryableHTTPService) SetRetryConfig(config RetryConfig) {
	e.retryConfig = config
	e.httpService.logger.Debug("Retry configuration updated",
		"max_retries", config.MaxRetries,
		"base_delay", config.BaseDelay,
		"max_delay", config.MaxDelay,
	)
}

// GracefulDegradation provides graceful degradation when API is unreachable
type GracefulDegradation struct {
	logger      shared.Logger
	classifier  *NetworkErrorClassifier
	diagnostics *NetworkDiagnostics
}

// NewGracefulDegradation creates a new graceful degradation handler
func NewGracefulDegradation(logger shared.Logger) *GracefulDegradation {
	return &GracefulDegradation{
		logger:      logger,
		classifier:  NewNetworkErrorClassifier(logger),
		diagnostics: NewNetworkDiagnostics(logger),
	}
}

// FallbackMode determines if we should enter fallback mode
type FallbackMode int

const (
	FallbackModeNone FallbackMode = iota
	FallbackModeOffline
	FallbackModeLimited
	FallbackModeCached
)

// EvaluateFallback evaluates if we should enter fallback mode based on errors
func (gd *GracefulDegradation) EvaluateFallback(err error) FallbackMode {
	if err == nil {
		return FallbackModeNone
	}

	errorType := gd.classifier.ClassifyError(err)

	switch errorType {
	case ErrorTypeConnectionRefused, ErrorTypeNetworkUnreachable:
		gd.logger.Warn("Network unreachable - entering offline mode", "error", err)
		return FallbackModeOffline

	case ErrorTypeTimeout:
		gd.logger.Warn("Network timeout - entering limited mode", "error", err)
		return FallbackModeLimited

	case ErrorTypeDNSResolution:
		gd.logger.Warn("DNS resolution failed - entering offline mode", "error", err)
		return FallbackModeOffline

	case ErrorTypeConnectionReset, ErrorTypeTemporaryFailure:
		gd.logger.Warn("Temporary network issue - using cached mode", "error", err)
		return FallbackModeCached

	default:
		gd.logger.Debug("Unknown error - no fallback mode", "error", err)
		return FallbackModeNone
	}
}

// GetFallbackMessage returns a user-friendly message for the fallback mode
func (gd *GracefulDegradation) GetFallbackMessage(mode FallbackMode) string {
	switch mode {
	case FallbackModeOffline:
		return "API is currently unreachable. Working in offline mode with limited functionality."
	case FallbackModeLimited:
		return "API connectivity is degraded. Some operations may be slow or unavailable."
	case FallbackModeCached:
		return "Using cached data due to temporary connectivity issues."
	default:
		return ""
	}
}
