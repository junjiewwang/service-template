package plugin

// PluginDomain 插件子域领域模型
type PluginDomain struct {
	// 所有插件共享的安装目录
	InstallDir string `yaml:"install_dir"`

	// 插件列表
	Items []Plugin `yaml:"items"`
}

// Plugin 插件配置
type Plugin struct {
	Name           string         `yaml:"name"`
	Description    string         `yaml:"description"`
	DownloadURL    interface{}    `yaml:"download_url"` // string or map[string]string
	InstallCommand string         `yaml:"install_command"`
	RuntimeEnv     []PluginEnvVar `yaml:"runtime_env,omitempty"`
	Required       bool           `yaml:"required"`
}

// PluginEnvVar 插件环境变量
type PluginEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// GetDownloadURL 获取指定架构的下载URL
func (p *Plugin) GetDownloadURL(arch string) string {
	// 如果是字符串，直接返回
	if url, ok := p.DownloadURL.(string); ok {
		return url
	}

	// 如果是map，根据架构查找
	if urlMap, ok := p.DownloadURL.(map[string]interface{}); ok {
		// 尝试精确匹配
		if url, ok := urlMap[arch].(string); ok {
			return url
		}

		// 架构别名映射
		archAliases := map[string][]string{
			"amd64":   {"x86_64", "amd64"},
			"arm64":   {"aarch64", "arm64"},
			"x86_64":  {"amd64", "x86_64"},
			"aarch64": {"arm64", "aarch64"},
		}

		// 尝试别名匹配
		if aliases, ok := archAliases[arch]; ok {
			for _, alias := range aliases {
				if url, ok := urlMap[alias].(string); ok {
					return url
				}
			}
		}

		// 尝试默认值
		if url, ok := urlMap["default"].(string); ok {
			return url
		}
	}

	return ""
}

// HasRuntimeEnv 是否有运行时环境变量
func (p *Plugin) HasRuntimeEnv() bool {
	return len(p.RuntimeEnv) > 0
}

// GetRequiredPlugins 获取必需的插件
func (d *PluginDomain) GetRequiredPlugins() []Plugin {
	var required []Plugin
	for _, plugin := range d.Items {
		if plugin.Required {
			required = append(required, plugin)
		}
	}
	return required
}

// HasPlugins 是否有插件
func (d *PluginDomain) HasPlugins() bool {
	return len(d.Items) > 0
}
