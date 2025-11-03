package generator

import (
	"path/filepath"

	"github.com/junjiewwang/service-template/pkg/config"
)

const (
	// 默认 CI 路径配置
	DefaultCIScriptDir      = "bk-ci/tcs"
	DefaultCIBuildConfigDir = "bk-ci/tcs/build"

	// 容器内路径（固定）
	ContainerCIScriptDir = "/opt/bk-ci/tcs"
)

// CIPaths 管理所有 CI 相关路径
type CIPaths struct {
	// 主机路径（相对于项目根目录）
	ScriptDir      string // bk-ci/tcs
	BuildConfigDir string // bk-ci/tcs/build

	// 容器内路径（绝对路径）
	ContainerScriptDir string // /opt/bk-ci/tcs

	// 脚本文件名
	BuildScript       string // build.sh
	DepsInstallScript string // build_deps_install.sh
	RtPrepareScript   string // rt_prepare.sh
	EntrypointScript  string // entrypoint.sh
	HealthcheckScript string // healthchk.sh
}

// NewCIPaths 创建 CI 路径管理器
func NewCIPaths(cfg *config.ServiceConfig) *CIPaths {
	scriptDir := DefaultCIScriptDir
	buildConfigDir := DefaultCIBuildConfigDir

	// 如果配置了自定义路径，使用自定义路径
	if cfg.CI.ScriptDir != "" {
		scriptDir = cfg.CI.ScriptDir
	}
	if cfg.CI.BuildConfigDir != "" {
		buildConfigDir = cfg.CI.BuildConfigDir
	}

	return &CIPaths{
		ScriptDir:          scriptDir,
		BuildConfigDir:     buildConfigDir,
		ContainerScriptDir: ContainerCIScriptDir,
		BuildScript:        "build.sh",
		DepsInstallScript:  "build_deps_install.sh",
		RtPrepareScript:    "rt_prepare.sh",
		EntrypointScript:   "entrypoint.sh",
		HealthcheckScript:  "healthchk.sh",
	}
}

// GetScriptPath 获取脚本的完整路径（主机）
func (p *CIPaths) GetScriptPath(scriptName string) string {
	return filepath.Join(p.ScriptDir, scriptName)
}

// GetContainerScriptPath 获取脚本的容器内路径
func (p *CIPaths) GetContainerScriptPath(scriptName string) string {
	return filepath.Join(p.ContainerScriptDir, scriptName)
}

// GetAllScriptPaths 获取所有脚本路径（用于生成）
func (p *CIPaths) GetAllScriptPaths() map[string]string {
	return map[string]string{
		"build-script":        p.GetScriptPath(p.BuildScript),
		"deps-install-script": p.GetScriptPath(p.DepsInstallScript),
		"rt-prepare-script":   p.GetScriptPath(p.RtPrepareScript),
		"entrypoint-script":   p.GetScriptPath(p.EntrypointScript),
		"healthcheck-script":  p.GetScriptPath(p.HealthcheckScript),
	}
}

// ToTemplateVars 转换为模板变量
func (p *CIPaths) ToTemplateVars() map[string]interface{} {
	return map[string]interface{}{
		"CI_SCRIPT_DIR":       p.ScriptDir,
		"CI_BUILD_CONFIG_DIR": p.BuildConfigDir,
		"CI_CONTAINER_DIR":    p.ContainerScriptDir,

		// 脚本文件名
		"BUILD_SCRIPT":        p.BuildScript,
		"DEPS_INSTALL_SCRIPT": p.DepsInstallScript,
		"RT_PREPARE_SCRIPT":   p.RtPrepareScript,
		"ENTRYPOINT_SCRIPT":   p.EntrypointScript,
		"HEALTHCHECK_SCRIPT":  p.HealthcheckScript,

		// 完整路径（主机）
		"BUILD_SCRIPT_PATH":        p.GetScriptPath(p.BuildScript),
		"DEPS_INSTALL_SCRIPT_PATH": p.GetScriptPath(p.DepsInstallScript),
		"RT_PREPARE_SCRIPT_PATH":   p.GetScriptPath(p.RtPrepareScript),
		"ENTRYPOINT_SCRIPT_PATH":   p.GetScriptPath(p.EntrypointScript),
		"HEALTHCHECK_SCRIPT_PATH":  p.GetScriptPath(p.HealthcheckScript),

		// 完整路径（容器）
		"BUILD_SCRIPT_CONTAINER_PATH":        p.GetContainerScriptPath(p.BuildScript),
		"DEPS_INSTALL_SCRIPT_CONTAINER_PATH": p.GetContainerScriptPath(p.DepsInstallScript),
		"RT_PREPARE_SCRIPT_CONTAINER_PATH":   p.GetContainerScriptPath(p.RtPrepareScript),
		"ENTRYPOINT_SCRIPT_CONTAINER_PATH":   p.GetContainerScriptPath(p.EntrypointScript),
		"HEALTHCHECK_SCRIPT_CONTAINER_PATH":  p.GetContainerScriptPath(p.HealthcheckScript),
	}
}
