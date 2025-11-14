package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

type sourceRepository struct {
	httpClient HTTPClient
	logger     Logger
}

// NewSourceRepository creates a new source repository
func NewSourceRepository(httpClient HTTPClient, logger Logger) SourceRepository {
	return &sourceRepository{
		httpClient: httpClient,
		logger:     logger,
	}
}

// List implements existing SourceRepository interface
func (s *sourceRepository) List(ctx context.Context, limit, offset int) ([]*models.SourceListResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("limit", fmt.Sprintf("%d", limit))
	queryParams.Set("offset", fmt.Sprintf("%d", offset))

	endpoint := "/api/sources?" + queryParams.Encode()
	resp, err := s.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}

	var result models.SourcesListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse sources response: %w", err)
	}

	// Convert to pointer slice
	sources := make([]*models.SourceListResponse, len(result.Sources))
	for i := range result.Sources {
		sources[i] = &result.Sources[i]
	}

	s.logger.Info("Retrieved sources", "count", len(sources))
	return sources, nil
}

// Create implements existing SourceRepository interface
func (s *sourceRepository) Create(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	resp, err := s.httpClient.Post(ctx, "/api/sources", source)
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	var result models.Source
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse source response: %w", err)
	}

	s.logger.Info("Created source", "id", *result.ID, "type", source.Type)
	return &result, nil
}

// CreateFromJSON implements existing SourceRepository interface
func (s *sourceRepository) CreateFromJSON(ctx context.Context, source *models.SourceCreate) (*models.Source, error) {
	return s.Create(ctx, source)
}

// Get implements existing SourceRepository interface
func (s *sourceRepository) Get(ctx context.Context, id string) (*models.Source, error) {
	endpoint := fmt.Sprintf("/api/sources/%s", id)
	resp, err := s.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get source %s: %w", id, err)
	}

	var result models.Source
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse source response: %w", err)
	}

	s.logger.Info("Retrieved source", "id", id)
	return &result, nil
}

// Update implements existing SourceRepository interface
func (s *sourceRepository) Update(ctx context.Context, id string, source *models.SourceUpdate) (*models.Source, error) {
	endpoint := fmt.Sprintf("/api/sources/%s", id)
	resp, err := s.httpClient.Put(ctx, endpoint, source)
	if err != nil {
		return nil, fmt.Errorf("failed to update source %s: %w", id, err)
	}

	var result models.Source
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	s.logger.Info("Updated source", "id", id)
	return &result, nil
}

// Delete implements existing SourceRepository interface
func (s *sourceRepository) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/sources/%s", id)
	_, err := s.httpClient.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete source %s: %w", id, err)
	}

	s.logger.Info("Deleted source", "id", id)
	return nil
}

// GetStatus implements existing SourceRepository interface
func (s *sourceRepository) GetStatus(ctx context.Context, id string) (*models.SourceStatusResponse, error) {
	endpoint := fmt.Sprintf("/api/sources/%s/status", id)
	resp, err := s.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get source status %s: %w", id, err)
	}

	var result models.SourceStatusResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	s.logger.Info("Retrieved source status", "id", id)
	return &result, nil
}

// Download implements existing SourceRepository interface
func (s *sourceRepository) Download(ctx context.Context, id string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/api/sources/%s/download", id)
	resp, err := s.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to download source %s: %w", id, err)
	}

	s.logger.Info("Downloaded source file", "id", id, "size", len(resp.Body))
	return io.NopCloser(bytes.NewReader(resp.Body)), nil
}

// Additional utility methods for file upload (not part of interface)

// CreateFromFile creates a source from file upload (extended method)
func (s *sourceRepository) CreateFromFile(ctx context.Context, filePath string, source *models.SourceCreate) (*models.Source, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add file part
	fileField := "file"
	_, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Add form fields
	if source.Title != nil {
		writer.WriteField("title", *source.Title)
	}
	if len(source.Transformations) > 0 {
		writer.WriteField("transformations", strings.Join(source.Transformations, ","))
	}
	writer.WriteField("type", string(source.Type))
	writer.WriteField("embed", fmt.Sprintf("%t", source.Embed))
	writer.WriteField("delete_source", fmt.Sprintf("%t", source.DeleteSource))
	writer.WriteField("async_processing", fmt.Sprintf("%t", source.AsyncProcessing))

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Send multipart request
	files := map[string]io.Reader{
		fileField: bytes.NewReader(requestBody.Bytes()),
	}
	fields := map[string]string{}

	resp, err := s.httpClient.PostMultipart(ctx, "/api/sources", fields, files)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	var result models.Source
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse upload response: %w", err)
	}

	s.logger.Info("Created source from file", "id", *result.ID, "file", filePath)
	return &result, nil
}