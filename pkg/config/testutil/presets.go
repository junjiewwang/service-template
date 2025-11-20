package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// ============================================
// 预设配置
// ============================================

// MinimalConfig 返回最小化配置（仅包含必需字段）
func MinimalConfig() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_default", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_default", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_default").
		WithRuntimeImage("@runtimes.alpine_default").
		WithBuildCommand("go build -o bin/app").
		WithStartupCommand("exec ./bin/${SERVICE_NAME}").
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// GoServiceConfig 返回 Go 服务的标准配置
func GoServiceConfig() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("go-service", "Go Service").
		WithPort("http", 8080, "TCP", true).
		WithPort("metrics", 9090, "TCP", false).
		WithLanguage("go").
		WithLanguageConfig("goproxy", "https://goproxy.cn,direct").
		WithLanguageConfig("gosumdb", "sum.golang.org").
		WithBuilder("go_1.21", "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0", "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0").
		WithRuntime("tencentos_minimal", "mirrors.tencent.com/tencentos/tencentos3-minimal:latest", "mirrors.tencent.com/tencentos/tencentos3-minimal:latest").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.tencentos_minimal").
		WithSystemPackages([]string{"git", "make", "gcc"}).
		WithBuildCommand("go build -o bin/${SERVICE_NAME} ./cmd/main.go").
		WithRuntimePackages([]string{"ca-certificates", "tzdata"}).
		WithStartupCommand("exec ./bin/${SERVICE_NAME}").
		WithHealthcheck(true, "default").
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// PythonServiceConfig 返回 Python 服务的标准配置
func PythonServiceConfig() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("python-service", "Python Service").
		WithPort("http", 8000, "TCP", true).
		WithLanguage("python").
		WithLanguageConfig("pip_index_url", "https://mirrors.tencent.com/pypi/simple").
		WithBuilder("python_3.11", "python:3.11-slim", "python:3.11-slim").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.python_3.11").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithSystemPackages([]string{"gcc", "python3-dev"}).
		WithBuildCommand("pip install -r requirements.txt").
		WithRuntimePackages([]string{"python3", "ca-certificates"}).
		WithStartupCommand("python3 app.py").
		WithHealthcheck(true, "default").
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// JavaServiceConfig 返回 Java 服务的标准配置
func JavaServiceConfig() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("java-service", "Java Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("java").
		WithBuilder("java_17", "openjdk:17-jdk-slim", "openjdk:17-jdk-slim").
		WithRuntime("java_17_jre", "openjdk:17-jre-slim", "openjdk:17-jre-slim").
		WithBuilderImage("@builders.java_17").
		WithRuntimeImage("@runtimes.java_17_jre").
		WithSystemPackages([]string{"maven"}).
		WithBuildCommand("mvn clean package -DskipTests").
		WithRuntimePackages([]string{"ca-certificates"}).
		WithStartupCommand("java -jar app.jar").
		WithHealthcheck(true, "default").
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// ConfigWithPlugins 返回带插件的配置
func ConfigWithPlugins() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("service-with-plugins", "Service With Plugins").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithStartupCommand("exec ./bin/${SERVICE_NAME}").
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
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// ConfigWithCustomHealthcheck 返回带自定义健康检查的配置
func ConfigWithCustomHealthcheck() *config.ServiceConfig {
	return NewConfigBuilder().
		WithService("service-with-custom-hc", "Service With Custom Healthcheck").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithStartupCommand("exec ./bin/${SERVICE_NAME}").
		WithCustomHealthcheck(`#!/bin/sh
curl -f http://localhost:8080/health || exit 1`).
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// ConfigWithMultiArchPlugin 返回带多架构插件的配置
func ConfigWithMultiArchPlugin() *config.ServiceConfig {
	downloadURLs := map[string]string{
		"x86_64":  "https://example.com/plugin-x86_64.tar.gz",
		"aarch64": "https://example.com/plugin-aarch64.tar.gz",
	}

	return NewConfigBuilder().
		WithService("service-with-multiarch-plugin", "Service With Multi-Arch Plugin").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithStartupCommand("exec ./bin/${SERVICE_NAME}").
		WithPluginInstallDir("/opt/plugins").
		WithPlugin(config.PluginConfig{
			Name:        "jre",
			Description: "Java Runtime Environment",
			DownloadURL: config.NewArchMappingDownloadURL(downloadURLs),
			InstallCommand: `echo "Installing JRE..."
curl -fsSL "${PLUGIN_DOWNLOAD_URL}" -o /tmp/jdk.tar.gz
mkdir -p "${PLUGIN_WORK_DIR}/${PLUGIN_NAME}"
tar -xzf /tmp/jdk.tar.gz -C "${PLUGIN_WORK_DIR}/${PLUGIN_NAME}" --strip-components=1`,
			RuntimeEnv: []config.EnvironmentVariable{
				{Name: "JAVA_HOME", Value: "${PLUGIN_INSTALL_DIR}/jre"},
				{Name: "PATH", Value: "${PLUGIN_INSTALL_DIR}/jre/bin:${PATH}"},
			},
			Required: false,
		}).
		WithMetadata("2.0.0", "svcgen").
		BuildWithDefaults()
}

// ============================================
// 预设基础镜像配置
// ============================================

// DefaultBaseImages 返回默认的基础镜像配置
func DefaultBaseImages() config.BaseImagesConfig {
	return config.BaseImagesConfig{
		Builders: map[string]config.ArchImageConfig{
			"go_1.21": {
				AMD64: "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0",
				ARM64: "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0",
			},
			"go_1.22": {
				AMD64: "golang:1.22",
				ARM64: "golang:1.22",
			},
			"python_3.11": {
				AMD64: "python:3.11-slim",
				ARM64: "python:3.11-slim",
			},
		},
		Runtimes: map[string]config.ArchImageConfig{
			"tencentos_minimal": {
				AMD64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest",
				ARM64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest",
			},
			"alpine_3.18": {
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
		},
	}
}
