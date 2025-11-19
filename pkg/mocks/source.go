package mocks

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// MockSourceRepository provides a mock implementation of SourceRepository
type MockSourceRepository struct {
	*MockBase
	sources map[string]*models.Source
	insights map[string][]*models.SourceInsightResponse
}

// NewMockSourceRepository creates a new mock source repository
func NewMockSourceRepository() *MockSourceRepository {
	return &MockSourceRepository{
		MockBase:  NewMockBase(0),
		sources:   make(map[string]*models.Source),
		insights:  make(map[string][]*models.SourceInsightResponse),
	}
}

// AddSource adds a source to the mock repository
func (m *MockSourceRepository) AddSource(source *models.Source) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if source.ID != nil {
		m.sources[*source.ID] = source
	}
}

// SetSources sets all sources in the mock repository
func (m *MockSourceRepository) SetSources(sources []*models.Source) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources = make(map[string]*models.Source)
	for _, src := range sources {
		if src.ID != nil {
			m.sources[*src.ID] = src
		}
	}
}

// List implements SourceRepository interface
func (m *MockSourceRepository) List(ctx context.Context, limit, offset int) ([]*models.SourceListResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("List", []interface{}{ctx, limit, offset}, nil, err)
		return nil, err
	}

	if err := m.GetError("List"); err != nil {
		m.RecordCall("List", []interface{}{ctx, limit, offset}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	allSources := make([]*models.Source, 0, len(m.sources))
	for _, src := range m.sources {
		srcCopy := *src
		allSources = append(allSources, &srcCopy)
	}
	m.mu.RUnlock()

	// Apply pagination
	total := len(allSources)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	var paginatedSources []*models.Source
	if start < end {
		paginatedSources = allSources[start:end]
	}

	result := make([]*models.SourceListResponse, len(paginatedSources))
	for i, src := range paginatedSources {
		result[i] = &models.SourceListResponse{
			ID:             src.ID,
			Title:          src.Title,
			Topics:         src.Topics,
			Asset:          src.Asset,
			Embedded:       src.Embedded,
			EmbeddedChunks: src.EmbeddedChunks,
			Created:        src.Created,
			Updated:        src.Updated,
			Status:         src.Status,
		}
	}

	m.RecordCall("List", []interface{}{ctx, limit, offset}, result, nil)
	return result, nil
}

// Create implements SourceRepository interface
func (m *MockSourceRepository) Create(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Create", []interface{}{ctx, source}, nil, err)
		return nil, err
	}

	if err := m.GetError("Create"); err != nil {
		m.RecordCall("Create", []interface{}{ctx, source}, nil, err)
		return nil, err
	}

	id := "mock-source-" + generateID()
	title := source.Title
	if title == nil {
		title = &id
	}

	status := models.SourceStatusPending

	newSource := &models.Source{
		ID:       &id,
		Title:    title,
		Topics:   []string{"mock"},
		FullText: source.Content,
		Embedded: false,
		Status:   &status,
		Created:  currentTime().Format(time.RFC3339),
		Updated:  currentTime().Format(time.RFC3339),
	}

	m.AddSource(newSource)
	m.RecordCall("Create", []interface{}{ctx, source}, newSource, nil)
	return newSource, nil
}

// CreateFromJSON implements SourceRepository interface
func (m *MockSourceRepository) CreateFromJSON(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	// For mock purposes, same as Create
	return m.Create(ctx, source)
}

// Get implements SourceRepository interface
func (m *MockSourceRepository) Get(ctx context.Context, id string) (*models.Source, error) {
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
	src, exists := m.sources[id]
	m.mu.RUnlock()

	if !exists {
		err := errors.New("source not found")
		m.RecordCall("Get", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	srcCopy := *src
	m.RecordCall("Get", []interface{}{ctx, id}, &srcCopy, nil)
	return &srcCopy, nil
}

// Update implements SourceRepository interface
func (m *MockSourceRepository) Update(ctx context.Context, id string, source *models.SourceUpdate) (*models.Source, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Update", []interface{}{ctx, id, source}, nil, err)
		return nil, err
	}

	if err := m.GetError("Update"); err != nil {
		m.RecordCall("Update", []interface{}{ctx, id, source}, nil, err)
		return nil, err
	}

	m.mu.Lock()
	src, exists := m.sources[id]
	if !exists {
		m.mu.Unlock()
		err := errors.New("source not found")
		m.RecordCall("Update", []interface{}{ctx, id, source}, nil, err)
		return nil, err
	}

	// Update fields
	if source.Title != nil && *source.Title != "" {
		src.Title = source.Title
	}
	if source.Topics != nil {
		src.Topics = source.Topics
	}
	src.Updated = currentTime().Format(time.RFC3339)

	srcCopy := *src
	m.mu.Unlock()

	m.RecordCall("Update", []interface{}{ctx, id, source}, &srcCopy, nil)
	return &srcCopy, nil
}

// Delete implements SourceRepository interface
func (m *MockSourceRepository) Delete(ctx context.Context, id string) error {
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

	if _, exists := m.sources[id]; !exists {
		err := errors.New("source not found")
		m.RecordCall("Delete", []interface{}{ctx, id}, nil, err)
		return err
	}

	delete(m.sources, id)
	delete(m.insights, id) // Also clean up insights
	m.RecordCall("Delete", []interface{}{ctx, id}, nil, nil)
	return nil
}

// GetStatus implements SourceRepository interface
func (m *MockSourceRepository) GetStatus(ctx context.Context, id string) (*models.SourceStatusResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("GetStatus", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	if err := m.GetError("GetStatus"); err != nil {
		m.RecordCall("GetStatus", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	src, err := m.Get(ctx, id)
	if err != nil {
		m.RecordCall("GetStatus", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	status := &models.SourceStatusResponse{
		Status:  src.Status,
		Message: "Mock status response",
	}

	m.RecordCall("GetStatus", []interface{}{ctx, id}, status, nil)
	return status, nil
}

// Download implements SourceRepository interface
func (m *MockSourceRepository) Download(ctx context.Context, id string) (io.ReadCloser, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Download", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	if err := m.GetError("Download"); err != nil {
		m.RecordCall("Download", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	src, err := m.Get(ctx, id)
	if err != nil {
		m.RecordCall("Download", []interface{}{ctx, id}, nil, err)
		return nil, err
	}

	// Return content as a reader
	var content string
	if src.FullText != nil {
		content = *src.FullText
	}
	if content == "" && src.Title != nil {
		content = "Mock source content for " + *src.Title
	}
	if content == "" {
		content = "Mock source content"
	}

	reader := io.NopCloser(strings.NewReader(content))
	m.RecordCall("Download", []interface{}{ctx, id}, reader, nil)
	return reader, nil
}

// GetInsights implements SourceRepository interface
func (m *MockSourceRepository) GetInsights(ctx context.Context, sourceID string) ([]*models.SourceInsightResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("GetInsights", []interface{}{ctx, sourceID}, nil, err)
		return nil, err
	}

	if err := m.GetError("GetInsights"); err != nil {
		m.RecordCall("GetInsights", []interface{}{ctx, sourceID}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	insights := m.insights[sourceID]
	if insights == nil {
		insights = []*models.SourceInsightResponse{}
	}
	// Deep copy
	insightsCopy := make([]*models.SourceInsightResponse, len(insights))
	copy(insightsCopy, insights)
	m.mu.RUnlock()

	m.RecordCall("GetInsights", []interface{}{ctx, sourceID}, insightsCopy, nil)
	return insightsCopy, nil
}

// CreateInsight implements SourceRepository interface
func (m *MockSourceRepository) CreateInsight(ctx context.Context, sourceID string, req *models.CreateSourceInsightRequest) (*models.SourceInsightResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("CreateInsight", []interface{}{ctx, sourceID, req}, nil, err)
		return nil, err
	}

	if err := m.GetError("CreateInsight"); err != nil {
		m.RecordCall("CreateInsight", []interface{}{ctx, sourceID, req}, nil, err)
		return nil, err
	}

	insight := &models.SourceInsightResponse{
		ID:          "mock-insight-" + generateID(),
		SourceID:    sourceID,
		InsightType: models.InsightTypeSummary,
		Content:     "Mock insight content for transformation: " + req.TransformationID,
		Created:     currentTime().Format(time.RFC3339),
		Updated:     currentTime().Format(time.RFC3339),
	}

	m.mu.Lock()
	m.insights[sourceID] = append(m.insights[sourceID], insight)
	m.mu.Unlock()

	m.RecordCall("CreateInsight", []interface{}{ctx, sourceID, req}, insight, nil)
	return insight, nil
}

// AddInsight adds an insight to the mock repository
func (m *MockSourceRepository) AddInsight(sourceID string, insight *models.SourceInsightResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.insights[sourceID] = append(m.insights[sourceID], insight)
}

// MockSourceService provides a mock implementation of SourceService
type MockSourceService struct {
	*MockBase
	repository *MockSourceRepository
}

// NewMockSourceService creates a new mock source service
func NewMockSourceService() *MockSourceService {
	return &MockSourceService{
		MockBase:   NewMockBase(0),
		repository: NewMockSourceRepository(),
	}
}

// Repository implements SourceService interface
func (m *MockSourceService) Repository() interface{} {
	return m.repository
}

// List implements SourceService interface
func (m *MockSourceService) List(ctx context.Context, limit, offset int) ([]*models.SourceListResponse, error) {
	return m.repository.List(ctx, limit, offset)
}

// AddSourceFromLink implements SourceService interface
func (m *MockSourceService) AddSourceFromLink(ctx context.Context, link string, options *models.SourceOptions) (*models.Source, error) {
	title := "Mock Source from Link"
	createReq := &models.SourceCreate{
		Title:           &title,
		URL:             &link,
		Type:            models.SourceTypeLink,
		Notebooks:       []string{"mock"},
		Embed:           true,
		DeleteSource:    false,
		AsyncProcessing: false,
	}
	if options != nil {
		createReq.Transformations = options.Transformations
		createReq.Embed = options.Embed
		createReq.DeleteSource = options.DeleteSource
		createReq.AsyncProcessing = options.AsyncProcessing
	}
	return m.repository.Create(ctx, createReq)
}

// AddSourceFromUpload implements SourceService interface
func (m *MockSourceService) AddSourceFromUpload(ctx context.Context, filename string, file io.Reader, options *models.SourceOptions) (*models.Source, error) {
	// For mock purposes, ignore file content and create a mock source
	createReq := &models.SourceCreate{
		Title:           &filename,
		FilePath:        &filename,
		Type:            models.SourceTypeUpload,
		Notebooks:       []string{"mock"},
		Embed:           true,
		DeleteSource:    false,
		AsyncProcessing: false,
	}
	if options != nil {
		createReq.Transformations = options.Transformations
		createReq.Embed = options.Embed
		createReq.DeleteSource = options.DeleteSource
		createReq.AsyncProcessing = options.AsyncProcessing
	}
	return m.repository.Create(ctx, createReq)
}

// AddSourceFromText implements SourceService interface
func (m *MockSourceService) AddSourceFromText(ctx context.Context, text, title string, options *models.SourceOptions) (*models.Source, error) {
	createReq := &models.SourceCreate{
		Title:           &title,
		Content:         &text,
		Type:            models.SourceTypeText,
		Notebooks:       []string{"mock"},
		Embed:           true,
		DeleteSource:    false,
		AsyncProcessing: false,
	}
	if options != nil {
		createReq.Transformations = options.Transformations
		createReq.Embed = options.Embed
		createReq.DeleteSource = options.DeleteSource
		createReq.AsyncProcessing = options.AsyncProcessing
	}
	return m.repository.Create(ctx, createReq)
}

// Get implements SourceService interface
func (m *MockSourceService) Get(ctx context.Context, id string) (*models.Source, error) {
	return m.repository.Get(ctx, id)
}

// Update implements SourceService interface
func (m *MockSourceService) Update(ctx context.Context, id string, title string, topics []string) (*models.Source, error) {
	updateReq := &models.SourceUpdate{
		Title:  &title,
		Topics: topics,
	}
	return m.repository.Update(ctx, id, updateReq)
}

// Delete implements SourceService interface
func (m *MockSourceService) Delete(ctx context.Context, id string) error {
	return m.repository.Delete(ctx, id)
}

// GetStatus implements SourceService interface
func (m *MockSourceService) GetStatus(ctx context.Context, id string) (*models.SourceStatusResponse, error) {
	return m.repository.GetStatus(ctx, id)
}

// Download implements SourceService interface
func (m *MockSourceService) Download(ctx context.Context, id string) (io.ReadCloser, error) {
	return m.repository.Download(ctx, id)
}

// Create implements SourceService interface
func (m *MockSourceService) Create(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	return m.repository.Create(ctx, source)
}

// CreateFromJSON implements SourceService interface
func (m *MockSourceService) CreateFromJSON(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	return m.repository.CreateFromJSON(ctx, source)
}

// GetInsights implements SourceService interface
func (m *MockSourceService) GetInsights(ctx context.Context, sourceID string) ([]*models.SourceInsightResponse, error) {
	return m.repository.GetInsights(ctx, sourceID)
}

// CreateInsight implements SourceService interface
func (m *MockSourceService) CreateInsight(ctx context.Context, sourceID string, request *models.CreateSourceInsightRequest) (*models.SourceInsightResponse, error) {
	return m.repository.CreateInsight(ctx, sourceID, request)
}

// GetRepository returns the underlying mock repository for testing purposes
func (m *MockSourceService) GetRepository() *MockSourceRepository {
	return m.repository
}