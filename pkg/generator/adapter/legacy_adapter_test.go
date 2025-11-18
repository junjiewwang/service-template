package adapter

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestConfig creates a test configuration
func createTestConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:        "test-service",
			Description: "Test service",
			DeployDir:   "/opt",
		},
		Language: config.LanguageConfig{
			Type: "golang",
			Config: map[string]interface{}{
				"version": "1.21",
			},
		},
		Build: config.BuildConfig{
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
			Commands: config.BuildCommandsConfig{
				Build: "go build -o bin/app ./cmd/app",
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/app",
			},
		},
	}
}

func TestNewLegacyGeneratorAdapter(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	assert.NotNil(t, adapter)
	assert.Equal(t, cfg, adapter.config)
	assert.Equal(t, "/tmp/output", adapter.outputDir)
	assert.NotNil(t, adapter.genCtx)
}

func TestGenerateDockerfile(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	tests := []struct {
		name string
		arch string
	}{
		{"amd64", "amd64"},
		{"arm64", "arm64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := adapter.GenerateDockerfile(tt.arch)
			require.NoError(t, err)
			assert.NotEmpty(t, content)
			assert.Contains(t, content, "FROM")
			assert.Contains(t, content, "WORKDIR")
		})
	}
}

func TestGenerateCompose(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateCompose()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "services:")
	assert.Contains(t, content, cfg.Service.Name)
}

func TestGenerateMakefile(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateMakefile()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, ".PHONY:")
	assert.Contains(t, content, "docker-build")
}

func TestGenerateDevOps(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateDevOps()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "tad:")
	assert.Contains(t, content, "export_envs:")
}

func TestGenerateBuildScript(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateBuildScript()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "#!/bin/bash")
	assert.Contains(t, content, "build")
}

func TestGenerateDepsInstallScript(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateDepsInstallScript()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "#!/bin/bash")
}

func TestGenerateEntrypointScript(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateEntrypointScript()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "#!/bin/")
	assert.Contains(t, content, "Entrypoint")
}

func TestGenerateHealthcheckScript(t *testing.T) {
	cfg := createTestConfig()
	cfg.Runtime.Healthcheck = config.HealthcheckConfig{
		Enabled: true,
		Type:    "default",
	}
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateHealthcheckScript()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "#!/bin/")
}

func TestGenerateRtPrepareScript(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateRtPrepareScript()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "#!/bin/")
}

func TestGenerateBuildPluginsScript(t *testing.T) {
	cfg := createTestConfig()
	cfg.Plugins = config.PluginsConfig{
		InstallDir: "/opt/plugins",
		Items: []config.PluginConfig{
			{
				Name:           "test-plugin",
				Description:    "Test plugin",
				InstallCommand: "echo 'install plugin'",
				Required:       true,
			},
		},
	}
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.GenerateBuildPluginsScript()
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "#!/bin/bash")
	assert.Contains(t, content, "test-plugin")
}

func TestGenerateByType_NotFound(t *testing.T) {
	cfg := createTestConfig()
	adapter := NewLegacyGeneratorAdapter(cfg, "/tmp/output")

	content, err := adapter.generateByType("non-existent-type")
	assert.Error(t, err)
	assert.Empty(t, content)
	assert.Contains(t, err.Error(), "not found")
}
