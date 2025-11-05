package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposeGenerator_Generate(t *testing.T) {
	// Arrange: Setup test configuration with comprehensive settings
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
			Ports: []config.PortConfig{
				{Name: "http", Port: 8080, Protocol: "TCP"},
				{Name: "metrics", Port: 9090, Protocol: "TCP"},
			},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
		Build: config.BuildConfig{},
		Runtime: config.RuntimeConfig{
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "default",
			},
			Startup: config.StartupConfig{
				Command: "./app",
				Env: []config.EnvConfig{
					{Name: "GO_ENV", Value: "production"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		LocalDev: config.LocalDevConfig{
			Compose: config.ComposeConfig{
				Resources: config.ResourcesConfig{
					Limits: config.ResourceLimits{
						CPUs:   "0.5",
						Memory: "1G",
					},
					Reservations: config.ResourceLimits{
						CPUs:   "0.25",
						Memory: "512M",
					},
				},
				Volumes: []config.VolumeConfig{
					{
						Source: "./config.yaml",
						Target: "/app/config.yaml",
						Type:   "bind",
					},
				},
				Healthcheck: config.ComposeHealthConfig{
					Interval:    "30s",
					Timeout:     "10s",
					Retries:     3,
					StartPeriod: "40s",
				},
				Labels: map[string]string{
					"kompose.image-pull-policy": "IfNotPresent",
				},
			},
		},
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
	}

	// Act: Generate compose file content
	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	generator := NewComposeGenerator(cfg, engine, vars)

	content, err := generator.Generate()

	// Assert: Verify generation succeeded
	require.NoError(t, err, "Generate() should not return an error")
	require.NotEmpty(t, content, "Generated content should not be empty")

	t.Logf("Generated compose.yaml content length: %d bytes", len(content))

	// Assert: Verify all expected sections are present
	expectedSections := map[string]string{
		"services":             "services:",
		"service_name":         "test-service:",
		"image":                "image:",
		"ports":                "ports:",
		"port_8080":            "\"8080\"",
		"port_9090":            "\"9090\"",
		"environment":          "environment:",
		"env_go_env":           "GO_ENV=production",
		"env_log_level":        "LOG_LEVEL=info",
		"volumes":              "volumes:",
		"volume_config":        "./config.yaml:/app/config.yaml",
		"deploy":               "deploy:",
		"resources":            "resources:",
		"limits":               "limits:",
		"cpus":                 "cpus:",
		"memory":               "memory:",
		"healthcheck":          "healthcheck:",
		"healthcheck_interval": "interval:",
		"healthcheck_timeout":  "timeout:",
		"healthcheck_retries":  "retries:",
		"labels":               "labels:",
		"restart":              "restart:",
	}

	for name, section := range expectedSections {
		assert.Contains(t, content, section,
			"Generated compose.yaml should contain %s section: %s", name, section)
	}

	// Assert: Verify specific configuration values
	assert.Contains(t, content, "0.5", "Should contain CPU limit")
	assert.Contains(t, content, "1G", "Should contain memory limit")
	assert.Contains(t, content, "0.25", "Should contain CPU reservation")
	assert.Contains(t, content, "512M", "Should contain memory reservation")
	assert.Contains(t, content, "30s", "Should contain healthcheck interval")
	assert.Contains(t, content, "10s", "Should contain healthcheck timeout")
	assert.Contains(t, content, "3", "Should contain healthcheck retries")
	assert.Contains(t, content, "40s", "Should contain healthcheck start_period")
	assert.Contains(t, content, "IfNotPresent", "Should contain image pull policy label")
}

func TestComposeGenerator_GenerateMinimal(t *testing.T) {
	// Arrange: Setup minimal configuration
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "minimal-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
		Build: config.BuildConfig{},
		Runtime: config.RuntimeConfig{
			Healthcheck: config.HealthcheckConfig{Enabled: false},
			Startup:     config.StartupConfig{Command: "./app"},
		},
		LocalDev: config.LocalDevConfig{
			Compose: config.ComposeConfig{},
		},
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
	}

	// Act: Generate compose file with minimal configuration
	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	generator := NewComposeGenerator(cfg, engine, vars)

	content, err := generator.Generate()

	// Assert: Verify generation succeeded
	require.NoError(t, err, "Generate() should not return an error")
	require.NotEmpty(t, content, "Generated content should not be empty")

	t.Logf("Generated minimal compose.yaml content length: %d bytes", len(content))

	// Assert: Verify minimal required sections are present
	requiredSections := map[string]string{
		"services":     "services:",
		"service_name": "minimal-service:",
		"image":        "image:",
		"ports":        "ports:",
		"port_8080":    "\"8080\"",
		"restart":      "restart:",
	}

	for name, section := range requiredSections {
		assert.Contains(t, content, section,
			"Generated minimal compose.yaml should contain required %s section: %s", name, section)
	}

	// Assert: Verify optional sections are NOT present when not configured
	optionalSections := []string{
		"healthcheck:",
		"deploy:",
		"resources:",
	}

	for _, section := range optionalSections {
		if assert.NotContains(t, content, section,
			"Minimal compose.yaml should not contain optional section: %s", section) {
			t.Logf("Correctly omitted optional section: %s", section)
		}
	}
}
