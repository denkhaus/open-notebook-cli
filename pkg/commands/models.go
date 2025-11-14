package commands

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// ModelsCommand returns the models command
func ModelsCommand() *cli.Command {
	return &cli.Command{
		Name:  "models",
		Usage: "AI model management commands",
		Description: "Manage AI models for OpenNotebook including listing, adding, removing models,\n" +
			"configuring defaults, and checking provider availability.\n\n" +
			"OpenNotebook supports multiple AI providers and model types:\n" +
			"• Providers: ollama, openai, groq, xai, vertex, google, openrouter, anthropic, elevenlabs, voyage, azure, mistral, deepseek, openai-compatible\n" +
			"• Types: language, embedding, text_to_speech, speech_to_text\n\n" +
			"Examples:\n" +
			"  onb models list                           # List all available models\n" +
			"  onb models add --name gpt-4 --provider openai --type language # Add new model\n" +
			"  onb models defaults                       # Show default model assignments\n" +
			"  onb models defaults --set chat=gpt-4      # Set default chat model\n" +
			"  onb models providers                      # Check provider availability\n" +
			"  onb models delete model-id                # Remove a model",
		Subcommands: []*cli.Command{
			modelsListCommand(),
			modelsAddCommand(),
			modelsShowCommand(),
			modelsDeleteCommand(),
			modelsDefaultsCommand(),
			modelsProvidersCommand(),
		},
	}
}

// modelsListCommand lists available models
func modelsListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all available AI models",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Filter by model type (language, embedding, text_to_speech, speech_to_text)",
			},
			&cli.StringFlag{
				Name:    "provider",
				Aliases: []string{"p"},
				Usage:   "Filter by provider (openai, ollama, groq, etc.)",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of models to return",
				Value:   50,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of models to skip",
				Value:   0,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("🤖 Listing AI models...")

			// Parse filter parameters
			if c.IsSet("type") {
				modelType := c.String("type")
				fmt.Printf("   Filter by type: %s\n", modelType)
			}
			if c.IsSet("provider") {
				provider := c.String("provider")
				fmt.Printf("   Filter by provider: %s\n", provider)
			}
			if c.IsSet("limit") {
				limit := c.Int("limit")
				fmt.Printf("   Limit: %d\n", limit)
			}
			if c.IsSet("offset") {
				offset := c.Int("offset")
				fmt.Printf("   Offset: %d\n", offset)
			}

			// TODO: Implement model listing with repository
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
	}
}

// modelsAddCommand adds a new model
func modelsAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a new AI model",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Model name (e.g., gpt-4, llama2)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "provider",
				Aliases:  []string{"p"},
				Usage:    "Model provider (openai, ollama, groq, etc.)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "Model type (language, embedding, text_to_speech, speech_to_text)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			provider := c.String("provider")
			modelType := c.String("type")

			fmt.Println("➕ Adding new AI model...")
			fmt.Printf("   Name: %s\n", name)
			fmt.Printf("   Provider: %s\n", provider)
			fmt.Printf("   Type: %s\n", modelType)

			// Validate model type
			validTypes := map[string]bool{
				"language":       true,
				"embedding":      true,
				"text_to_speech": true,
				"speech_to_text": true,
			}
			if !validTypes[modelType] {
				return fmt.Errorf("❌ Error: Invalid model type '%s'. Valid types: language, embedding, text_to_speech, speech_to_text", modelType)
			}

			// TODO: Implement model creation with repository
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
	}
}

// modelsShowCommand shows detailed information about a model
func modelsShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a specific model",
		Args:  true,
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing model ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only model ID")
			}

			modelID := c.Args().First()
			fmt.Printf("🔍 Showing model details: %s\n", modelID)

			// TODO: Implement model details retrieval
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
	}
}

// modelsDeleteCommand deletes a model
func modelsDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete an AI model",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("❌ Error: Missing model ID")
			}
			if c.NArg() > 1 {
				return fmt.Errorf("❌ Error: Too many arguments. Expected only model ID")
			}

			modelID := c.Args().First()
			force := c.Bool("force")

			if !force {
				fmt.Printf("⚠️  Are you sure you want to delete model '%s'? [y/N]: ", modelID)
				var response string
				fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "yes" {
					fmt.Println("❌ Deletion cancelled")
					return nil
				}
			}

			fmt.Printf("🗑️  Deleting model: %s\n", modelID)

			// TODO: Implement model deletion
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
	}
}

// modelsDefaultsCommand manages default model assignments
func modelsDefaultsCommand() *cli.Command {
	return &cli.Command{
		Name:  "defaults",
		Usage: "Manage default model assignments",
		Subcommands: []*cli.Command{
			modelsDefaultsShowCommand(),
			modelsDefaultsSetCommand(),
		},
	}
}

// modelsDefaultsShowCommand shows current default models
func modelsDefaultsShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show current default model assignments",
		Action: func(c *cli.Context) error {
			fmt.Println("🎯 Current default models:")

			// TODO: Implement default models retrieval
			fmt.Println("   Default chat model: (not yet implemented)")
			fmt.Println("   Default embedding model: (not yet implemented)")
			fmt.Println("   Default transformation model: (not yet implemented)")
			fmt.Println("   Large context model: (not yet implemented)")
			fmt.Println("   Default text-to-speech model: (not yet implemented)")
			fmt.Println("   Default speech-to-text model: (not yet implemented)")
			fmt.Println("   Default embedding model: (not yet implemented)")
			fmt.Println("   Default tools model: (not yet implemented)")
			return nil
		},
	}
}

// modelsDefaultsSetCommand sets default models
func modelsDefaultsSetCommand() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Set default models",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "chat",
				Usage: "Default chat model ID",
			},
			&cli.StringFlag{
				Name:  "embedding",
				Usage: "Default embedding model ID",
			},
			&cli.StringFlag{
				Name:  "transformation",
				Usage: "Default transformation model ID",
			},
			&cli.StringFlag{
				Name:  "large-context",
				Usage: "Default large context model ID",
			},
			&cli.StringFlag{
				Name:  "tts",
				Usage: "Default text-to-speech model ID",
			},
			&cli.StringFlag{
				Name:  "stt",
				Usage: "Default speech-to-text model ID",
			},
			&cli.StringFlag{
				Name:  "tools",
				Usage: "Default tools model ID",
			},
		},
		Action: func(c *cli.Context) error {
			if !c.IsSet("chat") && !c.IsSet("embedding") && !c.IsSet("transformation") &&
				!c.IsSet("large-context") && !c.IsSet("tts") && !c.IsSet("stt") && !c.IsSet("tools") {
				return fmt.Errorf("❌ Error: You must specify at least one default model to set")
			}

			fmt.Println("⚙️  Setting default models...")

			if c.IsSet("chat") {
				fmt.Printf("   Chat: %s\n", c.String("chat"))
			}
			if c.IsSet("embedding") {
				fmt.Printf("   Embedding: %s\n", c.String("embedding"))
			}
			if c.IsSet("transformation") {
				fmt.Printf("   Transformation: %s\n", c.String("transformation"))
			}
			if c.IsSet("large-context") {
				fmt.Printf("   Large Context: %s\n", c.String("large-context"))
			}
			if c.IsSet("tts") {
				fmt.Printf("   Text-to-Speech: %s\n", c.String("tts"))
			}
			if c.IsSet("stt") {
				fmt.Printf("   Speech-to-Text: %s\n", c.String("stt"))
			}
			if c.IsSet("tools") {
				fmt.Printf("   Tools: %s\n", c.String("tools"))
			}

			// TODO: Implement default models setting
			fmt.Println("   (Repository not yet implemented)")
			return nil
		},
	}
}

// modelsProvidersCommand checks provider availability
func modelsProvidersCommand() *cli.Command {
	return &cli.Command{
		Name:  "providers",
		Usage: "Check AI model provider availability",
		Action: func(c *cli.Context) error {
			fmt.Println("🔌 Checking provider availability...")

			// List of known providers
			providers := []string{
				"ollama", "openai", "groq", "xai", "vertex", "google",
				"openrouter", "anthropic", "elevenlabs", "voyage",
				"azure", "mistral", "deepseek", "openai-compatible",
			}

			fmt.Println("   Known providers:")
			for _, provider := range providers {
				fmt.Printf("   • %s\n", provider)
			}

			// TODO: Implement provider availability checking
			fmt.Println("\n   Provider status: (Repository not yet implemented)")
			return nil
		},
	}
}
