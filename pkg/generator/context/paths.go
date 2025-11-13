package context

import (
	"fmt"
	"path/filepath"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Paths manages all path-related information
type Paths struct {
	// CI paths
	CI *CIPaths

	// Service paths
	ServiceRoot string
	DeployDir   string
	ConfigDir   string
	BinDir      string
}

// NewPaths creates a new Paths instance
func NewPaths(cfg *config.ServiceConfig) *Paths {
	serviceName := cfg.Service.Name
	deployDir := cfg.Service.DeployDir
	serviceRoot := fmt.Sprintf("%s/%s", deployDir, serviceName)

	return &Paths{
		CI:          NewCIPaths(cfg),
		ServiceRoot: serviceRoot,
		DeployDir:   deployDir,
		ConfigDir:   fmt.Sprintf("%s/%s", serviceRoot, ConfigDirName),
		BinDir:      fmt.Sprintf("%s/%s", serviceRoot, BinDirName),
	}
}

// ToTemplateVars converts paths to template variables
func (p *Paths) ToTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Add CI paths
	if p.CI != nil {
		for k, v := range p.CI.ToTemplateVars() {
			vars[k] = v
		}
	}

	// Add service paths
	vars["SERVICE_ROOT"] = p.ServiceRoot
	vars["DEPLOY_DIR"] = p.DeployDir
	vars["CONFIG_DIR"] = p.ConfigDir
	vars["BIN_DIR"] = p.BinDir

	return vars
}

// CIPaths 管理所有 CI 相关路径
type CIPaths struct {
	// 主机路径（相对于项目根目录）
	ScriptDir         string // .tad/build/{service-name}
	BuildConfigDir    string // {script_dir}/build
	ConfigTemplateDir string // {script_dir}/config_template

	// 容器内路径（绝对路径）
	ContainerScriptDir string // /opt/.tad/build/{service-name}

	// 脚本文件名
	BuildScript        string // build.sh
	DepsInstallScript  string // build_deps_install.sh
	RtPrepareScript    string // rt_prepare.sh
	EntrypointScript   string // entrypoint.sh
	HealthcheckScript  string // healthchk.sh
	BuildPluginsScript string // build_plugins.sh
}

// NewCIPaths 创建 CI 路径管理器
func NewCIPaths(cfg *config.ServiceConfig) *CIPaths {
	serviceName := cfg.Service.Name

	// 计算默认路径
	defaultScriptDir := fmt.Sprintf(DefaultCIScriptDirPattern, serviceName)

	scriptDir := defaultScriptDir
	buildConfigDir := ""
	configTemplateDir := ""

	// 如果配置了自定义 script_dir，使用自定义路径
	if cfg.CI.ScriptDir != "" {
		scriptDir = cfg.CI.ScriptDir
	}

	// 计算 build_config_dir 默认值
	if cfg.CI.BuildConfigDir != "" {
		buildConfigDir = cfg.CI.BuildConfigDir
	} else {
		buildConfigDir = fmt.Sprintf(DefaultCIBuildConfigDirPattern, scriptDir)
	}

	// 计算 config_template_dir 默认值
	if cfg.CI.ConfigTemplateDir != "" {
		configTemplateDir = cfg.CI.ConfigTemplateDir
	} else {
		configTemplateDir = fmt.Sprintf(DefaultConfigTemplateDirPattern, scriptDir)
	}

	// 容器内路径：基于主机路径动态计算
	// 因为 Dockerfile 中 COPY . /opt/ 会保持目录结构
	// 所以容器内路径 = /opt + 主机相对路径
	containerScriptDir := filepath.Join(ContainerProjectRoot, scriptDir)

	return &CIPaths{
		ScriptDir:          scriptDir,
		BuildConfigDir:     buildConfigDir,
		ConfigTemplateDir:  configTemplateDir,
		ContainerScriptDir: containerScriptDir, // 动态计算的容器内路径
		BuildScript:        BuildScriptName,
		DepsInstallScript:  DepsInstallScriptName,
		RtPrepareScript:    RtPrepareScriptName,
		EntrypointScript:   EntrypointScriptName,
		HealthcheckScript:  HealthcheckScriptName,
		BuildPluginsScript: BuildPluginsScriptName,
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
		"build-script":         p.GetScriptPath(p.BuildScript),
		"deps-install-script":  p.GetScriptPath(p.DepsInstallScript),
		"rt-prepare-script":    p.GetScriptPath(p.RtPrepareScript),
		"entrypoint-script":    p.GetScriptPath(p.EntrypointScript),
		"healthcheck-script":   p.GetScriptPath(p.HealthcheckScript),
		"build-plugins-script": p.GetScriptPath(p.BuildPluginsScript),
	}
}

// ToTemplateVars 转换为模板变量
func (p *CIPaths) ToTemplateVars() map[string]interface{} {
	return map[string]interface{}{
		"CI_SCRIPT_DIR":          p.ScriptDir,
		"CI_BUILD_CONFIG_DIR":    p.BuildConfigDir,
		"CI_CONFIG_TEMPLATE_DIR": p.ConfigTemplateDir,
		"CI_CONTAINER_DIR":       p.ContainerScriptDir,

		// 脚本文件名
		"BUILD_SCRIPT":         p.BuildScript,
		"DEPS_INSTALL_SCRIPT":  p.DepsInstallScript,
		"RT_PREPARE_SCRIPT":    p.RtPrepareScript,
		"ENTRYPOINT_SCRIPT":    p.EntrypointScript,
		"HEALTHCHECK_SCRIPT":   p.HealthcheckScript,
		"BUILD_PLUGINS_SCRIPT": p.BuildPluginsScript,

		// 完整路径（主机）
		"BUILD_SCRIPT_PATH":         p.GetScriptPath(p.BuildScript),
		"DEPS_INSTALL_SCRIPT_PATH":  p.GetScriptPath(p.DepsInstallScript),
		"RT_PREPARE_SCRIPT_PATH":    p.GetScriptPath(p.RtPrepareScript),
		"ENTRYPOINT_SCRIPT_PATH":    p.GetScriptPath(p.EntrypointScript),
		"HEALTHCHECK_SCRIPT_PATH":   p.GetScriptPath(p.HealthcheckScript),
		"BUILD_PLUGINS_SCRIPT_PATH": p.GetScriptPath(p.BuildPluginsScript),

		// 完整路径（容器）
		"BUILD_SCRIPT_CONTAINER_PATH":         p.GetContainerScriptPath(p.BuildScript),
		"DEPS_INSTALL_SCRIPT_CONTAINER_PATH":  p.GetContainerScriptPath(p.DepsInstallScript),
		"RT_PREPARE_SCRIPT_CONTAINER_PATH":    p.GetContainerScriptPath(p.RtPrepareScript),
		"ENTRYPOINT_SCRIPT_CONTAINER_PATH":    p.GetContainerScriptPath(p.EntrypointScript),
		"HEALTHCHECK_SCRIPT_CONTAINER_PATH":   p.GetContainerScriptPath(p.HealthcheckScript),
		"BUILD_PLUGINS_SCRIPT_CONTAINER_PATH": p.GetContainerScriptPath(p.BuildPluginsScript),
	}
}
