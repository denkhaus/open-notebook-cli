package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/do/v2"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// Private search repository implementation
type searchRepository struct {
	httpClient HTTPClient
	logger     Logger
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(injector do.Injector) (SearchRepository, error) {
	httpClient := do.MustInvoke[HTTPClient](injector)
	logger := do.MustInvoke[Logger](injector)

	return &searchRepository{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// Search implements SearchRepository interface
func (r *searchRepository) Search(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error) {
	r.logger.Info("Performing search", "query", req.Query, "type", string(req.Type), "limit", req.Limit)

	resp, err := r.httpClient.Post(ctx, "/api/search", req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	var result models.SearchResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	r.logger.Info("Search completed", "results_count", len(result.Results), "total_count", result.TotalCount)
	return &result, nil
}

// Ask implements SearchRepository interface with streaming response
func (r *searchRepository) Ask(ctx context.Context, req *models.AskRequest) (<-chan *models.StreamChunk, error) {
	r.logger.Info("Starting AI ask", "question", req.Question, "strategy_model", req.StrategyModel, "answer_model", req.AnswerModel)

	// Use the streaming HTTP client for streaming responses
	streamChan, err := r.httpClient.Stream(ctx, "/api/search/ask", req)
	if err != nil {
		return nil, fmt.Errorf("failed to start ask stream: %w", err)
	}

	// Create output channel for stream chunks
	chunkChan := make(chan *models.StreamChunk, 100)

	// Start goroutine to process stream
	go func() {
		defer close(chunkChan)

		for {
			select {
			case <-ctx.Done():
				r.logger.Debug("Ask context cancelled")
				return
			case data, ok := <-streamChan:
				if !ok {
					// Stream ended
					r.logger.Info("Ask stream completed")
					return
				}

				// Parse streaming data as JSON line
				var chunk models.StreamChunk
				if err := json.Unmarshal(data, &chunk); err != nil {
					// If JSON parsing fails, treat as plain text content
					chunk = models.StreamChunk{
						Content: string(data),
						Done:    false,
					}
				}

				// Send chunk to output channel
				select {
				case chunkChan <- &chunk:
				case <-ctx.Done():
					return
				}

				// If this is the final chunk, end the stream
				if chunk.Done {
					r.logger.Info("Ask stream final chunk received")
					return
				}
			}
		}
	}()

	return chunkChan, nil
}

// AskSimple implements SearchRepository interface for non-streaming responses
func (r *searchRepository) AskSimple(ctx context.Context, req *models.AskRequest) (*models.AskResponse, error) {
	r.logger.Info("Starting simple AI ask", "question", req.Question, "strategy_model", req.StrategyModel, "answer_model", req.AnswerModel)

	resp, err := r.httpClient.Post(ctx, "/api/search/ask/simple", req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform simple ask: %w", err)
	}

	var result models.AskResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ask response: %w", err)
	}

	r.logger.Info("Simple ask completed", "answer_length", len(result.Answer))
	return &result, nil
}

// Private search service implementation
type searchService struct {
	repo SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(injector do.Injector) (SearchService, error) {
	repo := do.MustInvoke[SearchRepository](injector)

	return &searchService{
		repo: repo,
	}, nil
}

// Repository returns the underlying repository
func (s *searchService) Repository() SearchRepository {
	return s.repo
}

// Search performs a search with options
func (s *searchService) Search(ctx context.Context, query string, options *models.SearchOptions) (*models.SearchResponse, error) {
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

	return s.repo.Search(ctx, req)
}

// Ask performs an AI ask with streaming response
func (s *searchService) Ask(ctx context.Context, question string, options *models.AskOptions) (<-chan *models.StreamChunk, error) {
	if options == nil {
		options = &models.AskOptions{}
	}

	req := &models.AskRequest{
		Question:         question,
		StrategyModel:    options.StrategyModel,
		AnswerModel:      options.AnswerModel,
		FinalAnswerModel: options.FinalAnswerModel,
	}

	return s.repo.Ask(ctx, req)
}

// AskSimple performs an AI ask with non-streaming response
func (s *searchService) AskSimple(ctx context.Context, question string, options *models.AskOptions) (*models.AskResponse, error) {
	if options == nil {
		options = &models.AskOptions{}
	}

	req := &models.AskRequest{
		Question:         question,
		StrategyModel:    options.StrategyModel,
		AnswerModel:      options.AnswerModel,
		FinalAnswerModel: options.FinalAnswerModel,
	}

	return s.repo.AskSimple(ctx, req)
}