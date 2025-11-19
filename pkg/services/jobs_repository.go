package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

type jobRepository struct {
	httpClient shared.HTTPClient
	logger     shared.Logger
}

// NewJobRepository creates a new job repository
func NewJobRepository(injector do.Injector) (shared.JobRepository, error) {
	httpClient := do.MustInvoke[shared.HTTPClient](injector)
	logger := do.MustInvoke[shared.Logger](injector)

	return &jobRepository{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// List implements JobRepository interface
func (j *jobRepository) List(ctx context.Context) (*models.JobsListResponse, error) {
	endpoint := "/commands/jobs"
	resp, err := j.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, errors.FailedToList("jobs", err)
	}

	var result models.JobsListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, errors.FailedToDecode("jobs response", err)
	}

	j.logger.Info("Retrieved jobs", "count", len(result.Jobs))
	return &result, nil
}

// GetStatus implements JobRepository interface
func (j *jobRepository) GetStatus(ctx context.Context, jobID string) (*models.JobStatus, error) {
	endpoint := fmt.Sprintf("/commands/jobs/%s", jobID)
	resp, err := j.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, errors.FailedToGet("job status", err)
	}

	var result models.JobStatus
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, errors.FailedToDecode("job status response", err)
	}

	j.logger.Info("Retrieved job status", "job_id", jobID, "status", result.Status)
	return &result, nil
}

// Cancel implements JobRepository interface
func (j *jobRepository) Cancel(ctx context.Context, jobID string) error {
	endpoint := fmt.Sprintf("/commands/jobs/%s", jobID)
	_, err := j.httpClient.Delete(ctx, endpoint)
	if err != nil {
		return errors.FailedToCancel(fmt.Sprintf("job %s", jobID), err)
	}

	j.logger.Info("Cancelled job", "job_id", jobID)
	return nil
}

// Extended methods not part of interface

// ListWithStatus filters jobs by status
func (j *jobRepository) ListWithStatus(ctx context.Context, status string) (*models.JobsListResponse, error) {
	queryParams := url.Values{}
	if status != "" {
		queryParams.Set("status", status)
	}

	endpoint := "/commands/jobs?" + queryParams.Encode()
	resp, err := j.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs by status %s: %w", status, err)
	}

	var result models.JobsListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse jobs response: %w", err)
	}

	j.logger.Info("Retrieved jobs by status", "status", status, "count", len(result.Jobs))
	return &result, nil
}
