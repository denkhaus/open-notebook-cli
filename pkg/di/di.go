package di

import (
	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/services"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// Bootstrap initializes the dependency injection container with all services
func Bootstrap(cliCtx *cli.Context) do.Injector {
	injector := do.New()

	// Inject CLI context for configuration service and other CLI-dependent services
	do.ProvideValue(injector, cliCtx)

	// Core infrastructure services
	do.Provide(injector, config.NewConfig)
	do.Provide(injector, services.NewLogger)
	do.Provide(injector, services.NewRetryableHTTPClient)
	do.Provide(injector, services.NewAuth)

	// Repository layer (only implemented ones)
	do.Provide(injector, services.NewNotebookRepository)
	do.Provide(injector, services.NewNoteRepository)

	// Service layer (only implemented ones)
	do.Provide(injector, services.NewNotebookService)

	return injector
}

// Service getter helpers for easy access (only implemented services)

func GetConfig(injector do.Injector) config.Service {
	return do.MustInvoke[config.Service](injector)
}

func GetLogger(injector do.Injector) services.Logger {
	return do.MustInvoke[services.Logger](injector)
}

func GetAuth(injector do.Injector) services.Auth {
	return do.MustInvoke[services.Auth](injector)
}

func GetHTTPClient(injector do.Injector) services.HTTPClient {
	return do.MustInvoke[services.HTTPClient](injector)
}

func GetNotebookService(injector do.Injector) services.NotebookService {
	return do.MustInvoke[services.NotebookService](injector)
}

// Repository getters (only implemented ones)
func GetNotebookRepository(injector do.Injector) services.NotebookRepository {
	return do.MustInvoke[services.NotebookRepository](injector)
}

func GetNoteRepository(injector do.Injector) services.NoteRepository {
	return do.MustInvoke[services.NoteRepository](injector)
}
