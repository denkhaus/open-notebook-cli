package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
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
				Name:  "notebook",
				Aliases: []string{"n"},
				Usage: "Filter by notebook ID",
			},
			&cli.StringFlag{
				Name:  "type",
				Aliases: []string{"t"},
				Usage: "Filter by source type (text, link, upload)",
			},
			&cli.StringFlag{
				Name:  "status",
				Aliases: []string{"s"},
				Usage: "Filter by processing status (pending, running, completed, failed)",
			},
			&cli.IntFlag{
				Name:  "limit",
				Aliases: []string{"l"},
				Usage: "Maximum number of sources to return",
				Value: 20,
			},
			&cli.IntFlag{
				Name:  "offset",
				Aliases: []string{"o"},
				Usage: "Number of sources to skip",
				Value: 0,
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
		Action: func(c *cli.Context) error {
			// TODO: Implement source listing with repository
			fmt.Println("🔍 Listing sources...")

			// Parse filter parameters
			params := &models.SourcesListParams{}
			if c.IsSet("notebook") {
				notebookID := c.String("notebook")
				params.NotebookID = &notebookID
			}
			if c.IsSet("limit") {
				limit := c.Int("limit")
				params.Limit = &limit
			}
			if c.IsSet("offset") {
				offset := c.Int("offset")
				params.Offset = &offset
			}
			if c.IsSet("sort") {
				sortBy := c.String("sort")
				params.SortBy = &sortBy
			}
			if c.IsSet("order") {
				sortOrder := c.String("order")
				params.SortOrder = &sortOrder
			}
			if c.IsSet("type") {
				sourceTypeStr := c.String("type")
				sourceType := models.SourceType(sourceTypeStr)
				params.Type = &sourceType
			}
			if c.IsSet("status") {
				statusStr := c.String("status")
				status := models.ProcessingStatus(statusStr)
				params.Status = &status
			}

			// TODO: Call repository to list sources
			fmt.Printf("   Filters: %+v\n", params)
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
				Name:  "link",
				Aliases: []string{"url"},
				Usage: "URL or web link to add as source",
			},
			&cli.StringFlag{
				Name:  "file",
				Aliases: []string{"f"},
				Usage: "Local file path to upload as source",
			},
			&cli.StringFlag{
				Name:  "title",
				Aliases: []string{"t"},
				Usage: "Title for the source",
			},
			&cli.StringSliceFlag{
				Name:  "topic",
				Aliases: []string{"topics"},
				Usage: "Topics for categorization (can be specified multiple times)",
			},
			&cli.StringSliceFlag{
				Name:  "notebook",
				Aliases: []string{"notebooks", "n"},
				Usage: "Notebook IDs to associate with (can be specified multiple times)",
			},
			&cli.BoolFlag{
				Name:  "async",
				Usage: "Process source asynchronously (default: true)",
				Value: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Validate that exactly one source type is provided
			textSet := c.IsSet("text")
			linkSet := c.IsSet("link")
			fileSet := c.IsSet("file")

			sourceTypes := 0
			if textSet {
				sourceTypes++
			}
			if linkSet {
				sourceTypes++
			}
			if fileSet {
				sourceTypes++
			}

			if sourceTypes == 0 {
				return fmt.Errorf("❌ Error: You must specify one of --text, --link, or --file")
			}
			if sourceTypes > 1 {
				return fmt.Errorf("❌ Error: You can only specify one of --text, --link, or --file at a time")
			}

			// Parse common parameters
			title := c.String("title")
			topics := c.StringSlice("topic")
			notebooks := c.StringSlice("notebook")
			async := c.Bool("async")

			var err error
			if textSet {
				err = addTextSource(c.String("text"), title, topics, notebooks, async)
			} else if linkSet {
				err = addLinkSource(c.String("link"), title, topics, notebooks, async)
			} else if fileSet {
				err = addFileSource(c.String("file"), title, topics, notebooks, async)
			}

			if err != nil {
				return fmt.Errorf("❌ Failed to add source: %w", err)
			}

			return nil
		},
	}
}

// sourcesShowCommand shows source details
func sourcesShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a source",
		Args:  true,
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			fmt.Printf("📄 Showing source details: %s\n", sourceID)

			// TODO: Implement source details retrieval
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
				Name:  "title",
				Aliases: []string{"t"},
				Usage: "New title for the source",
			},
			&cli.StringSliceFlag{
				Name:  "topic",
				Aliases: []string{"topics"},
				Usage: "New topics for categorization (replaces existing)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			if !c.IsSet("title") && !c.IsSet("topic") {
				return fmt.Errorf("❌ Error: You must specify at least one field to update (--title or --topic)")
			}

			sourceID := c.Args().First()
			title := c.String("title")
			topics := c.StringSlice("topic")

			fmt.Printf("✏️  Updating source: %s\n", sourceID)
			if title != "" {
				fmt.Printf("   Title: %s\n", title)
			}
			if len(topics) > 0 {
				fmt.Printf("   Topics: %s\n", strings.Join(topics, ", "))
			}

			// TODO: Implement source update
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
				Name:  "force",
				Aliases: []string{"f"},
				Usage: "Force deletion without confirmation",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			force := c.Bool("force")

			if !force {
				fmt.Printf("⚠️  Are you sure you want to delete source '%s'? [y/N]: ", sourceID)
				var response string
				fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "yes" {
					fmt.Println("❌ Deletion cancelled")
					return nil
				}
			}

			fmt.Printf("🗑️  Deleting source: %s\n", sourceID)

			// TODO: Implement source deletion
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
				Name:  "output",
				Aliases: []string{"o"},
				Usage: "Output file path (default: current directory with original filename)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			outputPath := c.String("output")

			fmt.Printf("⬇️  Downloading source file: %s\n", sourceID)
			if outputPath != "" {
				fmt.Printf("   Output: %s\n", outputPath)
			}

			// TODO: Implement source file download
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
				Name:  "watch",
				Aliases: []string{"w"},
				Usage: "Watch status updates continuously",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			watch := c.Bool("watch")

			fmt.Printf("📊 Getting source status: %s\n", sourceID)

			// TODO: Implement source status checking
			fmt.Println("   (Repository not yet implemented)")

			if watch {
				fmt.Println("   (Watch mode not yet implemented)")
			}

			return nil
		},
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
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			force := c.Bool("force")

			fmt.Printf("🔄 Retrying source processing: %s\n", sourceID)
			if force {
				fmt.Println("   Force reprocessing enabled")
			}

			// TODO: Implement source retry
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
		Name:  "list",
		Usage: "List insights for a source",
		Args:  true,
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			fmt.Printf("💡 Listing insights for source: %s\n", sourceID)

			// TODO: Implement insights listing
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
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
				Name:  "content",
				Aliases: []string{"c"},
				Usage: "Insight content",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing source ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only source ID")
			}

			sourceID := c.Args().First()
			content := c.String("content")

			fmt.Printf("💭 Creating insight for source: %s\n", sourceID)
			fmt.Printf("   Content: %s\n", content)

			// TODO: Implement insight creation
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
	}
}

// Helper functions for different source types

func addTextSource(text, title string, topics, notebooks []string, async bool) error {
	fmt.Println("📝 Adding text source...")
	fmt.Printf("   Text: %s\n", truncateString(text, 100))
	if title != "" {
		fmt.Printf("   Title: %s\n", title)
	}
	if len(topics) > 0 {
		fmt.Printf("   Topics: %s\n", strings.Join(topics, ", "))
	}
	if len(notebooks) > 0 {
		fmt.Printf("   Notebooks: %s\n", strings.Join(notebooks, ", "))
	}
	fmt.Printf("   Async: %t\n", async)

	// TODO: Implement text source creation
	fmt.Println("   (Repository not yet implemented)")
	return nil
}

func addLinkSource(link, title string, topics, notebooks []string, async bool) error {
	fmt.Println("🔗 Adding link source...")
	fmt.Printf("   URL: %s\n", link)
	if title != "" {
		fmt.Printf("   Title: %s\n", title)
	}
	if len(topics) > 0 {
		fmt.Printf("   Topics: %s\n", strings.Join(topics, ", "))
	}
	if len(notebooks) > 0 {
		fmt.Printf("   Notebooks: %s\n", strings.Join(notebooks, ", "))
	}
	fmt.Printf("   Async: %t\n", async)

	// TODO: Implement link source creation
	fmt.Println("   (Repository not yet implemented)")
	return nil
}

func addFileSource(filePath, title string, topics, notebooks []string, async bool) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	fmt.Println("📁 Adding file source...")
	fmt.Printf("   File: %s\n", filePath)
	if title != "" {
		fmt.Printf("   Title: %s\n", title)
	}
	if len(topics) > 0 {
		fmt.Printf("   Topics: %s\n", strings.Join(topics, ", "))
	}
	if len(notebooks) > 0 {
		fmt.Printf("   Notebooks: %s\n", strings.Join(notebooks, ", "))
	}
	fmt.Printf("   Async: %t\n", async)

	// TODO: Implement file source upload
	fmt.Println("   (Repository not yet implemented)")
	return nil
}

// Utility function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Utility function to parse comma-separated strings
func parseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}