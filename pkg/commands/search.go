package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// SearchCommand returns the search command and its subcommands
func SearchCommand() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Search commands",
		Subcommands: []*cli.Command{
			{
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
				Action: func(ctx *cli.Context) error {
					query := ctx.String("query")
					searchType := ctx.String("type")
					limit := ctx.Int("limit")
					sources := ctx.String("sources")
					notes := ctx.String("notes")
					minScore := ctx.Float64("minimum-score")
					// TODO: Implement search query using DI injector
					fmt.Printf("Search query '%s' with type '%s', limit %d\n", query, searchType, limit)
					fmt.Printf("Sources: %s, Notes: %s, Min Score: %f\n", sources, notes, minScore)
					fmt.Println("Search query not yet implemented")
					return nil
				},
			},
			{
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
				Action: func(ctx *cli.Context) error {
					question := ctx.String("question")
					streaming := ctx.Bool("streaming")
					answerModel := ctx.String("answer-model")
					strategyModel := ctx.String("strategy-model")
					finalModel := ctx.String("final-model")
					// TODO: Implement search ask using DI injector
					fmt.Printf("Ask AI: '%s'\n", question)
					fmt.Printf("Streaming: %t, Answer Model: %s, Strategy Model: %s, Final Model: %s\n",
						streaming, answerModel, strategyModel, finalModel)
					fmt.Println("Ask AI not yet implemented")
					return nil
				},
			},
			{
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
				Action: func(ctx *cli.Context) error {
					question := ctx.String("question")
					// TODO: Implement search ask-simple using DI injector
					fmt.Printf("Ask AI (simple): '%s'\n", question)
					fmt.Println("Ask AI simple not yet implemented")
					return nil
				},
			},
		},
	}
}