package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
)

type modelRepository struct {
	httpClient HTTPClient
	logger     Logger
}

// NewModelRepository creates a new model repository
func NewModelRepository(injector do.Injector) (ModelRepository, error) {
	httpClient := do.MustInvoke[HTTPClient](injector)
	logger := do.MustInvoke[Logger](injector)

	return &modelRepository{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// List implements ModelRepository interface
func (m *modelRepository) List(ctx context.Context) ([]*models.Model, error) {
	endpoint := "/models"
	resp, err := m.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	var result []*models.Model
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse models response: %w", err)
	}

	m.logger.Info("Retrieved models", "count", len(result))
	return result, nil
}

// Create implements ModelRepository interface
func (m *modelRepository) Create(ctx context.Context, model *models.ModelCreate) (*models.Model, error) {
	resp, err := m.httpClient.Post(ctx, "/models", model)
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	var result models.Model
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse model response: %w", err)
	}

	m.logger.Info("Created model", "id", result.ID, "name", result.Name, "provider", result.Provider)
	return &result, nil
}

// Delete implements ModelRepository interface
func (m *modelRepository) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/models/%s", id)
	_, err := m.httpClient.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete model %s: %w", id, err)
	}

	m.logger.Info("Deleted model", "id", id)
	return nil
}

// GetDefaults implements ModelRepository interface
func (m *modelRepository) GetDefaults(ctx context.Context) (*models.DefaultModelsResponse, error) {
	endpoint := "/models/defaults"
	resp, err := m.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get default models: %w", err)
	}

	var result models.DefaultModelsResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse defaults response: %w", err)
	}

	m.logger.Info("Retrieved default models")
	return &result, nil
}

// SetDefaults implements ModelRepository interface
func (m *modelRepository) SetDefaults(ctx context.Context, defaults *models.DefaultModelsResponse) error {
	endpoint := "/models/defaults"
	resp, err := m.httpClient.Put(ctx, endpoint, defaults)
	if err != nil {
		return fmt.Errorf("failed to set default models: %w", err)
	}

	var result models.DefaultModelsResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse set defaults response: %w", err)
	}

	m.logger.Info("Updated default models")
	return nil
}

// GetProviders implements ModelRepository interface
func (m *modelRepository) GetProviders(ctx context.Context) (*models.ProviderAvailabilityResponse, error) {
	endpoint := "/models/providers"
	resp, err := m.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get providers: %w", err)
	}

	var result models.ProviderAvailabilityResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse providers response: %w", err)
	}

	m.logger.Info("Retrieved providers", "available", len(result.Available), "unavailable", len(result.Unavailable))
	return &result, nil
}

// Extended methods not part of interface

// ListByType filters models by type
func (m *modelRepository) ListByType(ctx context.Context, modelType models.ModelType) ([]*models.Model, error) {
	queryParams := url.Values{}
	queryParams.Set("type", string(modelType))

	endpoint := "/models?" + queryParams.Encode()
	resp, err := m.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list models by type %s: %w", modelType, err)
	}

	var result []*models.Model
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse models response: %w", err)
	}

	m.logger.Info("Retrieved models by type", "type", modelType, "count", len(result))
	return result, nil
}

// ListByProvider filters models by provider
func (m *modelRepository) ListByProvider(ctx context.Context, provider string) ([]*models.Model, error) {
	queryParams := url.Values{}
	queryParams.Set("provider", provider)

	endpoint := "/models?" + queryParams.Encode()
	resp, err := m.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list models by provider %s: %w", provider, err)
	}

	var result []*models.Model
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse models response: %w", err)
	}

	m.logger.Info("Retrieved models by provider", "provider", provider, "count", len(result))
	return result, nil
}

// ListWithPagination implements pagination support
func (m *modelRepository) ListWithPagination(ctx context.Context, limit, offset int) (*models.ModelsListResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("limit", strconv.Itoa(limit))
	queryParams.Set("offset", strconv.Itoa(offset))

	endpoint := "/models?" + queryParams.Encode()
	resp, err := m.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list models with pagination: %w", err)
	}

	var result models.ModelsListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse models response: %w", err)
	}

	m.logger.Info("Retrieved models with pagination", "count", len(result.Models), "total", result.Total)
	return &result, nil
}