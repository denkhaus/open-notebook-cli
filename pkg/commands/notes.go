package commands

import (
	"context"
	"fmt"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// NotesCommand returns the notes command and its subcommands
func NotesCommand() *cli.Command {
	return &cli.Command{
		Name:  "notes",
		Usage: "Note management commands",
		Subcommands: []*cli.Command{
			notesListCommand(),
			notesAddCommand(),
			notesShowCommand(),
			notesUpdateCommand(),
			notesDeleteCommand(),
			notesSearchCommand(),
		},
	}
}

// notesListCommand implements notes list functionality
func notesListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List notes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "notebook",
				Aliases: []string{"n"},
				Usage:   "Filter by notebook ID",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of notes to return",
				Value:   50,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of notes to skip",
				Value:   0,
			},
		},
		Action: handleNotesList,
	}
}

// notesAddCommand implements notes add functionality
func notesAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a new note",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "content",
				Aliases:  []string{"c"},
				Usage:    "Note content",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "notebook",
				Aliases:  []string{"n"},
				Usage:    "Notebook ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "Note title",
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"ty"},
				Usage:   "Note type (text, markdown, etc.)",
				Value:   "text",
			},
		},
		Action: handleNotesAdd,
	}
}

// notesShowCommand implements notes show functionality
func notesShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a specific note",
		Args:  true,
		Action: handleNotesShow,
	}
}

// notesUpdateCommand implements notes update functionality
func notesUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update an existing note",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "content",
				Aliases: []string{"c"},
				Usage:   "Updated note content",
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "Updated note title",
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"ty"},
				Usage:   "Updated note type",
			},
		},
		Action: handleNotesUpdate,
	}
}

// notesDeleteCommand implements notes delete functionality
func notesDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a note",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: handleNotesDelete,
	}
}

// notesSearchCommand implements notes search functionality
func notesSearchCommand() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Search notes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "notebook",
				Aliases:  []string{"n"},
				Usage:    "Search within specific notebook ID",
			},
			&cli.StringFlag{
				Name:     "query",
				Aliases:  []string{"q"},
				Usage:    "Search query",
				Required: true,
			},
		},
		Action: handleNotesSearch,
	}
}

// Handler functions with proper separation of concerns

// handleNotesList handles the notes list command
func handleNotesList(c *cli.Context) error {
	injector := getInjector(c)
	noteRepo := di.GetNoteRepository(injector)

	notebookID := c.String("notebook")
	limit := c.Int("limit")
	offset := c.Int("offset")

	ctx := context.Background()
	notes, err := noteRepo.List(ctx, notebookID, limit, offset)
	if err != nil {
		return fmt.Errorf("❌ Failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("📝 No notes found")
		return nil
	}

	fmt.Printf("📝 Found %d notes:\n\n", len(notes))
	for _, note := range notes {
		title := "Untitled"
		if note.Title != nil {
			title = *note.Title
		}
		noteType := "text"
		if note.NoteType != nil {
			noteType = string(*note.NoteType)
		}

		fmt.Printf("🔹 %s (%s)\n", safeDeref(note.ID), title)
		fmt.Printf("   Type: %s\n", noteType)
		fmt.Printf("   Created: %s\n", note.Created)
		fmt.Println()
	}

	return nil
}

// handleNotesAdd handles the notes add command
func handleNotesAdd(c *cli.Context) error {
	injector := getInjector(c)
	noteRepo := di.GetNoteRepository(injector)

	content := c.String("content")
	notebookID := c.String("notebook")
	title := c.String("title")
	noteType := c.String("type")

	fmt.Println("📝 Creating new note...")

	// Validate note type
	validTypes := map[string]bool{
		"text":     true,
		"markdown": true,
		"code":     true,
	}
	if !validTypes[noteType] {
		return fmt.Errorf("❌ Invalid note type '%s'. Valid types: text, markdown, code", noteType)
	}

	noteTypeTyped := models.NoteType(noteType)
	noteCreate := &models.NoteCreate{
		Title:      &title,
		Content:    content,
		NotebookID: &notebookID,
		NoteType:   &noteTypeTyped,
	}

	ctx := context.Background()
	note, err := noteRepo.Create(ctx, noteCreate)
	if err != nil {
		return fmt.Errorf("❌ Failed to create note: %w", err)
	}

	fmt.Printf("✅ Successfully created note: %s\n", safeDeref(note.ID))
	fmt.Printf("   Title: %s\n", safeDeref(note.Title))
	if note.NoteType != nil {
		fmt.Printf("   Type: %s\n", string(*note.NoteType))
	}

	return nil
}

// handleNotesShow handles the notes show command
func handleNotesShow(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("❌ Error: Missing note ID")
	}
	if c.NArg() > 1 {
		return fmt.Errorf("❌ Error: Too many arguments. Expected only note ID")
	}

	injector := getInjector(c)
	noteRepo := di.GetNoteRepository(injector)

	noteID := c.Args().First()
	ctx := context.Background()

	note, err := noteRepo.Get(ctx, noteID)
	if err != nil {
		return fmt.Errorf("❌ Failed to get note %s: %w", noteID, err)
	}

	fmt.Printf("📝 Note Details: %s\n", safeDeref(note.ID))
	fmt.Printf("   Title: %s\n", safeDeref(note.Title))
	if note.NoteType != nil {
		fmt.Printf("   Type: %s\n", string(*note.NoteType))
	}
	fmt.Printf("   Created: %s\n", note.Created)
	fmt.Printf("   Updated: %s\n", note.Updated)
	fmt.Println()
	fmt.Println("Content:")
	fmt.Println("--------")
	fmt.Println(note.Content)

	return nil
}

// handleNotesUpdate handles the notes update command
func handleNotesUpdate(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("❌ Error: Missing note ID")
	}
	if c.NArg() > 1 {
		return fmt.Errorf("❌ Error: Too many arguments. Expected only note ID")
	}

	injector := getInjector(c)
	noteRepo := di.GetNoteRepository(injector)

	noteID := c.Args().First()
	content := c.String("content")
	title := c.String("title")
	noteType := c.String("type")

	fmt.Printf("📝 Updating note: %s\n", noteID)

	// Build update request with only provided fields
	update := &models.NoteUpdate{}
	if c.IsSet("title") {
		update.Title = &title
	}
	if c.IsSet("content") {
		update.Content = &content
	}
	if c.IsSet("type") {
		// Validate note type
		validTypes := map[string]bool{
			"text":     true,
			"markdown": true,
			"code":     true,
		}
		if !validTypes[noteType] {
			return fmt.Errorf("❌ Invalid note type '%s'. Valid types: text, markdown, code", noteType)
		}
		typedNoteType := models.NoteType(noteType)
		update.NoteType = &typedNoteType
	}

	ctx := context.Background()
	note, err := noteRepo.Update(ctx, noteID, update)
	if err != nil {
		return fmt.Errorf("❌ Failed to update note %s: %w", noteID, err)
	}

	fmt.Printf("✅ Successfully updated note: %s\n", safeDeref(note.ID))
	fmt.Printf("   Title: %s\n", safeDeref(note.Title))
	if note.NoteType != nil {
		fmt.Printf("   Type: %s\n", string(*note.NoteType))
	}
	fmt.Printf("   Updated: %s\n", note.Updated)

	return nil
}

// handleNotesDelete handles the notes delete command
func handleNotesDelete(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("❌ Error: Missing note ID")
	}
	if c.NArg() > 1 {
		return fmt.Errorf("❌ Error: Too many arguments. Expected only note ID")
	}

	injector := getInjector(c)
	noteRepo := di.GetNoteRepository(injector)

	noteID := c.Args().First()
	force := c.Bool("force")

	if !force {
		fmt.Printf("⚠️  Are you sure you want to delete note '%s'? [y/N]: ", noteID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "yes" {
			fmt.Println("❌ Deletion cancelled")
			return nil
		}
	}

	fmt.Printf("🗑️  Deleting note: %s\n", noteID)

	ctx := context.Background()
	err := noteRepo.Delete(ctx, noteID)
	if err != nil {
		return fmt.Errorf("❌ Failed to delete note %s: %w", noteID, err)
	}

	fmt.Printf("✅ Successfully deleted note: %s\n", noteID)

	return nil
}

// handleNotesSearch handles the notes search command
func handleNotesSearch(c *cli.Context) error {
	injector := getInjector(c)
	noteRepo := di.GetNoteRepository(injector)

	notebookID := c.String("notebook")
	query := c.String("query")

	fmt.Printf("🔍 Searching notes with query: %s\n", query)
	if notebookID != "" {
		fmt.Printf("   Notebook: %s\n", notebookID)
	}

	ctx := context.Background()
	notes, err := noteRepo.Search(ctx, notebookID, query)
	if err != nil {
		return fmt.Errorf("❌ Failed to search notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("📝 No notes found matching your search")
		return nil
	}

	fmt.Printf("📝 Found %d notes:\n\n", len(notes))
	for _, note := range notes {
		title := "Untitled"
		if note.Title != nil {
			title = *note.Title
		}
		noteType := "text"
		if note.NoteType != nil {
			noteType = string(*note.NoteType)
		}

		fmt.Printf("🔹 %s (%s)\n", safeDeref(note.ID), title)
		fmt.Printf("   Type: %s\n", noteType)
		fmt.Printf("   Created: %s\n", note.Created)
		fmt.Println()
	}

	return nil
}

// safeDeref safely dereferences a string pointer
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getInjector extracts the DI injector from CLI context
func getInjector(c *cli.Context) do.Injector {
	if injector, ok := c.App.Metadata["injector"].(do.Injector); ok {
		return injector
	}
	panic("injector not found in CLI context - make sure DI is properly initialized")
}