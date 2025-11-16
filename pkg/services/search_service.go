package services

import (
	"context"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
)

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
