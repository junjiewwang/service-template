package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// ConfigBuilder 主配置构建器
type ConfigBuilder struct {
	config *config.ServiceConfig
}

// NewConfigBuilder 创建一个空的配置构建器
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &config.ServiceConfig{},
	}
}

// NewMinimalConfig 创建最小可用配置
// 适用于不关心配置细节的测试
func NewMinimalConfig() *ConfigBuilder {
	return &ConfigBuilder{
		config: &config.ServiceConfig{
			BaseImages: MinimalBaseImages.ToConfig(),
			Service: config.ServiceInfo{
				Name:      "test-service",
				DeployDir: "/app",
			},
			Build: config.BuildConfig{
				BuilderImage: "@builders.builder",
				RuntimeImage: "@runtimes.runtime",
				Commands: config.BuildCommandsConfig{
					Build: "echo build",
				},
				DependencyFiles: config.DependencyFilesConfig{
					AutoDetect: true,
				},
			},
			Language: config.LanguageConfig{
				Type: "go",
			},
			Runtime: config.RuntimeConfig{
				Startup: config.StartupConfig{
					Command: "./app",
				},
			},
			Plugins: config.PluginsConfig{
				InstallDir: "/tce",
				Items:      []config.PluginConfig{},
			},
		},
	}
}

// NewDefaultConfig 创建默认完整配置
// 适用于需要完整配置的测试
func NewDefaultConfig() *ConfigBuilder {
	return &ConfigBuilder{
		config: &config.ServiceConfig{
			BaseImages: DefaultBaseImages.ToConfig(),
			Service: config.ServiceInfo{
				Name:        "test-service",
				Description: "Test service",
				DeployDir:   "/opt/services",
				Ports: []config.PortConfig{
					{Port: 8080, Protocol: "http", Name: "http"},
				},
			},
			Build: config.BuildConfig{
				BuilderImage: "@builders.test_builder",
				RuntimeImage: "@runtimes.test_runtime",
				Commands: config.BuildCommandsConfig{
					Build: "go build -o bin/app",
				},
				DependencyFiles: config.DependencyFilesConfig{
					AutoDetect: true,
				},
			},
			Language: config.LanguageConfig{
				Type: "go",
			},
			Runtime: config.RuntimeConfig{
				Startup: config.StartupConfig{
					Command: "./bin/app",
				},
			},
			Plugins: config.PluginsConfig{
				InstallDir: "/tce",
				Items:      []config.PluginConfig{},
			},
			Makefile: config.MakefileConfig{
				CustomTargets: []config.CustomTarget{},
			},
			Metadata: config.MetadataConfig{
				Generator: "svcgen-test",
			},
		},
	}
}

// WithService 配置服务信息
func (b *ConfigBuilder) WithService(fn func(*ServiceBuilder)) *ConfigBuilder {
	sb := &ServiceBuilder{info: &b.config.Service}
	fn(sb)
	return b
}

// WithBuild 配置构建信息
func (b *ConfigBuilder) WithBuild(fn func(*BuildConfigBuilder)) *ConfigBuilder {
	bb := &BuildConfigBuilder{build: &b.config.Build}
	fn(bb)
	return b
}

// WithPlugins 配置插件信息
func (b *ConfigBuilder) WithPlugins(fn func(*PluginBuilder)) *ConfigBuilder {
	pb := &PluginBuilder{plugins: &b.config.Plugins}
	fn(pb)
	return b
}

// WithLanguage 设置语言类型
func (b *ConfigBuilder) WithLanguage(langType string) *ConfigBuilder {
	b.config.Language.Type = langType
	return b
}

// WithBaseImages 设置基础镜像配置
func (b *ConfigBuilder) WithBaseImages(images config.BaseImagesConfig) *ConfigBuilder {
	b.config.BaseImages = images
	return b
}

// WithBaseImagesPreset 使用预设的基础镜像
func (b *ConfigBuilder) WithBaseImagesPreset(preset BaseImagesPreset) *ConfigBuilder {
	b.config.BaseImages = preset.ToConfig()
	return b
}

// Build 构建配置对象
func (b *ConfigBuilder) Build() *config.ServiceConfig {
	return b.config
}

// MustBuild 构建配置对象（如果验证失败则 panic）
func (b *ConfigBuilder) MustBuild() *config.ServiceConfig {
	cfg := b.config
	validator := config.NewValidator(cfg)
	if err := validator.Validate(); err != nil {
		panic(err)
	}
	return cfg
}
