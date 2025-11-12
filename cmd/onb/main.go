package main

import (
	"fmt"
	"os"

	"github.com/denkhaus/open-notebook-cli/pkg/commands"
	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/urfave/cli/v2"
)

// Build information (set by Makefile)
var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	app := &cli.App{
		Name:    "onb",
		Usage:   "OpenNotebook CLI - Manage your knowledge bases from the command line",
		Version: fmt.Sprintf("%s (built %s)", version, buildTime),
		Authors: []*cli.Author{
			{Name: "denkhaus", Email: "denkhaus@example.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-url",
				Aliases: []string{"u"},
				Usage:   "OpenNotebook API URL",
				EnvVars: []string{"OPEN_NOTEBOOK_API_URL"},
				Value:   "http://localhost:5055",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "OpenNotebook API password",
				EnvVars: []string{"OPEN_NOTEBOOK_PASSWORD"},
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "Request timeout in seconds",
				EnvVars: []string{"OPEN_NOTEBOOK_TIMEOUT"},
				Value:   300,
			},
			&cli.IntFlag{
				Name:    "retry-count",
				Aliases: []string{"r"},
				Usage:   "Number of retry attempts",
				EnvVars: []string{"OPEN_NOTEBOOK_RETRY_COUNT"},
				Value:   3,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				EnvVars: []string{"OPEN_NOTEBOOK_VERBOSE"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output format (json, table, yaml)",
				EnvVars: []string{"OPEN_NOTEBOOK_OUTPUT"},
				Value:   "table",
			},
			&cli.StringFlag{
				Name:    "config-dir",
				Aliases: []string{"c"},
				Usage:   "Configuration directory",
				EnvVars: []string{"OPEN_NOTEBOOK_CONFIG_DIR"},
			},
		},
		Commands: commands.RegisterCommands(),
		Before: func(ctx *cli.Context) error {
			// Initialize dependency injection container with all services
			injector := di.Bootstrap(ctx)

			// Store injector in context for commands to access
			ctx.App.Metadata = map[string]interface{}{
				"injector": injector,
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}