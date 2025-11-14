package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"flag"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// TestFileOperationsTests tests file-related operations with proper response handling
func TestFileOperationsTests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping file operations integration tests")
	}

	t.Run("Mock file upload workflow", func(t *testing.T) {
		testMockFileUpload(t)
	})

	t.Run("Text source creation", func(t *testing.T) {
		testTextSourceCreation(t)
	})

	t.Run("File error scenarios", func(t *testing.T) {
		testFileErrorScenarios(t)
	})
}

// testMockFileUpload tests file upload with mock server
func testMockFileUpload(t *testing.T) {
	// Create mock server that accepts file uploads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/sources") {
			// Return a mock source response
			mockSource := map[string]interface{}{
				"id":    "test-source-123",
				"title": "Test Upload File",
				"embedded": true,
				"embedded_chunks": 5,
				"created": time.Now().Format(time.RFC3339),
				"updated": time.Now().Format(time.RFC3339),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(mockSource)
			return
		}

		// Default response
		w.WriteHeader(404)
	}))
	defer server.Close()

	// Test with our HTTP client
	ctx, err := createFileOpsCLIContextWithServer(server.URL)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	// Create source request
	sourceCreate := &models.SourceCreate{
		Type:      models.SourceTypeText,
		Content:   stringPtr("This is test content for file operations testing"),
		Title:     stringPtr("Test Upload File"),
		Embed:     true,
	}

	body, err := json.Marshal(sourceCreate)
	require.NoError(t, err)

	resp, err := httpClient.Post(context.Background(), "/sources", bytes.NewReader(body))
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode, "Should successfully create source")

	// Parse response using our helper
	var sourceResult map[string]interface{}
	err = json.Unmarshal(resp.Body, &sourceResult)
	require.NoError(t, err)

	assert.Equal(t, "test-source-123", sourceResult["id"])
	assert.Equal(t, "Test Upload File", sourceResult["title"])

	t.Logf("✅ Mock file upload completed successfully: %v", sourceResult["id"])
}

// testTextSourceCreation tests text source creation
func testTextSourceCreation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/sources" {
			// Parse request to validate
			var sourceCreate models.SourceCreate
			err := json.NewDecoder(r.Body).Decode(&sourceCreate)
			if err != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
				return
			}

			// Validate content
			if sourceCreate.Content == nil || *sourceCreate.Content == "" {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]string{"error": "Content required"})
				return
			}

			// Return success response
			response := map[string]interface{}{
				"id":    "text-source-" + fmt.Sprintf("%d", time.Now().Unix()),
				"title": sourceCreate.Title,
				"created": time.Now().Format(time.RFC3339),
				"embedded": sourceCreate.Embed,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	ctx, err := createFileOpsCLIContextWithServer(server.URL)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	testContent := "This is a comprehensive test document for text source creation testing. It contains multiple sentences and enough content to validate that the text processing functionality works correctly."

	sourceCreate := &models.SourceCreate{
		Type:      models.SourceTypeText,
		Content:   &testContent,
		Title:     stringPtr("Comprehensive Test Document"),
		Embed:     true,
		DeleteSource: false,
		AsyncProcessing: false,
	}

	body, err := json.Marshal(sourceCreate)
	require.NoError(t, err)

	start := time.Now()
	resp, err := httpClient.Post(context.Background(), "/api/sources", bytes.NewReader(body))
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Less(t, duration, 5*time.Second, "Request should complete within 5 seconds")

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body, &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "Comprehensive Test Document", response["title"])
	assert.Equal(t, true, response["embedded"])

	t.Logf("✅ Text source created successfully in %v: %s", duration, response["id"])
}

// testFileErrorScenarios tests file operation error scenarios
func testFileErrorScenarios(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/sources" {
			// Read request body
			var sourceCreate models.SourceCreate
			err := json.NewDecoder(r.Body).Decode(&sourceCreate)
			if err != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Bad Request",
					"message": "Invalid JSON format",
				})
				return
			}

			// Simulate various error scenarios based on content
			if sourceCreate.Content != nil && strings.Contains(*sourceCreate.Content, "trigger-error") {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal Server Error",
					"message": "Simulated processing error",
				})
				return
			}

			if sourceCreate.Title != nil && strings.Contains(*sourceCreate.Title, "large-file") {
				w.WriteHeader(413)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Payload Too Large",
					"message": "File size exceeds limit",
				})
				return
			}

			// Default success
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "success-id",
				"status": "created",
			})
		} else {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Not Found",
				"message": "Endpoint not found",
			})
		}
	}))
	defer server.Close()

	ctx, err := createFileOpsCLIContextWithServer(server.URL)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	t.Run("Invalid JSON error", func(t *testing.T) {
		invalidJSON := []byte(`{"invalid": json content}`)
		resp, err := httpClient.Post(context.Background(), "/api/sources", bytes.NewReader(invalidJSON))

		// Should succeed with error response
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.Unmarshal(resp.Body, &errorResp)
		require.NoError(t, err)

		assert.Equal(t, "Bad Request", errorResp.Error)
		assert.Contains(t, errorResp.Message, "Invalid JSON")

		t.Logf("✅ Invalid JSON error handled correctly")
	})

	t.Run("Processing error", func(t *testing.T) {
		sourceCreate := &models.SourceCreate{
			Type:    models.SourceTypeText,
			Content: stringPtr("This content will trigger-error in processing"),
			Title:   stringPtr("Error Test Source"),
		}

		body, err := json.Marshal(sourceCreate)
		require.NoError(t, err)

		resp, err := httpClient.Post(context.Background(), "/api/sources", bytes.NewReader(body))
		require.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.Unmarshal(resp.Body, &errorResp)
		require.NoError(t, err)

		assert.Equal(t, "Internal Server Error", errorResp.Error)
		assert.Contains(t, errorResp.Message, "processing error")

		t.Logf("✅ Processing error handled correctly")
	})

	t.Run("File size limit error", func(t *testing.T) {
		sourceCreate := &models.SourceCreate{
			Type:    models.SourceTypeText,
			Content: stringPtr("Test content for large-file simulation"),
			Title:   stringPtr("Test with large-file trigger"),
		}

		body, err := json.Marshal(sourceCreate)
		require.NoError(t, err)

		resp, err := httpClient.Post(context.Background(), "/api/sources", bytes.NewReader(body))
		require.NoError(t, err)
		assert.Equal(t, 413, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.Unmarshal(resp.Body, &errorResp)
		require.NoError(t, err)

		assert.Equal(t, "Payload Too Large", errorResp.Error)

		t.Logf("✅ File size limit error handled correctly")
	})

	t.Run("Endpoint not found", func(t *testing.T) {
		sourceCreate := &models.SourceCreate{
			Type:    models.SourceTypeText,
			Content: stringPtr("Test content"),
		}

		body, err := json.Marshal(sourceCreate)
		require.NoError(t, err)

		resp, err := httpClient.Post(context.Background(), "/api/nonexistent", bytes.NewReader(body))
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)

		var errorResp models.ErrorResponse
		err = json.Unmarshal(resp.Body, &errorResp)
		require.NoError(t, err)

		assert.Equal(t, "Not Found", errorResp.Error)

		t.Logf("✅ Endpoint not found error handled correctly")
	})
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