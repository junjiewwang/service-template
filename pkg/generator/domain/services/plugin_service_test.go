package services

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginService_PrepareForBuildScript_StaticURL(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/tce",
			Items: []config.PluginConfig{
				{
					Name:           "selfMonitor",
					Description:    "TCE Self Monitor Tool",
					DownloadURL:    config.NewStaticDownloadURL("https://example.com/tool.sh"),
					InstallCommand: "echo 'Installing...'",
					Required:       true,
				},
			},
		},
	}

	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	// Act
	plugins := service.PrepareForBuildScript()

	// Assert
	require.Len(t, plugins, 1)
	plugin := plugins[0]

	assert.Equal(t, "selfMonitor", plugin.Name)
	assert.Equal(t, "TCE Self Monitor Tool", plugin.Description)
	assert.Equal(t, "https://example.com/tool.sh", plugin.DownloadURL)
	assert.Equal(t, "/tce", plugin.InstallDir)
	assert.Contains(t, plugin.URLResolverScript, `PLUGIN_DOWNLOAD_URL="https://example.com/tool.sh"`)
}

func TestPluginService_PrepareForBuildScript_ArchMapping(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/tce",
			Items: []config.PluginConfig{
				{
					Name:        "jre",
					Description: "Java Runtime Environment",
					DownloadURL: config.NewArchMappingDownloadURL(map[string]string{
						"x86_64":  "https://example.com/jdk-x86_64.tar.gz",
						"aarch64": "https://example.com/jdk-aarch64.tar.gz",
						"default": "https://example.com/jdk-generic.tar.gz",
					}),
					InstallCommand: "echo 'Installing JDK...'",
					Required:       false,
				},
			},
		},
	}

	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	// Act
	plugins := service.PrepareForBuildScript()

	// Assert
	require.Len(t, plugins, 1)
	plugin := plugins[0]

	assert.Equal(t, "jre", plugin.Name)
	assert.Equal(t, "Java Runtime Environment", plugin.Description)
	assert.Equal(t, "${PLUGIN_DOWNLOAD_URL}", plugin.DownloadURL) // Placeholder for arch mapping
	assert.Equal(t, "/tce", plugin.InstallDir)

	// Verify URL resolver script contains case statement
	assert.Contains(t, plugin.URLResolverScript, "ARCH=$(uname -m)")
	assert.Contains(t, plugin.URLResolverScript, "case \"${ARCH}\" in")
	assert.Contains(t, plugin.URLResolverScript, "https://example.com/jdk-x86_64.tar.gz")
	assert.Contains(t, plugin.URLResolverScript, "https://example.com/jdk-aarch64.tar.gz")
	assert.Contains(t, plugin.URLResolverScript, "https://example.com/jdk-generic.tar.gz")
}

func TestPluginService_GenerateURLResolverScript_StaticURL(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{Name: "test-service"},
	}
	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	urlConfig := config.NewStaticDownloadURL("https://example.com/plugin.tar.gz")

	// Act
	script := service.GenerateURLResolverScript(urlConfig)

	// Assert
	assert.Equal(t, `PLUGIN_DOWNLOAD_URL="https://example.com/plugin.tar.gz"`, script)
}

func TestPluginService_GenerateURLResolverScript_ArchMapping(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{Name: "test-service"},
	}
	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	urlConfig := config.NewArchMappingDownloadURL(map[string]string{
		"x86_64":  "https://example.com/plugin-x86_64.tar.gz",
		"aarch64": "https://example.com/plugin-aarch64.tar.gz",
	})

	// Act
	script := service.GenerateURLResolverScript(urlConfig)

	// Assert
	assert.Contains(t, script, "# Detect architecture and set download URL")
	assert.Contains(t, script, "ARCH=$(uname -m)")
	assert.Contains(t, script, "case \"${ARCH}\" in")
	assert.Contains(t, script, "x86_64|amd64)")
	assert.Contains(t, script, "https://example.com/plugin-x86_64.tar.gz")
	assert.Contains(t, script, "aarch64|arm64)")
	assert.Contains(t, script, "https://example.com/plugin-aarch64.tar.gz")
	assert.Contains(t, script, "echo \"ERROR: Unsupported architecture ${ARCH}\"")
	assert.Contains(t, script, "exit 1")
}

func TestPluginService_GenerateURLResolverScript_WithDefault(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{Name: "test-service"},
	}
	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	urlConfig := config.NewArchMappingDownloadURL(map[string]string{
		"x86_64":  "https://example.com/plugin-x86_64.tar.gz",
		"default": "https://example.com/plugin-generic.tar.gz",
	})

	// Act
	script := service.GenerateURLResolverScript(urlConfig)

	// Assert
	assert.Contains(t, script, "x86_64|amd64)")
	assert.Contains(t, script, "https://example.com/plugin-x86_64.tar.gz")
	assert.Contains(t, script, "*)")
	assert.Contains(t, script, "https://example.com/plugin-generic.tar.gz")
	assert.NotContains(t, script, "ERROR: Unsupported architecture")
}

func TestPluginService_NormalizeArchMapping(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{Name: "test-service"},
	}
	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "x86_64 and amd64 combined",
			input: map[string]string{
				"x86_64":  "https://example.com/x86.tar.gz",
				"aarch64": "https://example.com/arm.tar.gz",
			},
			expected: map[string]string{
				"x86_64|amd64":  "https://example.com/x86.tar.gz",
				"aarch64|arm64": "https://example.com/arm.tar.gz",
			},
		},
		{
			name: "with default",
			input: map[string]string{
				"x86_64":  "https://example.com/x86.tar.gz",
				"default": "https://example.com/generic.tar.gz",
			},
			expected: map[string]string{
				"x86_64|amd64": "https://example.com/x86.tar.gz",
				"default":      "https://example.com/generic.tar.gz",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := service.normalizeArchMapping(tt.input)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPluginService_HasPlugins(t *testing.T) {
	tests := []struct {
		name     string
		plugins  []config.PluginConfig
		expected bool
	}{
		{
			name:     "no plugins",
			plugins:  []config.PluginConfig{},
			expected: false,
		},
		{
			name: "has plugins",
			plugins: []config.PluginConfig{
				{Name: "plugin1"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ServiceConfig{
				Service: config.ServiceInfo{Name: "test-service"},
				Plugins: config.PluginsConfig{
					Items: tt.plugins,
				},
			}
			ctx := context.NewGeneratorContext(cfg, ".")
			engine := core.NewTemplateEngine()
			service := NewPluginService(ctx, engine)

			result := service.HasPlugins()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPluginService_GetInstallDir(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{Name: "test-service"},
		Plugins: config.PluginsConfig{
			InstallDir: "/custom/plugins",
		},
	}
	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	result := service.GetInstallDir()
	assert.Equal(t, "/custom/plugins", result)
}

func TestPluginService_PrepareForEntrypoint(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{Name: "test-service"},
		Plugins: config.PluginsConfig{
			InstallDir: "/tce",
			Items: []config.PluginConfig{
				{
					Name: "plugin1",
					RuntimeEnv: []config.EnvironmentVariable{
						{Name: "PLUGIN_PATH", Value: "${PLUGIN_INSTALL_DIR}/bin"},
					},
				},
				{
					Name:       "plugin2",
					RuntimeEnv: []config.EnvironmentVariable{},
				},
			},
		},
	}
	ctx := context.NewGeneratorContext(cfg, ".")
	engine := core.NewTemplateEngine()
	service := NewPluginService(ctx, engine)

	result := service.PrepareForEntrypoint()

	// Only plugin1 should be included (has runtime env)
	require.Len(t, result, 1)
	assert.Equal(t, "plugin1", result[0]["Name"])
	assert.Equal(t, "/tce", result[0]["InstallDir"])
}
