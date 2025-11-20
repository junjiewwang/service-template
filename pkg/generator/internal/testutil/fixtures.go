package testutil

import (
	"github.com/junjiewwang/service-template/pkg/config"
	configtestutil "github.com/junjiewwang/service-template/pkg/config/testutil"
)

// NewTestBaseImages creates test base images configuration
func NewTestBaseImages() config.BaseImagesConfig {
	return configtestutil.DefaultBaseImages()
}

// NewTestConfig creates a test service configuration
// 使用新的 testutil 包创建配置
func NewTestConfig() *config.ServiceConfig {
	return configtestutil.NewConfigBuilder().
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

// NewTestConfigWithPlugins creates a test config with plugins
func NewTestConfigWithPlugins() *config.ServiceConfig {
	return configtestutil.NewConfigBuilder().
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
