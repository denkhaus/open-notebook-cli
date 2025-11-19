package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// MockModelRepository provides a mock implementation of ModelRepository
type MockModelRepository struct {
	*MockBase
	models      map[string]*models.Model
	defaults    *models.DefaultModelsResponse
	providers   *models.ProviderAvailabilityResponse
}

// NewMockModelRepository creates a new mock model repository
func NewMockModelRepository() *MockModelRepository {
	return &MockModelRepository{
		MockBase:  NewMockBase(0),
		models:    make(map[string]*models.Model),
		defaults: &models.DefaultModelsResponse{
			DefaultChatModel:       stringPtr("gpt-3.5-turbo"),
			DefaultEmbeddingModel:  stringPtr("text-embedding-ada-002"),
		},
		providers: &models.ProviderAvailabilityResponse{
			Available: []string{"openai", "anthropic"},
		},
	}
}

// AddModel adds a model to the mock repository
func (m *MockModelRepository) AddModel(model *models.Model) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.models[model.ID] = model
}

// SetModels sets all models in the mock repository
func (m *MockModelRepository) SetModels(modelList []*models.Model) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.models = make(map[string]*models.Model)
	for _, model := range modelList {
		m.models[model.ID] = model
	}
}

// SetDefaults sets the default models
func (m *MockModelRepository) SetDefaults(defaults *models.DefaultModelsResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaults = defaults
}

// SetProviders sets the provider availability
func (m *MockModelRepository) SetProviders(providers *models.ProviderAvailabilityResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers = providers
}

// List implements ModelRepository interface
func (m *MockModelRepository) List(ctx context.Context) ([]*models.Model, error) {
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
	result := make([]*models.Model, 0, len(m.models))
	for _, model := range m.models {
		modelCopy := *model
		result = append(result, &modelCopy)
	}
	m.mu.RUnlock()

	m.RecordCall("List", []interface{}{ctx}, result, nil)
	return result, nil
}

// Create implements ModelRepository interface
func (m *MockModelRepository) Create(ctx context.Context, model *models.ModelCreate) (*models.Model, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Create", []interface{}{ctx, model}, nil, err)
		return nil, err
	}

	if err := m.GetError("Create"); err != nil {
		m.RecordCall("Create", []interface{}{ctx, model}, nil, err)
		return nil, err
	}

	newModel := &models.Model{
		ID:       "mock-model-" + generateID(),
		Name:     model.Name,
		Provider: model.Provider,
		Type:     model.Type,
		Created:  currentTime().Format(time.RFC3339),
		Updated:  currentTime().Format(time.RFC3339),
	}

	m.AddModel(newModel)
	m.RecordCall("Create", []interface{}{ctx, model}, newModel, nil)
	return newModel, nil
}

// Delete implements ModelRepository interface
func (m *MockModelRepository) Delete(ctx context.Context, id string) error {
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

	if _, exists := m.models[id]; !exists {
		err := errors.New("model not found")
		m.RecordCall("Delete", []interface{}{ctx, id}, nil, err)
		return err
	}

	delete(m.models, id)
	m.RecordCall("Delete", []interface{}{ctx, id}, nil, nil)
	return nil
}

// GetDefaults implements ModelRepository interface
func (m *MockModelRepository) GetDefaults(ctx context.Context) (*models.DefaultModelsResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("GetDefaults", []interface{}{ctx}, nil, err)
		return nil, err
	}

	if err := m.GetError("GetDefaults"); err != nil {
		m.RecordCall("GetDefaults", []interface{}{ctx}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	defaultsCopy := *m.defaults
	m.mu.RUnlock()

	m.RecordCall("GetDefaults", []interface{}{ctx}, &defaultsCopy, nil)
	return &defaultsCopy, nil
}

// SetDefaults implements ModelRepository interface
func (m *MockModelRepository) UpdateDefaults(ctx context.Context, defaults *models.DefaultModelsResponse) error {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("UpdateDefaults", []interface{}{ctx, defaults}, nil, err)
		return err
	}

	if err := m.GetError("UpdateDefaults"); err != nil {
		m.RecordCall("UpdateDefaults", []interface{}{ctx, defaults}, nil, err)
		return err
	}

	m.SetDefaults(defaults)
	m.RecordCall("UpdateDefaults", []interface{}{ctx, defaults}, nil, nil)
	return nil
}

// GetProviders implements ModelRepository interface
func (m *MockModelRepository) GetProviders(ctx context.Context) (*models.ProviderAvailabilityResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("GetProviders", []interface{}{ctx}, nil, err)
		return nil, err
	}

	if err := m.GetError("GetProviders"); err != nil {
		m.RecordCall("GetProviders", []interface{}{ctx}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	providersCopy := *m.providers
	m.mu.RUnlock()

	m.RecordCall("GetProviders", []interface{}{ctx}, &providersCopy, nil)
	return &providersCopy, nil
}

// MockModelService provides a mock implementation of ModelService
type MockModelService struct {
	*MockBase
	repository *MockModelRepository
}

// NewMockModelService creates a new mock model service
func NewMockModelService() *MockModelService {
	return &MockModelService{
		MockBase:   NewMockBase(0),
		repository: NewMockModelRepository(),
	}
}

// Repository implements ModelService interface
func (m *MockModelService) Repository() interface{} {
	return m.repository
}

// List implements ModelService interface
func (m *MockModelService) List(ctx context.Context) ([]*models.Model, error) {
	return m.repository.List(ctx)
}

// Create implements ModelService interface
func (m *MockModelService) Create(ctx context.Context, model *models.ModelCreate) (*models.Model, error) {
	return m.repository.Create(ctx, model)
}

// Delete implements ModelService interface
func (m *MockModelService) Delete(ctx context.Context, id string) error {
	return m.repository.Delete(ctx, id)
}

// GetDefaults implements ModelService interface
func (m *MockModelService) GetDefaults(ctx context.Context) (*models.DefaultModelsResponse, error) {
	return m.repository.GetDefaults(ctx)
}

// SetDefaults implements ModelService interface
func (m *MockModelService) SetDefaults(ctx context.Context, defaults *models.DefaultModelsResponse) error {
	return m.repository.UpdateDefaults(ctx, defaults)
}

// GetProviders implements ModelService interface
func (m *MockModelService) GetProviders(ctx context.Context) (*models.ProviderAvailabilityResponse, error) {
	return m.repository.GetProviders(ctx)
}

// GetRepository returns the underlying mock repository for testing purposes
func (m *MockModelService) GetRepository() *MockModelRepository {
	return m.repository
}

// Helper function to create string pointers for model fields
func stringPtr(s string) *string {
	return &s
}