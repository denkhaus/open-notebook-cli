package di

import (
	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// Bootstrap initializes the dependency injection container with all services
func Bootstrap(cliCtx *cli.Context) do.Injector {
	injector := do.New()

	// Inject CLI context for configuration service and other CLI-dependent services
	do.ProvideValue(injector, cliCtx)

	// Provide configuration service (depends on CLI context)
	do.Provide[config.Service](injector, config.NewConfig)

	// TODO: Add other service providers here as they are implemented
	// do.Provide(injector, logger.New)
	// do.Provide(injector, auth.NewService)
	// do.Provide(injector, client.NewHTTPClient)
	// do.Provide(injector, notebooks.NewService)
	// do.Provide(injector, notes.NewService)
	// do.Provide(injector, search.NewService)
	// do.Provide(injector, chat.NewService)
	// do.Provide(injector, sources.NewService)
	// do.Provide(injector, models.NewService)
	// do.Provide(injector, jobs.NewService)
	// do.Provide(injector, settings.NewService)

	return injector
}

// GetConfig retrieves the configuration service from the injector
func GetConfig(injector do.Injector) config.Service {
	return do.MustInvoke[config.Service](injector)
}
