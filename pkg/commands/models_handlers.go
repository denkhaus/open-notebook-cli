package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// ModelsServices holds all the services needed for model commands
type ModelsServices struct {
	ModelService shared.ModelService
	Config       config.Service
	Logger       shared.Logger
}

// getModelsServices retrieves all required services via dependency injection
func getModelsServices(ctx *cli.Context) (*ModelsServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &ModelsServices{
		ModelService: do.MustInvoke[shared.ModelService](injector),
		Config:       do.MustInvoke[config.Service](injector),
		Logger:       do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validateModelArgs validates common argument patterns for model commands
func validateModelArgs(ctx *cli.Context, requireModelID bool) (string, error) {
	if requireModelID {
		if ctx.NArg() < 1 {
			return "", fmt.Errorf("❌ Error: Missing model ID")
		}
		if ctx.NArg() > 1 {
			return "", fmt.Errorf("❌ Error: Too many arguments. Expected only model ID")
		}
	}

	modelID := ctx.Args().First()
	if requireModelID && modelID == "" {
		return "", errors.UsageError("Model ID is required",
			"Usage: open-notebook models <command> <model-id>")
	}

	return modelID, nil
}

// validateModelType validates model type parameter
func validateModelType(modelType string) error {
	validTypes := map[string]bool{
		"language":       true,
		"embedding":      true,
		"text_to_speech": true,
		"speech_to_text": true,
	}

	if !validTypes[modelType] {
		return errors.UsageError("Invalid model type",
			fmt.Sprintf("Valid types: language, embedding, text_to_speech, speech_to_text (got: %s)", modelType))
	}

	return nil
}

// printModelSuccess prints standardized success messages for model operations
func printModelSuccess(operation string, model *models.Model) {
	fmt.Printf("✅ Model %s successfully!\n", operation)
	fmt.Printf("  ID:       %s\n", model.ID)
	fmt.Printf("  Name:     %s\n", model.Name)
	fmt.Printf("  Provider: %s\n", model.Provider)
	fmt.Printf("  Type:     %s\n", string(model.Type))
}

// displayModelCapabilities displays model capabilities in a readable format
func displayModelCapabilities(capabilities map[string]any) {
	if capabilities == nil || len(capabilities) == 0 {
		return
	}

	fmt.Println("  Capabilities:")
	for key, value := range capabilities {
		fmt.Printf("    %s: %v\n", key, value)
	}
}

// handleModelsList handles the models list command
func handleModelsList(ctx *cli.Context) error {
	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Listing models...")

	modelList, err := services.ModelService.List(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to list models",
			"Check API connection and permissions")
	}

	if len(modelList) == 0 {
		fmt.Println("No models found.")
		return nil
	}

	// Filter models based on flags
	filteredModels := []*models.Model{}
	for _, model := range modelList {
		// Filter by type
		if ctx.IsSet("type") && string(model.Type) != ctx.String("type") {
			continue
		}
		// Filter by provider
		if ctx.IsSet("provider") && model.Provider != ctx.String("provider") {
			continue
		}
		filteredModels = append(filteredModels, model)
	}

	// Apply limit and offset
	limit := 50
	offset := 0
	if ctx.IsSet("limit") {
		limit = ctx.Int("limit")
	}
	if ctx.IsSet("offset") {
		offset = ctx.Int("offset")
	}

	// Apply pagination
	start := offset
	if start > len(filteredModels) {
		start = len(filteredModels)
	}
	end := start + limit
	if end > len(filteredModels) {
		end = len(filteredModels)
	}

	displayModels := filteredModels[start:end]

	// Display models in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPROVIDER\tTYPE")

	for _, model := range displayModels {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			model.ID,
			utils.TruncateString(model.Name, 25),
			model.Provider,
			string(model.Type))
	}

	w.Flush()

	fmt.Printf("\nShowing %d models (use --limit and --offset for pagination)\n", len(displayModels))
	return nil
}

// handleModelsShow handles model details display
func handleModelsShow(ctx *cli.Context) error {
	modelID, err := validateModelArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Showing model details", "model_id", modelID)

	// First get all models and find the one with matching ID
	modelList, err := services.ModelService.List(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to get model details",
			"Check API connection and permissions")
	}

	var targetModel *models.Model
	for _, model := range modelList {
		if model.ID == modelID {
			targetModel = model
			break
		}
	}

	if targetModel == nil {
		return errors.ValidationError("Model not found",
			fmt.Sprintf("Model with ID '%s' does not exist", modelID))
	}

	// Display model details
	fmt.Printf("Model Details:\n")
	fmt.Printf("  ID:           %s\n", targetModel.ID)
	fmt.Printf("  Name:         %s\n", targetModel.Name)
	fmt.Printf("  Provider:     %s\n", targetModel.Provider)
	fmt.Printf("  Type:         %s\n", string(targetModel.Type))
	fmt.Printf("  Created:      %s\n", utils.FormatTimestamp(targetModel.Created))
	fmt.Printf("  Updated:      %s\n", utils.FormatTimestamp(targetModel.Updated))

	return nil
}

// handleModelsAdd handles model creation
func handleModelsAdd(ctx *cli.Context) error {
	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")
	provider := ctx.String("provider")
	modelType := ctx.String("type")

	// Validate model type
	if err := validateModelType(modelType); err != nil {
		return err
	}

	services.Logger.Info("Creating model", "name", name, "provider", provider, "type", modelType)

	model := &models.ModelCreate{
		Name:     name,
		Provider: provider,
		Type:     models.ModelType(modelType),
	}

	createdModel, err := services.ModelService.Create(ctx.Context, model)
	if err != nil {
		return errors.APIError("Failed to create model",
			"Check input parameters and API permissions")
	}

	printModelSuccess("created", createdModel)
	return nil
}

// handleModelsDelete handles model deletion
func handleModelsDelete(ctx *cli.Context) error {
	modelID, err := validateModelArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	// Confirm deletion unless force flag is used
	if !ctx.Bool("force") {
		fmt.Printf("⚠️  Are you sure you want to delete model '%s'? [y/N]: ", modelID)
		var response string
		fmt.Scanln(&response)
		response = fmt.Sprintf("%s", response) // Format the response
		if response != "y" && response != "yes" {
			fmt.Println("❌ Deletion cancelled")
			return nil
		}
	}

	services.Logger.Info("Deleting model", "model_id", modelID)

	err = services.ModelService.Delete(ctx.Context, modelID)
	if err != nil {
		return errors.APIError("Failed to delete model",
			"Check model ID and permissions")
	}

	fmt.Printf("✅ Model '%s' deleted successfully!\n", modelID)
	return nil
}

// handleModelsDefaultsShow handles showing current default models
func handleModelsDefaultsShow(ctx *cli.Context) error {
	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Getting default models")

	defaults, err := services.ModelService.GetDefaults(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to get default models",
			"Check API connection and permissions")
	}

	fmt.Println("🎯 Current default models:")

	defaultTypes := []struct {
		field *string
		label string
	}{
		{defaults.DefaultChatModel, "Chat"},
		{defaults.DefaultTransformationModel, "Transformation"},
		{defaults.LargeContextModel, "Large Context"},
		{defaults.DefaultTextToSpeechModel, "Text-to-Speech"},
		{defaults.DefaultSpeechToTextModel, "Speech-to-Text"},
		{defaults.DefaultEmbeddingModel, "Embedding"},
		{defaults.DefaultToolsModel, "Tools"},
	}

	for _, dt := range defaultTypes {
		if dt.field != nil && *dt.field != "" {
			fmt.Printf("  %s: %s\n", dt.label, *dt.field)
		} else {
			fmt.Printf("  %s: (not set)\n", dt.label)
		}
	}

	return nil
}

// handleModelsDefaultsSet handles setting default models
func handleModelsDefaultsSet(ctx *cli.Context) error {
	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	// Collect all default assignments
	assignments := map[string]string{}

	if ctx.IsSet("chat") {
		assignments["chat"] = ctx.String("chat")
	}
	if ctx.IsSet("embedding") {
		assignments["embedding"] = ctx.String("embedding")
	}
	if ctx.IsSet("transformation") {
		assignments["transformation"] = ctx.String("transformation")
	}
	if ctx.IsSet("large-context") {
		assignments["large_context"] = ctx.String("large-context")
	}
	if ctx.IsSet("tts") {
		assignments["text_to_speech"] = ctx.String("tts")
	}
	if ctx.IsSet("stt") {
		assignments["speech_to_text"] = ctx.String("stt")
	}
	if ctx.IsSet("tools") {
		assignments["tools"] = ctx.String("tools")
	}

	if len(assignments) == 0 {
		return errors.UsageError("No defaults specified",
			"Use at least one of the --chat, --embedding, --transformation, --large-context, --tts, --stt, --tools flags")
	}

	services.Logger.Info("Setting default models", "assignments", assignments)

	// Convert assignments to DefaultModelsResponse
	defaults := &models.DefaultModelsResponse{}
	if ctx.IsSet("chat") {
		chat := ctx.String("chat")
		defaults.DefaultChatModel = &chat
	}
	if ctx.IsSet("embedding") {
		embedding := ctx.String("embedding")
		defaults.DefaultEmbeddingModel = &embedding
	}
	if ctx.IsSet("transformation") {
		transformation := ctx.String("transformation")
		defaults.DefaultTransformationModel = &transformation
	}
	if ctx.IsSet("large-context") {
		largeContext := ctx.String("large-context")
		defaults.LargeContextModel = &largeContext
	}
	if ctx.IsSet("tts") {
		tts := ctx.String("tts")
		defaults.DefaultTextToSpeechModel = &tts
	}
	if ctx.IsSet("stt") {
		stt := ctx.String("stt")
		defaults.DefaultSpeechToTextModel = &stt
	}
	if ctx.IsSet("tools") {
		tools := ctx.String("tools")
		defaults.DefaultToolsModel = &tools
	}

	err = services.ModelService.SetDefaults(ctx.Context, defaults)
	if err != nil {
		return errors.APIError("Failed to set default models",
			"Check model IDs and API permissions")
	}

	fmt.Println("✅ Default models updated successfully!")
	for key, value := range assignments {
		fmt.Printf("  %s: %s\n", key, value)
	}

	return nil
}

// handleModelsProviders handles checking provider availability
func handleModelsProviders(ctx *cli.Context) error {
	services, err := getModelsServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Checking provider availability")

	providers, err := services.ModelService.GetProviders(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to get provider status",
			"Check API connection and permissions")
	}

	fmt.Println("🔌 Provider availability:")

	// Display available providers
	if len(providers.Available) > 0 {
		fmt.Println("  Available:")
		for _, provider := range providers.Available {
			fmt.Printf("    ✅ %s\n", provider)
		}
	}

	// Display unavailable providers
	if len(providers.Unavailable) > 0 {
		fmt.Println("  Unavailable:")
		for _, provider := range providers.Unavailable {
			fmt.Printf("    ❌ %s\n", provider)
		}
	}

	// Display supported types
	if len(providers.SupportedTypes) > 0 {
		fmt.Println("  Supported types by provider:")
		for provider, types := range providers.SupportedTypes {
			fmt.Printf("    %s: %v\n", provider, types)
		}
	}

	return nil
}
