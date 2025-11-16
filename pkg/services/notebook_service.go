package services

import (
	"context"
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