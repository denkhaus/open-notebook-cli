package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// AuthCommand returns the auth command and its subcommands
func AuthCommand() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Authentication commands",
		Subcommands: []*cli.Command{
			{
				Name:  "check",
				Usage: "Check authentication status",
				Action: func(ctx *cli.Context) error {
					// TODO: Implement auth check using DI injector
					fmt.Println("Auth check not yet implemented")
					return nil
				},
			},
		},
	}
}