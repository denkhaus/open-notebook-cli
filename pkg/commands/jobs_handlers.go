package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/errors"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/denkhaus/open-notebook-cli/pkg/utils"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// JobsServices holds all the services needed for job commands
type JobsServices struct {
	JobService shared.JobRepository
	Config     config.Service
	Logger     shared.Logger
}

// getJobsServices retrieves all required services via dependency injection
func getJobsServices(ctx *cli.Context) (*JobsServices, error) {
	injector, ok := ctx.App.Metadata["injector"].(do.Injector)
	if !ok {
		return nil, errors.UsageError("Dependency injector not found",
			"This command requires proper DI setup")
	}

	return &JobsServices{
		JobService: do.MustInvoke[shared.JobRepository](injector),
		Config:     do.MustInvoke[config.Service](injector),
		Logger:     do.MustInvoke[shared.Logger](injector),
	}, nil
}

// validateJobArgs validates common argument patterns for job commands
func validateJobArgs(ctx *cli.Context, requireJobID bool) (string, error) {
	if requireJobID {
		if ctx.NArg() < 1 {
			return "", fmt.Errorf("❌ Error: Missing job ID")
		}
		if ctx.NArg() > 1 {
			return "", fmt.Errorf("❌ Error: Too many arguments. Expected only job ID")
		}
	}

	jobID := ctx.Args().First()
	if requireJobID && jobID == "" {
		return "", errors.UsageError("Job ID is required",
			"Usage: open-notebook jobs <command> <job-id>")
	}

	return jobID, nil
}

// getJobStatusIcon returns an appropriate icon for job status
func getJobStatusIcon(status string) string {
	switch status {
	case "queued":
		return "⏳"
	case "running":
		return "🔄"
	case "completed":
		return "✅"
	case "failed":
		return "❌"
	default:
		return "❓"
	}
}

// formatJobDuration formats duration between created and updated timestamps
func formatJobDuration(created, updated string) string {
	if created == "" || updated == "" {
		return "N/A"
	}

	// Parse timestamps
	createdTime, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return "N/A"
	}

	updatedTime, err := time.Parse(time.RFC3339, updated)
	if err != nil {
		return "N/A"
	}

	duration := updatedTime.Sub(createdTime)
	if duration < time.Second {
		return "< 1s"
	}
	if duration < time.Minute {
		return fmt.Sprintf("%.0fs", duration.Seconds())
	}
	if duration < time.Hour {
		return fmt.Sprintf("%.0fm", duration.Minutes())
	}
	return fmt.Sprintf("%.1fh", duration.Hours())
}

// handleJobsList handles the jobs list command
func handleJobsList(ctx *cli.Context) error {
	services, err := getJobsServices(ctx)
	if err != nil {
		return err
	}

	statusFilter := ctx.String("status")
	limit := 20
	offset := 0

	if ctx.IsSet("limit") {
		limit = ctx.Int("limit")
	}
	if ctx.IsSet("offset") {
		offset = ctx.Int("offset")
	}

	services.Logger.Info("Listing background jobs", "status_filter", statusFilter)

	response, err := services.JobService.List(ctx.Context)
	if err != nil {
		return errors.APIError("Failed to list jobs",
			"Check API connection and permissions")
	}

	if len(response.Jobs) == 0 {
		fmt.Println("No background jobs found.")
		return nil
	}

	// Filter jobs based on status filter
	filteredJobs := []models.JobStatus{}
	for _, job := range response.Jobs {
		if statusFilter == "" || job.Status == statusFilter {
			filteredJobs = append(filteredJobs, job)
		}
	}

	// Apply pagination
	start := offset
	if start > len(filteredJobs) {
		start = len(filteredJobs)
	}
	end := start + limit
	if end > len(filteredJobs) {
		end = len(filteredJobs)
	}

	displayJobs := filteredJobs[start:end]

	if len(displayJobs) == 0 {
		fmt.Printf("No jobs found matching criteria.\n")
		return nil
	}

	// Display jobs in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tPROGRESS\tDURATION\tMESSAGE")

	for _, job := range displayJobs {
		progress := "N/A"
		if job.Progress != nil {
			progress = fmt.Sprintf("%.0f%%", *job.Progress*100)
		}

		duration := formatJobDuration(job.Created, utils.SafeDereferenceString(job.Updated))

		message := "N/A"
		if job.Message != nil {
			message = *job.Message
		}

		fmt.Fprintf(w, "%s\t%s %s\t%s\t%s\t%s\n",
			job.ID,
			getJobStatusIcon(job.Status),
			job.Status,
			progress,
			duration,
			message)
	}

	w.Flush()

	fmt.Printf("\nShowing %d jobs (use --limit and --offset for pagination)\n", len(displayJobs))
	return nil
}

// handleJobsStatus handles detailed job status display
func handleJobsStatus(ctx *cli.Context) error {
	jobID, err := validateJobArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getJobsServices(ctx)
	if err != nil {
		return err
	}

	watch := ctx.Bool("watch")

	fmt.Printf("📊 Job Status: %s\n", jobID)
	services.Logger.Info("Getting job status", "job_id", jobID)

	job, err := services.JobService.GetStatus(ctx.Context, jobID)
	if err != nil {
		return errors.APIError("Failed to get job status",
			"Check job ID and permissions")
	}

	// Display job details
	fmt.Printf("  ID:       %s\n", job.ID)
	fmt.Printf("  Status:   %s %s\n", getJobStatusIcon(job.Status), job.Status)

	if job.Progress != nil {
		fmt.Printf("  Progress: %.0f%%\n", *job.Progress*100)
	}

	duration := formatJobDuration(job.Created, utils.SafeDereferenceString(job.Updated))
	fmt.Printf("  Duration: %s\n", duration)

	if job.Message != nil {
		fmt.Printf("  Message:  %s\n", *job.Message)
	}

	fmt.Printf("  Created:  %s\n", utils.FormatTimestamp(job.Created))
	if job.Updated != nil && *job.Updated != "" {
		fmt.Printf("  Updated:  %s\n", utils.FormatTimestamp(*job.Updated))
	}

	if watch {
		fmt.Println("   🔄 Watching for status updates... (Press Ctrl+C to stop)")
		// Simple polling implementation for watching
		for i := 0; i < 10; i++ { // Watch for 10 iterations
			time.Sleep(2 * time.Second)

			// Get updated status
			updatedJob, err := services.JobService.GetStatus(ctx.Context, jobID)
			if err != nil {
				fmt.Printf("   Error checking status: %v\n", err)
				break
			}

			if updatedJob.Status != job.Status {
				fmt.Printf("   Status changed: %s → %s\n", job.Status, updatedJob.Status)
				job = updatedJob
			} else {
				fmt.Printf("   Checking status... (%d/10)\n", i+1)
			}

			if job.Status == "completed" || job.Status == "failed" {
				fmt.Printf("   Job finished with status: %s\n", job.Status)
				break
			}
		}
		fmt.Println("   Watch completed")
	}

	return nil
}

// handleJobsCancel handles job cancellation
func handleJobsCancel(ctx *cli.Context) error {
	jobID, err := validateJobArgs(ctx, true)
	if err != nil {
		return err
	}

	services, err := getJobsServices(ctx)
	if err != nil {
		return err
	}

	force := ctx.Bool("force")

	// Confirm cancellation unless force flag is used
	if !force {
		fmt.Printf("⚠️  Are you sure you want to cancel job '%s'? [y/N]: ", jobID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Cancellation cancelled")
			return nil
		}
	}

	fmt.Printf("🛑 Cancelling job: %s\n", jobID)
	services.Logger.Info("Cancelling job", "job_id", jobID)

	err = services.JobService.Cancel(ctx.Context, jobID)
	if err != nil {
		return errors.APIError("Failed to cancel job",
			"Check job ID and permissions. Job may not be cancellable")
	}

	fmt.Printf("✅ Job '%s' cancellation requested!\n", jobID)
	fmt.Println("   Note: The job may take a moment to stop gracefully")

	return nil
}
