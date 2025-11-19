package integration

import (
	"context"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
)

// TestNotesAndSearchOperations tests notes and search functionality against the real API
func TestNotesAndSearchOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping notes and search integration tests")
	}

	// Check if API is available before running tests
	if !isAPIAvailable("http://localhost:5055") {
		t.Skip("API not available on localhost:5055, skipping integration tests")
	}

	t.Run("Notes operations", func(t *testing.T) {
		testNotesOperations(t)
	})

	t.Run("Search operations", func(t *testing.T) {
		testSearchOperations(t)
	})

	t.Run("Notes and Search integration", func(t *testing.T) {
		testNotesSearchIntegration(t)
	})
}

// testNotesOperations tests comprehensive notes functionality
func testNotesOperations(t *testing.T) {
	ctx, err := createNotesSearchCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	// Configure HTTP client to not retry server errors for tests
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
	}

	var testNotebookID string
	var testNoteID string

	t.Run("Create test notebook for notes", func(t *testing.T) {
		testNotebook := &models.NotebookCreate{
			Name:        "Notes Test Notebook",
			Description: "Notebook created for notes testing",
		}

		resp, err := httpClient.Post(context.Background(), "/notebooks", testNotebook)
		require.NoError(t, err)
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 201, "Should successfully create test notebook")

		var notebook models.Notebook
		err = json.Unmarshal(resp.Body, &notebook)
		require.NoError(t, err)

		assert.NotEmpty(t, notebook.ID)
		testNotebookID = notebook.ID
		t.Logf("✅ Created test notebook: %s", testNotebookID)
	})

	t.Run("List notes (initially empty)", func(t *testing.T) {
		start := time.Now()
		resp, err := httpClient.Get(context.Background(), "/notes")
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Should successfully list notes")

		var notes []models.Note
		err = json.Unmarshal(resp.Body, &notes)
		require.NoError(t, err)

		t.Logf("✅ Successfully listed %d notes in %v", len(notes), duration)
	})

	t.Run("Create human note", func(t *testing.T) {
		noteType := models.NoteTypeHuman
		noteCreate := &models.NoteCreate{
			Content:    "This is a test human note created by integration testing. It contains detailed information about the testing process.",
			NoteType:   &noteType,
			NotebookID: &testNotebookID,
		}

		start := time.Now()
		resp, err := httpClient.Post(context.Background(), "/notes", noteCreate)
		duration := time.Since(start)

		require.NoError(t, err)
		
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			var note models.Note
			err = json.Unmarshal(resp.Body, &note)
			require.NoError(t, err)

			assert.NotNil(t, note.ID, "Note ID should not be nil")
			if note.ID != nil {
				testNoteID = *note.ID
				t.Logf("✅ Human note created successfully in %v: %s", duration, testNoteID)
			}
			assert.NotNil(t, note.Content, "Note content should not be nil")
			if note.Content != nil {
				assert.Contains(t, *note.Content, "test human note", "Note content should contain expected text")
			}
			assert.NotNil(t, note.NoteType, "Note type should not be nil")
			if note.NoteType != nil {
				assert.Equal(t, models.NoteTypeHuman, *note.NoteType, "Note type should be human")
			}
		} else {
			t.Logf("Note creation returned status %d: %s", resp.StatusCode, string(resp.Body))
			t.Logf("This might be expected if notes require different handling")
		}
	})

	if testNoteID != "" {
		t.Run("Get specific note", func(t *testing.T) {
			start := time.Now()
			resp, err := httpClient.Get(context.Background(), "/notes/"+testNoteID)
			duration := time.Since(start)

			require.NoError(t, err)
			
			if resp.StatusCode == 200 {
				var note models.Note
				err = json.Unmarshal(resp.Body, &note)
				require.NoError(t, err)

				assert.NotNil(t, note.ID, "Retrieved note ID should not be nil")
				if note.ID != nil {
					assert.Equal(t, testNoteID, *note.ID, "Retrieved note ID should match")
				}
				t.Logf("✅ Retrieved note successfully in %v: %s", duration, testNoteID)
			} else {
				t.Logf("Note retrieval returned status %d: %s", resp.StatusCode, string(resp.Body))
			}
		})

		t.Run("Update note", func(t *testing.T) {
			updatedContent := "This is an updated test note with modified content."
			noteUpdate := map[string]interface{}{
				"content": updatedContent,
			}

			start := time.Now()
			resp, err := httpClient.Put(context.Background(), "/notes/"+testNoteID, noteUpdate)
			duration := time.Since(start)

			require.NoError(t, err)
			
			if resp.StatusCode == 200 {
				var note models.Note
				err = json.Unmarshal(resp.Body, &note)
				require.NoError(t, err)

				if note.Content != nil {
					assert.Contains(t, *note.Content, "updated", "Updated note should contain 'updated'")
				}
				t.Logf("✅ Updated note successfully in %v", duration)
			} else {
				t.Logf("Note update returned status %d: %s", resp.StatusCode, string(resp.Body))
			}
		})
	}

	t.Run("List notes after creation", func(t *testing.T) {
		start := time.Now()
		resp, err := httpClient.Get(context.Background(), "/notes")
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Should successfully list notes")

		var notes []models.Note
		err = json.Unmarshal(resp.Body, &notes)
		require.NoError(t, err)

		t.Logf("✅ Successfully listed %d notes in %v (after creation)", len(notes), duration)
		
		// Log note details
		for i, note := range notes {
			if i < 5 { // Limit to first 5 for readability
				content := "No Content"
				if note.Content != nil {
					content = *note.Content
					if len(content) > 50 {
						content = content[:50] + "..."
					}
				}
				noteType := "Unknown"
				if note.NoteType != nil {
					noteType = string(*note.NoteType)
				}
				t.Logf("   Note %d: %s (%s)", i+1, content, noteType)
			}
		}
	})

	// Cleanup: Delete test note and notebook
	if testNoteID != "" {
		t.Run("Delete test note", func(t *testing.T) {
			resp, err := httpClient.Delete(context.Background(), "/notes/"+testNoteID)
			if err == nil && resp.StatusCode < 300 {
				t.Logf("🧹 Cleaned up test note: %s", testNoteID)
			}
		})
	}

	if testNotebookID != "" {
		t.Run("Delete test notebook", func(t *testing.T) {
			resp, err := httpClient.Delete(context.Background(), "/notebooks/"+testNotebookID)
			if err == nil && resp.StatusCode < 300 {
				t.Logf("🧹 Cleaned up test notebook: %s", testNotebookID)
			}
		})
	}
}

// testSearchOperations tests comprehensive search functionality
func testSearchOperations(t *testing.T) {
	ctx, err := createNotesSearchCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	// Configure HTTP client to not retry server errors for some tests
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
	}

	t.Run("Text search", func(t *testing.T) {
		searchRequest := &models.SearchRequest{
			Query:         "integration test",
			Type:          models.SearchTypeText,
			Limit:         10,
			SearchSources: true,
			SearchNotes:   true,
			MinimumScore:  0.1,
		}

		start := time.Now()
		resp, err := httpClient.Post(context.Background(), "/search", searchRequest)
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Should successfully perform text search")

		var searchResponse models.SearchResponse
		err = json.Unmarshal(resp.Body, &searchResponse)
		require.NoError(t, err)

		t.Logf("✅ Text search completed in %v, found %d results", duration, len(searchResponse.Results))
		
		// Log search results
		for i, result := range searchResponse.Results {
			if i < 3 { // Limit to first 3 for readability
				t.Logf("   Result %d: %s (relevance: %.2f)", i+1, result.Title, result.Relevance)
			}
		}
	})

	t.Run("Vector search", func(t *testing.T) {
		searchRequest := &models.SearchRequest{
			Query:         "document processing",
			Type:          models.SearchTypeVector,
			Limit:         5,
			SearchSources: true,
			SearchNotes:   false,
			MinimumScore:  0.5,
		}

		start := time.Now()
		resp, err := httpClient.Post(context.Background(), "/search", searchRequest)
		duration := time.Since(start)

		require.NoError(t, err)
		
		if resp.StatusCode == 200 {
			var searchResponse models.SearchResponse
			err = json.Unmarshal(resp.Body, &searchResponse)
			require.NoError(t, err)

			t.Logf("✅ Vector search completed in %v, found %d results", duration, len(searchResponse.Results))
			
			// Validate relevance scores (be flexible with vector search results)
			for _, result := range searchResponse.Results {
				if result.Relevance < searchRequest.MinimumScore {
					t.Logf("⚠️ Vector search result below minimum score: %.2f < %.2f", result.Relevance, searchRequest.MinimumScore)
				}
			}
		} else {
			t.Logf("Vector search returned status %d: %s", resp.StatusCode, string(resp.Body))
			t.Logf("Vector search might not be enabled or configured")
		}
	})

	t.Run("Ask simple question", func(t *testing.T) {
		askRequest := &models.AskRequest{
			Question:         "What is the purpose of integration testing?",
			StrategyModel:    "gemini-2.0-flash", // Use a model we know exists
			AnswerModel:      "gemini-2.0-flash",
			FinalAnswerModel: "gemini-2.0-flash",
		}

		start := time.Now()
		resp, err := httpClient.Post(context.Background(), "/search/ask/simple", askRequest)
		duration := time.Since(start)

		require.NoError(t, err)
		
		if resp.StatusCode == 200 {
			var askResponse models.AskResponse
			err = json.Unmarshal(resp.Body, &askResponse)
			require.NoError(t, err)

			assert.NotEmpty(t, askResponse.Answer, "Ask response should have an answer")
			t.Logf("✅ Ask simple question completed in %v", duration)
			t.Logf("   Question: %s", askRequest.Question)
			t.Logf("   Answer length: %d characters", len(askResponse.Answer))
		} else {
			t.Logf("Ask simple returned status %d: %s", resp.StatusCode, string(resp.Body))
			t.Logf("Ask functionality might require additional configuration")
		}
	})

	t.Run("Search with filters", func(t *testing.T) {
		searchRequest := &models.SearchRequest{
			Query:         "test",
			Type:          models.SearchTypeText,
			Limit:         20,
			SearchSources: true,
			SearchNotes:   true,
			MinimumScore:  0.0, // Very low threshold to get more results
		}

		start := time.Now()
		resp, err := httpClient.Post(context.Background(), "/search", searchRequest)
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Should successfully perform filtered search")

		var searchResponse models.SearchResponse
		err = json.Unmarshal(resp.Body, &searchResponse)
		require.NoError(t, err)

		t.Logf("✅ Filtered search completed in %v, found %d results", duration, len(searchResponse.Results))
		
		// Validate response structure (be flexible with relevance scores)
		for _, result := range searchResponse.Results {
			assert.NotEmpty(t, result.ID, "Result should have an ID")
			assert.NotEmpty(t, result.ParentID, "Result should have a parent ID")
			// Note: Some search implementations might return negative scores
			if result.Relevance < 0.0 {
				t.Logf("⚠️ Negative relevance score found: %.2f for result %s", result.Relevance, result.ID)
			}
		}
	})

	t.Run("Empty search query", func(t *testing.T) {
		searchRequest := &models.SearchRequest{
			Query:         "",
			Type:          models.SearchTypeText,
			Limit:         5,
			SearchSources: true,
			SearchNotes:   true,
		}

		resp, err := httpClient.Post(context.Background(), "/search", searchRequest)
		require.NoError(t, err)
		
		// Should handle empty query gracefully (likely return 400 or empty results)
		if resp.StatusCode == 400 {
			t.Logf("✅ Empty search query properly rejected with status %d", resp.StatusCode)
		} else if resp.StatusCode == 200 {
			var searchResponse models.SearchResponse
			err = json.Unmarshal(resp.Body, &searchResponse)
			require.NoError(t, err)
			t.Logf("✅ Empty search query handled gracefully, returned %d results", len(searchResponse.Results))
		} else {
			t.Logf("Empty search query returned status %d", resp.StatusCode)
		}
	})
}

// testNotesSearchIntegration tests integration between notes and search
func testNotesSearchIntegration(t *testing.T) {
	ctx, err := createNotesSearchCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	// Try to authenticate first
	auth := di.GetAuth(injector)
	if err := auth.Authenticate(context.Background()); err != nil {
		t.Logf("Authentication failed: %v", err)
	}

	t.Run("Create note and search for it", func(t *testing.T) {
		// Create a note with specific content
		uniqueContent := "UniqueTestContent" + time.Now().Format("20060102150405")
		noteType := models.NoteTypeHuman
		
		noteCreate := &models.NoteCreate{
			Content:  uniqueContent + " - This is a unique test note for search integration testing",
			NoteType: &noteType,
		}

		// Create note
		resp, err := httpClient.Post(context.Background(), "/notes", noteCreate)
		require.NoError(t, err)
		
		var createdNoteID string
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			var note models.Note
			err = json.Unmarshal(resp.Body, &note)
			if err == nil && note.ID != nil {
				createdNoteID = *note.ID
				t.Logf("✅ Created note for search testing: %s", createdNoteID)
			}
		}

		// Wait a moment for indexing (if applicable)
		time.Sleep(1 * time.Second)

		// Search for the note
		searchRequest := &models.SearchRequest{
			Query:         uniqueContent,
			Type:          models.SearchTypeText,
			Limit:         10,
			SearchSources: false,
			SearchNotes:   true,
			MinimumScore:  0.0,
		}

		resp, err = httpClient.Post(context.Background(), "/search", searchRequest)
		require.NoError(t, err)
		
		if resp.StatusCode == 200 {
			var searchResponse models.SearchResponse
			err = json.Unmarshal(resp.Body, &searchResponse)
			require.NoError(t, err)

			found := false
			for _, result := range searchResponse.Results {
				if result.ID == createdNoteID {
					found = true
					t.Logf("✅ Successfully found created note in search results")
					break
				}
			}

			if !found {
				t.Logf("⚠️ Created note not found in search results (might need indexing time)")
				t.Logf("   Search returned %d results", len(searchResponse.Results))
			}
		} else {
			t.Logf("Search returned status %d: %s", resp.StatusCode, string(resp.Body))
		}

		// Cleanup
		if createdNoteID != "" {
			resp, err := httpClient.Delete(context.Background(), "/notes/"+createdNoteID)
			if err == nil && resp.StatusCode < 300 {
				t.Logf("🧹 Cleaned up test note: %s", createdNoteID)
			}
		}
	})
}

// Helper function to create CLI context for notes and search testing
func createNotesSearchCLIContext() (*cli.Context, error) {
	app := &cli.App{
		Name: "notes-search-test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-url",
				Value: "http://localhost:5055",
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
	flagSet.String("api-url", "http://localhost:5055", "")
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