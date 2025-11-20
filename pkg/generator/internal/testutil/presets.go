package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// BaseImagesPreset 基础镜像预设
type BaseImagesPreset struct {
	Name     string
	Builders map[string]config.ArchImageConfig
	Runtimes map[string]config.ArchImageConfig
}

// ToConfig 转换为配置对象
func (p BaseImagesPreset) ToConfig() config.BaseImagesConfig {
	return config.BaseImagesConfig{
		Builders: p.Builders,
		Runtimes: p.Runtimes,
	}
}

// 预定义的镜像预设
var (
	// MinimalBaseImages 最小镜像（用于快速测试）
	MinimalBaseImages = BaseImagesPreset{
		Name: "minimal",
		Builders: map[string]config.ArchImageConfig{
			"builder": {
				AMD64: "builder:test",
				ARM64: "builder:test",
			},
		},
		Runtimes: map[string]config.ArchImageConfig{
			"runtime": {
				AMD64: "runtime:test",
				ARM64: "runtime:test",
			},
		},
	}

	// DefaultBaseImages 默认镜像（用于大多数测试）
	DefaultBaseImages = BaseImagesPreset{
		Name: "default",
		Builders: map[string]config.ArchImageConfig{
			"test_builder": {
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
		},
		Runtimes: map[string]config.ArchImageConfig{
			"test_runtime": {
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
		},
	}

	// CustomBaseImages 自定义镜像（用于特殊测试）
	CustomBaseImages = BaseImagesPreset{
		Name: "custom",
		Builders: map[string]config.ArchImageConfig{
			"go_1.21": {
				AMD64: "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0",
				ARM64: "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0",
			},
			"python_3.11": {
				AMD64: "docker.io/python:3.11-slim",
				ARM64: "docker.io/python:3.11-slim",
			},
		},
		Runtimes: map[string]config.ArchImageConfig{
			"tencentos": {
				AMD64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest",
				ARM64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest",
			},
			"alpine": {
				AMD64: "docker.io/alpine:3.18",
				ARM64: "docker.io/alpine:3.18",
			},
		},
	}
)

// NewGoServicePreset 创建 Go 服务预设配置
func NewGoServicePreset() *ConfigBuilder {
	return NewDefaultConfig().
		WithLanguage("go").
		WithBuild(func(b *BuildConfigBuilder) {
			b.BuildCommand("go build -o bin/app ./cmd/main.go").
				PreBuildCommand("go mod download")
		})
}

// NewPythonServicePreset 创建 Python 服务预设配置
func NewPythonServicePreset() *ConfigBuilder {
	return NewDefaultConfig().
		WithLanguage("python").
		WithBuild(func(b *BuildConfigBuilder) {
			b.BuildCommand("pip install -r requirements.txt")
		})
}

// NewServiceWithPluginsPreset 创建带插件的服务预设配置
func NewServiceWithPluginsPreset() *ConfigBuilder {
	return NewDefaultConfig().
		WithPlugins(func(p *PluginBuilder) {
			p.InstallDir("/tce")
		})
}

// NewCustomServicePreset 创建自定义镜像的服务预设配置
func NewCustomServicePreset() *ConfigBuilder {
	return NewDefaultConfig().
		WithBaseImagesPreset(CustomBaseImages).
		WithBuild(func(b *BuildConfigBuilder) {
			b.BuilderImage("@builders.go_1.21").
				RuntimeImage("@runtimes.tencentos")
		})
}
