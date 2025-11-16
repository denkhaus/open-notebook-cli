package commands

import (
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
		Action: handleModelsList,
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
		Action: handleModelsAdd,
	}
}

// modelsShowCommand shows detailed information about a model
func modelsShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show detailed information about a specific model",
		Args:  true,
		Action: handleModelsShow,
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
		Action: handleModelsDelete,
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
		Action: handleModelsDefaultsShow,
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
		Action: handleModelsDefaultsSet,
	}
}

// modelsProvidersCommand checks provider availability
func modelsProvidersCommand() *cli.Command {
	return &cli.Command{
		Name:  "providers",
		Usage: "Check AI model provider availability",
		Action: handleModelsProviders,
	}
}
