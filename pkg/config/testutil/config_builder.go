package testutil

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// ConfigBuilder 配置构建器
// 使用 Builder Pattern 提供流式 API 来构建测试配置
type ConfigBuilder struct {
	cfg *config.ServiceConfig
}

// NewConfigBuilder 创建配置构建器
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &config.ServiceConfig{},
	}
}

// FromPreset 从预设配置开始构建
func FromPreset(preset *config.ServiceConfig) *ConfigBuilder {
	// 深拷贝预设配置
	cfg := *preset
	return &ConfigBuilder{cfg: &cfg}
}

// ============================================
// 基础信息配置
// ============================================

// WithService 设置服务信息
func (b *ConfigBuilder) WithService(name, description string) *ConfigBuilder {
	b.cfg.Service.Name = name
	b.cfg.Service.Description = description
	return b
}

// WithServiceName 设置服务名称
func (b *ConfigBuilder) WithServiceName(name string) *ConfigBuilder {
	b.cfg.Service.Name = name
	return b
}

// WithServiceDescription 设置服务描述
func (b *ConfigBuilder) WithServiceDescription(desc string) *ConfigBuilder {
	b.cfg.Service.Description = desc
	return b
}

// WithPort 添加端口配置
func (b *ConfigBuilder) WithPort(name string, port int, protocol string, expose bool) *ConfigBuilder {
	b.cfg.Service.Ports = append(b.cfg.Service.Ports, config.PortConfig{
		Name:     name,
		Port:     port,
		Protocol: protocol,
		Expose:   expose,
	})
	return b
}

// WithDeployDir 设置部署目录
func (b *ConfigBuilder) WithDeployDir(dir string) *ConfigBuilder {
	b.cfg.Service.DeployDir = dir
	return b
}

// ============================================
// 语言配置
// ============================================

// WithLanguage 设置语言类型
func (b *ConfigBuilder) WithLanguage(langType string) *ConfigBuilder {
	b.cfg.Language.Type = langType
	return b
}

// WithLanguageConfig 设置语言配置
func (b *ConfigBuilder) WithLanguageConfig(key, value string) *ConfigBuilder {
	if b.cfg.Language.Config == nil {
		b.cfg.Language.Config = make(map[string]interface{})
	}
	b.cfg.Language.Config[key] = value
	return b
}

// ============================================
// 基础镜像配置
// ============================================

// WithBaseImages 设置基础镜像配置
func (b *ConfigBuilder) WithBaseImages(baseImages config.BaseImagesConfig) *ConfigBuilder {
	b.cfg.BaseImages = baseImages
	return b
}

// WithBuilder 添加构建镜像预设
func (b *ConfigBuilder) WithBuilder(name, amd64, arm64 string) *ConfigBuilder {
	if b.cfg.BaseImages.Builders == nil {
		b.cfg.BaseImages.Builders = make(map[string]config.ArchImageConfig)
	}
	b.cfg.BaseImages.Builders[name] = config.ArchImageConfig{
		AMD64: amd64,
		ARM64: arm64,
	}
	return b
}

// WithRuntime 添加运行时镜像预设
func (b *ConfigBuilder) WithRuntime(name, amd64, arm64 string) *ConfigBuilder {
	if b.cfg.BaseImages.Runtimes == nil {
		b.cfg.BaseImages.Runtimes = make(map[string]config.ArchImageConfig)
	}
	b.cfg.BaseImages.Runtimes[name] = config.ArchImageConfig{
		AMD64: amd64,
		ARM64: arm64,
	}
	return b
}

// ============================================
// 构建配置
// ============================================

// WithBuilderImage 设置构建镜像引用
func (b *ConfigBuilder) WithBuilderImage(ref string) *ConfigBuilder {
	b.cfg.Build.BuilderImage = config.NewImageSpec(ref)
	return b
}

// WithRuntimeImage 设置运行时镜像引用
func (b *ConfigBuilder) WithRuntimeImage(ref string) *ConfigBuilder {
	b.cfg.Build.RuntimeImage = config.NewImageSpec(ref)
	return b
}

// WithBuildCommand 设置构建命令
func (b *ConfigBuilder) WithBuildCommand(cmd string) *ConfigBuilder {
	b.cfg.Build.Commands.Build = cmd
	return b
}

// WithPreBuildCommand 设置预构建命令
func (b *ConfigBuilder) WithPreBuildCommand(cmd string) *ConfigBuilder {
	b.cfg.Build.Commands.PreBuild = cmd
	return b
}

// WithPostBuildCommand 设置后构建命令
func (b *ConfigBuilder) WithPostBuildCommand(cmd string) *ConfigBuilder {
	b.cfg.Build.Commands.PostBuild = cmd
	return b
}

// WithSystemPackage 添加系统包依赖
func (b *ConfigBuilder) WithSystemPackage(pkg string) *ConfigBuilder {
	b.cfg.Build.Dependencies.SystemPkgs = append(b.cfg.Build.Dependencies.SystemPkgs, pkg)
	return b
}

// WithSystemPackages 设置系统包依赖列表
func (b *ConfigBuilder) WithSystemPackages(pkgs []string) *ConfigBuilder {
	b.cfg.Build.Dependencies.SystemPkgs = pkgs
	return b
}

// WithCustomPackage 添加自定义包依赖
func (b *ConfigBuilder) WithCustomPackage(pkg config.CustomPackage) *ConfigBuilder {
	b.cfg.Build.Dependencies.CustomPkgs = append(b.cfg.Build.Dependencies.CustomPkgs, pkg)
	return b
}

// WithDependencyFiles 设置依赖文件配置
func (b *ConfigBuilder) WithDependencyFiles(autoDetect bool, files []string) *ConfigBuilder {
	b.cfg.Build.DependencyFiles.AutoDetect = autoDetect
	b.cfg.Build.DependencyFiles.Files = files
	return b
}

// ============================================
// 插件配置
// ============================================

// WithPlugin 添加插件
func (b *ConfigBuilder) WithPlugin(plugin config.PluginConfig) *ConfigBuilder {
	b.cfg.Plugins.Items = append(b.cfg.Plugins.Items, plugin)
	return b
}

// WithPluginInstallDir 设置插件安装目录
func (b *ConfigBuilder) WithPluginInstallDir(dir string) *ConfigBuilder {
	b.cfg.Plugins.InstallDir = dir
	return b
}

// ============================================
// 运行时配置
// ============================================

// WithRuntimePackage 添加运行时系统包
func (b *ConfigBuilder) WithRuntimePackage(pkg string) *ConfigBuilder {
	b.cfg.Runtime.SystemDependencies.Packages = append(b.cfg.Runtime.SystemDependencies.Packages, pkg)
	return b
}

// WithRuntimePackages 设置运行时系统包列表
func (b *ConfigBuilder) WithRuntimePackages(pkgs []string) *ConfigBuilder {
	b.cfg.Runtime.SystemDependencies.Packages = pkgs
	return b
}

// WithStartupCommand 设置启动命令
func (b *ConfigBuilder) WithStartupCommand(cmd string) *ConfigBuilder {
	b.cfg.Runtime.Startup.Command = cmd
	return b
}

// WithStartupEnv 添加启动环境变量
func (b *ConfigBuilder) WithStartupEnv(name, value string) *ConfigBuilder {
	b.cfg.Runtime.Startup.Env = append(b.cfg.Runtime.Startup.Env, config.EnvConfig{
		Name:  name,
		Value: value,
	})
	return b
}

// WithHealthcheck 设置健康检查配置
func (b *ConfigBuilder) WithHealthcheck(enabled bool, hcType string) *ConfigBuilder {
	b.cfg.Runtime.Healthcheck.Enabled = enabled
	b.cfg.Runtime.Healthcheck.Type = hcType
	return b
}

// WithCustomHealthcheck 设置自定义健康检查脚本
func (b *ConfigBuilder) WithCustomHealthcheck(script string) *ConfigBuilder {
	b.cfg.Runtime.Healthcheck.Enabled = true
	b.cfg.Runtime.Healthcheck.Type = "custom"
	b.cfg.Runtime.Healthcheck.CustomScript = script
	return b
}

// ============================================
// 本地开发配置
// ============================================

// WithComposeResources 设置 Docker Compose 资源限制
func (b *ConfigBuilder) WithComposeResources(cpuLimit, memLimit, cpuReserve, memReserve string) *ConfigBuilder {
	b.cfg.LocalDev.Compose.Resources.Limits.CPUs = cpuLimit
	b.cfg.LocalDev.Compose.Resources.Limits.Memory = memLimit
	b.cfg.LocalDev.Compose.Resources.Reservations.CPUs = cpuReserve
	b.cfg.LocalDev.Compose.Resources.Reservations.Memory = memReserve
	return b
}

// WithComposeVolume 添加 Docker Compose 卷挂载
func (b *ConfigBuilder) WithComposeVolume(volume config.VolumeConfig) *ConfigBuilder {
	b.cfg.LocalDev.Compose.Volumes = append(b.cfg.LocalDev.Compose.Volumes, volume)
	return b
}

// WithComposeLabel 添加 Docker Compose 标签
func (b *ConfigBuilder) WithComposeLabel(key, value string) *ConfigBuilder {
	if b.cfg.LocalDev.Compose.Labels == nil {
		b.cfg.LocalDev.Compose.Labels = make(map[string]string)
	}
	b.cfg.LocalDev.Compose.Labels[key] = value
	return b
}

// ============================================
// CI/CD 配置
// ============================================

// WithCIScriptDir 设置 CI 脚本目录
func (b *ConfigBuilder) WithCIScriptDir(dir string) *ConfigBuilder {
	b.cfg.CI.ScriptDir = dir
	return b
}

// WithCIBuildConfigDir 设置构建配置目录
func (b *ConfigBuilder) WithCIBuildConfigDir(dir string) *ConfigBuilder {
	b.cfg.CI.BuildConfigDir = dir
	return b
}

// WithCIConfigTemplateDir 设置配置模板目录
func (b *ConfigBuilder) WithCIConfigTemplateDir(dir string) *ConfigBuilder {
	b.cfg.CI.ConfigTemplateDir = dir
	return b
}

// ============================================
// 元数据配置
// ============================================

// WithMetadata 设置元数据
func (b *ConfigBuilder) WithMetadata(version, generator string) *ConfigBuilder {
	b.cfg.Metadata.TemplateVersion = version
	b.cfg.Metadata.Generator = generator
	return b
}

// ============================================
// 构建方法
// ============================================

// Build 构建配置
func (b *ConfigBuilder) Build() *config.ServiceConfig {
	return b.cfg
}

// MustBuild 构建配置（验证失败则 panic）
func (b *ConfigBuilder) MustBuild() *config.ServiceConfig {
	validator := config.NewValidator(b.cfg)
	if err := validator.Validate(); err != nil {
		panic(err)
	}
	return b.cfg
}

// BuildWithDefaults 构建配置并填充默认值
func (b *ConfigBuilder) BuildWithDefaults() *config.ServiceConfig {
	// 填充默认值
	if b.cfg.Service.DeployDir == "" {
		b.cfg.Service.DeployDir = "/usr/local/services"
	}

	if b.cfg.Metadata.TemplateVersion == "" {
		b.cfg.Metadata.TemplateVersion = "2.0.0"
	}

	if b.cfg.Metadata.Generator == "" {
		b.cfg.Metadata.Generator = "svcgen"
	}

	return b.cfg
}
