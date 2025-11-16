package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
)

// TestAdditionalErrorScenarios enhances the existing error handling tests
// with additional edge cases and complex scenarios that might have been missed
func TestAdditionalErrorScenarios(t *testing.T) {
	t.Run("Rate limiting with retry-after header", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			if requestCount <= 2 {
				w.Header().Set("Retry-After", "5")
				w.Header().Set("X-RateLimit-Limit", "100")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", "1640995200")
				w.WriteHeader(429)
				w.Write([]byte(`{"error": "Too Many Requests", "message": "Rate limit exceeded"}`))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "success after rate limit"}`))
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")
		require.NoError(t, err)
		// The client retries, so we should get success on final attempt
		assert.Equal(t, 200, resp.StatusCode)
		// But we can check that requestCount was incremented multiple times due to retries
		assert.Greater(t, requestCount, 2, "Should have retried multiple times")

		t.Log("✅ Rate limiting with headers tested")
	})

	t.Run("Large response body handling", func(t *testing.T) {
		// Create a very large response (10MB)
		largeBody := strings.Repeat("x", 10*1024*1024)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10485760")
			w.WriteHeader(200)
			w.Write([]byte(largeBody))
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServerAndTimeout(server.URL, 30) // 30 second timeout
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 10485760, len(resp.Body))

		t.Log("✅ Large response body handling tested")
	})

	t.Run("Chunked transfer encoding", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Transfer-Encoding", "chunked")
			w.WriteHeader(200)

			// Write in chunks
			flusher, _ := w.(http.Flusher)
			w.Write([]byte("chunk1 "))
			if flusher != nil {
				flusher.Flush()
			}
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte("chunk2 "))
			if flusher != nil {
				flusher.Flush()
			}
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte("chunk3"))
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "chunk1 chunk2 chunk3", string(resp.Body))

		t.Log("✅ Chunked transfer encoding tested")
	})

	t.Run("Compression handling (gzip)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`{"compressed": true, "message": "This should be gzip compressed"}`))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(`{"compressed": false, "message": "No compression"}`))
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		t.Log("✅ Compression handling tested")
	})

	t.Run("Redirect handling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/redirect":
				w.Header().Set("Location", "/target")
				w.WriteHeader(301)
				w.Write([]byte(`{"redirect": true}`))
			case "/target":
				w.WriteHeader(200)
				w.Write([]byte(`{"message": "Redirect target reached"}`))
			default:
				w.WriteHeader(404)
				w.Write([]byte(`{"error": "Not Found"}`))
			}
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/redirect")
		require.NoError(t, err)
		// Note: Current implementation follows redirects, so we should get the target
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(resp.Body), "Redirect target reached")

		t.Log("✅ Redirect handling tested")
	})

	t.Run("DNS resolution failure", func(t *testing.T) {
		ctx, err := createCLIContextWithServer("http://nonexistent.invalid.domain.test")
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()
		resp, err := httpClient.Get(testCtx, "/test")

		// Should fail with DNS resolution error
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t,
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "lookup") ||
			strings.Contains(err.Error(), "dns"),
			"Error should contain DNS-related message")

		t.Log("✅ DNS resolution failure tested")
	})

	t.Run("Concurrent request error handling", func(t *testing.T) {
		errorCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if errorCount < 3 {
				errorCount++
				w.WriteHeader(503)
				w.Write([]byte(`{"error": "Service Unavailable"}`))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "success"}`))
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()

		// Make concurrent requests
		results := make(chan error, 5)
		for i := 0; i < 5; i++ {
			go func() {
				_, err := httpClient.Get(testCtx, "/test")
				results <- err
			}()
		}

		// Collect results
		errorResults := 0
		successResults := 0
		for i := 0; i < 5; i++ {
			if <-results != nil {
				errorResults++
			} else {
				successResults++
			}
		}

		// Some should fail, some should succeed
		assert.True(t, errorResults > 0, "Some requests should fail")
		assert.True(t, successResults > 0, "Some requests should succeed")

		t.Log("✅ Concurrent request error handling tested")
	})

	t.Run("Invalid JSON response handling", func(t *testing.T) {
		testCases := []struct {
			name         string
			responseBody string
			description  string
		}{
			{
				name:         "Truncated JSON",
				responseBody: `{"message": "incomplete json`,
				description:  "JSON with missing closing brace",
			},
			{
				name:         "Invalid JSON syntax",
				responseBody: `{"message": "invalid": "syntax"}`,
				description:  "JSON with invalid syntax",
			},
			{
				name:         "Non-JSON content type",
				responseBody: `<html><body>Error page</body></html>`,
				description:  "HTML response instead of JSON",
			},
			{
				name:         "Empty response",
				responseBody: ``,
				description:  "Completely empty response",
			},
			{
				name:         "Whitespace only",
				responseBody: `   \n\t  `,
				description:  "Response with only whitespace",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(200)
					w.Write([]byte(tc.responseBody))
				}))
				defer server.Close()

				ctx, err := createCLIContextWithServer(server.URL)
				require.NoError(t, err)

				injector := di.Bootstrap(ctx)
				httpClient := di.GetHTTPClient(injector)

				testCtx := context.Background()
				resp, err := httpClient.Get(testCtx, "/test")
				require.NoError(t, err)
				assert.Equal(t, 200, resp.StatusCode)

				// Should fail to parse invalid JSON - just check that JSON parsing would fail
				var result map[string]interface{}
				parseErr := json.Unmarshal(resp.Body, &result)
				assert.Error(t, parseErr, tc.description)

				t.Logf("✅ %s: %s", tc.name, tc.description)
			})
		}
	})

	t.Run("Edge case HTTP methods", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				w.Write([]byte(`{"method": "GET"}`))
			case http.MethodPost:
				w.WriteHeader(201)
				w.Write([]byte(`{"method": "POST"}`))
			case http.MethodPut:
				w.WriteHeader(200)
				w.Write([]byte(`{"method": "PUT"}`))
			case http.MethodDelete:
				w.WriteHeader(204)
			case http.MethodPatch:
				w.WriteHeader(200)
				w.Write([]byte(`{"method": "PATCH"}`))
			default:
				w.WriteHeader(405)
				w.Write([]byte(`{"error": "Method Not Allowed"}`))
			}
		}))
		defer server.Close()

		ctx, err := createCLIContextWithServer(server.URL)
		require.NoError(t, err)

		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		testCtx := context.Background()

		// Test each HTTP method
		resp, err := httpClient.Get(testCtx, "/test")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		resp, err = httpClient.Post(testCtx, "/test", map[string]string{"test": "data"})
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		resp, err = httpClient.Put(testCtx, "/test", map[string]string{"test": "data"})
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		resp, err = httpClient.Delete(testCtx, "/test")
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)

		t.Log("✅ Edge case HTTP methods tested")
	})
}