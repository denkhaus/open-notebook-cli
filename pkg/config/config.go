package config

import (
	"fmt"
	"os"

	"github.com/samber/do/v2"
	"github.com/urfave/cli/v2"
)

// Service interface for configuration
type Service interface {
	GetAPIURL() string
	GetPassword() string
	GetTimeout() int
	GetRetryCount() int
	IsVerbose() bool
	GetOutput() string
	GetConfigDir() string
	IsAuthenticated() bool
	Validate() error
}

// Config implements the configuration service
type Config struct {
	apiURL     string
	password   string
	timeout    int
	retryCount int
	verbose    bool
	output     string
	configDir  string
}

// NewConfig creates a new configuration service by injecting the CLI context
// and extracting all resolved CLI flags and environment variables
func NewConfig(injector do.Injector) (Service, error) {
	// Inject the CLI context from urfav/cli
	cliContext := do.MustInvoke[*cli.Context](injector)

	// Extract values from CLI context (urfav/cli already resolved environment vars)
	apiURL := cliContext.String("api-url")
	password := cliContext.String("password")
	timeout := cliContext.Int("timeout")
	retryCount := cliContext.Int("retry-count")
	verbose := cliContext.Bool("verbose")
	output := cliContext.String("output")
	configDir := cliContext.String("config-dir")

	// Set defaults if not provided
	if apiURL == "" {
		apiURL = "http://localhost:5055"
	}
	if timeout <= 0 {
		timeout = 300 // 5 minutes default
	}
	if retryCount <= 0 {
		retryCount = 3
	}
	if output == "" {
		output = "table"
	}
	if configDir == "" {
		configDir = getDefaultConfigDir()
	}

	config := &Config{
		apiURL:     apiURL,
		password:   password,
		timeout:    timeout,
		retryCount: retryCount,
		verbose:    verbose,
		output:     output,
		configDir:  configDir,
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Interface implementation
func (c *Config) GetAPIURL() string     { return c.apiURL }
func (c *Config) GetPassword() string   { return c.password }
func (c *Config) GetTimeout() int       { return c.timeout }
func (c *Config) GetRetryCount() int    { return c.retryCount }
func (c *Config) IsVerbose() bool       { return c.verbose }
func (c *Config) GetOutput() string     { return c.output }
func (c *Config) GetConfigDir() string  { return c.configDir }
func (c *Config) IsAuthenticated() bool { return c.password != "" }

func (c *Config) Validate() error {
	if c.apiURL == "" {
		return fmt.Errorf("API URL is required")
	}

	if c.timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.retryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}

	validOutputs := map[string]bool{
		"json":  true,
		"table": true,
		"yaml":  true,
	}
	if !validOutputs[c.output] {
		return fmt.Errorf("invalid output format: %s (must be json, table, or yaml)", c.output)
	}

	return nil
}

func getDefaultConfigDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir + "/.config/open-notebook-cli"
	}
	return "/tmp/open-notebook-cli"
}
