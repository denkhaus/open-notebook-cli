package commands

import (
	"fmt"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// Services holds all the services needed for auth commands
type AuthServices struct {
	Auth   shared.Auth
	Config config.Service
	Logger shared.Logger
}

// getAuthServices retrieves all required services via dependency injection
func getAuthServices(ctx *cli.Context) (*AuthServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &AuthServices{
		Auth:   do.MustInvoke[shared.Auth](injector),
		Config: do.MustInvoke[config.Service](injector),
		Logger: do.MustInvoke[shared.Logger](injector),
	}, nil
}

// handleAuthCheck handles the auth check command
func handleAuthCheck(ctx *cli.Context) error {
	services, err := getAuthServices(ctx)
	if err != nil {
		return err
	}

	services.Logger.Info("Checking authentication status")

	// Check authentication
	isAuth := services.Auth.IsAuthenticated(ctx.Context)

	if isAuth {
		fmt.Println("✅ Authenticated successfully")

		// Try to get token to confirm it's valid
		token, err := services.Auth.GetToken(ctx.Context)
		if err != nil {
			return errors.AuthError("Failed to get auth token", err.Error())
		}

		fmt.Printf("🔑 Token: %s...\n", token[:utils.Min(20, len(token))])
		services.Logger.Info("Authentication check completed successfully")
	} else {
		fmt.Println("❌ Not authenticated")

		// Get password from config if available
		if services.Config.IsAuthenticated() {
			fmt.Println("💾 Using configured password for authentication")
		}

		// Try to authenticate with configured password
		if err := services.Auth.Authenticate(ctx.Context); err != nil {
			return errors.AuthError("Authentication failed",
				"Please check your password or API connection")
		}

		fmt.Println("✅ Authentication completed successfully")
		services.Logger.Info("Authentication successful")
	}

	return nil
}

// handleAuthLogin handles the auth login command
func handleAuthLogin(ctx *cli.Context) error {
	services, err := getAuthServices(ctx)
	if err != nil {
		return err
	}

	password := ctx.String("password")
	if password != "" {
		// Set password if provided
		services.Auth.SetPassword(password)
		services.Logger.Info("Password set from command line")
	} else {
		// Use password from config
		if !services.Config.IsAuthenticated() {
			return errors.AuthError("No password provided",
				"Use --password flag or set OPEN_NOTEBOOK_PASSWORD environment variable")
		}
		services.Logger.Info("Using configured password")
	}

	// Attempt authentication
	services.Logger.Info("Attempting authentication")
	if err := services.Auth.Authenticate(ctx.Context); err != nil {
		return errors.AuthError("Login failed",
			"Check your password and API connection")
	}

	fmt.Println("✅ Login successful!")
	services.Logger.Info("Authentication completed successfully")

	return nil
}

// AuthCommand returns the auth command and its subcommands
func AuthCommand() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Authentication commands",
		Description: "Manage OpenNotebook authentication and access control.\n\n" +
			"Authentication secures your OpenNotebook instance:\n" +
			"• Check current authentication status\n" +
			"• Login with password or token\n" +
			"• Validate API access permissions\n" +
			"• Manage authentication sessions\n\n" +
			"Examples:\n" +
			"  onb auth check                           # Check if authenticated\n" +
			"  onb auth login --password mypassword     # Login with password\n" +
			"  onb auth login                           # Login with configured password",
		Subcommands: []*cli.Command{
			{
				Name:   "check",
				Usage:  "Check authentication status",
				Action: handleAuthCheck,
			},
			{
				Name:  "login",
				Usage: "Authenticate with password",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Usage:    "API password (overrides config)",
						Required: false,
					},
				},
				Action: handleAuthLogin,
			},
		},
	}
}
