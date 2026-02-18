package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// ConfigOption 配置选项函数类型
// 使用 Options Pattern 提供灵活的配置修改方式
type ConfigOption func(*config.ServiceConfig)

// ============================================
// 服务配置选项
// ============================================

// WithServiceOpt 设置服务信息
func WithServiceOpt(name, description string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Service.Name = name
		cfg.Service.Description = description
	}
}

// WithServiceNameOpt 设置服务名称
func WithServiceNameOpt(name string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Service.Name = name
	}
}

// WithPortOpt 添加端口配置
func WithPortOpt(name string, port int, protocol string, expose bool) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Service.Ports = append(cfg.Service.Ports, config.PortConfig{
			Name:     name,
			Port:     port,
			Protocol: protocol,
			Expose:   expose,
		})
	}
}

// ============================================
// 语言配置选项
// ============================================

// WithLanguageOpt 设置语言类型
func WithLanguageOpt(langType string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Language.Type = langType
	}
}

// WithLanguageConfigOpt 设置语言配置
func WithLanguageConfigOpt(key, value string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		if cfg.Language.Config == nil {
			cfg.Language.Config = make(map[string]interface{})
		}
		cfg.Language.Config[key] = value
	}
}

// ============================================
// 镜像配置选项
// ============================================

// WithBuilderOpt 添加构建镜像预设
func WithBuilderOpt(name, amd64, arm64 string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		if cfg.BaseImages.Builders == nil {
			cfg.BaseImages.Builders = make(map[string]config.ArchImageConfig)
		}
		cfg.BaseImages.Builders[name] = config.ArchImageConfig{
			AMD64: amd64,
			ARM64: arm64,
		}
	}
}

// WithRuntimeOpt 添加运行时镜像预设
func WithRuntimeOpt(name, amd64, arm64 string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		if cfg.BaseImages.Runtimes == nil {
			cfg.BaseImages.Runtimes = make(map[string]config.ArchImageConfig)
		}
		cfg.BaseImages.Runtimes[name] = config.ArchImageConfig{
			AMD64: amd64,
			ARM64: arm64,
		}
	}
}

// WithBuilderImageOpt 设置构建镜像引用
func WithBuilderImageOpt(ref string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.BuilderImage = config.NewImageSpec(ref)
	}
}

// WithRuntimeImageOpt 设置运行时镜像引用
func WithRuntimeImageOpt(ref string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.RuntimeImage = config.NewImageSpec(ref)
	}
}

// ============================================
// 构建配置选项
// ============================================

// WithBuildCommandOpt 设置构建命令
func WithBuildCommandOpt(cmd string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.Commands.Build = cmd
	}
}

// WithPreBuildCommandOpt 设置预构建命令
func WithPreBuildCommandOpt(cmd string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.Commands.PreBuild = cmd
	}
}

// WithSystemPackagesOpt 设置系统包依赖
func WithSystemPackagesOpt(pkgs []string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.Dependencies.SystemPkgs = pkgs
	}
}

// ============================================
// 插件配置选项
// ============================================

// WithPluginOpt 添加插件
func WithPluginOpt(plugin config.PluginConfig) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Plugins.Items = append(cfg.Plugins.Items, plugin)
	}
}

// WithPluginInstallDirOpt 设置插件安装目录
func WithPluginInstallDirOpt(dir string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Plugins.InstallDir = dir
	}
}

// ============================================
// 运行时配置选项
// ============================================

// WithStartupCommandOpt 设置启动命令
func WithStartupCommandOpt(cmd string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Runtime.Startup.Command = cmd
	}
}

// WithHealthcheckOpt 设置健康检查配置
func WithHealthcheckOpt(enabled bool, hcType string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Runtime.Healthcheck.Enabled = enabled
		cfg.Runtime.Healthcheck.Type = hcType
	}
}

// WithCustomHealthcheckOpt 设置自定义健康检查脚本
func WithCustomHealthcheckOpt(script string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Runtime.Healthcheck.Enabled = true
		cfg.Runtime.Healthcheck.Type = "custom"
		cfg.Runtime.Healthcheck.CustomScript = script
	}
}

// ============================================
// 辅助函数
// ============================================

// ApplyOptions 应用配置选项到配置对象
func ApplyOptions(cfg *config.ServiceConfig, opts ...ConfigOption) *config.ServiceConfig {
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// NewConfigWithOptions 从预设配置创建新配置并应用选项
func NewConfigWithOptions(preset *config.ServiceConfig, opts ...ConfigOption) *config.ServiceConfig {
	// 深拷贝预设配置
	cfg := *preset
	return ApplyOptions(&cfg, opts...)
}
