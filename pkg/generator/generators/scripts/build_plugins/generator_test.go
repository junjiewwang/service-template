package build_plugins

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.ServiceConfig
		wantErr bool
		checks  []string
	}{
		{
			name: "generate build_plugins.sh with single plugin",
			config: &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/usr/local/services",
				},
				Plugins: config.PluginsConfig{
					InstallDir: "/tce",
					Items: []config.PluginConfig{
						{
							Name:           "selfMonitor",
							Description:    "TCE Self Monitor Tool",
							DownloadURL:    "https://example.com/download.sh",
							InstallCommand: `curl -fsSL "${PLUGIN_DOWNLOAD_URL}" | bash -s "${PLUGIN_WORK_DIR}"`,
							RuntimeEnv: []config.EnvironmentVariable{
								{Name: "TOOL_PATH", Value: "${PLUGIN_INSTALL_DIR}"},
							},
						},
					},
				},
			},
			wantErr: false,
			checks: []string{
				"Plugin Build System",
				"Building plugin: ${PLUGIN_NAME}",
				"TCE Self Monitor Tool",
				"https://example.com/download.sh",
				"TOOL_PATH",
				"install.sh",
				"Total plugins built: 1",
			},
		},
		{
			name: "no plugins configured",
			config: &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/usr/local/services",
				},
				Plugins: config.PluginsConfig{
					Items: []config.PluginConfig{},
				},
			},
			wantErr: false,
			checks:  []string{}, // Should return empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewGeneratorContext(tt.config, "/tmp/test")
			gen, err := New(ctx)
			if err != nil {
				t.Fatalf("Failed to create generator: %v", err)
			}

			content, err := gen.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no plugins, content should be empty
			if len(tt.config.Plugins.Items) == 0 {
				if content != "" {
					t.Errorf("Expected empty content for no plugins, got: %s", content)
				}
				return
			}

			// Check for expected strings
			for _, check := range tt.checks {
				if !strings.Contains(content, check) {
					t.Errorf("Generated content missing expected string: %s\nGenerated content:\n%s", check, content)
				}
			}

			// Verify script starts with shebang
			if !strings.HasPrefix(content, "#!/bin/bash") {
				t.Error("Generated script should start with #!/bin/bash")
			}
		})
	}
}

func TestGenerator_MultiplePlugins(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/usr/local/services",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/tce",
			Items: []config.PluginConfig{
				{
					Name:        "plugin1",
					Description: "First Plugin",
					DownloadURL: "https://example.com/plugin1.sh",
				},
				{
					Name:        "plugin2",
					Description: "Second Plugin",
					DownloadURL: "https://example.com/plugin2.sh",
				},
			},
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/test")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Check both plugins are included
	if !strings.Contains(content, "plugin1") {
		t.Error("Generated content missing plugin1")
	}
	if !strings.Contains(content, "plugin2") {
		t.Error("Generated content missing plugin2")
	}
	if !strings.Contains(content, "Total plugins built: 2") {
		t.Error("Generated content should show 2 plugins built")
	}
}
