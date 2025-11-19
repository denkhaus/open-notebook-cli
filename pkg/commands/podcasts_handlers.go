package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// PodcastServices holds all the services needed for podcast commands
type PodcastServices struct {
	PodcastRepository shared.PodcastRepository
	Config            config.Service
	Logger            shared.Logger
}

// getPodcastServices retrieves all required services via dependency injection
func getPodcastServices(ctx *cli.Context) (*PodcastServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &PodcastServices{
		PodcastRepository: do.MustInvoke[shared.PodcastRepository](injector),
		Config:            do.MustInvoke[config.Service](injector),
		Logger:            do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validatePodcastGenerateArgs validates podcast generation arguments
func validatePodcastGenerateArgs(ctx *cli.Context) (*models.PodcastGenerationRequest, error) {
	sources := ctx.StringSlice("sources")
	notebooks := ctx.StringSlice("notebooks")
	query := ctx.String("query")

	// Validate that at least one content source is provided
	if len(sources) == 0 && len(notebooks) == 0 && query == "" {
		return nil, errors.UsageError("Content source required",
			"Provide at least one of: --query, --sources, or --notebooks")
	}

	// Validate language code
	language := ctx.String("language")
	if language != "" && len(language) != 2 {
		return nil, errors.UsageError("Invalid language code",
			"Language code must be 2 characters (e.g., en, es, fr)")
	}

	// Validate voice
	voice := ctx.String("voice")
	if voice == "" {
		voice = "female"
	}

	// Validate style
	style := ctx.String("style")
	if style == "" {
		style = "educational"
	}

	var modelID *string
	if ctx.IsSet("model") {
		m := ctx.String("model")
		modelID = &m
	}

	return &models.PodcastGenerationRequest{
		SourceIDs:   sources,
		NotebookIDs: notebooks,
		Query:       query,
		ModelID:     modelID,
		Voice:       voice,
		Language:    language,
		Style:       style,
	}, nil
}

// validateEpisodeArgs validates episode ID arguments
func validateEpisodeArgs(ctx *cli.Context, requireEpisodeID bool) (string, error) {
	if requireEpisodeID {
		if ctx.NArg() < 1 {
			return "", fmt.Errorf("❌ Error: Missing episode ID")
		}
		if ctx.NArg() > 1 {
			return "", fmt.Errorf("❌ Error: Too many arguments. Expected only episode ID")
		}
	}

	episodeID := ctx.Args().First()
	if requireEpisodeID && episodeID == "" {
		return "", errors.UsageError("Episode ID is required",
			"Usage: open-notebook podcast episodes <command> <episode-id>")
	}

	return episodeID, nil
}

// TODO: move to shared package to make it reusable
// formatDuration formats duration in seconds to human readable format
func formatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	}
	minutes := int(seconds / 60)
	remainingSeconds := int(seconds) % 60
	return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
}

// printPodcastSuccess prints standardized success messages for podcast operations
func printPodcastSuccess(operation string, details interface{}) {
	fmt.Printf("✅ Podcast %s successfully!\n", operation)
	switch v := details.(type) {
	case string:
		fmt.Printf("  %s: %s\n", "ID", v)
	case *models.PodcastGenerationResponse:
		fmt.Printf("  Job ID: %s\n", v.JobID)
		fmt.Printf("  Message: %s\n", v.Message)
	case *models.PodcastEpisodeResponse:
		fmt.Printf("  Episode ID: %s\n", v.ID)
		fmt.Printf("  Title: %s\n", v.Title)
		fmt.Printf("  Duration: %s\n", formatDuration(v.Duration))
	}
}

// handlePodcastGenerate handles podcast generation command
func handlePodcastGenerate(ctx *cli.Context) error {
	services, err := getPodcastServices(ctx)
	if err != nil {
		return err
	}

	req, err := validatePodcastGenerateArgs(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Generating podcast",
		"query", req.Query,
		"sources_count", len(req.SourceIDs),
		"notebooks_count", len(req.NotebookIDs),
		"voice", req.Voice,
		"language", req.Language,
		"style", req.Style)

	response, err := services.PodcastRepository.Generate(ctx.Context, req)
	if err != nil {
		return errors.APIError("Failed to generate podcast",
			"Check content sources and API permissions")
	}

	printPodcastSuccess("generation started", response)

	// Watch progress if requested
	if ctx.Bool("watch") {
		return watchPodcastGeneration(ctx.Context, services, response.JobID)
	}

	fmt.Printf("💡 Use 'onb jobs status %s' to check progress\n", response.JobID)
	return nil
}

// watchPodcastGeneration watches podcast generation progress
func watchPodcastGeneration(ctx context.Context, services *PodcastServices, jobID string) error {
	fmt.Printf("🔄 Watching podcast generation progress...\n")
	fmt.Printf("Press Ctrl+C to stop watching\n\n")

	for i := 0; i < 30; i++ { // Watch for up to 30 iterations (~1 minute)
		jobStatus, err := services.PodcastRepository.GetJobStatus(ctx, jobID)
		if err != nil {
			fmt.Printf("\n❌ Error checking job status: %v\n", err)
			return err
		}

		// Display progress
		progress := "N/A"
		if jobStatus.Progress != nil {
			progress = fmt.Sprintf("%.0f%%", *jobStatus.Progress*100)
		}

		fmt.Printf("\r📊 Status: %s | Progress: %s | Checking... (%d/30)",
			jobStatus.Status, progress, i+1)

		// Check if job is completed or failed
		if jobStatus.Status == "completed" {
			fmt.Printf("\n✅ Podcast generation completed!\n")
			if jobStatus.EpisodeID != nil {
				fmt.Printf("   Episode ID: %s\n", *jobStatus.EpisodeID)
				fmt.Printf("   💡 Use 'onb podcast episodes show %s' to view details\n", *jobStatus.EpisodeID)
			}
			break
		}
		if jobStatus.Status == "failed" {
			fmt.Printf("\n❌ Podcast generation failed!\n")
			if jobStatus.Message != nil {
				fmt.Printf("   Error: %s\n", *jobStatus.Message)
			}
			break
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\n🏁 Watch completed\n")
	return nil
}

// handlePodcastEpisodesList handles episodes list command
func handlePodcastEpisodesList(ctx *cli.Context) error {
	services, err := getPodcastServices(ctx)
	if err != nil {
		return err
	}

	limit := ctx.Int("limit")
	offset := ctx.Int("offset")
	language := ctx.String("language")

	services.Logger.Info("Listing podcast episodes",
		"limit", limit, "offset", offset, "language", language)

	var episodesList *models.PodcastEpisodesListResponse
	if language != "" {
		// Use extended method with language filter (if implemented)
		episodesList, err = services.PodcastRepository.ListEpisodes(ctx.Context, limit, offset)
		if err != nil {
			return errors.APIError("Failed to list podcast episodes",
				"Check API connection and permissions")
		}
		// Filter by language manually
		filteredEpisodes := []models.PodcastEpisodeResponse{}
		for _, episode := range episodesList.Episodes {
			if episode.Language == language {
				filteredEpisodes = append(filteredEpisodes, episode)
			}
		}
		episodesList.Episodes = filteredEpisodes
	} else {
		episodesList, err = services.PodcastRepository.ListEpisodes(ctx.Context, limit, offset)
		if err != nil {
			return errors.APIError("Failed to list podcast episodes",
				"Check API connection and permissions")
		}
	}

	if len(episodesList.Episodes) == 0 {
		if language != "" {
			fmt.Printf("No podcast episodes found for language '%s'.\n", language)
		} else {
			fmt.Println("No podcast episodes found.")
		}
		return nil
	}

	// Display episodes in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tDURATION\tLANGUAGE\tVOICE\tCREATED")

	for _, episode := range episodesList.Episodes {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			utils.TruncateString(episode.ID, 12),
			utils.TruncateString(episode.Title, 30),
			formatDuration(episode.Duration),
			episode.Language,
			episode.Voice,
			utils.FormatTimestamp(episode.Created))
	}

	w.Flush()

	fmt.Printf("\nShowing %d episodes (Total: %d)\n", len(episodesList.Episodes), episodesList.Total)
	return nil
}

// handlePodcastEpisodesShow handles episode details display
func handlePodcastEpisodesShow(ctx *cli.Context) error {
	episodeID, err := validateEpisodeArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getPodcastServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Getting podcast episode details", "episode_id", episodeID)

	episode, err := services.PodcastRepository.GetEpisode(ctx.Context, episodeID)
	if err != nil {
		return errors.APIError("Failed to get episode details",
			"Check episode ID and API permissions")
	}

	// Display episode details
	fmt.Printf("Podcast Episode Details:\n")
	fmt.Printf("  ID:          %s\n", episode.ID)
	fmt.Printf("  Title:       %s\n", episode.Title)
	fmt.Printf("  Description: %s\n", episode.Description)
	fmt.Printf("  Duration:    %s\n", formatDuration(episode.Duration))
	fmt.Printf("  Language:    %s\n", episode.Language)
	fmt.Printf("  Voice:       %s\n", episode.Voice)
	fmt.Printf("  Style:       %s\n", episode.Style)
	fmt.Printf("  Model ID:    %s\n", episode.ModelID)
	fmt.Printf("  Job ID:      %s\n", episode.JobID)
	fmt.Printf("  Audio URL:   %s\n", episode.AudioURL)
	fmt.Printf("  Created:     %s\n", utils.FormatTimestamp(episode.Created))
	fmt.Printf("  Updated:     %s\n", utils.FormatTimestamp(episode.Updated))

	return nil
}

// handlePodcastEpisodesDownload handles episode audio download
func handlePodcastEpisodesDownload(ctx *cli.Context) error {
	episodeID, err := validateEpisodeArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getPodcastServices(ctx)
	if err != nil {
		return err
	}

	outputPath := ctx.String("output")
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s.mp3", episodeID)
	}

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.ConfigError("Failed to create output directory",
			fmt.Sprintf("Directory: %s, Error: %v", dir, err))
	}

	services.Logger.Info("Downloading podcast episode audio",
		"episode_id", episodeID, "output_path", outputPath)

	fmt.Printf("📥 Downloading podcast episode: %s\n", episodeID)
	fmt.Printf("   Output: %s\n", outputPath)

	audioReader, err := services.PodcastRepository.DownloadEpisodeAudio(ctx.Context, episodeID)
	if err != nil {
		return errors.APIError("Failed to download episode audio",
			"Check episode ID and API permissions")
	}
	defer audioReader.Close()

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return errors.ConfigError("Failed to create output file",
			fmt.Sprintf("File: %s, Error: %v", outputPath, err))
	}
	defer outputFile.Close()

	// Copy audio data to file
	written, err := outputFile.ReadFrom(audioReader)
	if err != nil {
		return errors.NetworkError("Failed to save audio file",
			fmt.Sprintf("Error: %v", err))
	}

	fmt.Printf("✅ Download completed!\n")
	fmt.Printf("   File size: %.1f MB\n", float64(written)/1024/1024)
	fmt.Printf("   Location: %s\n", outputPath)

	return nil
}

// handlePodcastEpisodesDelete handles episode deletion
func handlePodcastEpisodesDelete(ctx *cli.Context) error {
	episodeID, err := validateEpisodeArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getPodcastServices(ctx)
	if err != nil {
		return err
	}

	// Get episode details first for confirmation
	episode, err := services.PodcastRepository.GetEpisode(ctx.Context, episodeID)
	if err != nil {
		return errors.APIError("Failed to get episode details",
			"Check episode ID and API permissions")
	}

	// Confirm deletion unless force flag is used
	if !ctx.Bool("force") {
		fmt.Printf("⚠️  Are you sure you want to delete episode '%s'? (%.0fs)\n",
			episode.Title, episode.Duration)
		fmt.Printf("   This will permanently delete the episode and its audio file.\n")
		fmt.Printf("   Continue? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Deletion cancelled")
			return nil
		}
	}

	services.Logger.Info("Deleting podcast episode", "episode_id", episodeID)

	err = services.PodcastRepository.DeleteEpisode(ctx.Context, episodeID)
	if err != nil {
		return errors.APIError("Failed to delete episode",
			"Check episode ID and permissions")
	}

	fmt.Printf("✅ Episode '%s' deleted successfully!\n", episode.Title)
	return nil
}
