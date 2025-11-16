package services

import (
	"context"
	"fmt"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
)

type modelService struct {
	repo ModelRepository
}

// NewModelService creates a new model service
func NewModelService(injector do.Injector) (ModelService, error) {
	repo := do.MustInvoke[ModelRepository](injector)

	return &modelService{
		repo: repo,
	}, nil
}

// Interface implementation

func (s *modelService) Repository() ModelRepository {
	return s.repo
}

func (s *modelService) List(ctx context.Context) ([]*models.Model, error) {
	return s.repo.List(ctx)
}

func (s *modelService) Create(ctx context.Context, model *models.ModelCreate) (*models.Model, error) {
	// Business logic validation
	if model.Name == "" {
		return nil, fmt.Errorf("model name is required")
	}
	if model.Provider == "" {
		return nil, fmt.Errorf("model provider is required")
	}
	if model.Type == "" {
		return nil, fmt.Errorf("model type is required")
	}

	// Additional business logic could be added here:
	// - Validate provider is supported
	// - Check for duplicate models
	// - Validate model name format
	// etc.

	return s.repo.Create(ctx, model)
}

func (s *modelService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("model ID is required")
	}

	// Business logic: could check if model is being used as default
	// before allowing deletion

	return s.repo.Delete(ctx, id)
}

func (s *modelService) GetDefaults(ctx context.Context) (*models.DefaultModelsResponse, error) {
	return s.repo.GetDefaults(ctx)
}

func (s *modelService) SetDefaults(ctx context.Context, defaults *models.DefaultModelsResponse) error {
	if defaults == nil {
		return fmt.Errorf("defaults configuration is required")
	}

	// Business logic: validate that all specified models exist
	// This could be implemented by checking against the model list

	return s.repo.SetDefaults(ctx, defaults)
}

func (s *modelService) GetProviders(ctx context.Context) (*models.ProviderAvailabilityResponse, error) {
	return s.repo.GetProviders(ctx)
}