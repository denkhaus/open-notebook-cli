package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// MockNotebookRepository provides a mock implementation of NotebookRepository
type MockNotebookRepository struct {
	*MockBase
	notebooks map[string]*models.Notebook
}

// NewMockNotebookRepository creates a new mock notebook repository
func NewMockNotebookRepository() *MockNotebookRepository {
	return &MockNotebookRepository{
		MockBase:  NewMockBase(0),
		notebooks: make(map[string]*models.Notebook),
	}
}

// AddNotebook adds a notebook to the mock repository
func (m *MockNotebookRepository) AddNotebook(notebook *models.Notebook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notebooks[notebook.ID] = notebook
}

// SetNotebooks sets all notebooks in the mock repository
func (m *MockNotebookRepository) SetNotebooks(notebooks []*models.Notebook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notebooks = make(map[string]*models.Notebook)
	for _, nb := range notebooks {
		m.notebooks[nb.ID] = nb
	}
}

// List implements NotebookRepository interface
func (m *MockNotebookRepository) List(ctx context.Context) ([]*models.Notebook, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("List", []interface{}{ctx}, nil, err)
		return nil, err
	}

	if err := m.GetError("List"); err != nil {
		m.RecordCall("List", []interface{}{ctx}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	result := make([]*models.Notebook, 0, len(m.notebooks))
	for _, nb := range m.notebooks {
		// Deep copy to avoid mutation
		nbCopy := *nb
		result = append(result, &nbCopy)
	}
	m.mu.RUnlock()

	m.RecordCall("List", []interface{}{ctx}, result, nil)
	return result, nil
}

// Create implements NotebookRepository interface
func (m *MockNotebookRepository) Create(ctx context.Context, notebook *models.NotebookCreate) (*models.Notebook, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Create", []interface{}{ctx, notebook}, nil, err)
		return nil, err
	}

	if err := m.GetError("Create"); err != nil {
		m.RecordCall("Create", []interface{}{ctx, notebook}, nil, err)
		return nil, err
	}

	newNotebook := &models.Notebook{
		ID:          "mock-notebook-" + generateID(),
		Name:        notebook.Name,
		Description: notebook.Description,
		Created:     currentTime().Format(time.RFC3339),
		Updated:     currentTime().Format(time.RFC3339),
	}

	m.AddNotebook(newNotebook)
	m.RecordCall("Create", []interface{}{ctx, notebook}, newNotebook, nil)
	return newNotebook, nil
}

// Get implements NotebookRepository interface
func (m *MockNotebookRepository) Get(ctx context.Context, id string) (*models.Notebook, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Get", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	if err := m.GetError("Get"); err != nil {
		m.RecordCall("Get", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	nb, exists := m.notebooks[id]
	m.mu.RUnlock()

	if !exists {
		err := errors.New("notebook not found")
		m.RecordCall("Get", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	nbCopy := *nb
	m.RecordCall("Get", []interface{}{ctx, id}, &nbCopy, nil)
	return &nbCopy, nil
}

// Update implements NotebookRepository interface
func (m *MockNotebookRepository) Update(ctx context.Context, id string, notebook *models.NotebookUpdate) (*models.Notebook, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Update", []interface{}{ctx, id, notebook}, nil, err)
		return nil, err
	}

	if err := m.GetError("Update"); err != nil {
		m.RecordCall("Update", []interface{}{ctx, id, notebook}, nil, err)
		return nil, err
	}

	m.mu.Lock()
	nb, exists := m.notebooks[id]
	if !exists {
		m.mu.Unlock()
		err := errors.New("notebook not found")
		m.RecordCall("Update", []interface{}{ctx, id, notebook}, nil, err)
		return nil, err
	}

	// Update fields
	if notebook.Name != nil {
		nb.Name = *notebook.Name
	}
	if notebook.Description != nil {
		nb.Description = *notebook.Description
	}
	if notebook.Archived != nil {
		nb.Archived = *notebook.Archived
	}
	nb.Updated = currentTime().Format(time.RFC3339)

	nbCopy := *nb
	m.mu.Unlock()

	m.RecordCall("Update", []interface{}{ctx, id, notebook}, &nbCopy, nil)
	return &nbCopy, nil
}

// Delete implements NotebookRepository interface
func (m *MockNotebookRepository) Delete(ctx context.Context, id string) error {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Delete", []interface{}{ctx, id}, nil, err)
		return err
	}

	if err := m.GetError("Delete"); err != nil {
		m.RecordCall("Delete", []interface{}{ctx, id}, nil, err)
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.notebooks[id]; !exists {
		err := errors.New("notebook not found")
		m.RecordCall("Delete", []interface{}{ctx, id}, nil, err)
		return err
	}

	delete(m.notebooks, id)
	m.RecordCall("Delete", []interface{}{ctx, id}, nil, nil)
	return nil
}

// AddSource implements NotebookRepository interface
func (m *MockNotebookRepository) AddSource(ctx context.Context, notebookID, sourceID string) error {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("AddSource", []interface{}{ctx, notebookID, sourceID}, nil, err)
		return err
	}

	if err := m.GetError("AddSource"); err != nil {
		m.RecordCall("AddSource", []interface{}{ctx, notebookID, sourceID}, nil, err)
		return err
	}

	// For mock purposes, just record the call and succeed
	m.RecordCall("AddSource", []interface{}{ctx, notebookID, sourceID}, nil, nil)
	return nil
}

// RemoveSource implements NotebookRepository interface
func (m *MockNotebookRepository) RemoveSource(ctx context.Context, notebookID, sourceID string) error {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("RemoveSource", []interface{}{ctx, notebookID, sourceID}, nil, err)
		return err
	}

	if err := m.GetError("RemoveSource"); err != nil {
		m.RecordCall("RemoveSource", []interface{}{ctx, notebookID, sourceID}, nil, err)
		return err
	}

	// For mock purposes, just record the call and succeed
	m.RecordCall("RemoveSource", []interface{}{ctx, notebookID, sourceID}, nil, nil)
	return nil
}

// MockNotebookService provides a mock implementation of NotebookService
type MockNotebookService struct {
	*MockBase
	repository *MockNotebookRepository
}

// NewMockNotebookService creates a new mock notebook service
func NewMockNotebookService() *MockNotebookService {
	return &MockNotebookService{
		MockBase:   NewMockBase(0),
		repository: NewMockNotebookRepository(),
	}
}

// Repository implements NotebookService interface
func (m *MockNotebookService) Repository() interface{} {
	return m.repository
}

// ListNotebooks implements NotebookService interface
func (m *MockNotebookService) ListNotebooks(ctx context.Context) ([]*models.Notebook, error) {
	return m.repository.List(ctx)
}

// CreateNotebook implements NotebookService interface
func (m *MockNotebookService) CreateNotebook(ctx context.Context, name, description string) (*models.Notebook, error) {
	createReq := &models.NotebookCreate{
		Name:        name,
		Description: description,
	}
	return m.repository.Create(ctx, createReq)
}

// GetNotebook implements NotebookService interface
func (m *MockNotebookService) GetNotebook(ctx context.Context, id string) (*models.Notebook, error) {
	return m.repository.Get(ctx, id)
}

// UpdateNotebook implements NotebookService interface
func (m *MockNotebookService) UpdateNotebook(ctx context.Context, id string, name, description *string, archived *bool) (*models.Notebook, error) {
	updateReq := &models.NotebookUpdate{
		Name:        name,
		Description: description,
		Archived:    archived,
	}
	return m.repository.Update(ctx, id, updateReq)
}

// DeleteNotebook implements NotebookService interface
func (m *MockNotebookService) DeleteNotebook(ctx context.Context, id string) error {
	return m.repository.Delete(ctx, id)
}

// AddSourceToNotebook implements NotebookService interface
func (m *MockNotebookService) AddSourceToNotebook(ctx context.Context, notebookID, sourceID string) error {
	return m.repository.AddSource(ctx, notebookID, sourceID)
}

// RemoveSourceFromNotebook implements NotebookService interface
func (m *MockNotebookService) RemoveSourceFromNotebook(ctx context.Context, notebookID, sourceID string) error {
	return m.repository.RemoveSource(ctx, notebookID, sourceID)
}

// GetRepository returns the underlying mock repository for testing purposes
func (m *MockNotebookService) GetRepository() *MockNotebookRepository {
	return m.repository
}