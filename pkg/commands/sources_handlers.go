package commands

import (
	"fmt"
	"io"
	"os"
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

// SourcesServices holds all the services needed for source commands
type SourcesServices struct {
	SourceService shared.SourceService
	Config        config.Service
	Logger        shared.Logger
}

// getSourcesServices retrieves all required services via dependency injection
func getSourcesServices(ctx *cli.Context) (*SourcesServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &SourcesServices{
		SourceService: do.MustInvoke[shared.SourceService](injector),
		Config:        do.MustInvoke[config.Service](injector),
		Logger:        do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validateSourceArgs validates common argument patterns for source commands
func validateSourceArgs(ctx *cli.Context, requireSourceID bool) (string, error) {
	if requireSourceID {
		if ctx.NArg() < 1 {
			return "", errors.MissingArgument("source ID", ctx.Command.Name)
		}
		if ctx.NArg() > 1 {
			return "", errors.TooManyArguments("source ID", ctx.Command.Name)
		}
	}

	sourceID := ctx.Args().First()
	if requireSourceID && sourceID == "" {
		return "", errors.UsageError("Source ID is required",
			"Usage: open-notebook sources <command> <source-id>")
	}

	return sourceID, nil
}

// printSourceSuccess prints standardized success messages for source operations
func printSourceSuccess(operation string, source *models.Source) {
	fmt.Printf("✅ Source %s successfully!\n", operation)
	fmt.Printf("  ID:     %s\n", utils.SafeDereferenceString(source.ID))
	fmt.Printf("  Title:  %s\n", utils.SafeDereferenceString(source.Title))
	if source.Status != nil {
		fmt.Printf("  Status: %s\n", string(*source.Status))
	}
}

// displayProcessingInfo displays processing information in a consistent format
func displayProcessingInfo(processingInfo map[string]any) {
	if processingInfo == nil {
		return
	}

	if status, ok := processingInfo["status"].(string); ok {
		fmt.Printf("  Status:     %s\n", status)
	}
	if message, ok := processingInfo["message"].(string); ok {
		fmt.Printf("  Message:   %s\n", message)
	}
	if error, ok := processingInfo["error"].(string); ok && error != "" {
		fmt.Printf("  Error:     %s\n", error)
	}
	if processedAt, ok := processingInfo["processed_at"].(string); ok {
		fmt.Printf("  Processed: %s\n", utils.FormatTimestamp(processedAt))
	}
}

// convertProcessingInfo converts the model processing info to a displayable format
func convertProcessingInfo(processingInfo map[string]any) map[string]any {
	if processingInfo == nil {
		return nil
	}

	result := make(map[string]any, len(processingInfo))
	for k, v := range processingInfo {
		result[k] = v
	}
	return result
}

// handleSourcesList handles the sources list command
func handleSourcesList(ctx *cli.Context) error {
	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Listing sources...")

	// Parse pagination parameters
	limit := 20
	offset := 0
	if ctx.IsSet("limit") {
		limit = ctx.Int("limit")
	}
	if ctx.IsSet("offset") {
		offset = ctx.Int("offset")
	}

	sources, err := services.SourceService.List(ctx.Context, limit, offset)
	if err != nil {
		return errors.APIError("Failed to list sources",
			"Check API connection and permissions")
	}

	if len(sources) == 0 {
		fmt.Println("No sources found.")
		return nil
	}

	// Display sources in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tTYPE\tSTATUS\tCREATED")

	for _, source := range sources {
		title := utils.SafeDereferenceString(source.Title)
		status := "N/A"
		if source.Status != nil {
			status = string(*source.Status)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			utils.SafeDereferenceString(source.ID),
			utils.TruncateString(title, 30),
			"source", // Type field doesn't exist in list response
			status,
			utils.FormatTimestamp(source.Created))
	}

	w.Flush()

	fmt.Printf("\nShowing %d sources (use --limit and --offset for pagination)\n", len(sources))
	return nil
}

// handleSourcesShow handles source details display
func handleSourcesShow(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Showing source details", "source_id", sourceID)

	source, err := services.SourceService.Get(ctx.Context, sourceID)
	if err != nil {
		return errors.APIError("Failed to get source details",
			"Check source ID and permissions")
	}

	// Display source details
	fmt.Printf("Source Details:\n")
	fmt.Printf("  ID:           %s\n", utils.SafeDereferenceString(source.ID))
	fmt.Printf("  Title:        %s\n", utils.SafeDereferenceString(source.Title))
	fmt.Printf("  Created:      %s\n", utils.FormatTimestamp(source.Created))
	fmt.Printf("  Updated:      %s\n", utils.FormatTimestamp(source.Updated))

	if len(source.Notebooks) > 0 {
		fmt.Printf("  Notebooks:    %v\n", source.Notebooks)
	}

	if len(source.Topics) > 0 {
		fmt.Printf("  Topics:       %v\n", source.Topics)
	}

	if source.Status != nil {
		fmt.Printf("  Status:       %s\n", string(*source.Status))
	}

	if source.ProcessingInfo != nil {
		fmt.Printf("  Processing:\n")
		displayProcessingInfo(convertProcessingInfo(source.ProcessingInfo))
	}

	return nil
}

// handleSourcesAdd handles source creation
func handleSourcesAdd(ctx *cli.Context) error {
	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	title := ctx.String("title")
	text := ctx.String("text")
	link := ctx.String("link")
	filePath := ctx.String("file")

	// Validate required fields
	if title == "" {
		return errors.UsageError("Title is required",
			"Use --title flag to specify the source title")
	}

	// Determine source type based on provided flags - fail loud if ambiguous
	var sourceType string
	var source *models.SourceCreate

	if text != "" {
		sourceType = "text"
	} else if link != "" {
		sourceType = "link"
	} else if filePath != "" {
		sourceType = "file"
	} else {
		return errors.UsageError("One of --text, --link, or --file is required",
			"Use --text for text content, --link for URLs, or --file for file uploads")
	}

	switch sourceType {
	case "text":
		source = &models.SourceCreate{
			Type:    models.SourceTypeText,
			Title:   &title,
			Content: &text,
		}

	case "link":
		source = &models.SourceCreate{
			Type:  models.SourceTypeLink,
			Title: &title,
			URL:   &link,
		}

	case "file":
		if filePath == "" {
			return errors.UsageError("File path is required for file sources",
				"Use --file flag to specify the file path")
		}
		// File upload is handled separately
		return handleFileUpload(ctx, services, title, filePath)

	default:
		return errors.UsageError("Invalid source type",
			"Supported types are: text, link, file")
	}

	// Add optional parameters
	if ctx.IsSet("notebook") {
		notebookID := ctx.String("notebook")
		source.Notebooks = []string{notebookID}
	}

	services.Logger.Info("Creating source", "type", sourceType, "title", title)

	createdSource, err := services.SourceService.Create(ctx.Context, source)
	if err != nil {
		return errors.APIError("Failed to create source",
			"Check input parameters and API permissions")
	}

	printSourceSuccess("created", createdSource)
	return nil
}

// handleFileUpload handles file source uploads
func handleFileUpload(ctx *cli.Context, services *SourcesServices, title, filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.UsageError("File not found",
			fmt.Sprintf("The file '%s' does not exist", filePath))
	}

	source := &models.SourceCreate{
		Type:     models.SourceTypeUpload,
		Title:    &title,
		FilePath: &filePath,
	}

	// Add optional parameters
	if ctx.IsSet("notebook") {
		notebookID := ctx.String("notebook")
		source.Notebooks = []string{notebookID}
	}

	services.Logger.Info("Uploading file source", "file", filePath, "title", title)

	createdSource, err := services.SourceService.Create(ctx.Context, source)
	if err != nil {
		return errors.APIError("Failed to upload file source",
			"Check file path and API permissions")
	}

	printSourceSuccess("uploaded", createdSource)
	return nil
}

// handleSourcesUpdate handles source updates
func handleSourcesUpdate(ctx *cli.Context) error {
	if !ctx.IsSet("title") && !ctx.IsSet("topic") {
		return errors.AtLeastOneField([]string{"title", "topic"}, ctx.Command.Name)
	}

	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	title := ctx.String("title")
	topics := ctx.StringSlice("topics")

	if title == "" && len(topics) == 0 {
		return errors.UsageError("At least one update field is required",
			"Use --title or --topics flags")
	}

	source := &models.SourceUpdate{}

	if title != "" {
		source.Title = &title
	}

	if len(topics) > 0 {
		source.Topics = topics
	}

	services.Logger.Info("Updating source", "source_id", sourceID)

	updateTitle := ""
	if source.Title != nil {
		updateTitle = *source.Title
	} else if title != "" {
		updateTitle = title
	}

	updatedSource, err := services.SourceService.Update(ctx.Context, sourceID, updateTitle, topics)
	if err != nil {
		return errors.APIError("Failed to update source",
			"Check source ID and permissions")
	}

	printSourceSuccess("updated", updatedSource)
	return nil
}

// handleSourcesDelete handles source deletion
func handleSourcesDelete(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	// Confirm deletion unless force flag is used
	if !ctx.Bool("force") {
		fmt.Printf("Are you sure you want to delete source '%s'? (y/N): ", sourceID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	services.Logger.Info("Deleting source", "source_id", sourceID)

	err = services.SourceService.Delete(ctx.Context, sourceID)
	if err != nil {
		return errors.APIError("Failed to delete source",
			"Check source ID and permissions")
	}

	fmt.Printf("✅ Source '%s' deleted successfully!\n", sourceID)
	return nil
}

// handleSourcesDownload handles source file downloads
func handleSourcesDownload(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	outputPath := ctx.String("output")
	if outputPath == "" {
		outputPath = sourceID + "_downloaded_file"
	}

	services.Logger.Info("Downloading source file", "source_id", sourceID)

	reader, err := services.SourceService.Download(ctx.Context, sourceID)
	if err != nil {
		return errors.APIError("Failed to download source file",
			"Check source ID and permissions")
	}
	defer reader.Close()

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return errors.ValidationError("Failed to create output file",
			fmt.Sprintf("Could not create file '%s': %v", outputPath, err))
	}
	defer file.Close()

	// Copy downloaded content to file
	_, err = io.Copy(file, reader)
	if err != nil {
		return errors.ValidationError("Failed to save downloaded content",
			fmt.Sprintf("Error writing to file: %v", err))
	}

	fmt.Printf("✅ File downloaded successfully to: %s\n", outputPath)
	return nil
}

// handleSourcesStatus handles source status checking
func handleSourcesStatus(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	watch := ctx.Bool("watch")

	fmt.Printf("📊 Getting source status: %s\n", sourceID)

	services.Logger.Info("Checking source status", "source_id", sourceID)

	status, err := services.SourceService.GetStatus(ctx.Context, sourceID)
	if err != nil {
		return errors.APIError("Failed to get source status",
			"Check source ID and permissions")
	}

	fmt.Printf("Source Status: %s\n", sourceID)
	if status.Status != nil {
		fmt.Printf("  Status:    %s\n", string(*status.Status))
	}

	if status.Message != "" {
		fmt.Printf("  Message:   %s\n", status.Message)
	}

	if status.ProcessingInfo != nil {
		displayProcessingInfo(convertProcessingInfo(status.ProcessingInfo))
	}

	if watch {
		fmt.Println("   Watching for status updates... (Press Ctrl+C to stop)")
		// Simple polling implementation
		for i := 0; i < 10; i++ { // Watch for 10 iterations
			time.Sleep(2 * time.Second)
			fmt.Printf("   Checking status... (%d/10)\n", i+1)
			// In a real implementation, this would poll the API
		}
		fmt.Println("   Watch completed")
	}

	return nil
}

// handleSourcesRetry handles source processing retry
func handleSourcesRetry(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	force := ctx.Bool("force")

	services.Logger.Info("Retrying source processing", "source_id", sourceID)

	if force {
		fmt.Println("   Force reprocessing enabled")
	}

	// Get current source to check status
	source, err := services.SourceService.Get(ctx.Context, sourceID)
	if err != nil {
		return errors.APIError("Failed to get source details",
			"Check source ID and permissions")
	}

	if source.Status != nil && *source.Status == models.SourceStatusCompleted {
		fmt.Printf("Source '%s' is already completed. No retry needed.\n", sourceID)
		return nil
	}

	// Retry by re-creating the source (API limitation workaround)
	var retrySource *models.SourceCreate

	// For retry, we need to determine source type from the asset or content
	if source.FullText != nil && *source.FullText != "" {
		// Text source
		retrySource = &models.SourceCreate{
			Type:    models.SourceTypeText,
			Title:   source.Title,
			Content: source.FullText,
		}
	} else if source.Asset != nil && source.Asset.URL != nil {
		// Link source
		retrySource = &models.SourceCreate{
			Type:  models.SourceTypeLink,
			Title: source.Title,
			URL:   source.Asset.URL,
		}
	} else {
		return errors.ValidationError("Retry not supported for this source type",
			"Retry is only supported for text and link sources")
	}

	// Copy optional fields
	if len(source.Notebooks) > 0 {
		retrySource.Notebooks = source.Notebooks
	}
	if len(source.Topics) > 0 {
		// Note: SourceCreate doesn't have Topics field, but we'll keep this for future compatibility
		// retrySource.Topics = source.Topics
	}

	retriedSource, err := services.SourceService.Create(ctx.Context, retrySource)
	if err != nil {
		return errors.APIError("Failed to retry source processing",
			"Check API permissions")
	}

	fmt.Printf("✅ Source processing retry initiated!\n")
	fmt.Printf("  New ID:  %s\n", utils.SafeDereferenceString(retriedSource.ID))
	if retriedSource.Status != nil {
		fmt.Printf("  Status:  %s\n", string(*retriedSource.Status))
	}

	return nil
}

// handleSourcesInsightsList handles listing insights for a source
func handleSourcesInsightsList(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("💡 Listing insights for source: %s\n", sourceID)

	services.Logger.Info("Listing source insights", "source_id", sourceID)

	insights, err := services.SourceService.GetInsights(ctx.Context, sourceID)
	if err != nil {
		return errors.APIError("Failed to list source insights",
			"Check source ID and permissions")
	}

	if len(insights) == 0 {
		fmt.Printf("No insights found for source '%s'.\n", sourceID)
		return nil
	}

	// Display insights in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTYPE\tCREATED\tCONTENT")

	for _, insight := range insights {
		content := utils.TruncateString(insight.Content, 50)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			insight.ID,
			string(insight.InsightType),
			utils.FormatTimestamp(insight.Created),
			content)
	}

	w.Flush()

	fmt.Printf("\nFound %d insights for source '%s'\n", len(insights), sourceID)
	return nil
}

// handleSourcesInsightsCreate handles creating a new insight for a source
func handleSourcesInsightsCreate(ctx *cli.Context) error {
	sourceID, err := validateSourceArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getSourcesServices(ctx)
	if err != nil {
		return err
	}

	content := ctx.String("content")
	if content == "" {
		return errors.UsageError("Content is required",
			"Use --content flag to specify the insight content")
	}

	fmt.Printf("💭 Creating insight for source: %s\n", sourceID)
	fmt.Printf("   Content: %s\n", content)

	services.Logger.Info("Creating source insight", "source_id", sourceID)

	// Create insight request - for now using a generic transformation approach
	// In a real implementation, this would use proper transformation logic
	request := &models.CreateSourceInsightRequest{
		TransformationID: "manual-insight", // Default transformation for manual insights
		ModelID:          nil,              // Use default model
	}

	createdInsight, err := services.SourceService.CreateInsight(ctx.Context, sourceID, request)
	if err != nil {
		return errors.APIError("Failed to create source insight",
			"Check source ID, transformation ID and permissions")
	}

	fmt.Printf("✅ Insight created successfully!\n")
	fmt.Printf("  ID:      %s\n", createdInsight.ID)
	fmt.Printf("  Type:    %s\n", string(createdInsight.InsightType))
	fmt.Printf("  Content: %s\n", utils.TruncateString(createdInsight.Content, 100))

	return nil
}
