package commands

import (
	"github.com/urfave/cli/v2"
)

// PodcastCommand returns the podcast command
func PodcastCommand() *cli.Command {
	return &cli.Command{
		Name:  "podcast",
		Usage: "Podcast generation and episode management commands",
		Description: "Generate podcast episodes from your knowledge base and manage episodes.\n\n" +
			"Podcast generation transforms your content into audio format:\n" +
			"• Generate podcasts from sources and notebooks\n" +
			"• Track generation progress with job management\n" +
			"• List, download, and manage podcast episodes\n" +
			"• Support for multiple voices and languages\n" +
			"• Customizable podcast styles\n\n" +
			"Examples:\n" +
			"  onb podcast generate --query 'AI trends' --voice female     # Generate podcast\n" +
			"  onb podcast episodes list                                  # List episodes\n" +
			"  onb podcast episodes show abc123                           # Show episode details\n" +
			"  onb podcast episodes download abc123                       # Download audio\n" +
			"  onb podcast episodes delete abc123 --force                # Delete episode",
		Subcommands: []*cli.Command{
			podcastGenerateCommand(),
			podcastEpisodesCommand(),
		},
	}
}

// podcastGenerateCommand generates a new podcast episode
func podcastGenerateCommand() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Generate a new podcast episode",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "Query or topic for podcast generation",
			},
			&cli.StringSliceFlag{
				Name:    "sources",
				Aliases: []string{"s"},
				Usage:   "Source IDs to use as content (can be specified multiple times)",
			},
			&cli.StringSliceFlag{
				Name:    "notebooks",
				Aliases: []string{"n"},
				Usage:   "Notebook IDs to use as content (can be specified multiple times)",
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "Model ID to use for podcast generation (optional, uses default)",
			},
			&cli.StringFlag{
				Name:    "voice",
				Aliases: []string{"v"},
				Usage:   "Voice for podcast generation (e.g., male, female)",
				Value:   "female",
			},
			&cli.StringFlag{
				Name:    "language",
				Aliases: []string{"l"},
				Usage:   "Language code (e.g., en, es, fr, de)",
				Value:   "en",
			},
			&cli.StringFlag{
				Name:  "style",
				Usage: "Podcast style (e.g., educational, conversational, news)",
				Value: "educational",
			},
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch generation progress in real-time",
				Value:   false,
			},
		},
		Action: handlePodcastGenerate,
	}
}

// podcastEpisodesCommand manages podcast episodes
func podcastEpisodesCommand() *cli.Command {
	return &cli.Command{
		Name:  "episodes",
		Usage: "Podcast episode management commands",
		Description: "Manage your podcast episodes.\n\n" +
			"Episode management provides complete control over your content:\n" +
			"• List all episodes with pagination\n" +
			"• View detailed episode information\n" +
			"• Download episode audio files\n" +
			"• Delete unwanted episodes\n\n" +
			"Examples:\n" +
			"  onb podcast episodes list                              # List all episodes\n" +
			"  onb podcast episodes show abc123                       # Show episode details\n" +
			"  onb podcast episodes download abc123                   # Download audio file\n" +
			"  onb podcast episodes delete abc123 --force            # Delete episode",
		Subcommands: []*cli.Command{
			podcastEpisodesListCommand(),
			podcastEpisodesShowCommand(),
			podcastEpisodesDownloadCommand(),
			podcastEpisodesDeleteCommand(),
		},
	}
}

// podcastEpisodesListCommand lists podcast episodes
func podcastEpisodesListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all podcast episodes",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of episodes to return",
				Value:   20,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of episodes to skip",
				Value:   0,
			},
			&cli.StringFlag{
				Name:    "language",
				Aliases: []string{"lang"},
				Usage:   "Filter by language code",
			},
		},
		Action: handlePodcastEpisodesList,
	}
}

// podcastEpisodesShowCommand shows episode details
func podcastEpisodesShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a podcast episode",
		Args:  true,
		Action: handlePodcastEpisodesShow,
	}
}

// podcastEpisodesDownloadCommand downloads episode audio
func podcastEpisodesDownloadCommand() *cli.Command {
	return &cli.Command{
		Name:  "download",
		Usage: "Download podcast episode audio file",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path (optional, defaults to episode ID)",
			},
		},
		Action: handlePodcastEpisodesDownload,
	}
}

// podcastEpisodesDeleteCommand deletes a podcast episode
func podcastEpisodesDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a podcast episode",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: handlePodcastEpisodesDelete,
	}
}