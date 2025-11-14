package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// NotesServices holds all the services needed for note commands
type NotesServices struct {
	NoteService services.NoteRepository
	Config      config.Service
	Logger      services.Logger
}

// getNotesServices retrieves all required services via dependency injection
func getNotesServices(ctx *cli.Context) (*NotesServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &NotesServices{
		NoteService: do.MustInvoke[services.NoteRepository](injector),
		Config:      do.MustInvoke[config.Service](injector),
		Logger:      do.MustInvoke[services.Logger](injector),
	}, nil
}

// validateNoteArgs validates common argument patterns for note commands
func validateNoteArgs(ctx *cli.Context, requireNoteID bool) (string, error) {
	if requireNoteID {
		if ctx.NArg() < 1 {
			return "", fmt.Errorf("❌ Error: Missing note ID")
		}
		if ctx.NArg() > 1 {
			return "", fmt.Errorf("❌ Error: Too many arguments. Expected only note ID")
		}
	}

	noteID := ctx.Args().First()
	if requireNoteID && noteID == "" {
		return "", errors.UsageError("Note ID is required",
			"Usage: open-notebook notes <command> <note-id>")
	}

	return noteID, nil
}

// printNoteSuccess prints standardized success messages for note operations
func printNoteSuccess(operation string, note *models.Note) {
	fmt.Printf("✅ Note %s successfully!\n", operation)
	fmt.Printf("  ID:     %s\n", utils.SafeDereferenceString(note.ID))
	fmt.Printf("  Title:  %s\n", utils.SafeDereferenceString(note.Title))
	if note.NoteType != nil {
		fmt.Printf("  Type:   %s\n", string(*note.NoteType))
	}
}

// Handler functions with proper separation of concerns

// handleNotesList handles the notes list command
func handleNotesList(ctx *cli.Context) error {
	services, err := getNotesServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Listing notes...")

	notebookID := ctx.String("notebook")
	limit := ctx.Int("limit")
	offset := ctx.Int("offset")

	notes, err := services.NoteService.List(ctx.Context, notebookID, limit, offset)
	if err != nil {
		return errors.APIError("Failed to list notes",
			"Check API connection and permissions")
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	// Display notes in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tTYPE\tCREATED")

	for _, note := range notes {
		title := "Untitled"
		if note.Title != nil {
			title = *note.Title
		}
		noteType := "text"
		if note.NoteType != nil {
			noteType = string(*note.NoteType)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			utils.SafeDereferenceString(note.ID),
			utils.TruncateString(title, 30),
			noteType,
			utils.FormatTimestamp(note.Created))
	}

	w.Flush()

	fmt.Printf("\nShowing %d notes (use --limit and --offset for pagination)\n", len(notes))
	return nil
}

// handleNotesAdd handles the notes add command
func handleNotesAdd(ctx *cli.Context) error {
	services, err := getNotesServices(ctx)
	if err != nil {
		return err
	}

	content := ctx.String("content")
	notebookID := ctx.String("notebook")
	title := ctx.String("title")
	noteType := ctx.String("type")

	// Validate required fields
	if content == "" {
		return errors.UsageError("Content is required",
			"Use --content flag to specify the note content")
	}

	if notebookID == "" {
		return errors.UsageError("Notebook ID is required",
			"Use --notebook flag to specify the notebook")
	}

	// Validate note type
	validTypes := map[string]bool{
		"text":     true,
		"markdown": true,
		"code":     true,
	}
	if !validTypes[noteType] {
		return errors.UsageError("Invalid note type",
			"Supported types are: text, markdown, code")
	}

	services.Logger.Info("Creating note", "notebook", notebookID, "title", title)

	noteTypeTyped := models.NoteType(noteType)
	noteCreate := &models.NoteCreate{
		Title:      &title,
		Content:    content,
		NotebookID: &notebookID,
		NoteType:   &noteTypeTyped,
	}

	note, err := services.NoteService.Create(ctx.Context, noteCreate)
	if err != nil {
		return errors.APIError("Failed to create note",
			"Check input parameters and API permissions")
	}

	printNoteSuccess("created", note)
	return nil
}

// handleNotesShow handles the notes show command
func handleNotesShow(ctx *cli.Context) error {
	noteID, err := validateNoteArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getNotesServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Showing note details", "note_id", noteID)

	note, err := services.NoteService.Get(ctx.Context, noteID)
	if err != nil {
		return errors.APIError("Failed to get note details",
			"Check note ID and permissions")
	}

	// Display note details
	fmt.Printf("Note Details:\n")
	fmt.Printf("  ID:           %s\n", utils.SafeDereferenceString(note.ID))
	fmt.Printf("  Title:        %s\n", utils.SafeDereferenceString(note.Title))
	fmt.Printf("  Created:      %s\n", utils.FormatTimestamp(note.Created))
	fmt.Printf("  Updated:      %s\n", utils.FormatTimestamp(note.Updated))

	if note.NoteType != nil {
		fmt.Printf("  Type:         %s\n", string(*note.NoteType))
	}

	if note.Content != nil {
		fmt.Printf("  Content:\n")
		fmt.Printf("  %s\n", *note.Content)
	}

	return nil
}

// handleNotesUpdate handles the notes update command
func handleNotesUpdate(ctx *cli.Context) error {
	noteID, err := validateNoteArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getNotesServices(ctx)
	if err != nil {
		return err
	}

	title := ctx.String("title")
	content := ctx.String("content")
	noteType := ctx.String("type")

	if title == "" && content == "" && noteType == "" {
		return errors.UsageError("At least one update field is required",
			"Use --title, --content, or --type flags")
	}

	// Build update request with only provided fields
	update := &models.NoteUpdate{}
	if ctx.IsSet("title") {
		update.Title = &title
	}
	if ctx.IsSet("content") {
		update.Content = &content
	}
	if ctx.IsSet("type") {
		// Validate note type
		validTypes := map[string]bool{
			"text":     true,
			"markdown": true,
			"code":     true,
		}
		if !validTypes[noteType] {
			return errors.UsageError("Invalid note type",
				"Supported types are: text, markdown, code")
		}
		typedNoteType := models.NoteType(noteType)
		update.NoteType = &typedNoteType
	}

	services.Logger.Info("Updating note", "note_id", noteID)

	updatedNote, err := services.NoteService.Update(ctx.Context, noteID, update)
	if err != nil {
		return errors.APIError("Failed to update note",
			"Check note ID and permissions")
	}

	printNoteSuccess("updated", updatedNote)
	return nil
}

// handleNotesDelete handles the notes delete command
func handleNotesDelete(ctx *cli.Context) error {
	noteID, err := validateNoteArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getNotesServices(ctx)
	if err != nil {
		return err
	}

	force := ctx.Bool("force")

	// Confirm deletion unless force flag is used
	if !force {
		fmt.Printf("Are you sure you want to delete note '%s'? (y/N): ", noteID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	services.Logger.Info("Deleting note", "note_id", noteID)

	err = services.NoteService.Delete(ctx.Context, noteID)
	if err != nil {
		return errors.APIError("Failed to delete note",
			"Check note ID and permissions")
	}

	fmt.Printf("✅ Note '%s' deleted successfully!\n", noteID)
	return nil
}

// handleNotesSearch handles the notes search command
func handleNotesSearch(ctx *cli.Context) error {
	services, err := getNotesServices(ctx)
	if err != nil {
		return err
	}

	notebookID := ctx.String("notebook")
	query := ctx.String("query")

	if query == "" {
		return errors.UsageError("Query is required",
			"Use --query flag to specify the search query")
	}

	services.Logger.Info("Searching notes", "query", query, "notebook", notebookID)

	notes, err := services.NoteService.Search(ctx.Context, notebookID, query)
	if err != nil {
		return errors.APIError("Failed to search notes",
			"Check API connection and permissions")
	}

	if len(notes) == 0 {
		fmt.Println("No notes found matching your search.")
		return nil
	}

	// Display search results in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tTYPE\tCREATED")

	for _, note := range notes {
		title := "Untitled"
		if note.Title != nil {
			title = *note.Title
		}
		noteType := "text"
		if note.NoteType != nil {
			noteType = string(*note.NoteType)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			utils.SafeDereferenceString(note.ID),
			utils.TruncateString(title, 30),
			noteType,
			utils.FormatTimestamp(note.Created))
	}

	w.Flush()

	fmt.Printf("\nFound %d notes matching '%s'\n", len(notes), query)
	return nil
}
