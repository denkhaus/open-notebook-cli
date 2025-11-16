package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// getNotebookServices retrieves all required services via dependency injection
func getNotebookServices(ctx *cli.Context) (*NotebookServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &NotebookServices{
		NotebookService: do.MustInvoke[services.NotebookService](injector),
		Config:          do.MustInvoke[config.Service](injector),
		Logger:          do.MustInvoke[services.Logger](injector),
	}, nil
}

// handleNotebooksList handles the notebooks list command
func handleNotebooksList(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Listing notebooks")

	notebooks, err := services.NotebookService.ListNotebooks(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to list notebooks",
			"Check API connection and permissions")
	}

	if len(notebooks) == 0 {
		fmt.Println("No notebooks found")
		services.Logger.Info("No notebooks found")
		return nil
	}

	// Display in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tSOURCES\tNOTES\tARCHIVED")
	fmt.Fprintln(w, "--\t----\t-----------\t-------\t-----\t--------")

	for _, nb := range notebooks {
		archived := "No"
		if nb.Archived {
			archived = "Yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\n",
			nb.ID, nb.Name, nb.Description, nb.SourceCount, nb.NoteCount, archived)
	}
	w.Flush()

	services.Logger.Info("Listed notebooks successfully", "count", len(notebooks))
	return nil
}

// handleNotebooksCreate handles the notebooks create command
func handleNotebooksCreate(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	name := ctx.String("name")
	description := ctx.String("description")

	services.Logger.Info("Creating notebook", "name", name)

	notebook, err := services.NotebookService.CreateNotebook(ctx.Context, name, description)
	if err != nil {
		return errors.APIError("Failed to create notebook",
			"Check name length and API connection")
	}

	fmt.Printf("✅ Created notebook: %s (ID: %s)\n", notebook.Name, notebook.ID)
	if notebook.Description != "" {
		fmt.Printf("📝 Description: %s\n", notebook.Description)
	}

	services.Logger.Info("Notebook created successfully", "id", notebook.ID, "name", name)
	return nil
}

// handleNotebooksShow handles the notebooks show command
func handleNotebooksShow(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	id := ctx.String("id")
	services.Logger.Info("Getting notebook details", "id", id)

	notebook, err := services.NotebookService.GetNotebook(ctx.Context, id)
	if err != nil {
		return errors.NotFoundError("Notebook not found",
			fmt.Sprintf("Notebook with ID %s does not exist", id))
	}

	fmt.Printf("📓 Notebook Details:\n\n")
	fmt.Printf("ID:          %s\n", notebook.ID)
	fmt.Printf("Name:        %s\n", notebook.Name)
	fmt.Printf("Description: %s\n", notebook.Description)
	fmt.Printf("Created:     %s\n", notebook.Created)
	fmt.Printf("Updated:     %s\n", notebook.Updated)
	fmt.Printf("Sources:      %d\n", notebook.SourceCount)
	fmt.Printf("Notes:        %d\n", notebook.NoteCount)
	fmt.Printf("Archived:    %t\n", notebook.Archived)

	services.Logger.Info("Notebook details displayed", "id", id)
	return nil
}

// handleNotebooksUpdate handles the notebooks update command
func handleNotebooksUpdate(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	id := ctx.String("id")
	name := ctx.String("name")
	description := ctx.String("description")
	archived := ctx.Bool("archived")

	services.Logger.Info("Updating notebook", "id", id)

	var namePtr, descPtr *string
	var archivedPtr *bool

	if ctx.IsSet("name") {
		namePtr = &name
	}
	if ctx.IsSet("description") {
		descPtr = &description
	}
	if ctx.IsSet("archived") {
		archivedPtr = &archived
	}

	notebook, err := services.NotebookService.UpdateNotebook(ctx.Context, id, namePtr, descPtr, archivedPtr)
	if err != nil {
		return errors.APIError("Failed to update notebook",
			"Check that the notebook exists and field values are valid")
	}

	fmt.Printf("✅ Updated notebook: %s\n", notebook.Name)
	services.Logger.Info("Notebook updated successfully", "id", id)
	return nil
}

// handleNotebooksDelete handles the notebooks delete command
func handleNotebooksDelete(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	id := ctx.String("id")
	confirm := ctx.Bool("confirm")

	services.Logger.Info("Deleting notebook", "id", id)

	// Get notebook details for confirmation
	notebook, err := services.NotebookService.GetNotebook(ctx.Context, id)
	if err != nil {
		return errors.NotFoundError("Notebook not found",
			fmt.Sprintf("Notebook with ID %s does not exist", id))
	}

	// Confirmation prompt
	if !confirm {
		fmt.Printf("⚠️  Are you sure you want to delete notebook '%s'? (ID: %s)\n", notebook.Name, notebook.ID)
		fmt.Printf("This will also delete %d sources and %d notes.\n", notebook.SourceCount, notebook.NoteCount)
		fmt.Print("Type 'yes' to confirm: ")

		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			fmt.Println("❌ Delete cancelled")
			services.Logger.Info("Notebook deletion cancelled by user", "id", id)
			return nil
		}
	}

	if err := services.NotebookService.DeleteNotebook(ctx.Context, id); err != nil {
		return errors.APIError("Failed to delete notebook",
			"Check permissions and that the notebook is not in use")
	}

	fmt.Printf("✅ Deleted notebook: %s\n", notebook.Name)
	services.Logger.Info("Notebook deleted successfully", "id", id)
	return nil
}

// handleNotebooksAddSource handles the notebooks add-source command
func handleNotebooksAddSource(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	notebookID := ctx.String("notebook")
	sourceID := ctx.String("source")

	services.Logger.Info("Adding source to notebook", "notebook_id", notebookID, "source_id", sourceID)

	if err := services.NotebookService.AddSourceToNotebook(ctx.Context, notebookID, sourceID); err != nil {
		return errors.APIError("Failed to add source to notebook",
			"Check that both notebook and source exist")
	}

	fmt.Printf("✅ Added source %s to notebook %s\n", sourceID, notebookID)
	services.Logger.Info("Source added to notebook successfully")
	return nil
}

// handleNotebooksRemoveSource handles the notebooks remove-source command
func handleNotebooksRemoveSource(ctx *cli.Context) error {
	services, err := getNotebookServices(ctx)
	if err != nil {
		return err
	}

	notebookID := ctx.String("notebook")
	sourceID := ctx.String("source")

	services.Logger.Info("Removing source from notebook", "notebook_id", notebookID, "source_id", sourceID)

	if err := services.NotebookService.RemoveSourceFromNotebook(ctx.Context, notebookID, sourceID); err != nil {
		return errors.APIError("Failed to remove source from notebook",
			"Check that both notebook and source exist")
	}

	fmt.Printf("✅ Removed source %s from notebook %s\n", sourceID, notebookID)
	services.Logger.Info("Source removed from notebook successfully")
	return nil
}
