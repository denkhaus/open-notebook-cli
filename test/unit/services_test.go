package unit

import (
	"context"
	"testing"

	"github.com/denkhaus/open-notebook-cli/pkg/mocks"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSourceService_MockOperations tests basic source service operations
func TestSourceService_MockOperations(t *testing.T) {
	// Arrange
	service := mocks.NewMockSourceService()
	ctx := context.Background()

	// Test AddSourceFromLink
	t.Run("AddSourceFromLink", func(t *testing.T) {
		link := "https://example.com/test.pdf"
		options := &models.SourceOptions{
			Embed:           true,
			DeleteSource:    false,
			AsyncProcessing: false,
		}

		// Act
		result, err := service.AddSourceFromLink(ctx, link, options)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.ID)
		assert.Equal(t, "Mock Source from Link", *result.Title)
	})

	// Test AddSourceFromText
	t.Run("AddSourceFromText", func(t *testing.T) {
		text := "Test content"
		title := "Test Document"
		options := &models.SourceOptions{
			Embed:           true,
			DeleteSource:    false,
			AsyncProcessing: false,
		}

		// Act
		result, err := service.AddSourceFromText(ctx, text, title, options)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.ID)
		assert.Equal(t, &title, result.Title)
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		// Act
		result, err := service.List(ctx, 10, 0)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// TestNotebookService_MockOperations tests basic notebook service operations
func TestNotebookService_MockOperations(t *testing.T) {
	// Arrange
	service := mocks.NewMockNotebookService()
	ctx := context.Background()

	// Add test notebook
	testNotebook := &models.Notebook{
		ID:          "test-123",
		Name:        "Test Notebook",
		Description: "A test notebook",
		Created:     "2024-01-01T00:00:00Z",
		Updated:     "2024-01-01T00:00:00Z",
	}
	service.GetRepository().AddNotebook(testNotebook)

	// Test Repository operations
	t.Run("Repository Operations", func(t *testing.T) {
		// Test that we can access the repository
		repo := service.GetRepository()
		require.NotNil(t, repo)

		// Test listing notebooks through repository
		result, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "test-123", result[0].ID)
	})
}

// TestSearchService_MockOperations tests basic search service operations
func TestSearchService_MockOperations(t *testing.T) {
	// Arrange
	service := mocks.NewMockSearchService()
	ctx := context.Background()

	// Test Repository operations
	t.Run("Repository Operations", func(t *testing.T) {
		// Test that we can access the repository
		repo := service.GetRepository()
		require.NotNil(t, repo)

		// Test search operations through repository
		req := &models.SearchRequest{
			Query: "test query",
			Type:  models.SearchTypeVector,
			Limit: 10,
		}

		result, err := repo.Search(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "vector", result.SearchType)
	})
}