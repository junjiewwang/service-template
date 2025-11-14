package context

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariablePool_GetSharedVariables(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	tests := []struct {
		name     string
		category string
		wantKeys []string
	}{
		{
			name:     "common variables",
			category: CategoryCommon,
			wantKeys: []string{VarServiceName, VarDeployDir, VarServiceRoot, VarGeneratedAt},
		},
		{
			name:     "build variables",
			category: CategoryBuild,
			wantKeys: []string{VarBuildCommand, VarPreBuildCommand, VarPostBuildCommand},
		},
		{
			name:     "runtime variables",
			category: CategoryRuntime,
			wantKeys: []string{"STARTUP_COMMAND", "ENV_VARS", "HEALTHCHECK_ENABLED"},
		},
		{
			name:     "plugin variables",
			category: CategoryPlugin,
			wantKeys: []string{VarPluginRootDir, "PLUGIN_INSTALL_DIR", "HAS_PLUGINS"},
		},
		{
			name:     "service variables",
			category: CategoryService,
			wantKeys: []string{"PORTS", VarServicePort},
		},
		{
			name:     "language variables",
			category: CategoryLanguage,
			wantKeys: []string{VarLanguage},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			shared := ctx.VariablePool.GetSharedVariables(tt.category)

			// Assert
			require.NotNil(t, shared, "Shared variables should not be nil")
			assert.True(t, shared.IsFrozen(), "Shared variables should be frozen")
			assert.Equal(t, tt.category, shared.Category(), "Category should match")

			vars := shared.ToMap()
			for _, key := range tt.wantKeys {
				assert.Contains(t, vars, key, "Should contain key: %s", key)
			}

			t.Logf("✓ Category %s has %d variables", tt.category, shared.Size())
		})
	}
}

func TestVariablePool_Caching(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	// Act: Get same category twice
	shared1 := ctx.VariablePool.GetSharedVariables(CategoryCommon)
	shared2 := ctx.VariablePool.GetSharedVariables(CategoryCommon)

	// Assert: Should return the same instance (cached)
	assert.Same(t, shared1, shared2, "Should return cached instance")
	t.Logf("✓ Variable pool caching works correctly")
}

func TestSharedVariables_Immutability(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	shared := ctx.VariablePool.GetSharedVariables(CategoryCommon)

	// Act: Get map and try to modify
	vars1 := shared.ToMap()
	vars1["NEW_KEY"] = "new_value"

	// Assert: Original should not be affected
	vars2 := shared.ToMap()
	assert.NotContains(t, vars2, "NEW_KEY", "Original should not be modified")
	t.Logf("✓ Shared variables are immutable")
}

func TestSharedVariables_Get(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	shared := ctx.VariablePool.GetSharedVariables(CategoryCommon)

	// Act & Assert
	val, exists := shared.Get(VarServiceName)
	assert.True(t, exists, "Should find existing key")
	assert.Equal(t, "test-service", val, "Should return correct value")

	val, exists = shared.Get("NON_EXISTENT_KEY")
	assert.False(t, exists, "Should not find non-existent key")
	assert.Nil(t, val, "Should return nil for non-existent key")

	t.Logf("✓ Get method works correctly")
}

// Helper function to create test config
func createTestConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/app",
			Ports: []config.PortConfig{
				{Port: 8080, Expose: true},
				{Port: 9090, Expose: false},
			},
		},
		Language: config.LanguageConfig{
			Type:   "go",
			Config: map[string]interface{}{"GO111MODULE": "on"},
		},
		Build: config.BuildConfig{
			Commands: config.BuildCommandsConfig{
				Build:     "go build -o bin/app",
				PreBuild:  "go mod download",
				PostBuild: "echo done",
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
			Dependencies: config.DependenciesConfig{
				SystemPkgs: []string{"ca-certificates"},
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/app",
				Env: []config.EnvConfig{
					{Name: "ENV", Value: "production"},
				},
			},
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "http",
			},
			SystemDependencies: config.RuntimeSystemDependenciesConfig{
				Packages: []string{"curl"},
			},
			GenerateScripts: true,
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/opt/plugins",
			Items: []config.PluginConfig{
				{
					Name:        "test-plugin",
					DownloadURL: config.NewStaticDownloadURL("https://example.com/plugin.sh"),
					RuntimeEnv: []config.EnvironmentVariable{
						{Name: "PLUGIN_HOME", Value: "/opt/plugins"},
					},
				},
			},
		},
		CI: config.CIConfig{
			ScriptDir: ".tad/build/test-service",
		},
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
	}
}
