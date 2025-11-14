package commands

import (
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
		Name:   "show",
		Usage:  "Show detailed information about a specific note",
		Args:   true,
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
				Name:    "notebook",
				Aliases: []string{"n"},
				Usage:   "Search within specific notebook ID",
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
