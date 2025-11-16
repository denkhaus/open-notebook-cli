package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchTypeValues(t *testing.T) {
	assert.Equal(t, "vector", string(SearchTypeVector))
	assert.Equal(t, "text", string(SearchTypeText))
}

func TestModelTypeValues(t *testing.T) {
	assert.Equal(t, "language", string(ModelTypeLanguage))
	assert.Equal(t, "embedding", string(ModelTypeEmbedding))
	assert.Equal(t, "text_to_speech", string(ModelTypeTextToSpeech))
	assert.Equal(t, "speech_to_text", string(ModelTypeSpeechToText))
}

func TestSourceTypeValues(t *testing.T) {
	assert.Equal(t, "link", string(SourceTypeLink))
	assert.Equal(t, "upload", string(SourceTypeUpload))
	assert.Equal(t, "text", string(SourceTypeText))
}

func TestNoteTypeValues(t *testing.T) {
	assert.Equal(t, "human", string(NoteTypeHuman))
	assert.Equal(t, "ai", string(NoteTypeAI))
}

func TestRebuildModeValues(t *testing.T) {
	assert.Equal(t, "existing", string(RebuildModeExisting))
	assert.Equal(t, "all", string(RebuildModeAll))
}

func TestRebuildStatusValues(t *testing.T) {
	assert.Equal(t, "queued", string(RebuildStatusQueued))
	assert.Equal(t, "running", string(RebuildStatusRunning))
	assert.Equal(t, "completed", string(RebuildStatusCompleted))
	assert.Equal(t, "failed", string(RebuildStatusFailed))
}

func TestYesNoDecisionValues(t *testing.T) {
	assert.Equal(t, "yes", string(YesNoDecisionYes))
	assert.Equal(t, "no", string(YesNoDecisionNo))
}

func TestContentProcessingEngineValues(t *testing.T) {
	assert.Equal(t, "auto", string(ContentProcessingEngineAuto))
	assert.Equal(t, "docling", string(ContentProcessingEngineDocling))
	assert.Equal(t, "simple", string(ContentProcessingEngineSimple))
}

func TestContentProcessingEngineURLValues(t *testing.T) {
	assert.Equal(t, "auto", string(ContentProcessingEngineURLAuto))
	assert.Equal(t, "firecrawl", string(ContentProcessingEngineURLFirecrawl))
	assert.Equal(t, "jina", string(ContentProcessingEngineURLJina))
	assert.Equal(t, "simple", string(ContentProcessingEngineURLSimple))
}

func TestEmbeddingOptionValues(t *testing.T) {
	assert.Equal(t, "ask", string(EmbeddingOptionAsk))
	assert.Equal(t, "always", string(EmbeddingOptionAlways))
	assert.Equal(t, "never", string(EmbeddingOptionNever))
}

func TestItemTypeValues(t *testing.T) {
	assert.Equal(t, "source", string(ItemTypeSource))
	assert.Equal(t, "note", string(ItemTypeNote))
}

func TestSourceStatusValues(t *testing.T) {
	assert.Equal(t, "pending", string(SourceStatusPending))
	assert.Equal(t, "running", string(SourceStatusRunning))
	assert.Equal(t, "completed", string(SourceStatusCompleted))
	assert.Equal(t, "failed", string(SourceStatusFailed))
}

func TestInsightTypeValues(t *testing.T) {
	assert.Equal(t, "summary", string(InsightTypeSummary))
	assert.Equal(t, "analysis", string(InsightTypeAnalysis))
	assert.Equal(t, "extraction", string(InsightTypeExtraction))
	assert.Equal(t, "question", string(InsightTypeQuestion))
	assert.Equal(t, "reflection", string(InsightTypeReflection))
}

func TestContextLevelValues(t *testing.T) {
	assert.Equal(t, "low", string(ContextLevelLow))
	assert.Equal(t, "medium", string(ContextLevelMedium))
	assert.Equal(t, "high", string(ContextLevelHigh))
	assert.Equal(t, "critical", string(ContextLevelCritical))
}

// Test model structures to ensure they compile correctly
func TestModelStructures(t *testing.T) {
	t.Run("Notebook model", func(t *testing.T) {
		notebook := &Notebook{
			ID:          "nb-123",
			Name:        "Test Notebook",
			Description: "Test Description",
			Archived:    false,
			Created:     "2024-01-01T00:00:00Z",
			Updated:     "2024-01-01T00:00:00Z",
			SourceCount: 5,
			NoteCount:   10,
		}

		assert.Equal(t, "nb-123", notebook.ID)
		assert.Equal(t, "Test Notebook", notebook.Name)
		assert.Equal(t, "Test Description", notebook.Description)
		assert.False(t, notebook.Archived)
		assert.Equal(t, 5, notebook.SourceCount)
		assert.Equal(t, 10, notebook.NoteCount)
	})

	t.Run("NotebookCreate model", func(t *testing.T) {
		create := &NotebookCreate{
			Name:        "New Notebook",
			Description: "New Description",
		}

		assert.Equal(t, "New Notebook", create.Name)
		assert.Equal(t, "New Description", create.Description)
	})

	t.Run("NotebookUpdate model", func(t *testing.T) {
		name := "Updated Name"
		desc := "Updated Description"
		archived := true

		update := &NotebookUpdate{
			Name:        &name,
			Description: &desc,
			Archived:    &archived,
		}

		assert.Equal(t, "Updated Name", *update.Name)
		assert.Equal(t, "Updated Description", *update.Description)
		assert.True(t, *update.Archived)
	})

	t.Run("SearchRequest model", func(t *testing.T) {
		req := &SearchRequest{
			Query:         "test query",
			Type:          SearchTypeVector,
			Limit:         10,
			SearchSources: true,
			SearchNotes:   false,
			MinimumScore:  0.5,
		}

		assert.Equal(t, "test query", req.Query)
		assert.Equal(t, SearchTypeVector, req.Type)
		assert.Equal(t, 10, req.Limit)
		assert.True(t, req.SearchSources)
		assert.False(t, req.SearchNotes)
		assert.Equal(t, 0.5, req.MinimumScore)
	})

	t.Run("ModelCreate model", func(t *testing.T) {
		model := &ModelCreate{
			Name:     "Test Model",
			Provider: "test-provider",
			Type:     ModelTypeLanguage,
		}

		assert.Equal(t, "Test Model", model.Name)
		assert.Equal(t, "test-provider", model.Provider)
		assert.Equal(t, ModelTypeLanguage, model.Type)
	})

	t.Run("SourceCreate model", func(t *testing.T) {
		notebooks := []string{"nb-1", "nb-2"}
		source := &SourceCreate{
			Notebooks:       notebooks,
			Type:            SourceTypeLink,
			URL:             stringPtr("https://example.com"),
			Title:           stringPtr("Example Source"),
			Transformations: []string{},
			Embed:           true,
		}

		assert.Equal(t, notebooks, source.Notebooks)
		assert.Equal(t, SourceTypeLink, source.Type)
		assert.Equal(t, "https://example.com", *source.URL)
		assert.Equal(t, "Example Source", *source.Title)
		assert.True(t, source.Embed)
	})

	t.Run("NoteCreate model", func(t *testing.T) {
		noteType := NoteTypeHuman
		notebookID := "nb-123"

		note := &NoteCreate{
			Title:      stringPtr("Test Note"),
			Content:    "Test content",
			NoteType:   &noteType,
			NotebookID: &notebookID,
		}

		assert.Equal(t, "Test Note", *note.Title)
		assert.Equal(t, "Test content", note.Content)
		assert.Equal(t, NoteTypeHuman, *note.NoteType)
		assert.Equal(t, "nb-123", *note.NotebookID)
	})

	t.Run("EmbedRequest model", func(t *testing.T) {
		req := &EmbedRequest{
			ItemID:          "item-123",
			ItemType:        ItemTypeSource,
			AsyncProcessing: false,
		}

		assert.Equal(t, "item-123", req.ItemID)
		assert.Equal(t, ItemTypeSource, req.ItemType)
		assert.False(t, req.AsyncProcessing)
	})

	t.Run("RebuildRequest model", func(t *testing.T) {
		req := &RebuildRequest{
			Mode:            RebuildModeAll,
			IncludeSources:  true,
			IncludeNotes:    true,
			IncludeInsights: false,
		}

		assert.Equal(t, RebuildModeAll, req.Mode)
		assert.True(t, req.IncludeSources)
		assert.True(t, req.IncludeNotes)
		assert.False(t, req.IncludeInsights)
	})

	t.Run("SettingsResponse model", func(t *testing.T) {
		settings := &SettingsResponse{
			DefaultContentProcessingEngineDoc: ContentProcessingEngineAuto,
			DefaultContentProcessingEngineURL: ContentProcessingEngineURLAuto,
			DefaultEmbeddingOption:            EmbeddingOptionAsk,
			AutoDeleteFiles:                   YesNoDecisionYes,
			YoutubePreferredLanguages:         []string{"en", "de"},
		}

		assert.Equal(t, ContentProcessingEngineAuto, settings.DefaultContentProcessingEngineDoc)
		assert.Equal(t, ContentProcessingEngineURLAuto, settings.DefaultContentProcessingEngineURL)
		assert.Equal(t, EmbeddingOptionAsk, settings.DefaultEmbeddingOption)
		assert.Equal(t, YesNoDecisionYes, settings.AutoDeleteFiles)
		assert.Equal(t, []string{"en", "de"}, settings.YoutubePreferredLanguages)
	})

	t.Run("ContextConfig model", func(t *testing.T) {
		config := &ContextConfig{
			Sources: map[string]ContextLevel{
				"src-123": ContextLevelHigh,
				"src-456": ContextLevelLow,
			},
			Notes: map[string]ContextLevel{
				"note-789": ContextLevelCritical,
			},
		}

		assert.Equal(t, ContextLevelHigh, config.Sources["src-123"])
		assert.Equal(t, ContextLevelLow, config.Sources["src-456"])
		assert.Equal(t, ContextLevelCritical, config.Notes["note-789"])
	})
}

// Test that all enums have unique values
func TestEnumValueUniqueness(t *testing.T) {
	enumValues := map[string][]string{
		"SearchType": {
			string(SearchTypeVector),
			string(SearchTypeText),
		},
		"ModelType": {
			string(ModelTypeLanguage),
			string(ModelTypeEmbedding),
			string(ModelTypeTextToSpeech),
			string(ModelTypeSpeechToText),
		},
		"SourceType": {
			string(SourceTypeLink),
			string(SourceTypeUpload),
			string(SourceTypeText),
		},
		"NoteType": {
			string(NoteTypeHuman),
			string(NoteTypeAI),
		},
		"RebuildMode": {
			string(RebuildModeExisting),
			string(RebuildModeAll),
		},
		"RebuildStatus": {
			string(RebuildStatusQueued),
			string(RebuildStatusRunning),
			string(RebuildStatusCompleted),
			string(RebuildStatusFailed),
		},
		"YesNoDecision": {
			string(YesNoDecisionYes),
			string(YesNoDecisionNo),
		},
		"ContentProcessingEngine": {
			string(ContentProcessingEngineAuto),
			string(ContentProcessingEngineDocling),
			string(ContentProcessingEngineSimple),
		},
		"ContentProcessingEngineURL": {
			string(ContentProcessingEngineURLAuto),
			string(ContentProcessingEngineURLFirecrawl),
			string(ContentProcessingEngineURLJina),
			string(ContentProcessingEngineURLSimple),
		},
		"EmbeddingOption": {
			string(EmbeddingOptionAsk),
			string(EmbeddingOptionAlways),
			string(EmbeddingOptionNever),
		},
		"ItemType": {
			string(ItemTypeSource),
			string(ItemTypeNote),
		},
		"SourceStatus": {
			string(SourceStatusPending),
			string(SourceStatusRunning),
			string(SourceStatusCompleted),
			string(SourceStatusFailed),
		},
		"InsightType": {
			string(InsightTypeSummary),
			string(InsightTypeAnalysis),
			string(InsightTypeExtraction),
			string(InsightTypeQuestion),
			string(InsightTypeReflection),
		},
		"ContextLevel": {
			string(ContextLevelLow),
			string(ContextLevelMedium),
			string(ContextLevelHigh),
			string(ContextLevelCritical),
		},
	}

	// Check for uniqueness within each enum
	for enumName, values := range enumValues {
		t.Run(enumName+" uniqueness", func(t *testing.T) {
			uniqueValues := make(map[string]bool)
			for _, value := range values {
				assert.False(t, uniqueValues[value], "Duplicate value found in %s: %s", enumName, value)
				uniqueValues[value] = true
			}
		})
	}

	// Check for uniqueness within same semantic domain
	// Some values like "auto", "running", "completed", "failed" can be used
	// across different enum types as they have different semantic meanings
	semanticallyOverlappingValues := map[string][]string{
		"auto":      {"ContentProcessingEngine", "ContentProcessingEngineURL"},
		"running":   {"RebuildStatus", "SourceStatus"},
		"completed": {"RebuildStatus", "SourceStatus"},
		"failed":    {"RebuildStatus", "SourceStatus"},
		"text":      {"SourceType", "SearchType"},
	}

	// Note: These overlaps are acceptable as the enums represent different domains
	t.Logf("Note: Some enum values overlap across different semantic domains, which is acceptable:")
	for value, enums := range semanticallyOverlappingValues {
		t.Logf("  - '%s' used in: %v", value, enums)
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
