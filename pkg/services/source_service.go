package services

import (
	"context"
	"fmt"
	"io"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

type sourceService struct {
	repo shared.SourceRepository
}

// NewSourceService creates a new source service
func NewSourceService(injector do.Injector) (shared.SourceService, error) {
	repo := do.MustInvoke[shared.SourceRepository](injector)

	return &sourceService{
		repo: repo,
	}, nil
}

// Interface implementation

func (s *sourceService) Repository() shared.SourceRepository {
	return s.repo
}

func (s *sourceService) List(ctx context.Context, limit, offset int) ([]*models.SourceListResponse, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *sourceService) AddSourceFromLink(ctx context.Context, link string, options *models.SourceOptions) (*models.Source, error) {
	if link == "" {
		return nil, fmt.Errorf("source link is required")
	}

	// Validate URL format
	if !isValidURL(link) {
		return nil, fmt.Errorf("invalid URL format")
	}

	// Create source with business logic defaults
	source := &models.SourceCreate{
		URL:             &link,
		Embed:           options.Embed,
		DeleteSource:    options.DeleteSource,
		AsyncProcessing: options.AsyncProcessing,
	}

	// Add transformations if provided
	if len(options.Transformations) > 0 {
		source.Transformations = options.Transformations
	}

	return s.repo.Create(ctx, source)
}

func (s *sourceService) AddSourceFromUpload(ctx context.Context, filename string, file io.Reader, options *models.SourceOptions) (*models.Source, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename is required")
	}
	if file == nil {
		return nil, fmt.Errorf("file is required")
	}

	// Business logic for file validation
	if !isValidFileFormat(filename) {
		return nil, fmt.Errorf("unsupported file format")
	}

	// Create source with upload-specific business logic
	source := &models.SourceCreate{
		FilePath:        &filename,
		Type:            models.SourceTypeUpload,
		Embed:           options.Embed,
		DeleteSource:    options.DeleteSource,
		AsyncProcessing: options.AsyncProcessing,
	}

	// Add transformations if provided
	if len(options.Transformations) > 0 {
		source.Transformations = options.Transformations
	}

	// Use CreateFromJSON method for file uploads (if available)
	return s.repo.CreateFromJSON(ctx, source)
}

func (s *sourceService) AddSourceFromText(ctx context.Context, text, title string, options *models.SourceOptions) (*models.Source, error) {
	if text == "" {
		return nil, fmt.Errorf("text content is required")
	}
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Business logic validation
	if len(text) > 1000000 { // 1MB limit for text
		return nil, fmt.Errorf("text content too large (max 1MB)")
	}

	source := &models.SourceCreate{
		Content:         &text,
		Title:           &title,
		Type:            models.SourceTypeText,
		Embed:           options.Embed,
		DeleteSource:    options.DeleteSource,
		AsyncProcessing: options.AsyncProcessing,
	}

	// Add transformations if provided
	if len(options.Transformations) > 0 {
		source.Transformations = options.Transformations
	}

	return s.repo.CreateFromJSON(ctx, source)
}

func (s *sourceService) Get(ctx context.Context, id string) (*models.Source, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	return s.repo.Get(ctx, id)
}

func (s *sourceService) Update(ctx context.Context, id string, title string, topics []string) (*models.Source, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	// Business logic validation
	if len(title) > 500 {
		return nil, fmt.Errorf("title too long (max 500 characters)")
	}

	if len(topics) > 20 {
		return nil, fmt.Errorf("too many topics (max 20)")
	}

	update := &models.SourceUpdate{
		Title:  &title,
		Topics: topics,
	}

	return s.repo.Update(ctx, id, update)
}

func (s *sourceService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("source ID is required")
	}

	// Business logic: could check if source is referenced by notebooks
	// before allowing deletion

	return s.repo.Delete(ctx, id)
}

func (s *sourceService) GetStatus(ctx context.Context, id string) (*models.SourceStatusResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	return s.repo.GetStatus(ctx, id)
}

func (s *sourceService) Download(ctx context.Context, id string) (io.ReadCloser, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	return s.repo.Download(ctx, id)
}

func (s *sourceService) Create(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	if source == nil {
		return nil, fmt.Errorf("source is required")
	}

	return s.repo.Create(ctx, source)
}

func (s *sourceService) CreateFromJSON(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	if source == nil {
		return nil, fmt.Errorf("source is required")
	}

	return s.repo.CreateFromJSON(ctx, source)
}

func (s *sourceService) GetInsights(ctx context.Context, sourceID string) ([]*models.SourceInsightResponse, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	return s.repo.GetInsights(ctx, sourceID)
}

func (s *sourceService) CreateInsight(ctx context.Context, sourceID string, request *models.CreateSourceInsightRequest) (*models.SourceInsightResponse, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID is required")
	}
	if request == nil {
		return nil, fmt.Errorf("request is required")
	}

	return s.repo.CreateInsight(ctx, sourceID, request)
}

// Helper functions for business logic validation

func isValidURL(url string) bool {
	// Basic URL validation - could be enhanced
	return len(url) > 10 && (url[:7] == "http://" || url[:8] == "https://")
}

func isValidFileFormat(filename string) bool {
	// Check supported file formats
	supportedFormats := []string{
		".txt", ".md", ".pdf", ".doc", ".docx", ".html", ".htm",
		".json", ".csv", ".xml", ".rtf", ".epub", ".mobi",
	}

	for _, format := range supportedFormats {
		if len(filename) > len(format) && filename[len(filename)-len(format):] == format {
			return true
		}
	}

	return false
}
