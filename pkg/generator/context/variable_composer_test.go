package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableComposer_WithCategories(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	composer := ctx.GetVariableComposer()

	// Act
	vars := composer.
		WithCommon().
		WithBuild().
		WithRuntime().
		Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, VarBuildCommand, "Should contain build variables")
	assert.Contains(t, vars, "STARTUP_COMMAND", "Should contain runtime variables")

	t.Logf("✓ Composed %d variables from multiple categories", len(vars))
}

func TestVariableComposer_WithAll(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	composer := ctx.GetVariableComposer()

	// Act
	vars := composer.WithAll().Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, VarBuildCommand, "Should contain build variables")
	assert.Contains(t, vars, "STARTUP_COMMAND", "Should contain runtime variables")
	assert.Contains(t, vars, VarPluginRootDir, "Should contain plugin variables")
	assert.Contains(t, vars, "PORTS", "Should contain service variables")
	assert.Contains(t, vars, VarLanguage, "Should contain language variables")

	t.Logf("✓ WithAll() composed %d variables", len(vars))
}

func TestVariableComposer_WithArchitecture(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	tests := []struct {
		name string
		arch string
	}{
		{name: "amd64 architecture", arch: "amd64"},
		{name: "arm64 architecture", arch: "arm64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			vars := ctx.GetVariableComposer().
				WithBuild().
				WithArchitecture(tt.arch).
				Build()

			// Assert
			assert.Equal(t, tt.arch, vars[VarGOARCH], "Should set GOARCH")
			assert.Equal(t, tt.arch, vars["ARCH"], "Should set ARCH")
			assert.Equal(t, "linux", vars[VarGOOS], "Should set GOOS to linux")
			assert.Contains(t, vars, "BUILDER_IMAGE", "Should contain BUILDER_IMAGE")
			assert.Contains(t, vars, "RUNTIME_IMAGE", "Should contain RUNTIME_IMAGE")

			t.Logf("✓ Architecture %s variables set correctly", tt.arch)
		})
	}
}

func TestVariableComposer_WithCustom(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	// Act
	vars := ctx.GetVariableComposer().
		WithCommon().
		WithCustom("CUSTOM_KEY", "custom_value").
		WithCustom("ANOTHER_KEY", 123).
		Build()

	// Assert
	assert.Equal(t, "custom_value", vars["CUSTOM_KEY"], "Should contain custom string")
	assert.Equal(t, 123, vars["ANOTHER_KEY"], "Should contain custom int")

	t.Logf("✓ Custom variables added successfully")
}

func TestVariableComposer_WithCustomMap(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	customVars := map[string]interface{}{
		"KEY1": "value1",
		"KEY2": 42,
		"KEY3": true,
	}

	// Act
	vars := ctx.GetVariableComposer().
		WithCommon().
		WithCustomMap(customVars).
		Build()

	// Assert
	assert.Equal(t, "value1", vars["KEY1"])
	assert.Equal(t, 42, vars["KEY2"])
	assert.Equal(t, true, vars["KEY3"])

	t.Logf("✓ Custom map variables added successfully")
}

func TestVariableComposer_Override(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	// Act
	vars := ctx.GetVariableComposer().
		WithCommon().
		Override(VarServiceName, "overridden-service").
		Build()

	// Assert
	assert.Equal(t, "overridden-service", vars[VarServiceName], "Should override existing variable")

	t.Logf("✓ Variable override works correctly")
}

func TestVariableComposer_HasAndGet(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	composer := ctx.GetVariableComposer().WithCommon()

	// Act & Assert
	assert.True(t, composer.Has(VarServiceName), "Should have SERVICE_NAME")
	assert.False(t, composer.Has("NON_EXISTENT"), "Should not have non-existent key")

	val, exists := composer.Get(VarServiceName)
	assert.True(t, exists, "Should get existing variable")
	assert.Equal(t, "test-service", val, "Should return correct value")

	val, exists = composer.Get("NON_EXISTENT")
	assert.False(t, exists, "Should not get non-existent variable")
	assert.Nil(t, val, "Should return nil for non-existent variable")

	t.Logf("✓ Has() and Get() methods work correctly")
}

func TestVariableComposer_Clone(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	original := ctx.GetVariableComposer().WithCommon()

	// Act
	cloned := original.Clone()
	cloned.WithCustom("NEW_KEY", "new_value")

	// Assert
	assert.False(t, original.Has("NEW_KEY"), "Original should not have new key")
	assert.True(t, cloned.Has("NEW_KEY"), "Clone should have new key")

	t.Logf("✓ Clone creates independent copy")
}

func TestVariableComposer_Size(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	// Act
	composer := ctx.GetVariableComposer().WithCommon()
	size := composer.Size()

	// Assert
	assert.Greater(t, size, 0, "Should have variables")
	t.Logf("✓ Composer has %d variables", size)
}

func TestVariableComposer_NoOverwriteOnMerge(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")

	// Act: Add custom variable first, then merge common
	vars := ctx.GetVariableComposer().
		WithCustom(VarServiceName, "custom-name").
		WithCommon(). // This should not overwrite the custom value
		Build()

	// Assert
	assert.Equal(t, "custom-name", vars[VarServiceName], "Should not overwrite existing variable")

	t.Logf("✓ Merge does not overwrite existing variables")
}

func TestVariablePreset_ForDockerfile(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	preset := ctx.GetVariablePreset()

	// Act
	vars := preset.ForDockerfile("amd64").Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, VarBuildCommand, "Should contain build variables")
	assert.Contains(t, vars, "STARTUP_COMMAND", "Should contain runtime variables")
	assert.Contains(t, vars, VarPluginRootDir, "Should contain plugin variables")
	assert.Contains(t, vars, "PORTS", "Should contain service variables")
	assert.Contains(t, vars, VarLanguage, "Should contain language variables")
	assert.Equal(t, "amd64", vars[VarGOARCH], "Should set architecture")

	t.Logf("✓ Dockerfile preset contains %d variables", len(vars))
}

func TestVariablePreset_ForBuildScript(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	preset := ctx.GetVariablePreset()

	// Act
	vars := preset.ForBuildScript().Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, VarBuildCommand, "Should contain build variables")
	assert.Contains(t, vars, VarPluginRootDir, "Should contain plugin variables")
	assert.Contains(t, vars, "CI_SCRIPT_DIR", "Should contain CI path variables")

	t.Logf("✓ Build script preset contains %d variables", len(vars))
}

func TestVariablePreset_ForCompose(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	preset := ctx.GetVariablePreset()

	// Act
	vars := preset.ForCompose().Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, "STARTUP_COMMAND", "Should contain runtime variables")
	assert.Contains(t, vars, "PORTS", "Should contain service variables")

	t.Logf("✓ Compose preset contains %d variables", len(vars))
}

func TestVariablePreset_ForMakefile(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	preset := ctx.GetVariablePreset()

	// Act
	vars := preset.ForMakefile().Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, "PORTS", "Should contain service variables")
	assert.Contains(t, vars, "CI_SCRIPT_DIR", "Should contain CI path variables")

	t.Logf("✓ Makefile preset contains %d variables", len(vars))
}

func TestVariablePreset_ForDevOps(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	preset := ctx.GetVariablePreset()

	// Act
	vars := preset.ForDevOps().Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, VarBuildCommand, "Should contain build variables")
	assert.Contains(t, vars, VarLanguage, "Should contain language variables")

	t.Logf("✓ DevOps preset contains %d variables", len(vars))
}

func TestVariablePreset_ForScript(t *testing.T) {
	// Arrange
	cfg := createTestConfig()
	ctx := NewGeneratorContext(cfg, "/tmp/output")
	preset := ctx.GetVariablePreset()

	// Act
	vars := preset.ForScript().Build()

	// Assert
	assert.Contains(t, vars, VarServiceName, "Should contain common variables")
	assert.Contains(t, vars, "STARTUP_COMMAND", "Should contain runtime variables")
	assert.Contains(t, vars, "PORTS", "Should contain service variables")
	assert.Contains(t, vars, "CI_SCRIPT_DIR", "Should contain CI path variables")

	t.Logf("✓ Script preset contains %d variables", len(vars))
}
