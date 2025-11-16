package commands

import (
	"github.com/urfave/cli/v2"
)

// SettingsCommand returns the settings command
func SettingsCommand() *cli.Command {
	return &cli.Command{
		Name:  "settings",
		Usage: "Application settings management commands",
		Description: "Manage OpenNotebook application settings and preferences.\n\n" +
			"Settings control various aspects of application behavior:\n" +
			"• Content processing engines (docling, simple, auto)\n" +
			"• URL processing options (firecrawl, jina, simple)\n" +
			"• Embedding preferences and auto-deletion\n" +
			"• YouTube processing language preferences\n\n" +
			"Examples:\n" +
			"  onb settings get                           # Show all current settings\n" +
			"  onb settings set --engine docling          # Set document processing engine\n" +
			"  onb settings set --embed always           # Set embedding to always\n" +
			"  onb settings set --auto-delete yes        # Enable auto file deletion",
		Subcommands: []*cli.Command{
			settingsGetCommand(),
			settingsSetCommand(),
		},
	}
}

// settingsGetCommand gets current settings
func settingsGetCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Show current application settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format (table, json)",
				Value: "table",
			},
		},
		Action: handleSettingsGet,
	}
}

// settingsSetCommand updates settings
func settingsSetCommand() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Update application settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "engine",
				Usage: "Document processing engine (auto, docling, simple)",
			},
			&cli.StringFlag{
				Name:  "url-engine",
				Usage: "URL processing engine (auto, firecrawl, jina, simple)",
			},
			&cli.StringFlag{
				Name:  "embed",
				Usage: "Embedding preference (ask, always, never)",
			},
			&cli.StringFlag{
				Name:  "auto-delete",
				Usage: "Auto-delete files after processing (yes, no)",
			},
			&cli.StringSliceFlag{
				Name:  "youtube-languages",
				Usage: "YouTube preferred languages (can be specified multiple times)",
			},
		},
		Action: handleSettingsSet,
	}
}
