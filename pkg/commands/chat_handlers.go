package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// ChatServices holds all the services needed for chat commands
type ChatServices struct {
	ChatService shared.ChatRepository
	Config      config.Service
	Logger      shared.Logger
}

// getChatServices retrieves all required services via dependency injection
func getChatServices(ctx *cli.Context) (*ChatServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &ChatServices{
		ChatService: do.MustInvoke[shared.ChatRepository](injector),
		Config:      do.MustInvoke[config.Service](injector),
		Logger:      do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validateChatArgs validates common argument patterns for chat commands
func validateChatArgs(ctx *cli.Context, requireSessionID bool) (string, error) {
	if requireSessionID {
		if ctx.NArg() < 1 {
			return "", errors.MissingArgument("session ID", ctx.Command.Name)
		}
		if ctx.NArg() > 1 {
			return "", errors.TooManyArguments("session ID", ctx.Command.Name)
		}
	}

	sessionID := ctx.Args().First()
	if requireSessionID && sessionID == "" {
		return "", errors.UsageError("Session ID is required",
			"Usage: open-notebook chat sessions <command> <session-id>")
	}

	return sessionID, nil
}

// printChatSession prints formatted chat session information
func printChatSession(session *models.ChatSession) {
	status := "🔴 Inactive"
	if session.IsActive {
		status = "🟢 Active"
	}

	fmt.Printf("  ID:           %s\n", session.ID)
	fmt.Printf("  Title:        %s\n", session.Title)
	fmt.Printf("  Model:        %s\n", session.ModelID)
	fmt.Printf("  Status:       %s\n", status)
	fmt.Printf("  Messages:     %d\n", session.MessageCount)
	fmt.Printf("  Created:      %s\n", utils.FormatTimestamp(session.Created))
	fmt.Printf("  Updated:      %s\n", utils.FormatTimestamp(session.Updated))
}

// handleChatSessionsList handles the chat sessions list command
func handleChatSessionsList(ctx *cli.Context) error {
	services, err := getChatServices(ctx)
	if err != nil {
		return err
	}

	notebookID := ctx.String("notebook")
	services.Logger.Info("Listing chat sessions...", "notebook_id", notebookID)

	// API requires notebook_id parameter - fail loud if missing
	if notebookID == "" {
		return errors.RequiredField("Notebook ID", ctx.Command.Name)
	}

	response, err := services.ChatService.ListSessionsForNotebook(ctx.Context, notebookID)
	if err != nil {
		return errors.APIError("Failed to list chat sessions",
			"Check API connection and permissions")
	}

	if len(*response) == 0 {
		fmt.Println("No chat sessions found.")
		return nil
	}

	// Filter sessions
	filteredSessions := []models.ChatSession{}
	activeOnly := ctx.Bool("active-only")

	for _, session := range *response {
		if activeOnly && !session.IsActive {
			continue
		}
		filteredSessions = append(filteredSessions, session)
	}

	if len(filteredSessions) == 0 {
		fmt.Println("No sessions found matching criteria.")
		return nil
	}

	// Apply pagination
	limit := ctx.Int("limit")
	offset := ctx.Int("offset")

	start := offset
	if start > len(filteredSessions) {
		start = len(filteredSessions)
	}
	end := start + limit
	if end > len(filteredSessions) {
		end = len(filteredSessions)
	}

	displaySessions := filteredSessions[start:end]

	// Display sessions in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tMODEL\tMESSAGES\tSTATUS\tCREATED")

	for _, session := range displaySessions {
		status := "🔴"
		if session.IsActive {
			status = "🟢"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
			session.ID,
			utils.TruncateString(session.Title, 25),
			session.ModelID,
			session.MessageCount,
			status,
			utils.FormatTimestamp(session.Created))
	}

	w.Flush()

	fmt.Printf("\nShowing %d sessions (use --limit and --offset for pagination)\n", len(displaySessions))
	return nil
}

// handleChatSessionsCreate handles creating a new chat session
func handleChatSessionsCreate(ctx *cli.Context) error {
	services, err := getChatServices(ctx)
	if err != nil {
		return err
	}

	title := ctx.String("title")
	modelID := ctx.String("model")

	services.Logger.Info("Creating chat session", "title", title, "model", modelID)

	request := &models.ChatCreateRequest{
		Title: title,
	}

	if modelID != "" {
		request.ModelID = &modelID
	}

	session, err := services.ChatService.CreateSession(ctx.Context, request)
	if err != nil {
		return errors.APIError("Failed to create chat session",
			"Check input parameters and API permissions")
	}

	fmt.Printf("✅ Chat session created successfully!\n")
	printChatSession(session)

	fmt.Printf("\nStart chatting with:\n")
	fmt.Printf("  onb chat start --session %s \"Your message here\"\n", session.ID)

	return nil
}

// handleChatSessionsDelete handles deleting a chat session
func handleChatSessionsDelete(ctx *cli.Context) error {
	sessionID, err := validateChatArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getChatServices(ctx)
	if err != nil {
		return err
	}

	// Confirm deletion unless force flag is used
	if !ctx.Bool("force") {
		fmt.Printf("⚠️  Are you sure you want to delete chat session '%s'? [y/N]: ", sessionID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Deletion cancelled")
			return nil
		}
	}

	fmt.Printf("🗑️  Deleting chat session: %s\n", sessionID)
	services.Logger.Info("Deleting chat session", "session_id", sessionID)

	err = services.ChatService.DeleteSession(ctx.Context, sessionID)
	if err != nil {
		return errors.APIError("Failed to delete chat session",
			"Check session ID and permissions")
	}

	fmt.Printf("✅ Chat session '%s' deleted successfully!\n", sessionID)
	return nil
}

// handleChatSessionsShow handles showing chat session details
func handleChatSessionsShow(ctx *cli.Context) error {
	sessionID, err := validateChatArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getChatServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Getting chat session details", "session_id", sessionID)

	session, err := services.ChatService.GetSession(ctx.Context, sessionID)
	if err != nil {
		return errors.APIError("Failed to get chat session details",
			"Check session ID and permissions")
	}

	fmt.Printf("💬 Chat Session Details:\n")
	printChatSession(session)

	return nil
}

// handleChatStart handles starting or continuing a chat conversation
func handleChatStart(ctx *cli.Context) error {
	services, err := getChatServices(ctx)
	if err != nil {
		return err
	}

	if ctx.NArg() < 1 {
		return errors.UsageError("Message is required",
			"Usage: open-notebook chat start [--session <id>] \"Your message\"")
	}

	message := ctx.Args().First()
	sessionID := ctx.String("session")
	modelID := ctx.String("model")
	notebookID := ctx.String("notebook")
	sources := ctx.StringSlice("source")
	maxTokens := ctx.Int("max-tokens")
	stream := ctx.Bool("stream")

	services.Logger.Info("Starting chat", "session_id", sessionID, "message", utils.TruncateString(message, 50))

	// Build context request if any context options are provided
	var context *models.ChatContextRequest
	if notebookID != "" || len(sources) > 0 || maxTokens > 0 {
		context = &models.ChatContextRequest{
			NotebookID: notebookID,
			Sources:    sources,
		}
		if maxTokens > 0 {
			context.MaxTokens = &maxTokens
		}
	}

	// Build chat execution request
	request := &models.ChatExecuteRequest{
		Message: message,
		Stream:  stream,
		Context: context,
	}

	if sessionID != "" {
		request.SessionID = sessionID
	}

	if modelID != "" {
		request.ModelID = &modelID
	}

	fmt.Printf("💬 Starting chat...\n")
	if sessionID != "" {
		fmt.Printf("  Session: %s\n", sessionID)
	}
	if modelID != "" {
		fmt.Printf("  Model:   %s\n", modelID)
	}
	fmt.Printf("  Message: %s\n", utils.TruncateString(message, 100))

	if context != nil {
		fmt.Printf("  Context: ")
		if notebookID != "" {
			fmt.Printf("Notebook: %s ", notebookID)
		}
		if len(sources) > 0 {
			fmt.Printf("Sources: %v ", sources)
		}
		if maxTokens > 0 {
			fmt.Printf("Max Tokens: %d ", maxTokens)
		}
		fmt.Printf("\n")
	}

	if stream {
		return handleStreamingChat(services, ctx, request)
	} else {
		return handleSimpleChat(services, ctx, request)
	}
}

// handleStreamingChat handles streaming chat responses
func handleStreamingChat(services *ChatServices, ctx *cli.Context, request *models.ChatExecuteRequest) error {
	fmt.Printf("🔄 Assistant (streaming):\n")

	chunkChan, err := services.ChatService.StreamChat(ctx.Context, request)
	if err != nil {
		return errors.APIError("Failed to start chat stream",
			"Check connection and permissions")
	}

	// Stream chunks
	for chunk := range chunkChan {
		if chunk.Error != "" {
			return errors.APIError("Chat streaming error", chunk.Error)
		}

		fmt.Printf("%s", chunk.Content)

		if chunk.Done {
			break
		}
	}

	fmt.Printf("\n\n✅ Chat completed\n")
	return nil
}

// handleSimpleChat handles simple (non-streaming) chat responses
func handleSimpleChat(services *ChatServices, ctx *cli.Context, request *models.ChatExecuteRequest) error {
	fmt.Printf("🤖 Thinking...\n")

	response, err := services.ChatService.ExecuteChat(ctx.Context, request)
	if err != nil {
		return errors.APIError("Failed to execute chat",
			"Check connection and permissions")
	}

	fmt.Printf("💬 Assistant:\n%s\n\n", response.Content)
	fmt.Printf("✅ Chat completed (Session: %s, Message: %s)\n", response.SessionID, response.MessageID)
	return nil
}

// handleChatHistory handles showing chat history
func handleChatHistory(ctx *cli.Context) error {
	sessionID, err := validateChatArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getChatServices(ctx)
	if err != nil {
		return err
	}

	limit := ctx.Int("limit")
	reverse := ctx.Bool("reverse")

	services.Logger.Info("Getting chat history", "session_id", sessionID, "limit", limit)

	messages, err := services.ChatService.GetMessages(ctx.Context, sessionID)
	if err != nil {
		return errors.APIError("Failed to get chat history",
			"Check session ID and permissions")
	}

	if len(messages) == 0 {
		fmt.Printf("No messages found in session '%s'.\n", sessionID)
		return nil
	}

	// Apply limit
	if limit > 0 && len(messages) > limit {
		if reverse {
			messages = messages[len(messages)-limit:]
		} else {
			messages = messages[:limit]
		}
	}

	fmt.Printf("💬 Chat History (Session: %s)\n", sessionID)
	fmt.Printf("   Showing %d messages\n\n", len(messages))

	// Display messages
	if reverse {
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			displayChatMessage(msg)
		}
	} else {
		for _, msg := range messages {
			displayChatMessage(msg)
		}
	}

	return nil
}

// displayChatMessage displays a single chat message
func displayChatMessage(msg *models.ChatMessage) {
	roleIcon := "❓"
	switch msg.Role {
	case "user":
		roleIcon = "👤"
	case "assistant":
		roleIcon = "🤖"
	case "system":
		roleIcon = "⚙️"
	}

	fmt.Printf("%s %s (%s):\n", roleIcon, msg.Role, utils.FormatTimestamp(msg.Created))
	fmt.Printf("   %s\n\n", msg.Content)
}
