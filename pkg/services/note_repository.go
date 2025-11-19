package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

type noteRepository struct {
	httpClient shared.HTTPClient
	logger     shared.Logger
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(injector do.Injector) (shared.NoteRepository, error) {
	httpClient := do.MustInvoke[shared.HTTPClient](injector)
	logger := do.MustInvoke[shared.Logger](injector)

	return &noteRepository{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// List implements NoteRepository interface
func (n *noteRepository) List(ctx context.Context, notebookID string, limit, offset int) ([]*models.Note, error) {
	queryParams := url.Values{}
	if notebookID != "" {
		queryParams.Set("notebook_id", notebookID)
	}
	queryParams.Set("limit", strconv.Itoa(limit))
	queryParams.Set("offset", strconv.Itoa(offset))

	endpoint := "/api/notes?" + queryParams.Encode()
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	var result []*models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse notes response: %w", err)
	}

	n.logger.Info("Retrieved notes", "count", len(result))
	return result, nil
}

// Search implements NoteRepository interface
func (n *noteRepository) Search(ctx context.Context, notebookID, query string) ([]*models.Note, error) {
	queryParams := url.Values{}
	queryParams.Set("notebook_id", notebookID)
	queryParams.Set("query", query)

	endpoint := "/api/notes/search?" + queryParams.Encode()
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}

	var result []*models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	n.logger.Info("Searched notes", "notebook_id", notebookID, "query", query, "count", len(result))
	return result, nil
}

// Create implements NoteRepository interface
func (n *noteRepository) Create(ctx context.Context, note *models.NoteCreate) (*models.Note, error) {
	resp, err := n.httpClient.Post(ctx, "/api/notes", note)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	var result models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse note response: %w", err)
	}

	n.logger.Info("Created note", "id", result.ID)
	return &result, nil
}

// Get implements NoteRepository interface
func (n *noteRepository) Get(ctx context.Context, id string) (*models.Note, error) {
	endpoint := fmt.Sprintf("/api/notes/%s", id)
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get note %s: %w", id, err)
	}

	var result models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse note response: %w", err)
	}

	n.logger.Info("Retrieved note", "id", id)
	return &result, nil
}

// Update implements NoteRepository interface
func (n *noteRepository) Update(ctx context.Context, id string, note *models.NoteUpdate) (*models.Note, error) {
	endpoint := fmt.Sprintf("/api/notes/%s", id)
	resp, err := n.httpClient.Put(ctx, endpoint, note)
	if err != nil {
		return nil, fmt.Errorf("failed to update note %s: %w", id, err)
	}

	var result models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	n.logger.Info("Updated note", "id", id)
	return &result, nil
}

// Delete implements NoteRepository interface
func (n *noteRepository) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/notes/%s", id)
	_, err := n.httpClient.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete note %s: %w", id, err)
	}

	n.logger.Info("Deleted note", "id", id)
	return nil
}

// Extended methods not part of interface

// ListByNotebook filters notes by notebook
func (n *noteRepository) ListByNotebook(ctx context.Context, notebookID string, limit, offset int) ([]*models.Note, error) {
	queryParams := url.Values{}
	queryParams.Set("notebook_id", notebookID)
	queryParams.Set("limit", strconv.Itoa(limit))
	queryParams.Set("offset", strconv.Itoa(offset))

	endpoint := "/api/notes?" + queryParams.Encode()
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes by notebook %s: %w", notebookID, err)
	}

	var result []*models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse notes response: %w", err)
	}

	n.logger.Info("Retrieved notes by notebook", "notebook_id", notebookID, "count", len(result))
	return result, nil
}

// ListByType filters notes by type
func (n *noteRepository) ListByType(ctx context.Context, noteType models.NoteType, limit, offset int) ([]*models.Note, error) {
	queryParams := url.Values{}
	queryParams.Set("note_type", string(noteType))
	queryParams.Set("limit", strconv.Itoa(limit))
	queryParams.Set("offset", strconv.Itoa(offset))

	endpoint := "/api/notes?" + queryParams.Encode()
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes by type %s: %w", noteType, err)
	}

	var result []*models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse notes response: %w", err)
	}

	n.logger.Info("Retrieved notes by type", "type", noteType, "count", len(result))
	return result, nil
}

// SearchWithFilters searches notes with multiple filters
func (n *noteRepository) SearchWithFilters(ctx context.Context, filters map[string]string) ([]*models.Note, error) {
	queryParams := url.Values{}
	for key, value := range filters {
		queryParams.Set(key, value)
	}

	endpoint := "/api/notes/search?" + queryParams.Encode()
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes with filters: %w", err)
	}

	var result []*models.Note
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	n.logger.Info("Searched notes with filters", "filters", filters, "count", len(result))
	return result, nil
}

// GetNotesCount gets total count of notes
func (n *noteRepository) GetNotesCount(ctx context.Context, notebookID string) (int, error) {
	queryParams := url.Values{}
	if notebookID != "" {
		queryParams.Set("notebook_id", notebookID)
	}
	queryParams.Set("count_only", "true")

	endpoint := "/api/notes?" + queryParams.Encode()
	resp, err := n.httpClient.Get(ctx, endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get notes count: %w", err)
	}

	var result struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse count response: %w", err)
	}

	return result.Count, nil
}
