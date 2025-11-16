package services

import (
	"context"
	"fmt"
	"io"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
)

type podcastService struct {
	repo PodcastRepository
}

// NewPodcastService creates a new podcast service
func NewPodcastService(injector do.Injector) (PodcastService, error) {
	repo := do.MustInvoke[PodcastRepository](injector)

	return &podcastService{
		repo: repo,
	}, nil
}

// Interface implementation

func (s *podcastService) Repository() PodcastRepository {
	return s.repo
}

func (s *podcastService) ListEpisodes(ctx context.Context, limit, offset int) (*models.PodcastEpisodesListResponse, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0 // Default offset
	}

	return s.repo.ListEpisodes(ctx, limit, offset)
}

func (s *podcastService) GeneratePodcast(ctx context.Context, query string, sources, notebooks []string, modelID *string, voice, language, style string) (*models.PodcastGenerationResponse, error) {
	// Business logic validation
	if query == "" && len(sources) == 0 && len(notebooks) == 0 {
		return nil, fmt.Errorf("at least one of query, sources, or notebooks must be provided")
	}

	// Validate voice parameter
	validVoices := []string{"male", "female", "neutral"}
	if voice == "" {
		voice = "female" // Default
	} else {
		isValidVoice := false
		for _, v := range validVoices {
			if v == voice {
				isValidVoice = true
				break
			}
		}
		if !isValidVoice {
			return nil, fmt.Errorf("invalid voice. Valid options: %v", validVoices)
		}
	}

	// Validate language code
	if language == "" {
		language = "en" // Default
	} else if len(language) != 2 {
		return nil, fmt.Errorf("language code must be 2 characters (e.g., en, es, fr)")
	}

	// Validate style
	validStyles := []string{"educational", "conversational", "news", "storytelling"}
	if style == "" {
		style = "educational" // Default
	} else {
		isValidStyle := false
		for _, s := range validStyles {
			if s == style {
				isValidStyle = true
				break
			}
		}
		if !isValidStyle {
			return nil, fmt.Errorf("invalid style. Valid options: %v", validStyles)
		}
	}

	// Additional business logic
	if len(sources) > 10 {
		return nil, fmt.Errorf("too many sources (max 10)")
	}
	if len(notebooks) > 5 {
		return nil, fmt.Errorf("too many notebooks (max 5)")
	}

	req := &models.PodcastGenerationRequest{
		Query:       query,
		SourceIDs:   sources,
		NotebookIDs: notebooks,
		ModelID:     modelID,
		Voice:       voice,
		Language:    language,
		Style:       style,
	}

	return s.repo.Generate(ctx, req)
}

func (s *podcastService) GetEpisode(ctx context.Context, episodeID string) (*models.PodcastEpisodeResponse, error) {
	if episodeID == "" {
		return nil, fmt.Errorf("episode ID is required")
	}

	return s.repo.GetEpisode(ctx, episodeID)
}

func (s *podcastService) DownloadEpisodeAudio(ctx context.Context, episodeID string) (io.ReadCloser, error) {
	if episodeID == "" {
		return nil, fmt.Errorf("episode ID is required")
	}

	// Business logic: could check episode status before download
	episode, err := s.repo.GetEpisode(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get episode details: %w", err)
	}

	if episode.AudioURL == "" {
		return nil, fmt.Errorf("no audio file available for this episode")
	}

	return s.repo.DownloadEpisodeAudio(ctx, episodeID)
}

func (s *podcastService) DeleteEpisode(ctx context.Context, episodeID string) error {
	if episodeID == "" {
		return fmt.Errorf("episode ID is required")
	}

	// Business logic: get episode details for logging
	episode, err := s.repo.GetEpisode(ctx, episodeID)
	if err != nil {
		return fmt.Errorf("failed to get episode details: %w", err)
	}

	// Additional business logic: could check if episode is referenced
	// or has dependencies before allowing deletion

	err = s.repo.DeleteEpisode(ctx, episodeID)
	if err != nil {
		return fmt.Errorf("failed to delete episode '%s': %w", episode.Title, err)
	}

	return nil
}

func (s *podcastService) GetJobStatus(ctx context.Context, jobID string) (*models.PodcastJobStatus, error) {
	if jobID == "" {
		return nil, fmt.Errorf("job ID is required")
	}

	return s.repo.GetJobStatus(ctx, jobID)
}

// Additional business logic methods for enhanced functionality

func (s *podcastService) GetEpisodeByTitle(ctx context.Context, title string) (*models.PodcastEpisodeResponse, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Get all episodes and search by title
	episodes, err := s.repo.ListEpisodes(ctx, 100, 0) // Large limit for search
	if err != nil {
		return nil, fmt.Errorf("failed to search episodes: %w", err)
	}

	// Simple title matching - could be enhanced with fuzzy search
	for _, episode := range episodes.Episodes {
		if episode.Title == title {
			return &episode, nil
		}
	}

	return nil, fmt.Errorf("no episode found with title: %s", title)
}

func (s *podcastService) ListEpisodesByLanguage(ctx context.Context, language string, limit, offset int) (*models.PodcastEpisodesListResponse, error) {
	if language == "" {
		return nil, fmt.Errorf("language is required")
	}
	if len(language) != 2 {
		return nil, fmt.Errorf("language code must be 2 characters")
	}

	// Get all episodes and filter by language
	episodes, err := s.repo.ListEpisodes(ctx, 1000, offset) // Large limit for filtering
	if err != nil {
		return nil, fmt.Errorf("failed to get episodes: %w", err)
	}

	// Filter episodes by language
	var filteredEpisodes []models.PodcastEpisodeResponse
	for _, episode := range episodes.Episodes {
		if episode.Language == language {
			filteredEpisodes = append(filteredEpisodes, episode)
		}
	}

	// Apply pagination to filtered results
	start := offset
	if start > len(filteredEpisodes) {
		start = len(filteredEpisodes)
	}
	end := start + limit
	if end > len(filteredEpisodes) {
		end = len(filteredEpisodes)
	}

	var pagedEpisodes []models.PodcastEpisodeResponse
	if start < len(filteredEpisodes) {
		pagedEpisodes = filteredEpisodes[start:end]
	}

	return &models.PodcastEpisodesListResponse{
		Episodes: pagedEpisodes,
		Total:    len(filteredEpisodes),
	}, nil
}