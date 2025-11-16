package commands

import (
	"github.com/urfave/cli/v2"
)

// JobsCommand returns the jobs command
func JobsCommand() *cli.Command {
	return &cli.Command{
		Name:  "jobs",
		Usage: "Background job management commands",
		Description: "Manage and monitor background jobs in OpenNotebook.\n\n" +
			"Background jobs handle long-running operations like:\n" +
			"• File processing and embedding\n" +
			"• Content analysis and transformation\n" +
			"• Knowledge base rebuilding\n" +
			"• Batch operations on multiple items\n\n" +
			"Examples:\n" +
			"  onb jobs list                           # List all background jobs\n" +
			"  onb jobs status <job-id>                # Check job status\n" +
			"  onb jobs cancel <job-id>                # Cancel a running job\n" +
			"  onb jobs list --status running          # Show only running jobs",
		Subcommands: []*cli.Command{
			jobsListCommand(),
			jobsStatusCommand(),
			jobsCancelCommand(),
		},
	}
}

// jobsListCommand lists background jobs
func jobsListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List background jobs with optional filtering",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Filter by job status (queued, running, completed, failed)",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of jobs to return",
				Value:   20,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"o"},
				Usage:   "Number of jobs to skip",
				Value:   0,
			},
		},
		Action: handleJobsList,
	}
}

// jobsStatusCommand shows detailed job status
func jobsStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show detailed status of a specific job",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch job status updates continuously",
				Value:   false,
			},
		},
		Action: handleJobsStatus,
	}
}

// jobsCancelCommand cancels a running job
func jobsCancelCommand() *cli.Command {
	return &cli.Command{
		Name:  "cancel",
		Usage: "Cancel a running or queued job",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force cancellation without confirmation",
				Value:   false,
			},
		},
		Action: handleJobsCancel,
	}
}