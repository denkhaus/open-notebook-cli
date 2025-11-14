package commands

import (
	"github.com/urfave/cli/v2"
)

// SourcesCommand returns the sources command
func SourcesCommand() *cli.Command {
	return &cli.Command{
		Name:  "sources",
		Usage: "Source management commands",
		Description: "Manage knowledge sources including text content, file uploads, and web links.\n\n" +
			"Sources are the foundation of your knowledge base. You can add:\n" +
			"• Text content directly\n" +
			"• Local files (PDF, DOCX, TXT, etc.)\n" +
			"• Web links and URLs\n" +
			"• Track processing status and insights\n\n" +
			"Examples:\n" +
			"  onb sources list                          # List all sources\n" +
			"  onb sources add --text \"My note\"         # Add text content\n" +
			"  onb sources add --link https://example.com # Add web link\n" +
			"  onb sources add --file document.pdf      # Upload file\n" +
			"  onb sources show <source-id>              # Show source details\n" +
			"  onb sources status <source-id>            # Check processing status",
		Subcommands: []*cli.Command{
			sourcesListCommand(),
			sourcesAddCommand(),
			sourcesShowCommand(),
			sourcesUpdateCommand(),
			sourcesDeleteCommand(),
			sourcesDownloadCommand(),
			sourcesStatusCommand(),
			sourcesRetryCommand(),
			sourcesInsightsCommand(),
		},
	}
}

// sourcesListCommand lists sources
func sourcesListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List sources with optional filtering",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "notebook",
				Aliases: []string{"n"},
				Usage:   "Filter by notebook ID",
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Filter by source type (text, link, upload)",
			},
			&cli.StringFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Filter by processing status (pending, running, completed, failed)",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of sources to return",
				Value:   20,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of sources to skip",
				Value:   0,
			},
			&cli.StringFlag{
				Name:  "sort",
				Usage: "Sort field (created, updated)",
				Value: "created",
			},
			&cli.StringFlag{
				Name:  "order",
				Usage: "Sort order (asc, desc)",
				Value: "desc",
			},
		},
		Action: handleSourcesList,
	}
}

// sourcesAddCommand adds a new source
func sourcesAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a new source (text, link, or file)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "text",
				Usage: "Text content to add as source",
			},
			&cli.StringFlag{
				Name:    "link",
				Aliases: []string{"url"},
				Usage:   "URL or web link to add as source",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Local file path to upload as source",
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "Title for the source",
			},
			&cli.StringSliceFlag{
				Name:    "topic",
				Aliases: []string{"topics"},
				Usage:   "Topics for categorization (can be specified multiple times)",
			},
			&cli.StringSliceFlag{
				Name:    "notebook",
				Aliases: []string{"notebooks", "n"},
				Usage:   "Notebook IDs to associate with (can be specified multiple times)",
			},
			&cli.BoolFlag{
				Name:  "async",
				Usage: "Process source asynchronously (default: true)",
				Value: true,
			},
		},
		Action: handleSourcesAdd,
	}
}

// sourcesShowCommand shows source details
func sourcesShowCommand() *cli.Command {
	return &cli.Command{
		Name:   "show",
		Usage:  "Show detailed information about a source",
		Args:   true,
		Action: handleSourcesShow,
	}
}

// sourcesUpdateCommand updates a source
func sourcesUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update source metadata (title, topics)",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "New title for the source",
			},
			&cli.StringSliceFlag{
				Name:    "topic",
				Aliases: []string{"topics"},
				Usage:   "New topics for categorization (replaces existing)",
			},
		},
		Action: handleSourcesUpdate,
	}
}

// sourcesDeleteCommand deletes a source
func sourcesDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a source",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: handleSourcesDelete,
	}
}

// sourcesDownloadCommand downloads source file
func sourcesDownloadCommand() *cli.Command {
	return &cli.Command{
		Name:  "download",
		Usage: "Download source file to local system",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path (default: current directory with original filename)",
			},
		},
		Action: handleSourcesDownload,
	}
}

// sourcesStatusCommand shows source processing status
func sourcesStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Check source processing status",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch status updates continuously",
				Value:   false,
			},
		},
		Action: handleSourcesStatus,
	}
}

// sourcesRetryCommand retries source processing
func sourcesRetryCommand() *cli.Command {
	return &cli.Command{
		Name:  "retry",
		Usage: "Retry source processing",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force reprocessing even if already completed",
				Value: false,
			},
		},
		Action: handleSourcesRetry,
	}
}

// sourcesInsightsCommand manages source insights
func sourcesInsightsCommand() *cli.Command {
	return &cli.Command{
		Name:  "insights",
		Usage: "Manage source insights",
		Subcommands: []*cli.Command{
			sourcesInsightsListCommand(),
			sourcesInsightsCreateCommand(),
		},
	}
}

// sourcesInsightsListCommand lists source insights
func sourcesInsightsListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List insights for a source",
		Args:   true,
		Action: handleSourcesInsightsList,
	}
}

// sourcesInsightsCreateCommand creates a new insight
func sourcesInsightsCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new insight for a source",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "content",
				Aliases:  []string{"c"},
				Usage:    "Insight content",
				Required: true,
			},
		},
		Action: handleSourcesInsightsCreate,
	}
}
