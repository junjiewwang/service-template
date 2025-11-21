package testutil

import (
	"github.com/junjiewwang/service-template/pkg/config"
	configtestutil "github.com/junjiewwang/service-template/pkg/config/testutil"
)

// ============================================
// 导出核心层API (Re-export Core APIs)
// ============================================
// 这个包作为适配层，为生成器测试提供便捷的API
// 所有核心功能都委托给 pkg/config/testutil

// 类型别名 - 导出核心层的类型
type (
	// ConfigBuilder 配置构建器（来自核心层）
	ConfigBuilder = configtestutil.ConfigBuilder

	// ConfigOption 配置选项函数（来自核心层）
	ConfigOption = configtestutil.ConfigOption
)

// ============================================
// 构建器函数 - 导出核心层的构建器
// ============================================

var (
	// NewConfigBuilder 创建新的配置构建器
	// 推荐使用此方法创建配置
	NewConfigBuilder = configtestutil.NewConfigBuilder

	// FromPreset 从预设配置开始构建
	FromPreset = configtestutil.FromPreset

	// NewConfigWithOptions 使用选项模式创建配置
	NewConfigWithOptions = configtestutil.NewConfigWithOptions

	// ApplyOptions 应用选项到现有配置
	ApplyOptions = configtestutil.ApplyOptions
)

// ============================================
// 预设配置 - 导出核心层的预设
// ============================================

var (
	// MinimalConfig 最小化配置（仅包含必需字段）
	MinimalConfig = configtestutil.MinimalConfig

	// GoServiceConfig Go 服务的标准配置
	GoServiceConfig = configtestutil.GoServiceConfig

	// PythonServiceConfig Python 服务的标准配置
	PythonServiceConfig = configtestutil.PythonServiceConfig

	// JavaServiceConfig Java 服务的标准配置
	JavaServiceConfig = configtestutil.JavaServiceConfig

	// ConfigWithPlugins 带插件的配置
	ConfigWithPlugins = configtestutil.ConfigWithPlugins

	// ConfigWithCustomHealthcheck 带自定义健康检查的配置
	ConfigWithCustomHealthcheck = configtestutil.ConfigWithCustomHealthcheck

	// ConfigWithMultiArchPlugin 带多架构插件的配置
	ConfigWithMultiArchPlugin = configtestutil.ConfigWithMultiArchPlugin
)

// ============================================
// 选项函数 - 导出核心层的选项
// ============================================

var (
	// WithServiceOpt 设置服务信息（名称和描述）
	WithServiceOpt = configtestutil.WithServiceOpt

	// WithServiceNameOpt 设置服务名称
	WithServiceNameOpt = configtestutil.WithServiceNameOpt

	// WithPortOpt 添加端口配置
	WithPortOpt = configtestutil.WithPortOpt

	// WithLanguageOpt 设置语言类型
	WithLanguageOpt = configtestutil.WithLanguageOpt

	// WithBuilderOpt 添加构建镜像预设
	WithBuilderOpt = configtestutil.WithBuilderOpt

	// WithRuntimeOpt 添加运行时镜像预设
	WithRuntimeOpt = configtestutil.WithRuntimeOpt

	// WithBuilderImageOpt 设置构建镜像引用
	WithBuilderImageOpt = configtestutil.WithBuilderImageOpt

	// WithRuntimeImageOpt 设置运行时镜像引用
	WithRuntimeImageOpt = configtestutil.WithRuntimeImageOpt

	// WithBuildCommandOpt 设置构建命令
	WithBuildCommandOpt = configtestutil.WithBuildCommandOpt

	// WithPluginOpt 添加插件
	WithPluginOpt = configtestutil.WithPluginOpt

	// WithSystemPackagesOpt 设置系统包依赖
	WithSystemPackagesOpt = configtestutil.WithSystemPackagesOpt

	// WithHealthcheckOpt 设置健康检查
	WithHealthcheckOpt = configtestutil.WithHealthcheckOpt

	// WithCustomHealthcheckOpt 设置自定义健康检查
	WithCustomHealthcheckOpt = configtestutil.WithCustomHealthcheckOpt

	// WithStartupCommandOpt 设置启动命令
	WithStartupCommandOpt = configtestutil.WithStartupCommandOpt

	// WithPluginInstallDirOpt 设置插件安装目录
	WithPluginInstallDirOpt = configtestutil.WithPluginInstallDirOpt

	// WithPreBuildCommandOpt 设置预构建命令
	WithPreBuildCommandOpt = configtestutil.WithPreBuildCommandOpt

	// WithLanguageConfigOpt 设置语言配置
	WithLanguageConfigOpt = configtestutil.WithLanguageConfigOpt
)

// ============================================
// 向后兼容API (Backward Compatibility)
// ============================================
// 这些函数保持旧的API，以便现有测试代码无需修改

// ServiceConfig 类型别名，指向 config.ServiceConfig
// Deprecated: 直接使用 config.ServiceConfig
type ServiceConfig = config.ServiceConfig

// NewBuilder 创建配置构建器（兼容旧API）
// Deprecated: 使用 NewConfigBuilder 代替
//
// 示例迁移：
//
//	旧代码: cfg := testutil.NewBuilder().WithServiceName("test").Build()
//	新代码: cfg := testutil.NewConfigBuilder().WithServiceName("test").Build()
func NewBuilder() *configtestutil.ConfigBuilder {
	return configtestutil.NewConfigBuilder().
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuildCommand("go build -o bin/app").
		WithCIScriptDir(".tad/build/test-service").
		WithCIBuildConfigDir(".tad/build/test-service/build")
}

// NewMinimal 创建最小配置（兼容旧API）
// Deprecated: 使用 MinimalConfig 代替
//
// 示例迁移：
//
//	旧代码: cfg := testutil.NewMinimal("my-service")
//	新代码: cfg := testutil.MinimalConfig()
func NewMinimal(serviceName string) *config.ServiceConfig {
	cfg := configtestutil.MinimalConfig()
	cfg.Service.Name = serviceName
	return cfg
}

// ============================================
// 生成器特定的便捷函数
// ============================================
// 这些函数是生成器测试专用的，提供常用的配置组合

// NewGeneratorTestConfig 创建生成器测试专用的标准配置
// 包含生成器测试所需的所有默认值
func NewGeneratorTestConfig() *config.ServiceConfig {
	return configtestutil.NewConfigBuilder().
		WithService("generator-test", "Generator Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithCIScriptDir(".tad/build/generator-test").
		WithCIBuildConfigDir(".tad/build/generator-test/build").
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// NewGeneratorTestConfigWithPlugins 创建带插件的生成器测试配置
func NewGeneratorTestConfigWithPlugins() *config.ServiceConfig {
	return configtestutil.NewConfigBuilder().
		WithService("generator-test", "Generator Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPluginInstallDir("/tce").
		WithPlugin(config.PluginConfig{
			Name:        "selfMonitor",
			Description: "TCE Self Monitor Tool",
			DownloadURL: config.NewStaticDownloadURL("https://mirrors.tencent.com/repository/generic/selfMonitor/download_tool.sh"),
			InstallCommand: `echo "Installing selfMonitor..."
curl -fsSL "${PLUGIN_DOWNLOAD_URL}" | bash -s "${PLUGIN_WORK_DIR}"`,
			RuntimeEnv: []config.EnvironmentVariable{
				{Name: "TCESTAURY_TOOL_PATH", Value: "${PLUGIN_INSTALL_DIR}"},
			},
			Required: true,
		}).
		WithCIScriptDir(".tad/build/generator-test").
		WithCIBuildConfigDir(".tad/build/generator-test/build").
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// ============================================
// 基础镜像预设 - 导出核心层的预设
// ============================================

// DefaultBaseImages 返回默认的基础镜像配置
// Deprecated: 使用 configtestutil.DefaultBaseImages() 代替
func DefaultBaseImages() config.BaseImagesConfig {
	return configtestutil.DefaultBaseImages()
}
