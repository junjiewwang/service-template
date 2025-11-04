package generator

import (
	"fmt"
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Variables holds all template variables
type Variables struct {
	// Service variables
	ServiceName string
	ServicePort int
	ServiceRoot string
	DeployDir   string

	// Config variables
	ConfigDir     string
	ServiceBinDir string

	// Language variables
	Language        string
	LanguageVersion string
	LanguageConfig  map[string]string

	// Architecture variables
	GOARCH string
	GOOS   string

	// Ports
	Ports     []config.PortConfig
	PortsList string

	// Plugin variables
	PluginName        string
	PluginDescription string
	PluginDownloadURL string
	PluginInstallDir  string
	PluginRootDir     string // 新增：插件根目录 /plugins

	// CI paths
	CIPaths *CIPaths

	// All config for template access
	Config *config.ServiceConfig
}

// NewVariables creates a new Variables instance from config
func NewVariables(cfg *config.ServiceConfig) *Variables {
	vars := &Variables{
		ServiceName:     cfg.Service.Name,
		ServiceRoot:     fmt.Sprintf("%s/%s", cfg.Service.DeployDir, cfg.Service.Name),
		DeployDir:       cfg.Service.DeployDir,
		Language:        cfg.Language.Type,
		LanguageVersion: cfg.Language.Version,
		LanguageConfig:  cfg.Language.Config,
		Ports:           cfg.Service.Ports,
		Config:          cfg,
		CIPaths:         NewCIPaths(cfg), // 初始化 CI 路径管理器
	}

	// Set main service port (first port)
	if len(cfg.Service.Ports) > 0 {
		vars.ServicePort = cfg.Service.Ports[0].Port
	}

	// Build ports list
	var portsList []string
	for _, port := range cfg.Service.Ports {
		portsList = append(portsList, fmt.Sprintf("%d", port.Port))
	}
	vars.PortsList = strings.Join(portsList, ",")

	// Set directory paths
	vars.ConfigDir = fmt.Sprintf("%s/configs", vars.ServiceRoot)
	vars.ServiceBinDir = fmt.Sprintf("%s/bin", vars.ServiceRoot)

	// 新增：设置插件根目录
	vars.PluginRootDir = "/plugins"

	return vars
}

// WithArchitecture sets architecture-specific variables
func (v *Variables) WithArchitecture(arch string) *Variables {
	newVars := *v
	newVars.GOARCH = arch
	if arch == "amd64" {
		newVars.GOOS = "linux"
	} else if arch == "arm64" {
		newVars.GOOS = "linux"
	}
	return &newVars
}

// WithPlugin sets plugin-specific variables
func (v *Variables) WithPlugin(plugin config.PluginConfig) *Variables {
	newVars := *v
	newVars.PluginName = plugin.Name
	newVars.PluginDescription = plugin.Description
	newVars.PluginDownloadURL = plugin.DownloadURL
	newVars.PluginInstallDir = plugin.InstallDir
	return &newVars
}

// ToMap converts Variables to a map for template execution
func (v *Variables) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"ServiceName":       v.ServiceName,
		"ServicePort":       v.ServicePort,
		"ServiceRoot":       v.ServiceRoot,
		"DeployDir":         v.DeployDir,
		"ConfigDir":         v.ConfigDir,
		"ServiceBinDir":     v.ServiceBinDir,
		"Language":          v.Language,
		"LanguageVersion":   v.LanguageVersion,
		"LanguageConfig":    v.LanguageConfig,
		"GOARCH":            v.GOARCH,
		"GOOS":              v.GOOS,
		"Ports":             v.Ports,
		"PortsList":         v.PortsList,
		"PluginName":        v.PluginName,
		"PluginDescription": v.PluginDescription,
		"PluginDownloadURL": v.PluginDownloadURL,
		"PluginInstallDir":  v.PluginInstallDir,
		"PluginRootDir":     v.PluginRootDir,
		"Config":            v.Config,

		// Convenience functions
		"SERVICE_NAME":        v.ServiceName,
		"SERVICE_PORT":        v.ServicePort,
		"SERVICE_ROOT":        v.ServiceRoot,
		"DEPLOY_DIR":          v.DeployDir,
		"CONFIG_DIR":          v.ConfigDir,
		"SERVICE_BIN_DIR":     v.ServiceBinDir,
		"PLUGIN_NAME":         v.PluginName,
		"PLUGIN_DESCRIPTION":  v.PluginDescription,
		"PLUGIN_DOWNLOAD_URL": v.PluginDownloadURL,
		"PLUGIN_INSTALL_DIR":  v.PluginInstallDir,
		"PLUGIN_ROOT_DIR":     v.PluginRootDir,
	}

	// 合并 CI 路径变量
	if v.CIPaths != nil {
		for k, val := range v.CIPaths.ToTemplateVars() {
			result[k] = val
		}
	}

	return result
}
