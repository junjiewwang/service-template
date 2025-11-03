package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposeGenerator_Generate(t *testing.T) {
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
			Type:    "go",
			Version: "1.23",
		},
		Build: config.BuildConfig{
			OutputDir: "dist",
		},
		Runtime: config.RuntimeConfig{
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "http",
				HTTP: config.HTTPHealthConfig{
					Path:    "/health",
					Port:    8080,
					Timeout: 3,
				},
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

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	generator := NewComposeGenerator(cfg, engine, vars)

	content, err := generator.Generate()
	require.NoError(t, err, "Generate() should not return an error")
	assert.NotEmpty(t, content, "Generated content should not be empty")

	// Check that content contains expected sections
	expectedSections := []string{
		"services:",
		"test-service:",
		"image:",
		"ports:",
		"\"8080\"",
		"\"9090\"",
		"environment:",
		"GO_ENV=production",
		"LOG_LEVEL=info",
		"volumes:",
		"./config.yaml:/app/config.yaml",
		"deploy:",
		"resources:",
		"limits:",
		"cpus:",
		"memory:",
		"healthcheck:",
		"interval:",
		"timeout:",
		"retries:",
		"labels:",
		"restart:",
	}

	for _, section := range expectedSections {
		assert.Contains(t, content, section, "Generated compose.yaml should contain section: %s", section)
	}
}

func TestComposeGenerator_GenerateMinimal(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "minimal-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{Type: "go", Version: "1.23"},
		Build:    config.BuildConfig{OutputDir: "dist"},
		Runtime: config.RuntimeConfig{
			Healthcheck: config.HealthcheckConfig{Enabled: false},
			Startup:     config.StartupConfig{Command: "./app"},
		},
		LocalDev: config.LocalDevConfig{
			Compose: config.ComposeConfig{},
		},
		Metadata: config.MetadataConfig{GeneratedAt: "2024-01-01T00:00:00Z"},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	generator := NewComposeGenerator(cfg, engine, vars)

	content, err := generator.Generate()
	require.NoError(t, err, "Generate() should not return an error")
	assert.NotEmpty(t, content, "Generated content should not be empty")

	// Check minimal required sections
	requiredSections := []string{
		"services:",
		"minimal-service:",
		"image:",
		"ports:",
		"restart:",
	}

	for _, section := range requiredSections {
		assert.Contains(t, content, section, "Generated compose.yaml should contain required section: %s", section)
	}
}
