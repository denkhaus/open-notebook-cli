package services

import (
	"context"
	"fmt"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
)

type jobService struct {
	repo JobRepository
}

// NewJobService creates a new job service
func NewJobService(injector do.Injector) (JobService, error) {
	repo := do.MustInvoke[JobRepository](injector)

	return &jobService{
		repo: repo,
	}, nil
}

// Interface implementation

func (s *jobService) Repository() JobRepository {
	return s.repo
}

func (s *jobService) ListJobs(ctx context.Context) (*models.JobsListResponse, error) {
	return s.repo.List(ctx)
}

func (s *jobService) GetJobStatus(ctx context.Context, jobID string) (*models.JobStatus, error) {
	if jobID == "" {
		return nil, fmt.Errorf("job ID is required")
	}

	return s.repo.GetStatus(ctx, jobID)
}

func (s *jobService) CancelJob(ctx context.Context, jobID string) error {
	if jobID == "" {
		return fmt.Errorf("job ID is required")
	}

	// Business logic: check job status before cancellation
	jobStatus, err := s.repo.GetStatus(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job status: %w", err)
	}

	// Only allow cancellation of jobs that are not already completed or failed
	switch jobStatus.Status {
	case "completed":
		return fmt.Errorf("cannot cancel job '%s' - already completed", jobID)
	case "failed":
		return fmt.Errorf("cannot cancel job '%s' - already failed", jobID)
	case "cancelled":
		return fmt.Errorf("job '%s' is already cancelled", jobID)
	}

	return s.repo.Cancel(ctx, jobID)
}

// Additional business logic methods for enhanced job management

func (s *jobService) WatchJobStatus(ctx context.Context, jobID string, interval time.Duration, maxIterations int) (*models.JobStatus, error) {
	if jobID == "" {
		return nil, fmt.Errorf("job ID is required")
	}
	if interval <= 0 {
		interval = 2 * time.Second // Default 2 seconds
	}
	if maxIterations <= 0 {
		maxIterations = 30 // Default 30 iterations (1 minute)
	}

	var currentStatus *models.JobStatus
	var err error

	for i := 0; i < maxIterations; i++ {
		currentStatus, err = s.repo.GetStatus(ctx, jobID)
		if err != nil {
			return nil, fmt.Errorf("failed to get job status on iteration %d: %w", i+1, err)
		}

		// Check if job is in a terminal state
		if currentStatus.Status == "completed" || currentStatus.Status == "failed" || currentStatus.Status == "cancelled" {
			break
		}

		// Sleep before next iteration (unless it's the last iteration)
		if i < maxIterations-1 {
			select {
			case <-ctx.Done():
				return currentStatus, ctx.Err()
			case <-time.After(interval):
				// Continue with next iteration
			}
		}
	}

	return currentStatus, nil
}

func (s *jobService) GetJobsByStatus(ctx context.Context, status string) (*models.JobsListResponse, error) {
	if status == "" {
		return nil, fmt.Errorf("status filter is required")
	}

	// Validate status
	validStatuses := []string{"queued", "running", "completed", "failed", "cancelled"}
	isValidStatus := false
	for _, s := range validStatuses {
		if s == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return nil, fmt.Errorf("invalid status. Valid options: %v", validStatuses)
	}

	// Get all jobs and filter by status
	allJobs, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all jobs: %w", err)
	}

	// Filter jobs by status
	var filteredJobs []models.JobStatus
	for _, job := range allJobs.Jobs {
		if job.Status == status {
			filteredJobs = append(filteredJobs, job)
		}
	}

	return &models.JobsListResponse{
		Jobs: filteredJobs,
	}, nil
}

func (s *jobService) GetRunningJobs(ctx context.Context) (*models.JobsListResponse, error) {
	return s.GetJobsByStatus(ctx, "running")
}

func (s *jobService) GetPendingJobs(ctx context.Context) (*models.JobsListResponse, error) {
	return s.GetJobsByStatus(ctx, "queued")
}

func (s *jobService) GetFailedJobs(ctx context.Context) (*models.JobsListResponse, error) {
	return s.GetJobsByStatus(ctx, "failed")
}

func (s *jobService) CancelAllRunningJobs(ctx context.Context) (int, error) {
	runningJobs, err := s.GetRunningJobs(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get running jobs: %w", err)
	}

	if len(runningJobs.Jobs) == 0 {
		return 0, nil
	}

	successCount := 0
	var lastError error

	for _, job := range runningJobs.Jobs {
		err := s.CancelJob(ctx, job.ID)
		if err != nil {
			lastError = err
			continue
		}
		successCount++
	}

	if lastError != nil && successCount == 0 {
		return successCount, fmt.Errorf("failed to cancel any jobs. Last error: %w", lastError)
	}

	if lastError != nil {
		return successCount, fmt.Errorf("cancelled %d jobs with some errors. Last error: %w", successCount, lastError)
	}

	return successCount, nil
}

func (s *jobService) GetJobStatistics(ctx context.Context) (map[string]int, error) {
	allJobs, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs for statistics: %w", err)
	}

	stats := map[string]int{
		"queued":    0,
		"running":   0,
		"completed": 0,
		"failed":    0,
		"cancelled": 0,
	}

	for _, job := range allJobs.Jobs {
		if count, exists := stats[job.Status]; exists {
			stats[job.Status] = count + 1
		} else {
			stats[job.Status] = 1
		}
	}

	return stats, nil
}