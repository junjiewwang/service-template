package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCIPaths_DefaultValues(t *testing.T) {
	// Arrange: Setup configuration with no CI settings (use defaults)
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
		},
		CI: config.CIConfig{
			// 不设置任何值，使用默认值
		},
	}

	// Act: Create CI paths
	paths := context.NewCIPaths(cfg)

	// Assert: Verify default values are used
	assert.Equal(t, ".tad/build/test-service", paths.ScriptDir)
	assert.Equal(t, ".tad/build/test-service/build", paths.BuildConfigDir)
	assert.Equal(t, ".tad/build/test-service/config_template", paths.ConfigTemplateDir)
	assert.Equal(t, "/opt/.tad/build/test-service", paths.ContainerScriptDir, "Container path should be /opt + script_dir")
	assert.Equal(t, "build.sh", paths.BuildScript)
	assert.Equal(t, "build_deps_install.sh", paths.DepsInstallScript)
	assert.Equal(t, "rt_prepare.sh", paths.RtPrepareScript)
	assert.Equal(t, "entrypoint.sh", paths.EntrypointScript)
	assert.Equal(t, "healthchk.sh", paths.HealthcheckScript)

	t.Logf("✓ Default CI paths: ScriptDir=%s, BuildConfigDir=%s, ConfigTemplateDir=%s",
		paths.ScriptDir, paths.BuildConfigDir, paths.ConfigTemplateDir)
}

func TestNewCIPaths_CustomValues(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
		},
		CI: config.CIConfig{
			ScriptDir:         "custom/ci/scripts",
			BuildConfigDir:    "custom/ci/config",
			ConfigTemplateDir: "custom/ci/templates",
		},
	}

	paths := context.NewCIPaths(cfg)

	assert.Equal(t, "custom/ci/scripts", paths.ScriptDir)
	assert.Equal(t, "custom/ci/config", paths.BuildConfigDir)
	assert.Equal(t, "custom/ci/templates", paths.ConfigTemplateDir)
	assert.Equal(t, "/opt/custom/ci/scripts", paths.ContainerScriptDir, "Container path should be /opt + script_dir")
}

func TestCIPaths_GetScriptPath(t *testing.T) {
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			ScriptDir: "ci/scripts",
		},
	}

	paths := context.NewCIPaths(cfg)

	tests := []struct {
		name       string
		scriptName string
		expected   string
	}{
		{"build script", "build.sh", "ci/scripts/build.sh"},
		{"deps install script", "build_deps_install.sh", "ci/scripts/build_deps_install.sh"},
		{"rt prepare script", "rt_prepare.sh", "ci/scripts/rt_prepare.sh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetScriptPath(tt.scriptName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCIPaths_GetContainerScriptPath(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
		},
	}
	paths := context.NewCIPaths(cfg)

	tests := []struct {
		name       string
		scriptName string
		expected   string
	}{
		{"build script", "build.sh", "/opt/.tad/build/test-service/build.sh"},
		{"deps install script", "build_deps_install.sh", "/opt/.tad/build/test-service/build_deps_install.sh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetContainerScriptPath(tt.scriptName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCIPaths_GetAllScriptPaths(t *testing.T) {
	// Arrange: Setup configuration with custom script directory
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			ScriptDir: "ci/scripts",
		},
	}

	paths := context.NewCIPaths(cfg)

	// Act: Get all script paths
	allPaths := paths.GetAllScriptPaths()

	// Assert: Verify all script paths are present
	require.NotNil(t, allPaths)
	assert.Len(t, allPaths, 6, "Should have 6 script paths")

	// 验证所有脚本路径
	expectedPaths := map[string]string{
		"build-script":         "ci/scripts/build.sh",
		"deps-install-script":  "ci/scripts/build_deps_install.sh",
		"rt-prepare-script":    "ci/scripts/rt_prepare.sh",
		"entrypoint-script":    "ci/scripts/entrypoint.sh",
		"healthcheck-script":   "ci/scripts/healthchk.sh",
		"build-plugins-script": "ci/scripts/build_plugins.sh",
	}

	for key, expectedPath := range expectedPaths {
		assert.Equal(t, expectedPath, allPaths[key],
			"Script path for %s should match", key)
	}

	t.Logf("✓ Verified all %d script paths", len(allPaths))
}

func TestCIPaths_ToTemplateVars(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
		},
		CI: config.CIConfig{
			ScriptDir:         "ci/scripts",
			BuildConfigDir:    "ci/config",
			ConfigTemplateDir: "ci/templates",
		},
	}

	paths := context.NewCIPaths(cfg)
	vars := paths.ToTemplateVars()

	require.NotNil(t, vars)

	// 验证基本路径变量
	assert.Equal(t, "ci/scripts", vars["CI_SCRIPT_DIR"])
	assert.Equal(t, "ci/config", vars["CI_BUILD_CONFIG_DIR"])
	assert.Equal(t, "ci/templates", vars["CI_CONFIG_TEMPLATE_DIR"])
	assert.Equal(t, "/opt/ci/scripts", vars["CI_CONTAINER_DIR"], "Container dir should be /opt + script_dir")

	// 验证脚本文件名
	assert.Equal(t, "build.sh", vars["BUILD_SCRIPT"])
	assert.Equal(t, "build_deps_install.sh", vars["DEPS_INSTALL_SCRIPT"])
	assert.Equal(t, "rt_prepare.sh", vars["RT_PREPARE_SCRIPT"])
	assert.Equal(t, "entrypoint.sh", vars["ENTRYPOINT_SCRIPT"])
	assert.Equal(t, "healthchk.sh", vars["HEALTHCHECK_SCRIPT"])

	// 验证完整路径（主机）
	assert.Equal(t, "ci/scripts/build.sh", vars["BUILD_SCRIPT_PATH"])
	assert.Equal(t, "ci/scripts/build_deps_install.sh", vars["DEPS_INSTALL_SCRIPT_PATH"])

	// 验证完整路径（容器）
	assert.Equal(t, "/opt/ci/scripts/build.sh", vars["BUILD_SCRIPT_CONTAINER_PATH"])
	assert.Equal(t, "/opt/ci/scripts/build_deps_install.sh", vars["DEPS_INSTALL_SCRIPT_CONTAINER_PATH"])
}

func TestCIPaths_Integration(t *testing.T) {
	// Arrange: Setup complete integration scenario
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
		},
		CI: config.CIConfig{
			ScriptDir:      "devops/scripts",
			BuildConfigDir: "devops/config",
		},
	}

	paths := context.NewCIPaths(cfg)

	// Act & Assert: 验证可以正确生成所有路径
	allPaths := paths.GetAllScriptPaths()
	t.Logf("Generated %d script paths", len(allPaths))

	for generatorType, scriptPath := range allPaths {
		assert.NotEmpty(t, generatorType, "Generator type should not be empty")
		assert.NotEmpty(t, scriptPath, "Script path should not be empty")
		assert.Contains(t, scriptPath, "devops/scripts",
			"Script path should contain custom directory")
		t.Logf("✓ %s: %s", generatorType, scriptPath)
	}

	// 验证模板变量可以正确生成
	vars := paths.ToTemplateVars()
	assert.NotEmpty(t, vars, "Template vars should not be empty")
	assert.Equal(t, "devops/scripts", vars["CI_SCRIPT_DIR"])
	assert.Equal(t, "devops/config", vars["CI_BUILD_CONFIG_DIR"])

	t.Logf("✓ Template vars generated with %d keys", len(vars))
}
