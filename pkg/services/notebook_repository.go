package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/denkhaus/open-notebook-cli/pkg/models"

	"github.com/samber/do/v2"
)

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
	resp, err := r.http.Get(ctx, "/notebooks/"+id)
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
	resp, err := r.http.Put(ctx, "/notebooks/"+id, notebook)
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
	resp, err := r.http.Delete(ctx, "/notebooks/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete notebook: %w", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}

func (r *notebookRepository) AddSource(ctx context.Context, notebookID, sourceID string) error {
	// This might need a specific endpoint or be part of update
	// For now, assuming there's an endpoint for this
	payload := map[string]string{
		"source_id": sourceID,
	}

	resp, err := r.http.Post(ctx, "/notebooks/"+notebookID+"/sources", payload)
	if err != nil {
		return fmt.Errorf("failed to add source to notebook: %w", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}

func (r *notebookRepository) RemoveSource(ctx context.Context, notebookID, sourceID string) error {
	resp, err := r.http.Delete(ctx, "/notebooks/"+notebookID+"/sources/"+sourceID)
	if err != nil {
		return fmt.Errorf("failed to remove source from notebook: %w", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}
