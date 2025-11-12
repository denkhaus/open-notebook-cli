package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// NotesCommand returns the notes command and its subcommands
func NotesCommand() *cli.Command {
	return &cli.Command{
		Name:  "notes",
		Usage: "Note management commands",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List notes",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "notebook",
						Aliases: []string{"n"},
						Usage:   "Filter by notebook ID",
					},
				},
				Action: func(ctx *cli.Context) error {
					notebookID := ctx.Int("notebook")
					// TODO: Implement notes list using DI injector
					fmt.Printf("List notes for notebook %d not yet implemented\n", notebookID)
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "Add a new note",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "content",
						Aliases:  []string{"c"},
						Usage:    "Note content",
						Required: true,
					},
					&cli.IntFlag{
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
						Aliases: []string{"T"},
						Usage:   "Note type (human, ai)",
						Value:   "human",
					},
				},
				Action: func(ctx *cli.Context) error {
					content := ctx.String("content")
					notebookID := ctx.Int("notebook")
					title := ctx.String("title")
					noteType := ctx.String("type")
					// TODO: Implement notes add using DI injector
					fmt.Printf("Add note '%s' of type '%s' to notebook %d not yet implemented\n", title, noteType, notebookID)
					fmt.Printf("Content: %s\n", content)
					return nil
				},
			},
			{
				Name:  "show",
				Usage: "Show note details",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Note ID",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					id := ctx.Int("id")
					// TODO: Implement notes show using DI injector
					fmt.Printf("Show note %d not yet implemented\n", id)
					return nil
				},
			},
			{
				Name:  "edit",
				Usage: "Edit existing note",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Note ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "content",
						Aliases: []string{"c"},
						Usage:   "New note content",
					},
					&cli.StringFlag{
						Name:    "title",
						Aliases: []string{"t"},
						Usage:   "New note title",
					},
				},
				Action: func(ctx *cli.Context) error {
					id := ctx.Int("id")
					// TODO: Implement notes edit using DI injector
					fmt.Printf("Edit note %d not yet implemented\n", id)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete note",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Note ID",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					id := ctx.Int("id")
					// TODO: Implement notes delete using DI injector
					fmt.Printf("Delete note %d not yet implemented\n", id)
					return nil
				},
			},
		},
	}
}