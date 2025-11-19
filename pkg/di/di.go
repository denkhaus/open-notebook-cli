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
