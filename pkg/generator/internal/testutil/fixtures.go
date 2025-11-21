package testutil

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// ============================================
// 生成器测试固定数据
// ============================================
// 这个文件提供生成器测试专用的固定配置数据
// 所有实现都委托给核心层 (pkg/config/testutil)

// NewTestBaseImages 创建测试用的基础镜像配置
// Deprecated: 使用 testutil.DefaultBaseImages() 代替
func NewTestBaseImages() config.BaseImagesConfig {
	return DefaultBaseImages()
}

// NewTestConfig 创建标准的测试服务配置
// 这是生成器测试中最常用的配置
func NewTestConfig() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service ./cmd/test-service").
		WithDependencyFiles(false, []string{"go.mod", "go.sum"}).
		WithStartupCommand("./bin/test-service").
		WithHealthcheck(true, "default").
		WithCIScriptDir(".tad/build/test-service").
		WithCIBuildConfigDir(".tad/build/test-service/build").
		WithMetadata("1.0.0", "svcgen").
		BuildWithDefaults()
}

// NewTestConfigWithPlugins 创建带插件的测试配置
// 用于测试插件相关的生成器功能
func NewTestConfigWithPlugins() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service ./cmd/test-service").
		WithStartupCommand("./bin/test-service").
		WithPluginInstallDir("/tce").
		WithPlugin(config.PluginConfig{
			Name:        "selfMonitor",
			Description: "Self monitoring plugin",
			DownloadURL: config.NewStaticDownloadURL("https://example.com/selfMonitor.tar.gz"),
			RuntimeEnv: []config.EnvironmentVariable{
				{Name: "TOOL_PATH", Value: "${PLUGIN_INSTALL_DIR}"},
			},
			Required: true,
		}).
		WithMetadata("1.0.0", "svcgen").
		BuildWithDefaults()
}
