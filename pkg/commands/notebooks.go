package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// NotebooksCommand returns the notebooks command and its subcommands
func NotebooksCommand() *cli.Command {
	return &cli.Command{
		Name:  "notebooks",
		Usage: "Knowledge base management commands",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all notebooks",
				Action: func(ctx *cli.Context) error {
					// TODO: Implement notebooks list using DI injector
					fmt.Println("Notebooks list not yet implemented")
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "Create a new notebook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Notebook name",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "description",
						Aliases: []string{"d"},
						Usage:   "Notebook description",
					},
				},
				Action: func(ctx *cli.Context) error {
					name := ctx.String("name")
					description := ctx.String("description")
					// TODO: Implement notebooks create using DI injector
					fmt.Printf("Create notebook '%s' with description '%s' not yet implemented\n", name, description)
					return nil
				},
			},
			{
				Name:  "show",
				Usage: "Show notebook details",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					id := ctx.Int("id")
					// TODO: Implement notebooks show using DI injector
					fmt.Printf("Show notebook %d not yet implemented\n", id)
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "Update notebook",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Notebook name",
					},
					&cli.StringFlag{
						Name:    "description",
						Aliases: []string{"d"},
						Usage:   "Notebook description",
					},
					&cli.BoolFlag{
						Name:    "archived",
						Aliases: []string{"a"},
						Usage:   "Archive notebook",
					},
				},
				Action: func(ctx *cli.Context) error {
					id := ctx.Int("id")
					// TODO: Implement notebooks update using DI injector
					fmt.Printf("Update notebook %d not yet implemented\n", id)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete notebook",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					id := ctx.Int("id")
					// TODO: Implement notebooks delete using DI injector
					fmt.Printf("Delete notebook %d not yet implemented\n", id)
					return nil
				},
			},
			{
				Name:  "add-source",
				Usage: "Add source to notebook",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Usage:    "Source ID",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					notebookID := ctx.Int("id")
					sourceID := ctx.Int("source")
					// TODO: Implement add source using DI injector
					fmt.Printf("Add source %d to notebook %d not yet implemented\n", sourceID, notebookID)
					return nil
				},
			},
			{
				Name:  "remove-source",
				Usage: "Remove source from notebook",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Usage:    "Source ID",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					notebookID := ctx.Int("id")
					sourceID := ctx.Int("source")
					// TODO: Implement remove source using DI injector
					fmt.Printf("Remove source %d from notebook %d not yet implemented\n", sourceID, notebookID)
					return nil
				},
			},
		},
	}
}