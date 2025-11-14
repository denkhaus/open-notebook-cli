package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// SearchServices holds all the services needed for search commands
type SearchServices struct {
	SearchService services.SearchService
	Config        config.Service
	Logger        services.Logger
}

// getSearchServices retrieves all required services via dependency injection
func getSearchServices(ctx *cli.Context) (*SearchServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &SearchServices{
		SearchService: do.MustInvoke[services.SearchService](injector),
		Config:        do.MustInvoke[config.Service](injector),
		Logger:        do.MustInvoke[services.Logger](injector),
	}, nil
}

// validateSearchArgs validates common argument patterns for search commands
func validateSearchArgs(ctx *cli.Context, requireQuery bool) (string, error) {
	query := ctx.String("query")
	if requireQuery && query == "" {
		return "", errors.UsageError("Query is required",
			"Use --query flag to specify the search query")
	}
	return query, nil
}

// printSearchResults prints search results in a formatted table
func printSearchResults(results []models.SearchResult, searchType string) {
	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	// Display results in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tRELEVANCE\tTYPE")

	for _, result := range results {
		fmt.Fprintf(w, "%s\t%s\t%.3f\t%s\n",
			result.ID,
			utils.TruncateString(result.Title, 40),
			result.Relevance,
			result.ID[:4]) // Show first 4 chars as type indicator
	}

	w.Flush()
	fmt.Printf("\nFound %d results (%s search)\n", len(results), searchType)
}

// Handler functions with proper separation of concerns

// handleSearchQuery handles the search query command
func handleSearchQuery(ctx *cli.Context) error {
	services, err := getSearchServices(ctx)
	if err != nil {
		return err
	}

	query, err := validateSearchArgs(ctx, true)
	if err != nil {
		return err
	}

	searchType := ctx.String("type")
	limit := ctx.Int("limit")
	sources := ctx.String("sources")
	notes := ctx.String("notes")
	minScore := ctx.Float64("minimum-score")

	// Validate search type
	validTypes := map[string]bool{
		"vector": true,
		"text":   true,
	}
	if !validTypes[searchType] {
		return errors.UsageError("Invalid search type",
			"Supported types are: vector, text")
	}

	services.Logger.Info("Performing search query", "query", query, "type", searchType, "limit", limit)

	// Parse sources and notes if provided
	searchSources := true
	searchNotes := true
	if sources != "" || notes != "" {
		// If specific sources or notes are provided, only search those
		searchSources = sources != ""
		searchNotes = notes != ""
	}

	options := &models.SearchOptions{
		Type:          searchType,
		Limit:         limit,
		MinimumScore:  minScore,
		SearchSources: searchSources,
		SearchNotes:   searchNotes,
	}

	response, err := services.SearchService.Search(ctx.Context, query, options)
	if err != nil {
		return errors.APIError("Failed to perform search",
			"Check query parameters and API permissions")
	}

	printSearchResults(response.Results, response.SearchType)
	return nil
}

// handleSearchAsk handles the search ask command with streaming
func handleSearchAsk(ctx *cli.Context) error {
	services, err := getSearchServices(ctx)
	if err != nil {
		return err
	}

	question := ctx.String("question")
	streaming := ctx.Bool("streaming")
	answerModel := ctx.String("answer-model")
	strategyModel := ctx.String("strategy-model")
	finalModel := ctx.String("final-model")

	if question == "" {
		return errors.UsageError("Question is required",
			"Use --question flag to specify the question")
	}

	services.Logger.Info("Starting AI ask", "question", question, "streaming", streaming)

	options := &models.AskOptions{
		StrategyModel:    strategyModel,
		AnswerModel:      answerModel,
		FinalAnswerModel: finalModel,
	}

	fmt.Printf("🤖 Asking: %s\n", question)
	fmt.Println("─" + strings.Repeat("─", len(question)+10))
	fmt.Println()

	if streaming {
		// Streaming response
		chunkChan, err := services.SearchService.Ask(ctx.Context, question, options)
		if err != nil {
			return errors.APIError("Failed to start AI conversation",
				"Check API connection and model availability")
		}

		// Print streaming response
		for chunk := range chunkChan {
			if chunk.Error != "" {
				fmt.Fprintf(os.Stderr, "\n❌ Error: %s\n", chunk.Error)
				return fmt.Errorf("AI response error: %s", chunk.Error)
			}

			fmt.Print(chunk.Content)

			if chunk.Done {
				fmt.Println() // Add newline after completion
				break
			}
		}
	} else {
		// Non-streaming response
		response, err := services.SearchService.AskSimple(ctx.Context, question, options)
		if err != nil {
			return errors.APIError("Failed to get AI response",
				"Check API connection and model availability")
		}

		fmt.Println(response.Answer)
	}

	fmt.Println()
	fmt.Println("─" + strings.Repeat("─", 50))
	return nil
}

// handleSearchAskSimple handles the search ask-simple command
func handleSearchAskSimple(ctx *cli.Context) error {
	services, err := getSearchServices(ctx)
	if err != nil {
		return err
	}

	question := ctx.String("question")
	if question == "" {
		return errors.UsageError("Question is required",
			"Use --question flag to specify the question")
	}

	services.Logger.Info("Starting simple AI ask", "question", question)

	fmt.Printf("🤖 Asking (simple): %s\n", question)
	fmt.Println("─" + strings.Repeat("─", len(question)+18))
	fmt.Println()

	response, err := services.SearchService.AskSimple(ctx.Context, question, nil)
	if err != nil {
		return errors.APIError("Failed to get AI response",
			"Check API connection and model availability")
	}

	fmt.Println(response.Answer)
	fmt.Println()
	fmt.Println("─" + strings.Repeat("─", 50))
	return nil
}

