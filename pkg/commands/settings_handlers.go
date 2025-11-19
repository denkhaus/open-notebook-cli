package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// SettingsServices holds all the services needed for settings commands
type SettingsServices struct {
	SettingsService shared.SettingsRepository
	Config          config.Service
	Logger          shared.Logger
}

// getSettingsServices retrieves all required services via dependency injection
func getSettingsServices(ctx *cli.Context) (*SettingsServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &SettingsServices{
		SettingsService: do.MustInvoke[shared.SettingsRepository](injector),
		Config:          do.MustInvoke[config.Service](injector),
		Logger:          do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validateSettingsArgs validates that at least one setting is specified for updates
func validateSettingsArgs(ctx *cli.Context) error {
	if !ctx.IsSet("engine") && !ctx.IsSet("url-engine") && !ctx.IsSet("embed") &&
		!ctx.IsSet("auto-delete") && !ctx.IsSet("youtube-languages") {
		return errors.UsageError("No settings specified",
			"Use at least one of: --engine, --url-engine, --embed, --auto-delete, --youtube-languages")
	}
	return nil
}

// parseYesNoValue parses yes/no string values
func parseYesNoValue(value string) (models.YesNoDecision, error) {
	switch strings.ToLower(value) {
	case "yes", "y", "true", "1":
		return models.YesNoDecisionYes, nil
	case "no", "n", "false", "0":
		return models.YesNoDecisionNo, nil
	default:
		return "", errors.ValidationError("Invalid yes/no value",
			fmt.Sprintf("Expected 'yes' or 'no', got '%s'", value))
	}
}

// formatSettingsTable formats settings for table display
func formatSettingsTable(settings *models.SettingsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SETTING\tVALUE\tDESCRIPTION")

	fmt.Fprintf(w, "Document Engine\t%s\tContent processing engine for documents\n", string(settings.DefaultContentProcessingEngineDoc))
	fmt.Fprintf(w, "URL Engine\t%s\tContent processing engine for web pages\n", string(settings.DefaultContentProcessingEngineURL))
	fmt.Fprintf(w, "Embedding\t%s\tWhen to perform embeddings\n", string(settings.DefaultEmbeddingOption))
	fmt.Fprintf(w, "Auto Delete Files\t%s\tAutomatically delete files after processing\n", string(settings.AutoDeleteFiles))

	if len(settings.YoutubePreferredLanguages) > 0 {
		fmt.Fprintf(w, "YouTube Languages\t%s\tPreferred languages for YouTube content\n", strings.Join(settings.YoutubePreferredLanguages, ", "))
	} else {
		fmt.Fprintf(w, "YouTube Languages\t(N/A)\tNo preferred languages set\n")
	}

	w.Flush()
}

// formatSettingsJSON formats settings for JSON display
func formatSettingsJSON(settings *models.SettingsResponse) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return errors.ValidationError("Failed to format settings as JSON",
			fmt.Sprintf("JSON marshaling error: %v", err))
	}

	fmt.Println(string(data))
	return nil
}

// handleSettingsGet handles the settings get command
func handleSettingsGet(ctx *cli.Context) error {
	services, err := getSettingsServices(ctx)
	if err != nil {
		return err
	}

	format := ctx.String("format")

	services.Logger.Info("Getting application settings")

	settings, err := services.SettingsService.Get(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to get settings",
			"Check API connection and permissions")
	}

	fmt.Println("⚙️  Application Settings:")

	switch format {
	case "json":
		return formatSettingsJSON(settings)
	case "table":
		formatSettingsTable(settings)
	default:
		return errors.UsageError("Invalid format",
			fmt.Sprintf("Supported formats: table, json (got: %s)", format))
	}

	return nil
}

// handleSettingsSet handles the settings set command
func handleSettingsSet(ctx *cli.Context) error {
	services, err := getSettingsServices(ctx)
	if err != nil {
		return err
	}

	// Validate that at least one setting is being updated
	if err := validateSettingsArgs(ctx); err != nil {
		return err
	}

	services.Logger.Info("Updating application settings")

	// Build settings update request
	update := &models.SettingsUpdate{}

	if ctx.IsSet("engine") {
		engine := ctx.String("engine")
		processingEngine := models.ContentProcessingEngine(engine)
		update.DefaultContentProcessingEngineDoc = &processingEngine
	}

	if ctx.IsSet("url-engine") {
		urlEngine := ctx.String("url-engine")
		processingEngineURL := models.ContentProcessingEngineURL(urlEngine)
		update.DefaultContentProcessingEngineURL = &processingEngineURL
	}

	if ctx.IsSet("embed") {
		embed := ctx.String("embed")
		embeddingOption := models.EmbeddingOption(embed)
		update.DefaultEmbeddingOption = &embeddingOption
	}

	if ctx.IsSet("auto-delete") {
		autoDelete := ctx.String("auto-delete")
		yesNoDecision, err := parseYesNoValue(autoDelete)
		if err != nil {
			return err
		}
		update.AutoDeleteFiles = &yesNoDecision
	}

	if ctx.IsSet("youtube-languages") {
		languages := ctx.StringSlice("youtube-languages")
		update.YoutubePreferredLanguages = languages
	}

	// Update settings
	updatedSettings, err := services.SettingsService.Update(ctx.Context, update)
	if err != nil {
		return errors.APIError("Failed to update settings",
			"Check input parameters and API permissions")
	}

	fmt.Println("✅ Settings updated successfully!")

	// Show updated settings
	fmt.Println("\nUpdated Settings:")
	if update.DefaultContentProcessingEngineDoc != nil {
		fmt.Printf("  Document Engine: %s\n", string(updatedSettings.DefaultContentProcessingEngineDoc))
	}
	if update.DefaultContentProcessingEngineURL != nil {
		fmt.Printf("  URL Engine: %s\n", string(updatedSettings.DefaultContentProcessingEngineURL))
	}
	if update.DefaultEmbeddingOption != nil {
		fmt.Printf("  Embedding: %s\n", string(updatedSettings.DefaultEmbeddingOption))
	}
	if update.AutoDeleteFiles != nil {
		fmt.Printf("  Auto Delete Files: %s\n", string(updatedSettings.AutoDeleteFiles))
	}
	if update.YoutubePreferredLanguages != nil {
		fmt.Printf("  YouTube Languages: %v\n", updatedSettings.YoutubePreferredLanguages)
	}

	return nil
}
