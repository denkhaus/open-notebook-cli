package commands

import (
	"fmt"
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

// TransformationsServices holds all the services needed for transformation commands
type TransformationsServices struct {
	TransformationService shared.TransformationRepository
	Config                config.Service
	Logger                shared.Logger
}

// getTransformationsServices retrieves all required services via dependency injection
func getTransformationsServices(ctx *cli.Context) (*TransformationsServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &TransformationsServices{
		TransformationService: do.MustInvoke[shared.TransformationRepository](injector),
		Config:                do.MustInvoke[config.Service](injector),
		Logger:                do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validateTransformationArgs validates common argument patterns for transformation commands
func validateTransformationArgs(ctx *cli.Context, requireTransformationID bool) (string, error) {
	if requireTransformationID {
		if ctx.NArg() < 1 {
			return "", fmt.Errorf("❌ Error: Missing transformation ID")
		}
		if ctx.NArg() > 1 {
			return "", fmt.Errorf("❌ Error: Too many arguments. Expected only transformation ID")
		}
	}

	transformationID := ctx.Args().First()
	if requireTransformationID && transformationID == "" {
		return "", errors.UsageError("Transformation ID is required",
			"Usage: open-notebook transformations <command> <transformation-id>")
	}

	return transformationID, nil
}

// printTransformationSuccess prints standardized success messages for transformation operations
func printTransformationSuccess(operation string, transformation *models.Transformation) {
	fmt.Printf("✅ Transformation %s successfully!\n", operation)
	fmt.Printf("  ID:          %s\n", transformation.ID)
	fmt.Printf("  Name:        %s\n", transformation.Name)
	fmt.Printf("  Title:       %s\n", transformation.Title)
	if transformation.Description != "" {
		fmt.Printf("  Description: %s\n", transformation.Description)
	}
}

// handleTransformationsList handles the transformations list command
func handleTransformationsList(ctx *cli.Context) error {
	services, err := getTransformationsServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Listing transformations...")

	transformationList, err := services.TransformationService.List(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to list transformations",
			"Check API connection and permissions")
	}

	if len(transformationList) == 0 {
		fmt.Println("No transformations found.")
		return nil
	}

	// Apply pagination
	limit := 20
	offset := 0
	if ctx.IsSet("limit") {
		limit = ctx.Int("limit")
	}
	if ctx.IsSet("offset") {
		offset = ctx.Int("offset")
	}

	// Apply pagination
	start := offset
	if start > len(transformationList) {
		start = len(transformationList)
	}
	end := start + limit
	if end > len(transformationList) {
		end = len(transformationList)
	}

	displayTransformations := transformationList[start:end]

	// Display transformations in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTITLE\tDEFAULT\tCREATED")

	for _, transformation := range displayTransformations {
		defaultFlag := "No"
		if transformation.ApplyDefault {
			defaultFlag = "Yes"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			transformation.ID,
			transformation.Name,
			utils.TruncateString(transformation.Title, 25),
			defaultFlag,
			utils.FormatTimestamp(transformation.Created))
	}

	w.Flush()

	fmt.Printf("\nShowing %d transformations (use --limit and --offset for pagination)\n", len(displayTransformations))
	return nil
}

// handleTransformationsShow handles transformation details display
func handleTransformationsShow(ctx *cli.Context) error {
	transformationID, err := validateTransformationArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getTransformationsServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Showing transformation details", "transformation_id", transformationID)

	// Get all transformations and find the one with matching ID
	transformationList, err := services.TransformationService.List(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to get transformation details",
			"Check API connection and permissions")
	}

	var targetTransformation *models.Transformation
	for _, transformation := range transformationList {
		if transformation.ID == transformationID {
			targetTransformation = transformation
			break
		}
	}

	if targetTransformation == nil {
		return errors.ValidationError("Transformation not found",
			fmt.Sprintf("Transformation with ID '%s' does not exist", transformationID))
	}

	// Display transformation details
	fmt.Printf("Transformation Details:\n")
	fmt.Printf("  ID:           %s\n", targetTransformation.ID)
	fmt.Printf("  Name:         %s\n", targetTransformation.Name)
	fmt.Printf("  Title:        %s\n", targetTransformation.Title)
	fmt.Printf("  Description:  %s\n", targetTransformation.Description)
	fmt.Printf("  Created:      %s\n", utils.FormatTimestamp(targetTransformation.Created))
	fmt.Printf("  Updated:      %s\n", utils.FormatTimestamp(targetTransformation.Updated))
	fmt.Printf("  Default:      %t\n", targetTransformation.ApplyDefault)

	// Display prompt
	fmt.Printf("  Prompt:\n")
	fmt.Printf("    %s\n", targetTransformation.Prompt)

	return nil
}

// handleTransformationsCreate handles transformation creation
func handleTransformationsCreate(ctx *cli.Context) error {
	services, err := getTransformationsServices(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")
	title := ctx.String("title")
	description := ctx.String("description")
	prompt := ctx.String("prompt")
	applyDefault := ctx.Bool("apply-default")

	services.Logger.Info("Creating transformation", "name", name, "title", title)

	transformation := &models.TransformationCreate{
		Name:         name,
		Title:        title,
		Description:  description,
		Prompt:       prompt,
		ApplyDefault: applyDefault,
	}

	createdTransformation, err := services.TransformationService.Create(ctx.Context, transformation)
	if err != nil {
		return errors.APIError("Failed to create transformation",
			"Check input parameters and API permissions")
	}

	printTransformationSuccess("created", createdTransformation)
	return nil
}

// handleTransformationsUpdate handles transformation updates
func handleTransformationsUpdate(ctx *cli.Context) error {
	transformationID, err := validateTransformationArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getTransformationsServices(ctx)
	if err != nil {
		return err
	}

	// Check if any update fields are provided
	if !ctx.IsSet("name") && !ctx.IsSet("title") && !ctx.IsSet("description") &&
		!ctx.IsSet("prompt") && !ctx.IsSet("apply-default") {
		return errors.UsageError("No update fields specified",
			"Use at least one of: --name, --title, --description, --prompt, --apply-default")
	}

	services.Logger.Info("Updating transformation", "transformation_id", transformationID)

	transformation := &models.TransformationUpdate{}

	if ctx.IsSet("name") {
		name := ctx.String("name")
		transformation.Name = &name
	}
	if ctx.IsSet("title") {
		title := ctx.String("title")
		transformation.Title = &title
	}
	if ctx.IsSet("description") {
		description := ctx.String("description")
		transformation.Description = &description
	}
	if ctx.IsSet("prompt") {
		prompt := ctx.String("prompt")
		transformation.Prompt = &prompt
	}
	if ctx.IsSet("apply-default") {
		applyDefault := ctx.Bool("apply-default")
		transformation.ApplyDefault = &applyDefault
	}

	updatedTransformation, err := services.TransformationService.Update(ctx.Context, transformationID, transformation)
	if err != nil {
		return errors.APIError("Failed to update transformation",
			"Check transformation ID and permissions")
	}

	printTransformationSuccess("updated", updatedTransformation)
	return nil
}

// handleTransformationsDelete handles transformation deletion
func handleTransformationsDelete(ctx *cli.Context) error {
	transformationID, err := validateTransformationArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getTransformationsServices(ctx)
	if err != nil {
		return err
	}

	// Confirm deletion unless force flag is used
	if !ctx.Bool("force") {
		fmt.Printf("⚠️  Are you sure you want to delete transformation '%s'? [y/N]: ", transformationID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Deletion cancelled")
			return nil
		}
	}

	services.Logger.Info("Deleting transformation", "transformation_id", transformationID)

	err = services.TransformationService.Delete(ctx.Context, transformationID)
	if err != nil {
		return errors.APIError("Failed to delete transformation",
			"Check transformation ID and permissions")
	}

	fmt.Printf("✅ Transformation '%s' deleted successfully!\n", transformationID)
	return nil
}

// handleTransformationsExecute handles transformation execution
func handleTransformationsExecute(ctx *cli.Context) error {
	transformationID, err := validateTransformationArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getTransformationsServices(ctx)
	if err != nil {
		return err
	}

	inputText := ctx.String("text")
	modelID := ctx.String("model")
	stream := ctx.Bool("stream")

	if inputText == "" {
		return errors.UsageError("Input text is required",
			"Use --text flag to specify the text to transform")
	}

	services.Logger.Info("Executing transformation", "transformation_id", transformationID)

	// Build execution request
	request := &models.TransformationExecuteRequest{
		TransformationID: transformationID,
		InputText:        inputText,
	}

	if modelID != "" {
		request.ModelID = modelID
	}

	fmt.Printf("🔄 Executing transformation: %s\n", transformationID)
	if modelID != "" {
		fmt.Printf("  Using model: %s\n", modelID)
	}
	fmt.Printf("  Input text: %s\n", utils.TruncateString(inputText, 50))

	response, err := services.TransformationService.Execute(ctx.Context, request)
	if err != nil {
		return errors.APIError("Failed to execute transformation",
			"Check transformation ID, model ID and permissions")
	}

	if stream {
		fmt.Println("  🔄 Streaming output:")
		// Simulate streaming by printing chunks
		output := response.Output
		chunkSize := 20
		for i := 0; i < len(output); i += chunkSize {
			end := i + chunkSize
			if end > len(output) {
				end = len(output)
			}
			fmt.Printf("  %s\n", output[i:end])
			time.Sleep(100 * time.Millisecond) // Simulate streaming delay
		}
		return nil
	}

	fmt.Printf("✅ Transformation executed successfully!\n")
	fmt.Printf("  Output: %s\n", response.Output)
	if response.TransformationID != "" {
		fmt.Printf("  Used transformation: %s\n", response.TransformationID)
	}
	if response.ModelID != "" {
		fmt.Printf("  Used model: %s\n", response.ModelID)
	}

	return nil
}
