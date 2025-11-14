package commands

import (
	"github.com/urfave/cli/v2"
)

// SearchCommand returns the search command and its subcommands
func SearchCommand() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Search commands",
		Description: "Search your knowledge base and ask AI questions.\n\n" +
			"Powerful search functionality with vector and text search capabilities.\n" +
			"AI-powered question answering with streaming responses.\n\n" +
			"Examples:\n" +
			"  onb search query --query \"machine learning\"           # Vector search\n" +
			"  onb search query --query \"python\" --type text       # Text search\n" +
			"  onb search ask --question \"What is AI?\"             # Streaming AI response\n" +
			"  onb search ask-simple --question \"Explain ML\"       # Simple AI response",
		Subcommands: []*cli.Command{
			searchQueryCommand(),
			searchAskCommand(),
			searchAskSimpleCommand(),
		},
	}
}

// searchQueryCommand performs knowledge base search
func searchQueryCommand() *cli.Command {
	return &cli.Command{
		Name:  "query",
		Usage: "Search knowledge base",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "query",
				Aliases:  []string{"q"},
				Usage:    "Search query",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Search type (vector, text)",
				Value:   "vector",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Result limit",
				Value:   10,
			},
			&cli.StringFlag{
				Name:    "sources",
				Aliases: []string{"s"},
				Usage:   "Comma-separated source IDs to search",
			},
			&cli.StringFlag{
				Name:    "notes",
				Aliases: []string{"n"},
				Usage:   "Comma-separated note IDs to search",
			},
			&cli.Float64Flag{
				Name:    "minimum-score",
				Aliases: []string{"m"},
				Usage:   "Minimum similarity score for vector search",
			},
		},
		Action: handleSearchQuery,
	}
}

// searchAskCommand asks AI a question with streaming response
func searchAskCommand() *cli.Command {
	return &cli.Command{
		Name:  "ask",
		Usage: "Ask AI a question",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "question",
				Aliases:  []string{"q"},
				Usage:    "Question to ask",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "streaming",
				Aliases: []string{"s"},
				Usage:   "Enable streaming response",
				Value:   true,
			},
			&cli.StringFlag{
				Name:    "answer-model",
				Aliases: []string{"a"},
				Usage:   "Model for answer generation",
			},
			&cli.StringFlag{
				Name:    "strategy-model",
				Aliases: []string{"S"},
				Usage:   "Model for search strategy",
			},
			&cli.StringFlag{
				Name:    "final-model",
				Aliases: []string{"f"},
				Usage:   "Model for final response",
			},
		},
		Action: handleSearchAsk,
	}
}

// searchAskSimpleCommand asks AI a question (non-streaming)
func searchAskSimpleCommand() *cli.Command {
	return &cli.Command{
		Name:  "ask-simple",
		Usage: "Ask AI a question (non-streaming)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "question",
				Aliases:  []string{"q"},
				Usage:    "Question to ask",
				Required: true,
			},
		},
		Action: handleSearchAskSimple,
	}
}