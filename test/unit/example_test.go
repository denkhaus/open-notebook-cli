package unit

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/mocks"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNotebookService_MockUsage demonstrates how to use the mock notebook service
func TestNotebookService_MockUsage(t *testing.T) {
	// Create mock service with no delay for fast tests
	mockService := mocks.NewMockNotebookService()

	// Set up test data
	testNotebook := &models.Notebook{
		ID:          "test-123",
		Name:        "Test Notebook",
		Description: "A test notebook for demonstration",
		Created:     time.Now().Format(time.RFC3339),
		Updated:     time.Now().Format(time.RFC3339),
	}

	// Add test data to the mock
	mockService.GetRepository().AddNotebook(testNotebook)

	// Use the service in your test
	ctx := context.Background()
	notebooks, err := mockService.GetRepository().List(ctx)

	// Assert
	require.NoError(t, err)
	assert.Len(t, notebooks, 1)
	assert.Equal(t, "test-123", notebooks[0].ID)
	assert.Equal(t, "Test Notebook", notebooks[0].Name)

	// Check that List was called exactly once on the repository
	repo := mockService.GetRepository()
	assert.Equal(t, 1, repo.CallCount("List"))
}

// TestSourceService_MockUsage demonstrates how to use the mock source service
func TestSourceService_MockUsage(t *testing.T) {
	// Create mock service
	mockService := mocks.NewMockSourceService()

	// Add test source
	id := "source-123"
	title := "Test Source"
	status := models.SourceStatusPending
	testSource := &models.Source{
		ID:       &id,
		Title:    &title,
		Topics:   []string{"testing"},
		Embedded: false,
		Status:   &status,
		Created:  time.Now().Format(time.RFC3339),
		Updated:  time.Now().Format(time.RFC3339),
	}

	mockService.GetRepository().AddSource(testSource)

	// Test operations
	ctx := context.Background()

	// List sources
	_, _ = mockService.List(ctx, 10, 0)

	// Get source by ID
	_, _ = mockService.Get(ctx, "source-123")

	// Create new source
	_, _ = mockService.AddSourceFromText(ctx, "test content", "Test Title", nil)
}

// TestErrorInjection demonstrates how to test error scenarios
func TestErrorInjection(t *testing.T) {
	// Create mock service
	mockService := mocks.NewMockNotebookService()

	// Configure the mock to return an error on the next call
	mockService.SetError("ListNotebooks", assert.AnError)

	// Or configure all subsequent calls to fail
	mockService.SetFailure(assert.AnError)

	ctx := context.Background()

	// This call will return the configured error
	_, _ = mockService.ListNotebooks(ctx)

	// Clear failure mode
	mockService.ClearFailure()
}

// TestDelayedOperations demonstrates testing with simulated delays
func TestDelayedOperations(t *testing.T) {
	// Create mock with 100ms delay to simulate network latency
	mockService := mocks.NewMockNotebookService()

	// Use in tests that need to test timeouts or delayed responses
	ctx := context.Background()

	start := time.Now()
	_, _ = mockService.ListNotebooks(ctx)
	_ = time.Since(start)

	// You can also test timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// This will fail due to context timeout
	_, _ = mockService.ListNotebooks(ctx)
}

// TestCallVerification demonstrates how to verify method calls
func TestCallVerification(t *testing.T) {
	// Create mock service
	mockService := mocks.NewMockNotebookService()

	// Execute some operations
	ctx := context.Background()
	_, _ = mockService.GetRepository().List(ctx)

	// Verify method was called on repository
	repo := mockService.GetRepository()
	assert.Equal(t, 1, repo.CallCount("List"))

	// Get detailed call information
	calls := repo.GetCalls("List")
	assert.Len(t, calls, 1)
	assert.Equal(t, ctx, calls[0].Args[0])
}