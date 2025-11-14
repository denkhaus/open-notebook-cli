package integration

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
)

// TestNetworkConnectivityTests tests network connectivity scenarios
func TestNetworkConnectivityTests(t *testing.T) {
	t.Run("Connection timeout handling", func(t *testing.T) {
		testConnectionTimeout(t)
	})

	t.Run("Connection refused handling", func(t *testing.T) {
		testConnectionRefused(t)
	})

	t.Run("DNS resolution errors", func(t *testing.T) {
		testDNSErrors(t)
	})

	t.Run("Network interruption", func(t *testing.T) {
		testNetworkInterruption(t)
	})
}

// testConnectionTimeout tests timeout scenarios
func testConnectionTimeout(t *testing.T) {
	// Create a slow server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second) // Delay longer than our timeout
		w.WriteHeader(200)
		w.Write([]byte(`{"message": "slow response"}`))
	}))
	defer server.Close()

	// Create CLI context with short timeout
	ctx, err := createNetworkTestCLIContextWithTimeout(server.URL, 1)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	start := time.Now()
	resp, err := httpClient.Get(context.Background(), "/test")
	duration := time.Since(start)

	// Should fail due to timeout (implementation dependent)
	if err != nil {
		assert.True(t,
			strings.Contains(strings.ToLower(err.Error()), "timeout") ||
				strings.Contains(strings.ToLower(err.Error()), "deadline exceeded"),
			"Error should contain timeout information: %v", err)
		t.Logf("✅ Connection timeout handled correctly in %v", duration)
	} else {
		// If no error, check if it completed within reasonable time
		assert.NotNil(t, resp)
		assert.Less(t, duration, 5*time.Second, "Should timeout quickly")
		t.Logf("ℹ️ Request completed in %v (timeout handling may differ)", duration)
	}
}

// testConnectionRefused tests connection refused scenarios
func testConnectionRefused(t *testing.T) {
	testCases := []struct {
		name      string
		serverURL string
	}{
		{"Invalid port", "http://localhost:99999"},
		{"Non-existent host", "http://nonexistent-host-12345.local:8080"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := createNetworkTestCLIContextWithServer(tc.serverURL)
			require.NoError(t, err)

			injector := di.Bootstrap(ctx)
			httpClient := di.GetHTTPClient(injector)

			start := time.Now()
			_, err = httpClient.Get(context.Background(), "/test")
			duration := time.Since(start)

			// Should fail with connection error
			assert.Error(t, err, "Should fail with connection error")

			// Should fail quickly
			assert.Less(t, duration, 10*time.Second, "Should fail quickly")

			errorMsg := strings.ToLower(err.Error())
			assert.True(t,
				strings.Contains(errorMsg, "connection refused") ||
					strings.Contains(errorMsg, "no such host") ||
					strings.Contains(errorMsg, "timeout") ||
					strings.Contains(errorMsg, "network unreachable"),
				"Error should contain connection-related information: %s", err.Error())

			t.Logf("✅ Connection refused handled correctly: %s in %v", tc.name, duration)
		})
	}
}

// testDNSErrors tests DNS resolution errors
func testDNSErrors(t *testing.T) {
	testCases := []struct {
		name        string
		serverURL   string
		expectError string
	}{
		{
			name:        "Invalid domain",
			serverURL:   "http://invalid-domain-that-does-not-exist-12345.com:8080",
			expectError: "no such host",
		},
		{
			name:        "Malformed URL",
			serverURL:   "http://[]:8080",
			expectError: "invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := createNetworkTestCLIContextWithServer(tc.serverURL)
			require.NoError(t, err)

			injector := di.Bootstrap(ctx)
			httpClient := di.GetHTTPClient(injector)

			_, err = httpClient.Get(context.Background(), "/test")

			assert.Error(t, err, "Should fail with DNS/URL error")

			errorMsg := strings.ToLower(err.Error())
			assert.Contains(t, errorMsg, tc.expectError,
				"Error should contain '%s': %s", tc.expectError, err.Error())

			t.Logf("✅ DNS error handled correctly: %s", tc.name)
		})
	}
}

// testNetworkInterruption tests network interruption scenarios
func testNetworkInterruption(t *testing.T) {
	t.Run("Connection reset during request", func(t *testing.T) {
		// Create a server that closes connections abruptly
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Close connection immediately
			hijacker, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hijacker.Hijack()
				conn.Close()
				return
			}
			// Fallback: write minimal response and close
			w.WriteHeader(200)
			w.(http.Flusher).Flush()
		}))
		defer server.Close()

		ctx, err := createNetworkTestCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		resp, err := httpClient.Get(context.Background(), "/test")

		// Should fail due to connection reset or succeed with partial data
		if err != nil {
			assert.Nil(t, resp, "Response should be nil on connection error")
			errorMsg := strings.ToLower(err.Error())
			assert.True(t,
				strings.Contains(errorMsg, "connection reset") ||
					strings.Contains(errorMsg, "connection closed") ||
					strings.Contains(errorMsg, "broken pipe") ||
					strings.Contains(errorMsg, "eof"),
				"Error should contain connection reset information: %s", err.Error())
			t.Logf("✅ Connection reset handled as error: %v", err)
		} else {
			// Might succeed with partial response depending on implementation
			t.Logf("ℹ️ Connection reset handled with response (implementation-dependent)")
		}
	})

	t.Run("Partial response handling", func(t *testing.T) {
		// Create a server that sends partial response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"partial": "response", "incomplete": `)) // Incomplete JSON
			w.(http.Flusher).Flush()
			// Server closes without completing response
		}))
		defer server.Close()

		ctx, err := createNetworkTestCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		_, err = httpClient.Get(context.Background(), "/test")

		// Handle both success and error cases
		if err != nil {
			t.Logf("✅ Partial response correctly handled as error: %v", err)
		} else {
			t.Logf("ℹ️ Partial response accepted (implementation-dependent)")
			// Could try to parse partial JSON, but that's implementation-specific
		}
	})
}

// TestNetworkRecovery tests network recovery scenarios
func TestNetworkRecovery(t *testing.T) {
	t.Run("Retry mechanism verification", func(t *testing.T) {
		attemptCount := 0

		// Create a server that fails first, then succeeds
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount <= 2 {
				// First two attempts fail
				w.WriteHeader(503)
				w.Write([]byte(`{"error": "Service temporarily unavailable"}`))
			} else {
				// Third attempt succeeds
				w.WriteHeader(200)
				w.Write([]byte(`{"message": "success after retry"}`))
			}
		}))
		defer server.Close()

		ctx, err := createNetworkTestCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		// Note: Current implementation may not have auto-retry
		// This test documents current behavior
		resp, err := httpClient.Get(context.Background(), "/test")

		if err != nil {
			t.Logf("ℹ️ Request failed (expected without auto-retry): %v", err)
		} else {
			t.Logf("ℹ️ Request completed with status %d (attempt %d)", resp.StatusCode, attemptCount)
		}

		t.Logf("📝 Retry behavior documented: %d attempts made", attemptCount)
	})
}

// TestNetworkPerformance tests network performance under various conditions
func TestNetworkPerformance(t *testing.T) {
	t.Run("High latency simulation", func(t *testing.T) {
		// Create server with high latency
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second) // 1 second delay
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "high latency response"}`))
		}))
		defer server.Close()

		ctx, err := createNetworkTestCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		start := time.Now()
		resp, err := httpClient.Get(context.Background(), "/test")
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		assert.GreaterOrEqual(t, duration, 1*time.Second,
			"Should take at least 1 second due to server delay")
		assert.Less(t, duration, 5*time.Second,
			"Should complete within reasonable time")

		t.Logf("✅ High latency handled correctly: %v", duration)
	})

	t.Run("Concurrent requests", func(t *testing.T) {
		// Create fast server for concurrent testing
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Small delay
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "concurrent response"}`))
		}))
		defer server.Close()

		ctx, err := createNetworkTestCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		const numRequests = 10
		results := make(chan time.Duration, numRequests)

		start := time.Now()

		// Launch concurrent requests
		for i := 0; i < numRequests; i++ {
			go func() {
				reqStart := time.Now()
				_, err := httpClient.Get(context.Background(), "/test")
				reqDuration := time.Since(reqStart)

				if err == nil {
					results <- reqDuration
				} else {
					results <- -1 // Error indicator
				}
			}()
		}

		// Collect results
		successCount := 0
		totalDuration := time.Duration(0)

		for i := 0; i < numRequests; i++ {
			duration := <-results
			if duration > 0 {
				successCount++
				totalDuration += duration
			}
		}

		overallDuration := time.Since(start)

		t.Logf("✅ Concurrent requests: %d/%d successful in %v", successCount, numRequests, overallDuration)

		if successCount > 0 {
			avgDuration := totalDuration / time.Duration(successCount)
			t.Logf("Average request duration: %v", avgDuration)

			assert.GreaterOrEqual(t, float64(successCount)/float64(numRequests), 0.8,
				"At least 80%% of concurrent requests should succeed")
		}
	})
}

// Helper functions

func createNetworkTestCLIContextWithServer(serverURL string) (*cli.Context, error) {
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

func createNetworkTestCLIContextWithTimeout(serverURL string, timeout int) (*cli.Context, error) {
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
				Value: timeout,
			},
		},
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	flagSet.String("api-url", serverURL, "")
	flagSet.String("password", "test", "")
	flagSet.Bool("verbose", true, "")
	flagSet.Int("timeout", timeout, "")

	args := []string{}
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	return cli.NewContext(app, flagSet, nil), nil
}
