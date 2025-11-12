package testutil

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// NewTestConfig creates a test service configuration
func NewTestConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:        "test-service",
			Description: "Test service",
			DeployDir:   "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "http", Name: "http"},
			},
		},
		Language: config.LanguageConfig{
			Type:    "go",
			Version: "1.21",
			Config:  map[string]string{},
		},
		Build: config.BuildConfig{
			Commands: config.BuildCommandsConfig{
				Build: "go build -o bin/test-service ./cmd/test-service",
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
			DependencyFiles: config.DependencyFilesConfig{
				AutoDetect: true,
				Files:      []string{"go.mod", "go.sum"},
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/test-service",
			},
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "default",
			},
		},
		CI: config.CIConfig{
			ScriptDir:         ".tad/build/test-service",
			BuildConfigDir:    ".tad/build/test-service/build",
			ConfigTemplateDir: ".tad/build/test-service/config_template",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/tce",
			Items:      []config.PluginConfig{},
		},
		Metadata: config.MetadataConfig{
			Generator:       "svcgen",
			TemplateVersion: "1.0.0",
		},
	}
}

// NewTestConfigWithPlugins creates a test config with plugins
func NewTestConfigWithPlugins() *config.ServiceConfig {
	cfg := NewTestConfig()
	cfg.Plugins.Items = []config.PluginConfig{
		{
			Name:        "selfMonitor",
			Description: "Self monitoring plugin",
			DownloadURL: "https://example.com/selfMonitor.tar.gz",
			RuntimeEnv: []config.EnvironmentVariable{
				{Name: "TOOL_PATH", Value: "${PLUGIN_INSTALL_DIR}"},
			},
		},
	}
	return cfg
}
