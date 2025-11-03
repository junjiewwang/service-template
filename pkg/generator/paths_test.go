package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCIPaths_DefaultValues(t *testing.T) {
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			// 不设置任何值，使用默认值
		},
	}

	paths := NewCIPaths(cfg)

	assert.Equal(t, DefaultCIScriptDir, paths.ScriptDir)
	assert.Equal(t, DefaultCIBuildConfigDir, paths.BuildConfigDir)
	assert.Equal(t, ContainerCIScriptDir, paths.ContainerScriptDir)
	assert.Equal(t, "build.sh", paths.BuildScript)
	assert.Equal(t, "build_deps_install.sh", paths.DepsInstallScript)
	assert.Equal(t, "rt_prepare.sh", paths.RtPrepareScript)
	assert.Equal(t, "entrypoint.sh", paths.EntrypointScript)
	assert.Equal(t, "healthchk.sh", paths.HealthcheckScript)
}

func TestNewCIPaths_CustomValues(t *testing.T) {
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			ScriptDir:      "custom/ci/scripts",
			BuildConfigDir: "custom/ci/config",
		},
	}

	paths := NewCIPaths(cfg)

	assert.Equal(t, "custom/ci/scripts", paths.ScriptDir)
	assert.Equal(t, "custom/ci/config", paths.BuildConfigDir)
	assert.Equal(t, ContainerCIScriptDir, paths.ContainerScriptDir) // 容器路径固定
}

func TestCIPaths_GetScriptPath(t *testing.T) {
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			ScriptDir: "ci/scripts",
		},
	}

	paths := NewCIPaths(cfg)

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
	cfg := &config.ServiceConfig{}
	paths := NewCIPaths(cfg)

	tests := []struct {
		name       string
		scriptName string
		expected   string
	}{
		{"build script", "build.sh", "/opt/bk-ci/tcs/build.sh"},
		{"deps install script", "build_deps_install.sh", "/opt/bk-ci/tcs/build_deps_install.sh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetContainerScriptPath(tt.scriptName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCIPaths_GetAllScriptPaths(t *testing.T) {
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			ScriptDir: "ci/scripts",
		},
	}

	paths := NewCIPaths(cfg)
	allPaths := paths.GetAllScriptPaths()

	require.NotNil(t, allPaths)
	assert.Len(t, allPaths, 5)

	// 验证所有脚本路径
	assert.Equal(t, "ci/scripts/build.sh", allPaths["build-script"])
	assert.Equal(t, "ci/scripts/build_deps_install.sh", allPaths["deps-install-script"])
	assert.Equal(t, "ci/scripts/rt_prepare.sh", allPaths["rt-prepare-script"])
	assert.Equal(t, "ci/scripts/entrypoint.sh", allPaths["entrypoint-script"])
	assert.Equal(t, "ci/scripts/healthchk.sh", allPaths["healthcheck-script"])
}

func TestCIPaths_ToTemplateVars(t *testing.T) {
	cfg := &config.ServiceConfig{
		CI: config.CIConfig{
			ScriptDir:      "ci/scripts",
			BuildConfigDir: "ci/config",
		},
	}

	paths := NewCIPaths(cfg)
	vars := paths.ToTemplateVars()

	require.NotNil(t, vars)

	// 验证基本路径变量
	assert.Equal(t, "ci/scripts", vars["CI_SCRIPT_DIR"])
	assert.Equal(t, "ci/config", vars["CI_BUILD_CONFIG_DIR"])
	assert.Equal(t, ContainerCIScriptDir, vars["CI_CONTAINER_DIR"])

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
	assert.Equal(t, "/opt/bk-ci/tcs/build.sh", vars["BUILD_SCRIPT_CONTAINER_PATH"])
	assert.Equal(t, "/opt/bk-ci/tcs/build_deps_install.sh", vars["DEPS_INSTALL_SCRIPT_CONTAINER_PATH"])
}

func TestCIPaths_Integration(t *testing.T) {
	// 测试完整的集成场景
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

	paths := NewCIPaths(cfg)

	// 验证可以正确生成所有路径
	allPaths := paths.GetAllScriptPaths()
	for generatorType, scriptPath := range allPaths {
		assert.NotEmpty(t, generatorType)
		assert.NotEmpty(t, scriptPath)
		assert.Contains(t, scriptPath, "devops/scripts")
	}

	// 验证模板变量可以正确生成
	vars := paths.ToTemplateVars()
	assert.NotEmpty(t, vars)
	assert.Equal(t, "devops/scripts", vars["CI_SCRIPT_DIR"])
	assert.Equal(t, "devops/config", vars["CI_BUILD_CONFIG_DIR"])
}
