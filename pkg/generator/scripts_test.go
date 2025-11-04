package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptsGenerator_GenerateBuildScript(t *testing.T) {
	// Arrange: Setup service configuration with build commands
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
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

	// Act: Generate build script
	content, err := g.GenerateBuildScript()
	require.NoError(t, err, "GenerateBuildScript() should not return an error")
	require.NotEmpty(t, content, "Build script should not be empty")

	t.Logf("Generated build script: %d bytes", len(content))

	// Assert: Verify expected content
	expectedStrings := map[string]string{
		"shebang":      "#!/bin/bash",
		"header":       "TCS Service Build System",
		"service_name": "test-service",
	}

	for name, expected := range expectedStrings {
		assert.Contains(t, content, expected,
			"Build script should contain %s: %s", name, expected)
	}
	t.Logf("✓ Verified all %d expected sections present", len(expectedStrings))
}

func TestScriptsGenerator_GenerateDepsInstallScript(t *testing.T) {
	tests := []struct {
		name     string
		langType string
		check    string
	}{
		{
			name:     "golang dependencies",
			langType: "golang",
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
			// Arrange: Setup configuration for specific language
			cfg := &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/opt/services",
				},
				Build: config.BuildConfig{
				},
				Language: config.LanguageConfig{
					Type:    tt.langType,
					Version: "1.21",
				},
			}

			engine := NewTemplateEngine()
			vars := NewVariables(cfg)
			g := NewScriptsGenerator(cfg, engine, vars)

			// Act: Generate dependency install script
			content, err := g.GenerateDepsInstallScript()
			require.NoError(t, err, "GenerateDepsInstallScript() should not return an error")
			require.NotEmpty(t, content, "Deps install script should not be empty")

			t.Logf("Generated deps install script for %s: %d bytes", tt.langType, len(content))

			// Assert: Verify script content
			assert.Contains(t, content, "#!/bin/bash", "Script should contain shebang")
			assert.Contains(t, content, tt.check,
				"Script should contain expected command for %s: %s", tt.langType, tt.check)
			t.Logf("✓ Verified %s dependency command present", tt.langType)
		})
	}
}

func TestScriptsGenerator_GenerateRtPrepareScript(t *testing.T) {
	// Arrange: Setup configuration with plugins
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
		},
		Build: config.BuildConfig{
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

	// Act: Generate runtime preparation script
	content, err := g.GenerateRtPrepareScript()
	require.NoError(t, err, "GenerateRtPrepareScript() should not return an error")
	require.NotEmpty(t, content, "Rt prepare script should not be empty")

	t.Logf("Generated runtime preparation script: %d bytes", len(content))

	// Assert: Verify expected content
	expectedStrings := map[string]string{
		"shebang": "#!/bin/sh",
		"header":  "TCS Runtime Preparation",
	}

	for name, expected := range expectedStrings {
		assert.Contains(t, content, expected,
			"Rt prepare script should contain %s: %s", name, expected)
	}
	t.Logf("✓ Verified all %d expected sections present", len(expectedStrings))
}

func TestHealthcheckScriptGenerator_Generate(t *testing.T) {
	tests := []struct {
		name        string
		healthcheck config.HealthcheckConfig
		expectCheck string
	}{
		{
			name: "default process check when disabled",
			healthcheck: config.HealthcheckConfig{
				Enabled: false,
				Type:    "http",
			},
			expectCheck: "ps=$(ls -l /proc/*/exe",
		},
		{
			name: "default process check when type is http",
			healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "http",
				HTTP: config.HTTPHealthConfig{
					Path:    "/health",
					Port:    8080,
					Timeout: 3,
				},
			},
			expectCheck: "ps=$(ls -l /proc/*/exe",
		},
		{
			name: "default process check when type is custom but no custom_script",
			healthcheck: config.HealthcheckConfig{
				Enabled:      true,
				Type:         "custom",
				CustomScript: "",
			},
			expectCheck: "ps=$(ls -l /proc/*/exe",
		},
		{
			name: "custom script when type is custom with custom_script",
			healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "custom",
				CustomScript: `#!/bin/sh
# Custom health check
curl -f http://localhost:8080/health || exit 1`,
			},
			expectCheck: "curl -f http://localhost:8080/health",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup configuration with specific healthcheck
			cfg := &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/opt/services",
				},
				Runtime: config.RuntimeConfig{
					Healthcheck: tt.healthcheck,
				},
			}

			engine := NewTemplateEngine()
			vars := NewVariables(cfg)
			g := NewHealthcheckScriptTemplateGenerator(cfg, engine, vars)

			// Act: Generate healthcheck script
			content, err := g.Generate()
			require.NoError(t, err, "Generate() should not return an error")
			require.NotEmpty(t, content, "Healthcheck script should not be empty")

			t.Logf("Generated healthcheck script: %d bytes", len(content))

			// Assert: Verify expected check is present
			assert.Contains(t, content, tt.expectCheck,
				"Healthcheck script should contain expected check: %s", tt.expectCheck)
			t.Logf("✓ Verified healthcheck type: %s", tt.healthcheck.Type)
		})
	}
}
