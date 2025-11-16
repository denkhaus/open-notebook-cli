package commands

import (
	"github.com/urfave/cli/v2"
)

// TransformationsCommand returns the transformations command
func TransformationsCommand() *cli.Command {
	return &cli.Command{
		Name:  "transformations",
		Usage: "Transformation management commands",
		Description: "Manage text transformations for content processing and analysis.\n\n" +
			"Transformations allow you to define custom processing pipelines:\n" +
			"• Create reusable text processing templates\n" +
			"• Apply transformations to any text content\n" +
			"• Use specific models for different transformation types\n" +
			"• Execute transformations on-demand or set as defaults\n\n" +
			"Examples:\n" +
			"  onb transformations list                     # List all transformations\n" +
			"  onb transformations create --name summary   # Create new transformation\n" +
			"  onb transformations execute <id> --text \"sample text\" # Execute transformation\n" +
			"  onb transformations show <id>               # Show transformation details",
		Subcommands: []*cli.Command{
			transformationsListCommand(),
			transformationsCreateCommand(),
			transformationsShowCommand(),
			transformationsUpdateCommand(),
			transformationsDeleteCommand(),
			transformationsExecuteCommand(),
		},
	}
}

// transformationsListCommand lists transformations
func transformationsListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all available transformations",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of transformations to return",
				Value:   20,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of transformations to skip",
				Value:   0,
			},
		},
		Action: handleTransformationsList,
	}
}

// transformationsCreateCommand creates a new transformation
func transformationsCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new transformation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Transformation name (identifier)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Transformation title (display name)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "description",
				Aliases:  []string{"d"},
				Usage:    "Transformation description",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "prompt",
				Aliases:  []string{"p"},
				Usage:    "Transformation prompt template",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "apply-default",
				Usage: "Set as default transformation for sources",
				Value: false,
			},
		},
		Action: handleTransformationsCreate,
	}
}

// transformationsShowCommand shows transformation details
func transformationsShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a transformation",
		Args:  true,
		Action: handleTransformationsShow,
	}
}

// transformationsUpdateCommand updates a transformation
func transformationsUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update transformation metadata",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "New transformation name",
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "New transformation title",
			},
			&cli.StringFlag{
				Name:    "description",
				Aliases: []string{"d"},
				Usage:   "New transformation description",
			},
			&cli.StringFlag{
				Name:    "prompt",
				Aliases: []string{"p"},
				Usage:   "New transformation prompt template",
			},
			&cli.BoolFlag{
				Name:  "apply-default",
				Usage: "Set as default transformation for sources",
			},
		},
		Action: handleTransformationsUpdate,
	}
}

// transformationsDeleteCommand deletes a transformation
func transformationsDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a transformation",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: handleTransformationsDelete,
	}
}

// transformationsExecuteCommand executes a transformation
func transformationsExecuteCommand() *cli.Command {
	return &cli.Command{
		Name:  "execute",
		Usage: "Execute a transformation on text",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "text",
				Aliases:  []string{"t"},
				Usage:    "Input text to transform",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "Model ID to use for transformation (optional, uses default if not specified)",
			},
			&cli.BoolFlag{
				Name:  "stream",
				Usage: "Stream the transformation output in real-time",
				Value: false,
			},
		},
		Action: handleTransformationsExecute,
	}
}