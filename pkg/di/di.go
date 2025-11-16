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
	do.Provide(injector, services.NewSourceRepository)
	do.Provide(injector, services.NewModelRepository)
	do.Provide(injector, services.NewChatRepository)
	do.Provide(injector, services.NewNotebookRepository)
	do.Provide(injector, services.NewJobRepository)
	do.Provide(injector, services.NewPodcastRepository)
	do.Provide(injector, services.NewNoteRepository)
	do.Provide(injector, services.NewSearchRepository)

	// Service layer (only implemented ones)
	do.Provide(injector, services.NewNotebookService)
	do.Provide(injector, services.NewSearchService)
	do.Provide(injector, services.NewModelService)
	do.Provide(injector, services.NewSourceService)
	do.Provide(injector, services.NewPodcastService)
	do.Provide(injector, services.NewJobService)

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

func GetSearchService(injector do.Injector) services.SearchService {
	return do.MustInvoke[services.SearchService](injector)
}

// Repository getters (only implemented ones)
func GetNotebookRepository(injector do.Injector) services.NotebookRepository {
	return do.MustInvoke[services.NotebookRepository](injector)
}

func GetNoteRepository(injector do.Injector) services.NoteRepository {
	return do.MustInvoke[services.NoteRepository](injector)
}

func GetSearchRepository(injector do.Injector) services.SearchRepository {
	return do.MustInvoke[services.SearchRepository](injector)
}

func GetSourceRepository(injector do.Injector) services.SourceRepository {
	return do.MustInvoke[services.SourceRepository](injector)
}

func GetSourceService(injector do.Injector) services.SourceService {
	return do.MustInvoke[services.SourceService](injector)
}

func GetChatRepository(injector do.Injector) services.ChatRepository {
	return do.MustInvoke[services.ChatRepository](injector)
}

func GetModelRepository(injector do.Injector) services.ModelRepository {
	return do.MustInvoke[services.ModelRepository](injector)
}

func GetModelService(injector do.Injector) services.ModelService {
	return do.MustInvoke[services.ModelService](injector)
}

func GetJobRepository(injector do.Injector) services.JobRepository {
	return do.MustInvoke[services.JobRepository](injector)
}

func GetJobService(injector do.Injector) services.JobService {
	return do.MustInvoke[services.JobService](injector)
}

func GetPodcastRepository(injector do.Injector) services.PodcastRepository {
	return do.MustInvoke[services.PodcastRepository](injector)
}

func GetPodcastService(injector do.Injector) services.PodcastService {
	return do.MustInvoke[services.PodcastService](injector)
}
