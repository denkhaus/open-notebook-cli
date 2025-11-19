package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

type podcastRepository struct {
	httpClient shared.HTTPClient
	logger     shared.Logger
}

// NewPodcastRepository creates a new podcast repository
func NewPodcastRepository(injector do.Injector) (shared.PodcastRepository, error) {
	httpClient := do.MustInvoke[shared.HTTPClient](injector)
	logger := do.MustInvoke[shared.Logger](injector)

	return &podcastRepository{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// Generate implements PodcastRepository interface
func (p *podcastRepository) Generate(ctx context.Context, req *models.PodcastGenerationRequest) (*models.PodcastGenerationResponse, error) {
	resp, err := p.httpClient.Post(ctx, "/podcasts/generate", req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate podcast: %w", err)
	}

	var result models.PodcastGenerationResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse podcast generation response: %w", err)
	}

	p.logger.Info("Generated podcast", "job_id", result.JobID)
	return &result, nil
}

// GetJobStatus implements PodcastRepository interface
func (p *podcastRepository) GetJobStatus(ctx context.Context, jobID string) (*models.PodcastJobStatus, error) {
	endpoint := fmt.Sprintf("/podcasts/jobs/%s", jobID)
	resp, err := p.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get podcast job status: %w", err)
	}

	var result models.PodcastJobStatus
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse podcast job status response: %w", err)
	}

	p.logger.Info("Retrieved podcast job status", "job_id", jobID, "status", result.Status)
	return &result, nil
}

// ListEpisodes implements PodcastRepository interface
func (p *podcastRepository) ListEpisodes(ctx context.Context, limit, offset int) (*models.PodcastEpisodesListResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("limit", strconv.Itoa(limit))
	queryParams.Set("offset", strconv.Itoa(offset))

	endpoint := "/podcasts/episodes?" + queryParams.Encode()
	resp, err := p.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list podcast episodes: %w", err)
	}

	var result models.PodcastEpisodesListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse podcast episodes response: %w", err)
	}

	p.logger.Info("Retrieved podcast episodes", "count", len(result.Episodes), "total", result.Total)
	return &result, nil
}

// GetEpisode implements PodcastRepository interface
func (p *podcastRepository) GetEpisode(ctx context.Context, episodeID string) (*models.PodcastEpisodeResponse, error) {
	endpoint := fmt.Sprintf("/podcasts/episodes/%s", episodeID)
	resp, err := p.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get podcast episode: %w", err)
	}

	var result models.PodcastEpisodeResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse podcast episode response: %w", err)
	}

	p.logger.Info("Retrieved podcast episode", "episode_id", episodeID)
	return &result, nil
}

// DownloadEpisodeAudio implements PodcastRepository interface
func (p *podcastRepository) DownloadEpisodeAudio(ctx context.Context, episodeID string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/podcasts/episodes/%s/audio", episodeID)
	resp, err := p.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to download podcast audio: %w", err)
	}

	p.logger.Info("Downloaded podcast episode audio", "episode_id", episodeID, "size", len(resp.Body))
	return io.NopCloser(bytes.NewReader(resp.Body)), nil
}

// DeleteEpisode implements PodcastRepository interface
func (p *podcastRepository) DeleteEpisode(ctx context.Context, episodeID string) error {
	endpoint := fmt.Sprintf("/podcasts/episodes/%s", episodeID)
	_, err := p.httpClient.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete podcast episode %s: %w", episodeID, err)
	}

	p.logger.Info("Deleted podcast episode", "episode_id", episodeID)
	return nil
}

// Extended methods not part of interface

// ListEpisodesWithLanguage filters episodes by language
func (p *podcastRepository) ListEpisodesWithLanguage(ctx context.Context, language string, limit, offset int) (*models.PodcastEpisodesListResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("limit", strconv.Itoa(limit))
	queryParams.Set("offset", strconv.Itoa(offset))
	if language != "" {
		queryParams.Set("language", language)
	}

	endpoint := "/podcasts/episodes?" + queryParams.Encode()
	resp, err := p.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list podcast episodes by language %s: %w", language, err)
	}

	var result models.PodcastEpisodesListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse podcast episodes response: %w", err)
	}

	p.logger.Info("Retrieved podcast episodes by language", "language", language, "count", len(result.Episodes))
	return &result, nil
}
