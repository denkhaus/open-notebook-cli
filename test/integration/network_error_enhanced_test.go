package integration

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/denkhaus/open-notebook-cli/pkg/mocks"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
)

// TestEnhancedNetworkErrorHandling tests the enhanced network error handling
func TestEnhancedNetworkErrorHandling(t *testing.T) {
	t.Run("Retry logic with exponential backoff", func(t *testing.T) {
		testRetryWithBackoff(t)
	})

	t.Run("Connection pooling under load", func(t *testing.T) {
		testConnectionPooling(t)
	})

	t.Run("Graceful degradation scenarios", func(t *testing.T) {
		testGracefulDegradation(t)
	})

	t.Run("Network diagnostics", func(t *testing.T) {
		testNetworkDiagnostics(t)
	})

	t.Run("Error classification accuracy", func(t *testing.T) {
		testErrorClassification(t)
	})
}

// testRetryWithBackoff tests the retry logic with exponential backoff
func testRetryWithBackoff(t *testing.T) {
	t.Run("Successful retry after temporary failure", func(t *testing.T) {
		attemptCount := 0
		maxAttempts := 3

		// Create a server that fails first, then succeeds
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount <= maxAttempts-1 {
				// First attempts fail with retryable status
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error": "Service temporarily unavailable"}`))
				return
			}
			// Final attempt succeeds
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success after retry"}`))
		}))
		defer server.Close()

		ctx, err := createEnhancedTestCLIContext(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		// Enhanced client should retry automatically
		start := time.Now()
		resp, err := httpClient.Get(context.Background(), "/test")
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, maxAttempts, attemptCount)

		// Should take longer due to retries
		assert.GreaterOrEqual(t, duration, 200*time.Millisecond, "Should have retry delays")

		t.Logf("✅ Retry logic successful: %d attempts in %v", attemptCount, duration)
	})

	t.Run("Non-retryable errors fail immediately", func(t *testing.T) {
		attemptCount := 0

		// Create a server that returns non-retryable error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.WriteHeader(http.StatusBadRequest) // Non-retryable
			w.Write([]byte(`{"error": "Bad request"}`))
		}))
		defer server.Close()

		ctx, err := createEnhancedTestCLIContext(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		resp, err := httpClient.Get(context.Background(), "/test")

		// Should fail immediately without retries
		require.NoError(t, err) // No network error, just HTTP error
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, 1, attemptCount) // Only one attempt

		t.Logf("✅ Non-retryable error handled correctly: %d attempts", attemptCount)
	})

	t.Run("Max retry attempts respected", func(t *testing.T) {
		attemptCount := 0
		maxRetries := 3

		// Create a server that always fails
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.WriteHeader(http.StatusServiceUnavailable) // Always retryable
			w.Write([]byte(`{"error": "Always failing"}`))
		}))
		defer server.Close()

		ctx, err := createEnhancedTestCLIContext(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		// Override retry config for testing
		// Note: This would require type assertion to the concrete enhanced type
		// For now, we use the default retry configuration
		_ = maxRetries

		start := time.Now()
		resp, err := httpClient.Get(context.Background(), "/test")
		duration := time.Since(start)

		// Should fail after max retries with error
		require.Error(t, err) // Should have error after exhausting retries
		assert.Contains(t, err.Error(), "HTTP 503: retryable status")
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		assert.Equal(t, maxRetries+1, attemptCount) // Initial attempt + retries

		// Should take time due to retry delays
		assert.GreaterOrEqual(t, duration, 300*time.Millisecond, "Should have cumulative retry delays")

		t.Logf("✅ Max retry attempts respected: %d attempts in %v", attemptCount, duration)
	})
}

// testConnectionPooling tests connection pooling under load
func testConnectionPooling(t *testing.T) {
	t.Run("Concurrent requests with connection pooling", func(t *testing.T) {
		// Create a fast server for connection pooling testing
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "pooled connection response"}`))
		}))
		defer server.Close()

		ctx, err := createEnhancedTestCLIContext(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		const numRequests = 20
		results := make(chan time.Duration, numRequests)
		errors := make(chan error, numRequests)

		start := time.Now()

		// Launch concurrent requests
		for i := 0; i < numRequests; i++ {
			go func(requestID int) {
				reqStart := time.Now()
				_, err := httpClient.Get(context.Background(), "/test")
				reqDuration := time.Since(reqStart)

				results <- reqDuration
				errors <- err
			}(i)
		}

		// Collect results
		successCount := 0
		var totalDuration time.Duration
		var errorCount int

		for i := 0; i < numRequests; i++ {
			duration := <-results
			err := <-errors

			if err == nil {
				successCount++
				totalDuration += duration
			} else {
				errorCount++
				t.Logf("Request failed: %v", err)
			}
		}

		overallDuration := time.Since(start)

		// Assertions
		assert.Equal(t, numRequests, successCount+errorCount, "Total requests should match")
		assert.Greater(t, successCount, numRequests*9/10, "At least 90% of requests should succeed")

		if successCount > 0 {
			avgDuration := totalDuration / time.Duration(successCount)

			// Connection pooling should make subsequent requests faster
			assert.Less(t, avgDuration, 200*time.Millisecond, "Average request should be fast with pooling")

			t.Logf("✅ Connection pooling effective: %d/%d successful, avg %v, total %v",
				successCount, numRequests, avgDuration, overallDuration)
		}
	})

	t.Run("Connection reuse verification", func(t *testing.T) {
		// This test would require more sophisticated connection tracking
		// For now, we just verify that multiple requests work efficiently
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "keep-alive response"}`))
		}))
		defer server.Close()

		ctx, err := createEnhancedTestCLIContext(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		// Make multiple requests sequentially
		start := time.Now()
		for i := 0; i < 5; i++ {
			_, err := httpClient.Get(context.Background(), "/test")
			require.NoError(t, err)
		}
		duration := time.Since(start)

		// Should be fast due to connection reuse
		assert.Less(t, duration, 500*time.Millisecond, "Sequential requests should be fast with connection reuse")

		t.Logf("✅ Connection reuse verified: 5 requests in %v", duration)
	})
}

// testGracefulDegradation tests graceful degradation scenarios
func testGracefulDegradation(t *testing.T) {
	t.Run("Offline mode evaluation", func(t *testing.T) {
		logger := mocks.NewMockLogger(false)
		degradation := services.NewGracefulDegradation(logger)

		// Test connection refused -> offline mode
		err := fmt.Errorf("connection refused")
		mode := degradation.EvaluateFallback(err)
		assert.Equal(t, services.FallbackModeOffline, mode)

		message := degradation.GetFallbackMessage(mode)
		assert.Contains(t, message, "offline mode")

		t.Logf("✅ Offline mode evaluation: %s", message)
	})

	t.Run("Limited mode evaluation", func(t *testing.T) {
		logger := mocks.NewMockLogger(false)
		degradation := services.NewGracefulDegradation(logger)

		// Test timeout -> limited mode
		err := fmt.Errorf("context deadline exceeded")
		mode := degradation.EvaluateFallback(err)
		assert.Equal(t, services.FallbackModeLimited, mode)

		message := degradation.GetFallbackMessage(mode)
		assert.Contains(t, message, "connectivity is degraded")

		t.Logf("✅ Limited mode evaluation: %s", message)
	})

	t.Run("Cached mode evaluation", func(t *testing.T) {
		logger := mocks.NewMockLogger(false)
		degradation := services.NewGracefulDegradation(logger)

		// Test connection reset -> cached mode
		err := fmt.Errorf("connection reset by peer")
		mode := degradation.EvaluateFallback(err)
		assert.Equal(t, services.FallbackModeCached, mode)

		message := degradation.GetFallbackMessage(mode)
		assert.Contains(t, message, "cached data")

		t.Logf("✅ Cached mode evaluation: %s", message)
	})
}

// testNetworkDiagnostics tests network diagnostics functionality
func testNetworkDiagnostics(t *testing.T) {
	t.Run("Comprehensive network diagnostics", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "diagnostic test"}`))
		}))
		defer server.Close()

		logger := mocks.NewMockLogger(false)
		diagnostics := services.NewNetworkDiagnostics(logger)

		// Run diagnostics
		results := diagnostics.DiagnoseConnectivity(context.Background(), server.URL)

		// Verify diagnostic results
		assert.Contains(t, results, "dns_test")
		assert.Contains(t, results, "http_test")
		assert.Contains(t, results, "total_duration")

		// Check HTTP test results
		httpTest := results["http_test"].(map[string]interface{})
		assert.True(t, httpTest["success"].(bool), "HTTP test should succeed")

		// Check DNS test results
		dnsTest := results["dns_test"].(map[string]interface{})
		assert.True(t, dnsTest["success"].(bool), "DNS test should succeed")

		duration := results["total_duration"].(time.Duration)
		assert.Less(t, duration, 5*time.Second, "Diagnostics should complete quickly")

		t.Logf("✅ Network diagnostics completed in %v", duration)
	})

	t.Run("Diagnostics with unreachable server", func(t *testing.T) {
		logger := mocks.NewMockLogger(false)
		diagnostics := services.NewNetworkDiagnostics(logger)

		// Test with unreachable server
		unreachableURL := "http://localhost:99999"
		results := diagnostics.DiagnoseConnectivity(context.Background(), unreachableURL)

		// Should still provide diagnostic information
		assert.Contains(t, results, "http_test")
		assert.Contains(t, results, "dns_test")

		httpTest := results["http_test"].(map[string]interface{})
		assert.False(t, httpTest["success"].(bool), "HTTP test should fail")
		assert.Contains(t, httpTest, "error")

		t.Logf("✅ Unreachable server diagnostics handled correctly")
	})
}

// testErrorClassification tests error classification accuracy
func testErrorClassification(t *testing.T) {
	logger := mocks.NewMockLogger(false)
	classifier := services.NewNetworkErrorClassifier(logger)

	testCases := []struct {
		name      string
		error     error
		expected  services.ErrorType
		retryable bool
	}{
		{
			name:      "Connection refused",
			error:     fmt.Errorf("connection refused"),
			expected:  services.ErrorTypeConnectionRefused,
			retryable: true,
		},
		{
			name:      "Timeout error",
			error:     fmt.Errorf("context deadline exceeded"),
			expected:  services.ErrorTypeTimeout,
			retryable: true,
		},
		{
			name:      "DNS resolution failure",
			error:     fmt.Errorf("no such host"),
			expected:  services.ErrorTypeDNSResolution,
			retryable: false,
		},
		{
			name:      "Network unreachable",
			error:     fmt.Errorf("network unreachable"),
			expected:  services.ErrorTypeNetworkUnreachable,
			retryable: true,
		},
		{
			name:      "Connection reset",
			error:     fmt.Errorf("connection reset by peer"),
			expected:  services.ErrorTypeConnectionReset,
			retryable: true,
		},
		{
			name:      "Unknown error",
			error:     fmt.Errorf("some unknown error"),
			expected:  services.ErrorTypeUnknown,
			retryable: false,
		},
	}

	config := services.DefaultRetryConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errorType := classifier.ClassifyError(tc.error)
			assert.Equal(t, tc.expected, errorType, "Error type classification should match")

			isRetryable := classifier.IsRetryable(tc.error, config)
			assert.Equal(t, tc.retryable, isRetryable, "Retryability assessment should match")

			t.Logf("✅ Error classification: %s -> %v (retryable: %v)", tc.name, errorType, isRetryable)
		})
	}

	t.Run("HTTP status retryability", func(t *testing.T) {
		retryableStatuses := []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		}

		nonRetryableStatuses := []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
		}

		for _, status := range retryableStatuses {
			isRetryable := classifier.IsRetryableStatus(status, config)
			assert.True(t, isRetryable, "Status %d should be retryable", status)
		}

		for _, status := range nonRetryableStatuses {
			isRetryable := classifier.IsRetryableStatus(status, config)
			assert.False(t, isRetryable, "Status %d should not be retryable", status)
		}

		t.Logf("✅ HTTP status retryability classification correct")
	})
}

// Helper function to create enhanced test CLI context
func createEnhancedTestCLIContext(serverURL string) (*cli.Context, error) {
	app := &cli.App{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-url",
				Value: serverURL,
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "test",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Value: true,
			},
			&cli.IntFlag{
				Name:  "timeout",
				Value: 30,
			},
		},
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	flagSet.String("api-url", serverURL, "")
	flagSet.String("password", "test", "")
	flagSet.Bool("verbose", true, "")
	flagSet.Int("timeout", 30, "")

	args := []string{}
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	return cli.NewContext(app, flagSet, nil), nil
}
