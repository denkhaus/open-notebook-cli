package commands

import (
	"github.com/urfave/cli/v2"
)

// ChatCommand returns the chat command
func ChatCommand() *cli.Command {
	return &cli.Command{
		Name:  "chat",
		Usage: "Chat and conversation management commands",
		Description: "Interactive chat sessions with AI models for knowledge exploration.\n\n" +
			"Chat enables conversational access to your knowledge bases:\n" +
			"• Create and manage chat sessions\n" +
			"• Stream responses in real-time\n" +
			"• Use context from notebooks and sources\n" +
			"• Switch between different AI models\n" +
			"• View chat history and sessions\n\n" +
			"Examples:\n" +
			"  onb chat sessions list                  # List all chat sessions\n" +
			"  onb chat sessions create --title 'Q&A'  # Create new chat session\n" +
			"  onb chat start 'What is X?'              # Start new chat with default session\n" +
			"  onb chat start --session abc123 'How?'   # Continue existing session\n" +
			"  onb chat sessions delete abc123          # Delete a chat session",
		Subcommands: []*cli.Command{
			chatSessionsCommand(),
			chatStartCommand(),
			chatHistoryCommand(),
		},
	}
}

// chatSessionsCommand manages chat sessions
func chatSessionsCommand() *cli.Command {
	return &cli.Command{
		Name:  "sessions",
		Usage: "Chat session management commands",
		Description: "Manage your chat sessions.\n\n" +
			"Chat sessions maintain conversation context and history:\n" +
			"• List all your sessions with message counts\n" +
			"• Create new sessions with custom titles\n" +
			"• Delete unwanted sessions\n" +
			"• View session details and settings\n\n" +
			"Examples:\n" +
			"  onb chat sessions list                     # List all sessions\n" +
			"  onb chat sessions create --title 'Research' # Create new session\n" +
			"  onb chat sessions delete abc123             # Delete session",
		Subcommands: []*cli.Command{
			chatSessionsListCommand(),
			chatSessionsCreateCommand(),
			chatSessionsDeleteCommand(),
			chatSessionsShowCommand(),
		},
	}
}

// chatSessionsListCommand lists chat sessions
func chatSessionsListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all chat sessions",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of sessions to return",
				Value:   20,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of sessions to skip",
				Value:   0,
			},
			&cli.BoolFlag{
				Name:  "active-only",
				Usage: "Show only active sessions",
				Value: false,
			},
			&cli.StringFlag{
				Name:     "notebook",
				Aliases:  []string{"n"},
				Usage:    "Notebook ID to list chat sessions for (required)",
				Required: true,
			},
		},
		Action: handleChatSessionsList,
	}
}

// chatSessionsCreateCommand creates a new chat session
func chatSessionsCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new chat session",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Session title",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "Model ID to use for this session (optional, uses default)",
			},
		},
		Action: handleChatSessionsCreate,
	}
}

// chatSessionsDeleteCommand deletes a chat session
func chatSessionsDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a chat session",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: handleChatSessionsDelete,
	}
}

// chatSessionsShowCommand shows session details
func chatSessionsShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a chat session",
		Args:  true,
		Action: handleChatSessionsShow,
	}
}

// chatStartCommand starts or continues a chat conversation
func chatStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start a new chat or continue existing session",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Usage:   "Session ID to continue (creates new if not provided)",
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "Model ID to use for this conversation (optional, uses session default)",
			},
			&cli.StringFlag{
				Name:    "notebook",
				Aliases: []string{"n"},
				Usage:   "Notebook ID to use as context (optional)",
			},
			&cli.StringSliceFlag{
				Name:    "source",
				Aliases: []string{"sources"},
				Usage:   "Source IDs to include as context (can be specified multiple times)",
			},
			&cli.IntFlag{
				Name:    "max-tokens",
				Aliases: []string{"mt"},
				Usage:   "Maximum context tokens (optional)",
			},
			&cli.BoolFlag{
				Name:    "stream",
				Aliases: []string{"r"},
				Usage:   "Stream response in real-time (default: true)",
				Value:   true,
			},
		},
		Action: handleChatStart,
	}
}

// chatHistoryCommand shows chat history
func chatHistoryCommand() *cli.Command {
	return &cli.Command{
		Name:  "history",
		Usage: "Show chat history for a session",
		Args:  true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of messages to show",
				Value:   50,
			},
			&cli.BoolFlag{
				Name:    "reverse",
				Aliases: []string{"r"},
				Usage:   "Show messages in reverse order (newest first)",
				Value:   false,
			},
		},
		Action: handleChatHistory,
	}
}