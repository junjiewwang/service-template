package generator

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
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

	if err != nil {
		t.Fatalf("GenerateBuildScript() error = %v", err)
	}

	expectedStrings := []string{
		"#!/bin/bash",
		"Starting build process",
		"go build",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("Build script missing expected string: %s", expected)
		}
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

			if err != nil {
				t.Fatalf("GenerateDepsInstallScript() error = %v", err)
			}

			if !strings.Contains(content, "#!/bin/bash") {
				t.Error("Script missing shebang")
			}

			if !strings.Contains(content, tt.check) {
				t.Errorf("Script missing expected command: %s", tt.check)
			}
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

	if err != nil {
		t.Fatalf("GenerateRtPrepareScript() error = %v", err)
	}

	expectedStrings := []string{
		"#!/bin/bash",
		"runtime preparation",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("Rt prepare script missing expected string: %s", expected)
		}
	}
}
