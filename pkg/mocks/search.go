package mocks

import (
	"context"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// MockSearchRepository provides a mock implementation of SearchRepository
type MockSearchRepository struct {
	*MockBase
	searchResults map[string]*models.SearchResponse
	askResponses  map[string]*models.AskResponse
	streamChunks  map[string][]*models.StreamChunk
}

// NewMockSearchRepository creates a new mock search repository
func NewMockSearchRepository() *MockSearchRepository {
	return &MockSearchRepository{
		MockBase:      NewMockBase(0),
		searchResults: make(map[string]*models.SearchResponse),
		askResponses:  make(map[string]*models.AskResponse),
		streamChunks:  make(map[string][]*models.StreamChunk),
	}
}

// SetSearchResult sets a mock search response for a query
func (m *MockSearchRepository) SetSearchResult(query string, response *models.SearchResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.searchResults[query] = response
}

// SetAskResponse sets a mock ask response for a question
func (m *MockSearchRepository) SetAskResponse(question string, response *models.AskResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.askResponses[question] = response
}

// SetStreamChunks sets mock stream chunks for a question
func (m *MockSearchRepository) SetStreamChunks(question string, chunks []*models.StreamChunk) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streamChunks[question] = chunks
}

// Search implements SearchRepository interface
func (m *MockSearchRepository) Search(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Search", []interface{}{ctx, req}, nil, err)
		return nil, err
	}

	if err := m.GetError("Search"); err != nil {
		m.RecordCall("Search", []interface{}{ctx, req}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	response, exists := m.searchResults[req.Query]
	if !exists {
		// Return empty result if no mock is set
		response = &models.SearchResponse{
			Results:    []models.SearchResult{},
			TotalCount: 0,
			SearchType: string(req.Type),
		}
	}
	// Deep copy
	responseCopy := *response
	responseCopy.Results = make([]models.SearchResult, len(response.Results))
	copy(responseCopy.Results, response.Results)
	m.mu.RUnlock()

	m.RecordCall("Search", []interface{}{ctx, req}, &responseCopy, nil)
	return &responseCopy, nil
}

// Ask implements SearchRepository interface with streaming response
func (m *MockSearchRepository) Ask(ctx context.Context, req *models.AskRequest) (<-chan *models.StreamChunk, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("Ask", []interface{}{ctx, req}, nil, err)
		return nil, err
	}

	if err := m.GetError("Ask"); err != nil {
		m.RecordCall("Ask", []interface{}{ctx, req}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	chunks, exists := m.streamChunks[req.Question]
	if !exists {
		// Generate default mock chunks
		chunks = []*models.StreamChunk{
			{Content: "Mock response to: ", Done: false},
			{Content: req.Question, Done: false},
			{Content: "\nThis is a simulated AI response.", Done: false},
			{Content: "", Done: true},
		}
	}
	// Deep copy
	chunksCopy := make([]*models.StreamChunk, len(chunks))
	copy(chunksCopy, chunks)
	m.mu.RUnlock()

	// Create channel and stream chunks
	chunkChan := make(chan *models.StreamChunk, len(chunksCopy))
	go func() {
		defer close(chunkChan)
		for _, chunk := range chunksCopy {
			chunkCopy := *chunk
			chunkChan <- &chunkCopy
		}
	}()

	m.RecordCall("Ask", []interface{}{ctx, req}, chunkChan, nil)
	return chunkChan, nil
}

// AskSimple implements SearchRepository interface for non-streaming responses
func (m *MockSearchRepository) AskSimple(ctx context.Context, req *models.AskRequest) (*models.AskResponse, error) {
	m.simulateDelay()

	if err := m.checkFailure(); err != nil {
		m.RecordCall("AskSimple", []interface{}{ctx, req}, nil, err)
		return nil, err
	}

	if err := m.GetError("AskSimple"); err != nil {
		m.RecordCall("AskSimple", []interface{}{ctx, req}, nil, err)
		return nil, err
	}

	m.mu.RLock()
	response, exists := m.askResponses[req.Question]
	if !exists {
		// Generate default mock response
		response = &models.AskResponse{
			Answer:   "This is a mock AI response to: " + req.Question,
			Question: req.Question,
		}
	}
	// Deep copy
	responseCopy := *response
	m.mu.RUnlock()

	m.RecordCall("AskSimple", []interface{}{ctx, req}, &responseCopy, nil)
	return &responseCopy, nil
}

// MockSearchService provides a mock implementation of SearchService
type MockSearchService struct {
	*MockBase
	repository *MockSearchRepository
}

// NewMockSearchService creates a new mock search service
func NewMockSearchService() *MockSearchService {
	return &MockSearchService{
		MockBase:   NewMockBase(0),
		repository: NewMockSearchRepository(),
	}
}

// Repository implements SearchService interface
func (m *MockSearchService) Repository() interface{} {
	return m.repository
}

// Search implements SearchService interface
func (m *MockSearchService) Search(ctx context.Context, query string, options *models.SearchOptions) (*models.SearchResponse, error) {
	if options == nil {
		options = &models.SearchOptions{
			Type:          "vector",
			Limit:         10,
			MinimumScore:  0.0,
			SearchSources: true,
			SearchNotes:   true,
		}
	}

	// Convert string type to SearchType enum
	searchType := models.SearchTypeText
	if options.Type == "vector" {
		searchType = models.SearchTypeVector
	}

	req := &models.SearchRequest{
		Query:         query,
		Type:          searchType,
		Limit:         options.Limit,
		SearchSources: options.SearchSources,
		SearchNotes:   options.SearchNotes,
		MinimumScore:  options.MinimumScore,
	}

	return m.repository.Search(ctx, req)
}

// Ask implements SearchService interface
func (m *MockSearchService) Ask(ctx context.Context, question string, options *models.AskOptions) (<-chan *models.StreamChunk, error) {
	if options == nil {
		options = &models.AskOptions{}
	}

	req := &models.AskRequest{
		Question:         question,
		StrategyModel:    options.StrategyModel,
		AnswerModel:      options.AnswerModel,
		FinalAnswerModel: options.FinalAnswerModel,
	}

	return m.repository.Ask(ctx, req)
}

// AskSimple implements SearchService interface
func (m *MockSearchService) AskSimple(ctx context.Context, question string, options *models.AskOptions) (*models.AskResponse, error) {
	if options == nil {
		options = &models.AskOptions{}
	}

	req := &models.AskRequest{
		Question:         question,
		StrategyModel:    options.StrategyModel,
		AnswerModel:      options.AnswerModel,
		FinalAnswerModel: options.FinalAnswerModel,
	}

	return m.repository.AskSimple(ctx, req)
}

// GetRepository returns the underlying mock repository for testing purposes
func (m *MockSearchService) GetRepository() *MockSearchRepository {
	return m.repository
}

// Helper methods for test setup

// AddMockSearchResult adds a predefined search result
func (m *MockSearchService) AddMockSearchResult(query string, results []models.SearchResult) {
	response := &models.SearchResponse{
		Results:    results,
		TotalCount: len(results),
		SearchType: "vector",
	}
	m.repository.SetSearchResult(query, response)
}

// AddMockAskResponse adds a predefined ask response
func (m *MockSearchService) AddMockAskResponse(question string, answer string) {
	response := &models.AskResponse{
		Answer:   answer,
		Question: question,
	}
	m.repository.SetAskResponse(question, response)
}

// AddMockStreamChunks adds predefined stream chunks for a question
func (m *MockSearchService) AddMockStreamChunks(question string, contentChunks []string) {
	chunks := make([]*models.StreamChunk, len(contentChunks))
	for i, content := range contentChunks {
		chunks[i] = &models.StreamChunk{
			Content: content,
			Done:    i == len(contentChunks)-1, // Last chunk is done
		}
	}
	m.repository.SetStreamChunks(question, chunks)
}