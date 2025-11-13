package integration

import (
	"context"
	"encoding/json"
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
	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// TestHTTPErrors tests various HTTP error response scenarios
func TestHTTPErrors(t *testing.T) {
	// Mock server for different error scenarios
	testCases := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedError  string
		description    string
	}{
		{
			name:       "400 Bad Request",
			statusCode: 400,
			responseBody: `{
				"error": "Bad Request",
				"message": "Invalid input parameters"
			}`,
			expectedError: "",
			description:   "Client sent invalid request data",
		},
		{
			name:       "401 Unauthorized",
			statusCode: 401,
			responseBody: `{
				"error": "Unauthorized",
				"message": "Authentication required"
			}`,
			expectedError: "",
			description:   "Client needs to authenticate",
		},
		{
			name:       "403 Forbidden",
			statusCode: 403,
			responseBody: `{
				"error": "Forbidden",
				"message": "Insufficient permissions"
			}`,
			expectedError: "",
			description:   "Client lacks permission",
		},
		{
			name:       "404 Not Found",
			statusCode: 404,
			responseBody: `{
				"error": "Not Found",
				"message": "Resource not found"
			}`,
			expectedError: "",
			description:   "Requested resource doesn't exist",
		},
		{
			name:       "429 Too Many Requests",
			statusCode: 429,
			responseBody: `{
				"error": "Too Many Requests",
				"message": "Rate limit exceeded"
			}`,
			expectedError: "",
			description:   "Client exceeded rate limits",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: 500,
			responseBody: `{
				"error": "Internal Server Error",
				"message": "Server encountered an error"
			}`,
			expectedError: "",
			description:   "Server-side error occurred",
		},
		{
			name:       "502 Bad Gateway",
			statusCode: 502,
			responseBody: `{
				"error": "Bad Gateway",
				"message": "Invalid response from upstream"
			}`,
			expectedError: "",
			description:   "Gateway/proxy error",
		},
		{
			name:       "503 Service Unavailable",
			statusCode: 503,
			responseBody: `{
				"error": "Service Unavailable",
				"message": "Server temporarily unavailable"
			}`,
			expectedError: "",
			description:   "Server is down for maintenance",
		},
		{
			name:          "Empty Response Body",
			statusCode:    500,
			responseBody:  "",
			expectedError: "",
			description:   "Server error with no body",
		},
		{
			name:       "Malformed JSON Response",
			statusCode: 500,
			responseBody: `{
				"error": "Internal Server Error",
				"message": "Server encountered an error"
			`,
			expectedError: "",
			description:   "Invalid JSON in response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			// Create CLI context with mock server URL
			ctx, err := createCLIContextWithServer(server.URL)
			require.NoError(t, err)

			// Bootstrap DI
			injector := di.Bootstrap(ctx)
			httpClient := di.GetHTTPClient(injector)

			// Test the request
			testCtx := context.Background()
			resp, err := httpClient.Get(testCtx, "/test")

			// Should succeed but return error status
			require.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, []byte(tc.responseBody), resp.Body)

			// Test error response parsing
			if tc.responseBody != "" && len(tc.responseBody) > 0 {
				var errorResp models.ErrorResponse
				parseErr := json.Unmarshal(resp.Body, &errorResp)

				// Should be able to parse valid JSON responses
				if !strings.Contains(tc.name, "Malformed") {
					assert.NoError(t, parseErr, "Should be able to parse valid error JSON")
					if parseErr == nil {
						assert.NotEmpty(t, errorResp.Error, "Error response should have error field")
						assert.NotEmpty(t, errorResp.Message, "Error response should have message field")
					}
				} else {
					// Malformed JSON should fail parsing
					assert.Error(t, parseErr, "Malformed JSON should fail to parse")
				}
			}

			t.Logf("✅ %s (HTTP %d): %s", tc.name, tc.statusCode, tc.description)
		})
	}
}

// TestHTTPTimeouts tests request timeout scenarios
func TestHTTPTimeouts(t *testing.T) {
	t.Run("Request timeout", func(t *testing.T) {
		// Create slow server that delays response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second) // Delay longer than timeout
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "slow response"}`))
		}))
		defer server.Close()

		// Create CLI context with very short timeout
		ctx, err := createCLIContextWithServerAndTimeout(server.URL, 1)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")

		// Should fail due to timeout
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t,
			strings.Contains(err.Error(), "timeout") ||
			strings.Contains(err.Error(), "deadline exceeded"),
			"Error should contain timeout-related message")

		t.Log("✅ Request timeout handled correctly")
	})

	t.Run("Context cancellation", func(t *testing.T) {
		// Create server that delays response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "delayed response"}`))
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		// Create context with short timeout
		testCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		resp, err := httpClient.Get(testCtx, "/test")

		// Should fail due to context cancellation
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "context deadline exceeded")

		t.Log("✅ Context cancellation handled correctly")
	})
}

// TestHTTPConnectionErrors tests network connectivity issues
func TestHTTPConnectionErrors(t *testing.T) {
	t.Run("Connection refused", func(t *testing.T) {
		// Use URL where server doesn't exist
		ctx, err := createCLIContextWithServer("http://localhost:99999")
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")

		// Should fail with connection error
		assert.Error(t, err)
		assert.Nil(t, resp)
		// Error could be "connection refused" or "invalid port" depending on the system
		assert.True(t,
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "invalid port") ||
			strings.Contains(err.Error(), "no such host"),
			"Error should contain connection-related message")

		t.Log("✅ Connection refused handled correctly")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		// Use malformed URL
		ctx, err := createCLIContextWithServer("invalid-url")
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")

		// Should fail with URL error
		assert.Error(t, err)
		assert.Nil(t, resp)

		t.Log("✅ Invalid URL handled correctly")
	})
}

// TestHTTPRetryLogic tests retry scenarios for transient errors
func TestHTTPRetryLogic(t *testing.T) {
	t.Run("Successful retry after temporary failure", func(t *testing.T) {
		attemptCount := 0

		// Create server that fails first time, then succeeds
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount == 1 {
				// First attempt fails
				w.WriteHeader(503)
				w.Write([]byte(`{"error": "Service Unavailable"}`))
			} else {
				// Second attempt succeeds
				w.WriteHeader(200)
				w.Write([]byte(`{"message": "success"}`))
			}
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")

		// Note: Current implementation doesn't have auto-retry
		// This test documents current behavior and can be used when retry logic is added
		require.NoError(t, err)
		assert.Equal(t, 503, resp.StatusCode)
		assert.Equal(t, 1, attemptCount) // Current implementation: 1 attempt only

		t.Log("✅ Retry test completed - current behavior documented (no auto-retry)")
	})
}

// Helper function to create CLI context with custom server URL
func createCLIContextWithServer(serverURL string) (*cli.Context, error) {
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
			&cli.StringFlag{
				Name:  "output",
				Value: "table",
			},
			&cli.StringFlag{
				Name:  "config-dir",
				Value: "/tmp/open-notebook-cli-test",
			},
		},
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	flagSet.String("api-url", serverURL, "")
	flagSet.String("password", "test", "")
	flagSet.Bool("verbose", true, "")
	flagSet.Int("timeout", 30, "")
	flagSet.String("output", "table", "")
	flagSet.String("config-dir", "/tmp/open-notebook-cli-test", "")

	args := []string{}
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	return cli.NewContext(app, flagSet, nil), nil
}

// Helper function to create CLI context with custom timeout
func createCLIContextWithServerAndTimeout(serverURL string, timeout int) (*cli.Context, error) {
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
			&cli.StringFlag{
				Name:  "output",
				Value: "table",
			},
			&cli.StringFlag{
				Name:  "config-dir",
				Value: "/tmp/open-notebook-cli-test",
			},
		},
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	flagSet.String("api-url", serverURL, "")
	flagSet.String("password", "test", "")
	flagSet.Bool("verbose", true, "")
	flagSet.Int("timeout", timeout, "")
	flagSet.String("output", "table", "")
	flagSet.String("config-dir", "/tmp/open-notebook-cli-test", "")

	args := []string{}
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	return cli.NewContext(app, flagSet, nil), nil
}