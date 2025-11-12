package commands

import "github.com/urfave/cli/v2"

// RegisterCommands registers all CLI commands and returns them
func RegisterCommands() []*cli.Command {
	return []*cli.Command{
		AuthCommand(),
		NotebooksCommand(),
		NotesCommand(),
		SearchCommand(),
		// TODO: Add more commands as they are implemented
		// ChatCommand(),
		// SourcesCommand(),
		// ModelsCommand(),
		// JobsCommand(),
		// SettingsCommand(),
		// PodcastCommand(),
	}
}