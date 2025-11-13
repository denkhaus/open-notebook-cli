package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/do/v2"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// Private notebook service implementation
type notebookService struct {
	repo NotebookRepository
}

// NewNotebookService creates a new notebook service
func NewNotebookService(injector do.Injector) (NotebookService, error) {
	repo := do.MustInvoke[NotebookRepository](injector)

	return &notebookService{
		repo: repo,
	}, nil
}

// Interface implementation

func (s *notebookService) Repository() NotebookRepository {
	return s.repo
}

func (s *notebookService) ListNotebooks(ctx context.Context) ([]*models.Notebook, error) {
	return s.repo.List(ctx)
}

func (s *notebookService) CreateNotebook(ctx context.Context, name, description string) (*models.Notebook, error) {
	if name == "" {
		return nil, fmt.Errorf("notebook name is required")
	}

	create := &models.NotebookCreate{
		Name:        name,
		Description: description,
	}

	return s.repo.Create(ctx, create)
}

func (s *notebookService) GetNotebook(ctx context.Context, id string) (*models.Notebook, error) {
	if id == "" {
		return nil, fmt.Errorf("notebook ID is required")
	}

	return s.repo.Get(ctx, id)
}

func (s *notebookService) UpdateNotebook(ctx context.Context, id string, name, description *string, archived *bool) (*models.Notebook, error) {
	if id == "" {
		return nil, fmt.Errorf("notebook ID is required")
	}

	update := &models.NotebookUpdate{}
	if name != nil {
		update.Name = name
	}
	if description != nil {
		update.Description = description
	}
	if archived != nil {
		update.Archived = archived
	}

	return s.repo.Update(ctx, id, update)
}

func (s *notebookService) DeleteNotebook(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("notebook ID is required")
	}

	return s.repo.Delete(ctx, id)
}

func (s *notebookService) AddSourceToNotebook(ctx context.Context, notebookID, sourceID string) error {
	if notebookID == "" {
		return fmt.Errorf("notebook ID is required")
	}
	if sourceID == "" {
		return fmt.Errorf("source ID is required")
	}

	return s.repo.AddSource(ctx, notebookID, sourceID)
}

func (s *notebookService) RemoveSourceFromNotebook(ctx context.Context, notebookID, sourceID string) error {
	if notebookID == "" {
		return fmt.Errorf("notebook ID is required")
	}
	if sourceID == "" {
		return fmt.Errorf("source ID is required")
	}

	return s.repo.RemoveSource(ctx, notebookID, sourceID)
}

// Repository implementation
type notebookRepository struct {
	http HTTPClient
}

// NewNotebookRepository creates a new notebook repository
func NewNotebookRepository(injector do.Injector) (NotebookRepository, error) {
	http := do.MustInvoke[HTTPClient](injector)

	return &notebookRepository{
		http: http,
	}, nil
}

// Repository interface implementation

func (r *notebookRepository) List(ctx context.Context) ([]*models.Notebook, error) {
	resp, err := r.http.Get(ctx, "/notebooks")
	if err != nil {
		return nil, fmt.Errorf("failed to list notebooks: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	var notebooks []*models.Notebook
	if err := json.Unmarshal(resp.Body, &notebooks); err != nil {
		return nil, fmt.Errorf("failed to decode notebooks response: %w", err)
	}

	return notebooks, nil
}

func (r *notebookRepository) Create(ctx context.Context, notebook *models.NotebookCreate) (*models.Notebook, error) {
	resp, err := r.http.Post(ctx, "/notebooks", notebook)
	if err != nil {
		return nil, fmt.Errorf("failed to create notebook: %w", err)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	var createdNotebook models.Notebook
	if err := json.Unmarshal(resp.Body, &createdNotebook); err != nil {
		return nil, fmt.Errorf("failed to decode notebook response: %w", err)
	}

	return &createdNotebook, nil
}

func (r *notebookRepository) Get(ctx context.Context, id string) (*models.Notebook, error) {
	endpoint := fmt.Sprintf("/notebooks/%s", id)

	resp, err := r.http.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get notebook: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	var notebook models.Notebook
	if err := json.Unmarshal(resp.Body, &notebook); err != nil {
		return nil, fmt.Errorf("failed to decode notebook response: %w", err)
	}

	return &notebook, nil
}

func (r *notebookRepository) Update(ctx context.Context, id string, notebook *models.NotebookUpdate) (*models.Notebook, error) {
	endpoint := fmt.Sprintf("/notebooks/%s", id)

	resp, err := r.http.Put(ctx, endpoint, notebook)
	if err != nil {
		return nil, fmt.Errorf("failed to update notebook: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	var updatedNotebook models.Notebook
	if err := json.Unmarshal(resp.Body, &updatedNotebook); err != nil {
		return nil, fmt.Errorf("failed to decode notebook response: %w", err)
	}

	return &updatedNotebook, nil
}

func (r *notebookRepository) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/notebooks/%s", id)

	resp, err := r.http.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete notebook: %w", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}

func (r *notebookRepository) AddSource(ctx context.Context, notebookID, sourceID string) error {
	endpoint := fmt.Sprintf("/notebooks/%s/sources/%s", notebookID, sourceID)

	resp, err := r.http.Post(ctx, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to add source to notebook: %w", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}

func (r *notebookRepository) RemoveSource(ctx context.Context, notebookID, sourceID string) error {
	endpoint := fmt.Sprintf("/notebooks/%s/sources/%s", notebookID, sourceID)

	resp, err := r.http.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to remove source from notebook: %w", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}

// Mock implementations for testing

type mockNotebookService struct {
	notebooks map[string]*models.Notebook
}

func NewMockNotebookService() NotebookService {
	return &mockNotebookService{
		notebooks: make(map[string]*models.Notebook),
	}
}

func (m *mockNotebookService) Repository() NotebookRepository {
	return NewMockNotebookRepository()
}

func (m *mockNotebookService) ListNotebooks(ctx context.Context) ([]*models.Notebook, error) {
	var notebooks []*models.Notebook
	for _, nb := range m.notebooks {
		notebooks = append(notebooks, nb)
	}
	return notebooks, nil
}

func (m *mockNotebookService) CreateNotebook(ctx context.Context, name, description string) (*models.Notebook, error) {
	if name == "" {
		return nil, fmt.Errorf("notebook name is required")
	}

	notebook := &models.Notebook{
		ID:          fmt.Sprintf("nb-%d", len(m.notebooks)+1),
		Name:        name,
		Description: description,
		Archived:    false,
		Created:     "2024-01-01T00:00:00Z",
		Updated:     "2024-01-01T00:00:00Z",
		SourceCount: 0,
		NoteCount:   0,
	}

	m.notebooks[notebook.ID] = notebook
	return notebook, nil
}

func (m *mockNotebookService) GetNotebook(ctx context.Context, id string) (*models.Notebook, error) {
	notebook, exists := m.notebooks[id]
	if !exists {
		return nil, fmt.Errorf("notebook not found: %s", id)
	}
	return notebook, nil
}

func (m *mockNotebookService) UpdateNotebook(ctx context.Context, id string, name, description *string, archived *bool) (*models.Notebook, error) {
	notebook, exists := m.notebooks[id]
	if !exists {
		return nil, fmt.Errorf("notebook not found: %s", id)
	}

	if name != nil {
		notebook.Name = *name
	}
	if description != nil {
		notebook.Description = *description
	}
	if archived != nil {
		notebook.Archived = *archived
	}

	return notebook, nil
}

func (m *mockNotebookService) DeleteNotebook(ctx context.Context, id string) error {
	if _, exists := m.notebooks[id]; !exists {
		return fmt.Errorf("notebook not found: %s", id)
	}
	delete(m.notebooks, id)
	return nil
}

func (m *mockNotebookService) AddSourceToNotebook(ctx context.Context, notebookID, sourceID string) error {
	if _, exists := m.notebooks[notebookID]; !exists {
		return fmt.Errorf("notebook not found: %s", notebookID)
	}
	return nil // Mock implementation
}

func (m *mockNotebookService) RemoveSourceFromNotebook(ctx context.Context, notebookID, sourceID string) error {
	if _, exists := m.notebooks[notebookID]; !exists {
		return fmt.Errorf("notebook not found: %s", notebookID)
	}
	return nil // Mock implementation
}

type mockNotebookRepository struct {
	notebooks map[string]*models.Notebook
}

func NewMockNotebookRepository() NotebookRepository {
	return &mockNotebookRepository{
		notebooks: make(map[string]*models.Notebook),
	}
}

func (m *mockNotebookRepository) List(ctx context.Context) ([]*models.Notebook, error) {
	var notebooks []*models.Notebook
	for _, nb := range m.notebooks {
		notebooks = append(notebooks, nb)
	}
	return notebooks, nil
}

func (m *mockNotebookRepository) Create(ctx context.Context, notebook *models.NotebookCreate) (*models.Notebook, error) {
	// Create a notebook from the create request
	result := &models.Notebook{
		ID:          fmt.Sprintf("nb-%d", len(m.notebooks)+1),
		Name:        notebook.Name,
		Description: notebook.Description,
		Archived:    false,
		Created:     "2024-01-01T00:00:00Z",
		Updated:     "2024-01-01T00:00:00Z",
		SourceCount: 0,
		NoteCount:   0,
	}

	m.notebooks[result.ID] = result
	return result, nil
}

func (m *mockNotebookRepository) Get(ctx context.Context, id string) (*models.Notebook, error) {
	notebook, exists := m.notebooks[id]
	if !exists {
		return nil, fmt.Errorf("notebook not found: %s", id)
	}
	return notebook, nil
}

func (m *mockNotebookRepository) Update(ctx context.Context, id string, notebook *models.NotebookUpdate) (*models.Notebook, error) {
	existing, exists := m.notebooks[id]
	if !exists {
		return nil, fmt.Errorf("notebook not found: %s", id)
	}

	if notebook.Name != nil {
		existing.Name = *notebook.Name
	}
	if notebook.Description != nil {
		existing.Description = *notebook.Description
	}
	if notebook.Archived != nil {
		existing.Archived = *notebook.Archived
	}

	return existing, nil
}

func (m *mockNotebookRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.notebooks[id]; !exists {
		return fmt.Errorf("notebook not found: %s", id)
	}
	delete(m.notebooks, id)
	return nil
}

func (m *mockNotebookRepository) AddSource(ctx context.Context, notebookID, sourceID string) error {
	if _, exists := m.notebooks[notebookID]; !exists {
		return fmt.Errorf("notebook not found: %s", notebookID)
	}
	// Mock implementation - would update source count in real implementation
	return nil
}

func (m *mockNotebookRepository) RemoveSource(ctx context.Context, notebookID, sourceID string) error {
	if _, exists := m.notebooks[notebookID]; !exists {
		return fmt.Errorf("notebook not found: %s", notebookID)
	}
	// Mock implementation - would update source count in real implementation
	return nil
}