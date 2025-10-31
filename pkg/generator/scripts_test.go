package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptsGenerator_GenerateBuildScript(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
			Commands: config.BuildCommandsConfig{
				PreBuild:  "echo 'Pre-build'",
				Build:     "go build",
				PostBuild: "echo 'Post-build'",
			},
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	g := NewScriptsGenerator(cfg, engine, vars)
	content, err := g.GenerateBuildScript()
	require.NoError(t, err, "GenerateBuildScript() should not return an error")

	expectedStrings := []string{
		"#!/bin/bash",
		"TCS Service Build System",
		"test-service",
	}

	for _, expected := range expectedStrings {
		assert.Contains(t, content, expected, "Build script should contain expected string: %s", expected)
	}
}

func TestScriptsGenerator_GenerateDepsInstallScript(t *testing.T) {
	tests := []struct {
		name     string
		langType string
		check    string
	}{
		{
			name:     "golang dependencies",
			langType: "go",
			check:    "go mod download",
		},
		{
			name:     "python dependencies",
			langType: "python",
			check:    "pip install",
		},
		{
			name:     "nodejs dependencies",
			langType: "nodejs",
			check:    "npm install",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/opt/services",
				},
				Build: config.BuildConfig{
					OutputDir: "build",
				},
				Language: config.LanguageConfig{
					Type:    tt.langType,
					Version: "1.21",
				},
			}

			engine := NewTemplateEngine()
			vars := NewVariables(cfg)
			g := NewScriptsGenerator(cfg, engine, vars)
			content, err := g.GenerateDepsInstallScript()
			require.NoError(t, err, "GenerateDepsInstallScript() should not return an error")

			assert.Contains(t, content, "#!/bin/bash", "Script should contain shebang")
			assert.Contains(t, content, tt.check, "Script should contain expected command: %s", tt.check)
		})
	}
}

func TestScriptsGenerator_GenerateRtPrepareScript(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
		Plugins: []config.PluginConfig{
			{
				Name:           "test-plugin",
				Description:    "Test plugin",
				DownloadURL:    "https://example.com/plugin.tar.gz",
				InstallDir:     "/opt/plugins",
				InstallCommand: "echo 'Installing plugin'",
			},
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	g := NewScriptsGenerator(cfg, engine, vars)
	content, err := g.GenerateRtPrepareScript()
	require.NoError(t, err, "GenerateRtPrepareScript() should not return an error")

	expectedStrings := []string{
		"#!/bin/sh",
		"TCS Runtime Preparation",
	}

	for _, expected := range expectedStrings {
		assert.Contains(t, content, expected, "Rt prepare script should contain expected string: %s", expected)
	}
}
