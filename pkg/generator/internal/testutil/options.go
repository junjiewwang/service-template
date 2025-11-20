package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// ConfigOption 配置选项函数
type ConfigOption func(*config.ServiceConfig)

// WithServiceName 设置服务名称
func WithServiceName(name string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Service.Name = name
	}
}

// WithDeployDir 设置部署目录
func WithDeployDir(dir string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Service.DeployDir = dir
	}
}

// WithPort 添加端口配置
func WithPort(port int, protocol string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Service.Ports = append(cfg.Service.Ports, config.PortConfig{
			Port:     port,
			Protocol: protocol,
			Name:     protocol,
		})
	}
}

// WithLanguageType 设置语言类型
func WithLanguageType(langType string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Language.Type = langType
	}
}

// WithPlugin 添加插件
func WithPlugin(name, url string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Plugins.Items = append(cfg.Plugins.Items, config.PluginConfig{
			Name:           name,
			Description:    name,
			DownloadURL:    config.NewStaticDownloadURL(url),
			InstallCommand: "echo 'Installing...'",
			Required:       true,
		})
	}
}

// WithBuildCommand 设置构建命令
func WithBuildCommand(cmd string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.Commands.Build = cmd
	}
}

// WithBuilderImageRef 设置构建镜像引用
func WithBuilderImageRef(ref string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.BuilderImage = config.ImageRef(ref)
	}
}

// WithRuntimeImageRef 设置运行时镜像引用
func WithRuntimeImageRef(ref string) ConfigOption {
	return func(cfg *config.ServiceConfig) {
		cfg.Build.RuntimeImage = config.ImageRef(ref)
	}
}

// NewConfigWithOptions 使用选项创建配置
func NewConfigWithOptions(opts ...ConfigOption) *config.ServiceConfig {
	cfg := NewMinimalConfig().Build()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
