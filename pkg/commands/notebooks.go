package commands

import (
	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/urfave/cli/v2"
)

// NotebookServices holds all the services needed for notebook commands
type NotebookServices struct {
	NotebookService shared.NotebookService
	Config          config.Service
	Logger          shared.Logger
}

// NotebooksCommand returns the notebooks command and its subcommands
func NotebooksCommand() *cli.Command {
	return &cli.Command{
		Name:  "notebooks",
		Usage: "Knowledge base management commands",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all notebooks",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "archived",
						Usage: "Show archived notebooks",
						Value: false,
					},
				},
				Action: handleNotebooksList,
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
				Action: handleNotebooksCreate,
			},
			{
				Name:  "show",
				Usage: "Show notebook details",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
				},
				Action: handleNotebooksShow,
			},
			{
				Name:  "update",
				Usage: "Update notebook",
				Flags: []cli.Flag{
					&cli.StringFlag{
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
				Action: handleNotebooksUpdate,
			},
			{
				Name:  "delete",
				Usage: "Delete notebook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "Notebook ID",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "confirm",
						Aliases: []string{"y"},
						Usage:   "Skip confirmation prompt",
						Value:   false,
					},
				},
				Action: handleNotebooksDelete,
			},
			{
				Name:  "add-source",
				Usage: "Add source to notebook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "notebook",
						Aliases:  []string{"n"},
						Usage:    "Notebook ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Usage:    "Source ID",
						Required: true,
					},
				},
				Action: handleNotebooksAddSource,
			},
			{
				Name:  "remove-source",
				Usage: "Remove source from notebook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "notebook",
						Aliases:  []string{"n"},
						Usage:    "Notebook ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Usage:    "Source ID",
						Required: true,
					},
				},
				Action: handleNotebooksRemoveSource,
			},
		},
	}
}
