package integration

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
)

// TestFileOperationsTests tests file-related operations against the real API
func TestFileOperationsTests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping file operations integration tests")
	}

	// Check if API is available before running tests
	if !isAPIAvailable("http://localhost:5055") {
		t.Skip("API not available on localhost:5055, skipping integration tests")
	}

	t.Run("Text source creation", func(t *testing.T) {
		testTextSourceCreation(t)
	})

	t.Run("Source listing", func(t *testing.T) {
		testSourceListing(t)
	})

	t.Run("Error handling", func(t *testing.T) {
		testErrorHandling(t)
	})
}

// isAPIAvailable checks if the API server is running and reachable
func isAPIAvailable(apiURL string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(apiURL + "/api/health")
	if err != nil {
		// Try a simple ping to the base URL if health endpoint doesn't exist
		resp, err = client.Get(apiURL)
		if err != nil {
			return false
		}
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}

// testTextSourceCreation tests text source creation against real API
func testTextSourceCreation(t *testing.T) {
	ctx, err := createRealAPICLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)
	
	// Configure HTTP client to not retry server errors for this test
	if retryableClient, ok := httpClient.(interface{ SetRetryConfig(services.RetryConfig) }); ok {
		config := services.DefaultRetryConfig()
		// Remove HTTP 500 from retryable statuses for this test
		var filteredStatuses []int
		for _, status := range config.RetryableStatus {
			if status != 500 {
				filteredStatuses = append(filteredStatuses, status)
			}
		}
		config.RetryableStatus = filteredStatuses
		retryableClient.SetRetryConfig(config)
	}
	
	// Try to authenticate first
	auth := di.GetAuth(injector)
	if err := auth.Authenticate(context.Background()); err != nil {
		t.Logf("Authentication failed: %v", err)
		// Continue anyway as some endpoints might not require auth
	}

	// Test if we can reach the API with a simple GET request first
	testResp, err := httpClient.Get(context.Background(), "/sources")
	if err != nil {
		t.Logf("GET /sources failed: %v", err)
	} else {
		t.Logf("GET /sources successful: status %d", testResp.StatusCode)
	}

	testContent := "This is a comprehensive integration test document for text source creation. It contains multiple sentences and enough content to validate that the text processing functionality works correctly with the real API."

	sourceCreate := &models.SourceCreate{
		Type:            models.SourceTypeText,
		Content:         &testContent,
		Title:           StringPtr("Integration Test Document"),
		Embed:           true,
		DeleteSource:    false,
		AsyncProcessing: true, // Use async processing to avoid server-side asyncio issue
	}

	start := time.Now()
	resp, err := httpClient.Post(context.Background(), "/sources/json", sourceCreate)
	duration := time.Since(start)

	require.NoError(t, err)
	
	// Check for successful creation (200 or 201)
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		// Success case - validate the response
		var response models.Source
		err = json.Unmarshal(resp.Body, &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.ID, "Source ID should not be empty")
		if response.Title != nil {
			assert.Equal(t, "Integration Test Document", *response.Title)
		}
		assert.NotEmpty(t, response.Created, "Created timestamp should not be empty")

		sourceID := "unknown"
		if response.ID != nil {
			sourceID = *response.ID
		}
		
		// For async processing, check if we have command_id indicating background processing
		if response.CommandID != nil {
			t.Logf("✅ Text source queued for async processing in %v: %s (command: %s)", duration, sourceID, *response.CommandID)
		} else {
			t.Logf("✅ Text source created successfully in %v: %s", duration, sourceID)
		}

		// Optional: Clean up by deleting the created source
		if response.ID != nil {
			deleteResp, err := httpClient.Delete(context.Background(), "/sources/"+*response.ID)
			if err == nil && deleteResp.StatusCode < 300 {
				t.Logf("🧹 Cleaned up test source: %s", *response.ID)
			}
		}
	} else {
		// Log error for debugging
		t.Logf("Response status %d, body: %s", resp.StatusCode, string(resp.Body))
		t.Fatalf("Failed to create source - status: %d, response: %s", resp.StatusCode, string(resp.Body))
	}
}

// testSourceListing tests listing sources from real API
func testSourceListing(t *testing.T) {
	ctx, err := createRealAPICLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	start := time.Now()
	resp, err := httpClient.Get(context.Background(), "/sources")
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Should successfully list sources")
	assert.Less(t, duration, 10*time.Second, "Request should complete within 10 seconds")

	var sources models.SourcesListResponse
	err = json.Unmarshal(resp.Body, &sources)
	require.NoError(t, err)

	t.Logf("✅ Successfully listed %d sources in %v", len(sources), duration)
}

// testErrorHandling tests error handling with real API
func testErrorHandling(t *testing.T) {
	ctx, err := createRealAPICLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)
	
	// Configure HTTP client to not retry server errors for this test
	if retryableClient, ok := httpClient.(interface{ SetRetryConfig(services.RetryConfig) }); ok {
		config := services.DefaultRetryConfig()
		// Remove HTTP 500 from retryable statuses for this test
		var filteredStatuses []int
		for _, status := range config.RetryableStatus {
			if status != 500 {
				filteredStatuses = append(filteredStatuses, status)
			}
		}
		config.RetryableStatus = filteredStatuses
		retryableClient.SetRetryConfig(config)
	}

	t.Run("Invalid source creation - missing content", func(t *testing.T) {
		sourceCreate := &models.SourceCreate{
			Type:  models.SourceTypeText,
			Title: StringPtr("Test Source Without Content"),
			Embed: true,
			// Content is intentionally nil
		}

		resp, err := httpClient.Post(context.Background(), "/sources", sourceCreate)
		require.NoError(t, err)
		
		// Should return an error status (400 Bad Request is expected)
		assert.True(t, resp.StatusCode >= 400, "Should return error status for invalid request")
		
		t.Logf("✅ Invalid source creation properly rejected with status %d", resp.StatusCode)
	})

	t.Run("Nonexistent endpoint", func(t *testing.T) {
		resp, err := httpClient.Get(context.Background(), "/nonexistent-endpoint")
		require.NoError(t, err)
		
		assert.Equal(t, 404, resp.StatusCode, "Should return 404 for nonexistent endpoint")
		
		t.Logf("✅ Nonexistent endpoint properly returns 404")
	})

	t.Run("Invalid source ID", func(t *testing.T) {
		resp, err := httpClient.Get(context.Background(), "/sources/invalid-source-id-12345")
		require.NoError(t, err)
		
		// Should return 404 or similar error for invalid source ID
		assert.True(t, resp.StatusCode >= 400, "Should return error status for invalid source ID")
		
		t.Logf("✅ Invalid source ID properly handled with status %d", resp.StatusCode)
	})
}

// Helper function to create CLI context for real API
func createRealAPICLIContext() (*cli.Context, error) {
	return createFileOpsCLIContextWithServer("http://localhost:5055")
}

// Helper function to create CLI context with custom server URL
func createFileOpsCLIContextWithServer(serverURL string) (*cli.Context, error) {
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
