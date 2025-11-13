package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGetters(t *testing.T) {
	cfg := &Config{
		apiURL:     "http://test:5055",
		password:   "test-pass",
		timeout:    60,
		retryCount: 5,
		verbose:    true,
		output:     "yaml",
		configDir:  "/tmp/config",
	}

	assert.Equal(t, "http://test:5055", cfg.GetAPIURL())
	assert.Equal(t, "test-pass", cfg.GetPassword())
	assert.Equal(t, 60, cfg.GetTimeout())
	assert.Equal(t, 5, cfg.GetRetryCount())
	assert.True(t, cfg.IsVerbose())
	assert.Equal(t, "yaml", cfg.GetOutput())
	assert.Equal(t, "/tmp/config", cfg.GetConfigDir())
}

func TestConfigIsAuthenticated(t *testing.T) {
	cfg := &Config{password: ""}
	assert.False(t, cfg.IsAuthenticated())

	cfg.password = "admin"
	assert.True(t, cfg.IsAuthenticated())
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		expectErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				apiURL:  "http://localhost:5055",
				timeout: 30,
				output:  "table", // Add output format
			},
			expectErr: false,
		},
		{
			name: "invalid API URL",
			cfg: &Config{
				apiURL:  "",
				timeout: 30,
				output:  "table",
			},
			expectErr: true,
		},
		{
			name: "invalid timeout",
			cfg: &Config{
				apiURL:  "http://localhost:5055",
				timeout: 0,
				output: "table",
			},
			expectErr: true,
		},
		{
			name: "invalid retry count",
			cfg: &Config{
				apiURL:     "http://localhost:5055",
				timeout:    30,
				retryCount: -1,
			},
			expectErr: true,
		},
		{
			name: "invalid output format",
			cfg: &Config{
				apiURL:  "http://localhost:5055",
				timeout: 30,
				output:  "invalid",
			},
			expectErr: true,
		},
		{
			name: "valid output formats",
			cfg: &Config{
				apiURL:  "http://localhost:5055",
				timeout: 30,
				output:  "json",
			},
			expectErr: false,
		},
		{
			name: "valid output format table",
			cfg: &Config{
				apiURL:  "http://localhost:5055",
				timeout: 30,
				output:  "table",
			},
			expectErr: false,
		},
		{
			name: "valid output format yaml",
			cfg: &Config{
				apiURL:  "http://localhost:5055",
				timeout: 30,
				output:  "yaml",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDefaultConfigDir(t *testing.T) {
	dir := getDefaultConfigDir()
	assert.NotEmpty(t, dir)
	// The directory might contain "open-notebook-cli" instead of ".open-notebook"
	assert.True(t, strings.Contains(dir, "open-notebook"))
}