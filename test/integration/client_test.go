package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// TestModelTypes validates our type-safe enums
func TestModelTypes(t *testing.T) {
	// This test doesn't require API connection

	t.Run("SearchType enum", func(t *testing.T) {
		assert.Equal(t, "vector", string(models.SearchTypeVector))
		assert.Equal(t, "text", string(models.SearchTypeText))
	})

	t.Run("ModelType enum", func(t *testing.T) {
		assert.Equal(t, "language", string(models.ModelTypeLanguage))
		assert.Equal(t, "embedding", string(models.ModelTypeEmbedding))
		assert.Equal(t, "text_to_speech", string(models.ModelTypeTextToSpeech))
		assert.Equal(t, "speech_to_text", string(models.ModelTypeSpeechToText))
	})

	t.Run("SourceType enum", func(t *testing.T) {
		assert.Equal(t, "link", string(models.SourceTypeLink))
		assert.Equal(t, "upload", string(models.SourceTypeUpload))
		assert.Equal(t, "text", string(models.SourceTypeText))
	})

	t.Run("NoteType enum", func(t *testing.T) {
		assert.Equal(t, "human", string(models.NoteTypeHuman))
		assert.Equal(t, "ai", string(models.NoteTypeAI))
	})

	t.Run("RebuildMode enum", func(t *testing.T) {
		assert.Equal(t, "existing", string(models.RebuildModeExisting))
		assert.Equal(t, "all", string(models.RebuildModeAll))
	})

	t.Run("RebuildStatus enum", func(t *testing.T) {
		assert.Equal(t, "queued", string(models.RebuildStatusQueued))
		assert.Equal(t, "running", string(models.RebuildStatusRunning))
		assert.Equal(t, "completed", string(models.RebuildStatusCompleted))
		assert.Equal(t, "failed", string(models.RebuildStatusFailed))
	})

	t.Run("YesNoDecision enum", func(t *testing.T) {
		assert.Equal(t, "yes", string(models.YesNoDecisionYes))
		assert.Equal(t, "no", string(models.YesNoDecisionNo))
	})

	t.Run("ContentProcessingEngine enum", func(t *testing.T) {
		assert.Equal(t, "auto", string(models.ContentProcessingEngineAuto))
		assert.Equal(t, "docling", string(models.ContentProcessingEngineDocling))
		assert.Equal(t, "simple", string(models.ContentProcessingEngineSimple))
	})

	t.Run("ContentProcessingEngineURL enum", func(t *testing.T) {
		assert.Equal(t, "auto", string(models.ContentProcessingEngineURLAuto))
		assert.Equal(t, "firecrawl", string(models.ContentProcessingEngineURLFirecrawl))
		assert.Equal(t, "jina", string(models.ContentProcessingEngineURLJina))
		assert.Equal(t, "simple", string(models.ContentProcessingEngineURLSimple))
	})

	t.Run("EmbeddingOption enum", func(t *testing.T) {
		assert.Equal(t, "ask", string(models.EmbeddingOptionAsk))
		assert.Equal(t, "always", string(models.EmbeddingOptionAlways))
		assert.Equal(t, "never", string(models.EmbeddingOptionNever))
	})

	t.Run("ItemType enum", func(t *testing.T) {
		assert.Equal(t, "source", string(models.ItemTypeSource))
		assert.Equal(t, "note", string(models.ItemTypeNote))
	})

	t.Run("SourceStatus enum", func(t *testing.T) {
		assert.Equal(t, "pending", string(models.SourceStatusPending))
		assert.Equal(t, "running", string(models.SourceStatusRunning))
		assert.Equal(t, "completed", string(models.SourceStatusCompleted))
		assert.Equal(t, "failed", string(models.SourceStatusFailed))
	})

	t.Run("InsightType enum", func(t *testing.T) {
		assert.Equal(t, "summary", string(models.InsightTypeSummary))
		assert.Equal(t, "analysis", string(models.InsightTypeAnalysis))
		assert.Equal(t, "extraction", string(models.InsightTypeExtraction))
		assert.Equal(t, "question", string(models.InsightTypeQuestion))
		assert.Equal(t, "reflection", string(models.InsightTypeReflection))
	})

	t.Run("ContextLevel enum", func(t *testing.T) {
		assert.Equal(t, "low", string(models.ContextLevelLow))
		assert.Equal(t, "medium", string(models.ContextLevelMedium))
		assert.Equal(t, "high", string(models.ContextLevelHigh))
		assert.Equal(t, "critical", string(models.ContextLevelCritical))
	})
}

// TestNotebookModels tests notebook model structure
func TestNotebookModels(t *testing.T) {
	t.Run("NotebookCreate", func(t *testing.T) {
		nb := &models.NotebookCreate{
			Name:        "Test Notebook",
			Description: "A test notebook",
		}

		assert.Equal(t, "Test Notebook", nb.Name)
		assert.Equal(t, "A test notebook", nb.Description)
	})

	t.Run("Notebook", func(t *testing.T) {
		nb := &models.Notebook{
			ID:          "nb-123",
			Name:        "Test Notebook",
			Description: "A test notebook",
			Archived:    false,
			Created:     "2024-01-01T00:00:00Z",
			Updated:     "2024-01-01T00:00:00Z",
			SourceCount: 0,
			NoteCount:   0,
		}

		assert.Equal(t, "nb-123", nb.ID)
		assert.Equal(t, "Test Notebook", nb.Name)
		assert.Equal(t, "A test notebook", nb.Description)
		assert.False(t, nb.Archived)
		assert.Equal(t, 0, nb.SourceCount)
		assert.Equal(t, 0, nb.NoteCount)
	})

	t.Run("NotebookUpdate", func(t *testing.T) {
		name := "Updated Name"
		desc := "Updated Description"
		archived := true

		nb := &models.NotebookUpdate{
			Name:        &name,
			Description: &desc,
			Archived:    &archived,
		}

		assert.Equal(t, "Updated Name", *nb.Name)
		assert.Equal(t, "Updated Description", *nb.Description)
		assert.Equal(t, true, *nb.Archived)
	})
}

// TestSearchModels tests search model structure
func TestSearchModels(t *testing.T) {
	t.Run("SearchRequest", func(t *testing.T) {
		req := &models.SearchRequest{
			Query:         "test query",
			Type:          models.SearchTypeVector,
			Limit:         10,
			SearchSources: true,
			SearchNotes:   false,
			MinimumScore:  0.5,
		}

		assert.Equal(t, "test query", req.Query)
		assert.Equal(t, models.SearchTypeVector, req.Type)
		assert.Equal(t, 10, req.Limit)
		assert.True(t, req.SearchSources)
		assert.False(t, req.SearchNotes)
		assert.Equal(t, 0.5, req.MinimumScore)
	})

	t.Run("AskRequest", func(t *testing.T) {
		req := &models.AskRequest{
			Question:         "What is the meaning of life?",
			StrategyModel:    "gpt-4",
			AnswerModel:      "gpt-4",
			FinalAnswerModel: "gpt-4",
		}

		assert.Equal(t, "What is the meaning of life?", req.Question)
		assert.Equal(t, "gpt-4", req.StrategyModel)
		assert.Equal(t, "gpt-4", req.AnswerModel)
		assert.Equal(t, "gpt-4", req.FinalAnswerModel)
	})
}

// TestNoteModels tests note model structure
func TestNoteModels(t *testing.T) {
	t.Run("NoteCreate", func(t *testing.T) {
		noteType := models.NoteTypeHuman
		notebookID := "nb-123"

		note := &models.NoteCreate{
			Content:    "This is a test note",
			NoteType:   &noteType,
			NotebookID: &notebookID,
		}

		assert.Equal(t, "This is a test note", note.Content)
		assert.Equal(t, models.NoteTypeHuman, *note.NoteType)
		assert.Equal(t, "nb-123", *note.NotebookID)
	})

	t.Run("Note", func(t *testing.T) {
		noteType := models.NoteTypeAI

		note := &models.Note{
			ID:       nil, // Optional field
			Title:    nil, // Optional field
			Content:  nil, // Optional field
			NoteType: &noteType,
			Created:  "2024-01-01T00:00:00Z",
			Updated:  "2024-01-01T00:00:00Z",
		}

		assert.Equal(t, models.NoteTypeAI, *note.NoteType)
		assert.Equal(t, "2024-01-01T00:00:00Z", note.Created)
		assert.Equal(t, "2024-01-01T00:00:00Z", note.Updated)
	})
}

// TestSourceModels tests source model structure
func TestSourceModels(t *testing.T) {
	t.Run("SourceCreate", func(t *testing.T) {
		notebooks := []string{"nb-123", "nb-456"}

		source := &models.SourceCreate{
			NotebookID:      nil, // deprecated
			Notebooks:       notebooks,
			Type:            models.SourceTypeLink,
			URL:             stringPtr("https://example.com"),
			FilePath:        nil,
			Content:         nil,
			Title:           stringPtr("Example Source"),
			Transformations: []string{},
			Embed:           true,
			DeleteSource:    false,
			AsyncProcessing: false,
		}

		assert.Equal(t, notebooks, source.Notebooks)
		assert.Equal(t, models.SourceTypeLink, source.Type)
		assert.Equal(t, "https://example.com", *source.URL)
		assert.Equal(t, "Example Source", *source.Title)
		assert.True(t, source.Embed)
	})
}

// TestModelValidation tests model validation
func TestModelValidation(t *testing.T) {
	t.Run("Notebook validation", func(t *testing.T) {
		// Empty name should be handled by service layer, not model
		nb := &models.NotebookCreate{
			Name:        "",
			Description: "Test description",
		}
		assert.Equal(t, "", nb.Name) // Model doesn't validate, service does
	})

	t.Run("Search request validation", func(t *testing.T) {
		req := &models.SearchRequest{
			Query: "",
			Type:  models.SearchTypeVector,
			Limit: 0,
		}
		assert.Equal(t, "", req.Query) // Model doesn't validate, service does
		assert.Equal(t, models.SearchTypeVector, req.Type)
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

// TestAPIConnectivityLive tests against live API (if available)
func TestAPIConnectivityLive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping live API test")
	}

	// This test requires the API to be running
	// For now, we'll just create the test structure
	// In a real environment, this would make HTTP requests

	t.Run("Notebook endpoint test", func(t *testing.T) {
		// TODO: Implement actual HTTP request test
		// This would use our HTTP client to test against localhost:5055
		t.Skip("Live API testing requires running OpenNotebook instance")
	})

	t.Run("Authentication test", func(t *testing.T) {
		// TODO: Implement auth flow test
		t.Skip("Live API testing requires running OpenNotebook instance")
	})
}

// TestHTTPResponseStructure tests HTTP response model
func TestHTTPResponseStructure(t *testing.T) {
	t.Run("Response model", func(t *testing.T) {
		resp := &models.Response{
			StatusCode: 200,
			Body:       []byte(`{"message": "success"}`),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, []byte(`{"message": "success"}`), resp.Body)
		assert.Equal(t, "application/json", resp.Header["Content-Type"][0])
	})
}

// TestTimeoutAndCancellation tests context cancellation
func TestTimeoutAndCancellation(t *testing.T) {
	t.Run("Context timeout", func(t *testing.T) {
		// Create a context with very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait a bit to ensure timeout
		time.Sleep(2 * time.Millisecond)

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			assert.Equal(t, context.DeadlineExceeded, ctx.Err())
		default:
			t.Error("Context should have been cancelled due to timeout")
		}
	})
}
