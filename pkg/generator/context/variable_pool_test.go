package context

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
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

func TestVariablePool_EmptyPorts(t *testing.T) {
	// Arrange: 创建一个没有端口的配置
	cfg := testutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithDeployDir("/app").
		WithLanguage("go").
		// 不添加任何端口
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		Build()

	ctx := NewGeneratorContext(cfg, "/tmp/output")

	// Act: 获取 service 变量
	shared := ctx.VariablePool.GetSharedVariables(CategoryService)

	// Assert
	vars := shared.ToMap()

	// PORTS 应该是空数组
	ports, exists := vars["PORTS"]
	assert.True(t, exists, "PORTS should exist")
	assert.Empty(t, ports, "PORTS should be empty")

	// SERVICE_PORT 应该是空字符串
	servicePort, exists := vars[VarServicePort]
	assert.True(t, exists, "SERVICE_PORT should exist")
	assert.Equal(t, "", servicePort, "SERVICE_PORT should be empty string when no ports configured")

	// EXPOSE_PORTS 应该是空数组
	exposePorts, exists := vars["EXPOSE_PORTS"]
	assert.True(t, exists, "EXPOSE_PORTS should exist")
	assert.Empty(t, exposePorts, "EXPOSE_PORTS should be empty")

	t.Logf("✓ Empty ports handling works correctly")
}

// Helper function to create test config
func createTestConfig() *testutil.ServiceConfig {
	return testutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithDeployDir("/app").
		WithLanguage("go").
		WithPort("http", 8080, "TCP", true).
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPlugin(config.PluginConfig{
			Name:        "test-plugin",
			DownloadURL: config.NewStaticDownloadURL("https://example.com/plugin.sh"),
		}).
		Build()
}
