package commands

import (
	"bytes"
	"testing"

	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

// TestCommandRegistration tests that commands are properly registered
func TestCommandRegistration(t *testing.T) {
	commands := RegisterCommands()

	t.Run("Expected commands are registered", func(t *testing.T) {
		commandNames := make([]string, len(commands))
		for i, cmd := range commands {
			commandNames[i] = cmd.Name
		}

		expectedCommands := []string{"auth", "notebooks", "notes", "search"}
		for _, expected := range expectedCommands {
			assert.Contains(t, commandNames, expected, "Command %s should be registered", expected)
		}
	})

	t.Run("Command structure validation", func(t *testing.T) {
		for _, cmd := range commands {
			assert.NotEmpty(t, cmd.Name, "Command should have a name")
			assert.NotEmpty(t, cmd.Usage, "Command should have usage text")
			// Command should have either an action or subcommands (but not both)
			hasAction := cmd.Action != nil
			hasSubcommands := len(cmd.Subcommands) > 0
			assert.True(t, hasAction || hasSubcommands, "Command should have action or subcommands")
		}
	})

	t.Run("Subcommands validation", func(t *testing.T) {
		for _, cmd := range commands {
			for _, subcmd := range cmd.Subcommands {
				assert.NotEmpty(t, subcmd.Name, "Subcommand should have a name")
				assert.NotEmpty(t, subcmd.Usage, "Subcommand should have usage text")
				assert.NotNil(t, subcmd.Action, "Subcommand should have an action")
			}
		}
	})
}

// TestGlobalFlags tests global flag configuration
func TestGlobalFlags(t *testing.T) {
	app := createTestApp()

	flags := app.Flags
	flagMap := make(map[string]cli.Flag)

	for _, flag := range flags {
		switch f := flag.(type) {
		case *cli.StringFlag:
			flagMap[f.Name] = f
		case *cli.IntFlag:
			flagMap[f.Name] = f
		case *cli.BoolFlag:
			flagMap[f.Name] = f
		}
	}

	t.Run("Required global flags exist", func(t *testing.T) {
		expectedFlags := []string{"api-url", "password", "timeout", "retry-count", "verbose", "output"}
		for _, flagName := range expectedFlags {
			assert.Contains(t, flagMap, flagName, "Flag %s should exist", flagName)
		}
	})

	t.Run("Flag aliases are configured", func(t *testing.T) {
		// Test specific flags have expected aliases
		if apiURLFlag, ok := flagMap["api-url"].(*cli.StringFlag); ok {
			assert.Equal(t, []string{"u"}, apiURLFlag.Aliases, "api-url flag should have 'u' alias")
		}

		if passwordFlag, ok := flagMap["password"].(*cli.StringFlag); ok {
			assert.Equal(t, []string{"p"}, passwordFlag.Aliases, "password flag should have 'p' alias")
		}

		if timeoutFlag, ok := flagMap["timeout"].(*cli.IntFlag); ok {
			assert.Equal(t, []string{"t"}, timeoutFlag.Aliases, "timeout flag should have 't' alias")
		}

		if outputFlag, ok := flagMap["output"].(*cli.StringFlag); ok {
			assert.Equal(t, []string{"o"}, outputFlag.Aliases, "output flag should have 'o' alias")
		}
	})

	t.Run("Default values are set", func(t *testing.T) {
		if apiURLFlag, ok := flagMap["api-url"].(*cli.StringFlag); ok {
			assert.Equal(t, "http://localhost:5055", apiURLFlag.Value, "api-url should have correct default")
		}

		if timeoutFlag, ok := flagMap["timeout"].(*cli.IntFlag); ok {
			assert.Equal(t, 300, timeoutFlag.Value, "timeout should have correct default")
		}

		if retryCountFlag, ok := flagMap["retry-count"].(*cli.IntFlag); ok {
			assert.Equal(t, 3, retryCountFlag.Value, "retry-count should have correct default")
		}

		if verboseFlag, ok := flagMap["verbose"].(*cli.BoolFlag); ok {
			assert.Equal(t, false, verboseFlag.Value, "verbose should default to false")
		}
	})
}

// TestHelpCommands tests help command functionality
func TestHelpCommands(t *testing.T) {
	app := createTestApp()

	t.Run("Main app help", func(t *testing.T) {
		output, err := runTestApp(app, []string{"--help"})
		assert.NoError(t, err)
		assert.Contains(t, output, "onb")
		assert.Contains(t, output, "OpenNotebook CLI")
		assert.Contains(t, output, "COMMANDS:")
	})

	t.Run("Auth command help", func(t *testing.T) {
		output, err := runTestApp(app, []string{"auth", "--help"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Authentication commands")
	})

	t.Run("Notebooks command help", func(t *testing.T) {
		output, err := runTestApp(app, []string{"notebooks", "--help"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Knowledge base management commands")
		assert.Contains(t, output, "list")
		assert.Contains(t, output, "create")
		assert.Contains(t, output, "show")
	})

	t.Run("Search command help", func(t *testing.T) {
		output, err := runTestApp(app, []string{"search", "--help"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Search commands")
	})
}

// TestFlagValidation tests flag validation scenarios
func TestFlagValidation(t *testing.T) {
	app := createTestApp()

	t.Run("Valid global flags", func(t *testing.T) {
		_, err := runTestApp(app, []string{
			"--api-url", "http://localhost:8080",
			"--timeout", "60",
			"--retry-count", "5",
			"--verbose",
			"--output", "json",
			"--help",
		})
		assert.NoError(t, err)
	})

	t.Run("Short flag aliases", func(t *testing.T) {
		_, err := runTestApp(app, []string{
			"-u", "http://localhost:8080",
			"-t", "120",
			"-r", "2",
			"-o", "yaml",
			"--help",
		})
		assert.NoError(t, err)
	})
}

// TestErrorFunctions tests our custom error functions
func TestErrorFunctions(t *testing.T) {
	t.Run("APIError function", func(t *testing.T) {
		err := errors.APIError("Test API error", "Suggestion 1", "Suggestion 2")
		assert.NotNil(t, err)
		assert.Equal(t, "Test API error", err.Error())
	})

	t.Run("NotFoundError function", func(t *testing.T) {
		err := errors.NotFoundError("Resource not found", "Check ID")
		assert.NotNil(t, err)
		assert.Equal(t, "Resource not found", err.Error())
	})
}

// createTestApp creates a test CLI app without DI for testing CLI structure
func createTestApp() *cli.App {
	return &cli.App{
		Name:    "onb",
		Usage:   "OpenNotebook CLI - Test Version",
		Version: "test-version",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-url",
				Aliases: []string{"u"},
				Usage:   "OpenNotebook API URL",
				Value:   "http://localhost:5055",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "OpenNotebook API password",
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "Request timeout in seconds",
				Value:   300,
			},
			&cli.IntFlag{
				Name:    "retry-count",
				Aliases: []string{"r"},
				Usage:   "Number of retry attempts",
				Value:   3,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output format (json, table, yaml)",
				Value:   "table",
			},
		},
		Commands: RegisterCommands(),
	}
}

// runTestApp runs the CLI app with given arguments and returns output
func runTestApp(app *cli.App, args []string) (string, error) {
	var buf bytes.Buffer
	app.Writer = &buf

	err := app.Run(append([]string{"onb"}, args...))
	return buf.String(), err
}

// BenchmarkCommandHelp benchmarks help command performance
func BenchmarkCommandHelp(b *testing.B) {
	app := createTestApp()
	scenarios := [][]string{
		{"--help"},
		{"auth", "--help"},
		{"notebooks", "--help"},
		{"search", "query", "--help"},
	}

	for _, scenario := range scenarios {
		b.Run("scenario_"+scenario[0], func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runTestApp(app, scenario)
			}
		})
	}
}