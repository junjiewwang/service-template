package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================
// 端到端验收测试
// 验证 YAML 加载 → 验证 → 镜像解析 → 构建命令解析 的完整链路
// ============================================

// --- 用例1: 格式1 - 完全不填镜像和构建命令（自动推导） ---

func TestAcceptance_Format1_AutoInfer_Go(t *testing.T) {
	yaml := `
service:
  name: my-go-service
  ports:
    - name: http
      port: 8080
      protocol: TCP
      expose: true

language:
  type: go
  config:
    go_version: "1.21"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	// 验证通过（无需 base_images、无需 builder_image/runtime_image、无需 build command）
	validator := NewValidator(cfg)
	err = validator.Validate()
	require.NoError(t, err, "format1 Go config should pass validation")

	// 镜像自动推导
	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "golang:1.21-alpine", builderImg.AMD64, "Go builder should be golang:{version}-alpine")
	assert.Equal(t, "golang:1.21-alpine", builderImg.ARM64, "multi-arch: same image for both")

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "alpine:3.19", runtimeImg.AMD64, "Go runtime should be alpine (static binary)")
	assert.Equal(t, "alpine:3.19", runtimeImg.ARM64)

	// 构建命令自动推导
	buildCmd := ResolveBuildCommand(cfg)
	assert.Contains(t, buildCmd, "go build", "Go default build command should contain 'go build'")
}

func TestAcceptance_Format1_AutoInfer_Python(t *testing.T) {
	yaml := `
service:
  name: my-python-service
  ports:
    - name: http
      port: 8000
      protocol: TCP

language:
  type: python
  config:
    python_version: "3.11"

runtime:
  startup:
    command: "python app.py"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "python:3.11-slim", builderImg.AMD64)

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "python:3.11-slim", runtimeImg.AMD64, "Python runtime needs Python environment")

	buildCmd := ResolveBuildCommand(cfg)
	assert.Contains(t, buildCmd, "cp -r")
}

func TestAcceptance_Format1_AutoInfer_Java(t *testing.T) {
	yaml := `
service:
  name: my-java-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: java
  config:
    jdk_version: "17"
    build_tool: "gradle"
    gradle_version: "8"

runtime:
  startup:
    command: "java -jar app.jar"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "gradle:8-jdk17", builderImg.AMD64, "Java+Gradle builder image")

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "eclipse-temurin:17-jre-alpine", runtimeImg.AMD64, "Java runtime: JRE only")

	buildCmd := ResolveBuildCommand(cfg)
	assert.Contains(t, buildCmd, "gradle build", "Java+Gradle default build command")
}

func TestAcceptance_Format1_AutoInfer_NodeJS(t *testing.T) {
	yaml := `
service:
  name: my-node-service
  ports:
    - name: http
      port: 3000
      protocol: TCP

language:
  type: nodejs
  config:
    node_version: "18"

runtime:
  startup:
    command: "node index.js"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "node:18-alpine", builderImg.AMD64)

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "node:18-alpine", runtimeImg.AMD64, "Node.js runtime needs Node environment")

	buildCmd := ResolveBuildCommand(cfg)
	assert.Contains(t, buildCmd, "npm run build")
}

func TestAcceptance_Format1_AutoInfer_Rust(t *testing.T) {
	yaml := `
service:
  name: my-rust-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: rust

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "rust:1.78-alpine", builderImg.AMD64, "Rust default version")

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "alpine:3.19", runtimeImg.AMD64, "Rust static binary: alpine")

	buildCmd := ResolveBuildCommand(cfg)
	assert.Contains(t, buildCmd, "cargo build")
}

// --- 用例2: 格式2 - 直接指定镜像名（multi-arch） ---

func TestAcceptance_Format2_DirectImage(t *testing.T) {
	yaml := `
service:
  name: my-go-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image: "golang:1.23-alpine"
  runtime_image: "alpine:3.19"
  commands:
    build: "go build -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/server"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	// 验证：无需 base_images
	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	// 直接镜像名解析
	assert.Equal(t, ImageSpecDirect, cfg.Build.BuilderImage.Kind())
	assert.Equal(t, ImageSpecDirect, cfg.Build.RuntimeImage.Kind())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-alpine", builderImg.AMD64, "direct: same image for amd64")
	assert.Equal(t, "golang:1.23-alpine", builderImg.ARM64, "direct: same image for arm64")

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "alpine:3.19", runtimeImg.AMD64)
	assert.Equal(t, "alpine:3.19", runtimeImg.ARM64)

	// 用户显式命令优先
	assert.Equal(t, "go build -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/server", ResolveBuildCommand(cfg))
}

// --- 用例3: 格式3 - 预设引用 (@builders.xxx / @runtimes.xxx) ---

func TestAcceptance_Format3_PresetRef(t *testing.T) {
	yaml := `
base_images:
  builders:
    go_1.23:
      amd64: "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0"
      arm64: "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0"
  runtimes:
    tencentos_minimal:
      amd64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"
      arm64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"

service:
  name: my-go-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image: "@builders.go_1.23"
  runtime_image: "@runtimes.tencentos_minimal"
  commands:
    build: "go build -o app"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	assert.Equal(t, ImageSpecPreset, cfg.Build.BuilderImage.Kind())
	assert.Equal(t, ImageSpecPreset, cfg.Build.RuntimeImage.Kind())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0", builderImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0", builderImg.ARM64)

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/tencentos/tencentos3-minimal:latest", runtimeImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/tencentos/tencentos3-minimal:latest", runtimeImg.ARM64)
}

// --- 用例4: 格式4 - 按架构指定 ---

func TestAcceptance_Format4_PerArch(t *testing.T) {
	yaml := `
service:
  name: my-go-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image:
    amd64: "mirrors.tencent.com/builder:x86"
    arm64: "mirrors.tencent.com/builder:arm"
  runtime_image:
    amd64: "mirrors.tencent.com/runtime:x86"
    arm64: "mirrors.tencent.com/runtime:arm"
  commands:
    build: "go build -o app"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	assert.Equal(t, ImageSpecPerArch, cfg.Build.BuilderImage.Kind())
	assert.Equal(t, ImageSpecPerArch, cfg.Build.RuntimeImage.Kind())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/builder:x86", builderImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/builder:arm", builderImg.ARM64)

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/runtime:x86", runtimeImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/runtime:arm", runtimeImg.ARM64)
}

// --- 用例5: 混合格式（builder 用直接镜像，runtime 用预设引用） ---

func TestAcceptance_MixedFormats(t *testing.T) {
	yaml := `
base_images:
  runtimes:
    alpine:
      amd64: "alpine:3.19-amd64"
      arm64: "alpine:3.19-arm64"

service:
  name: mixed-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image: "golang:1.23-alpine"
  runtime_image: "@runtimes.alpine"
  commands:
    build: "go build"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	assert.Equal(t, ImageSpecDirect, cfg.Build.BuilderImage.Kind())
	assert.Equal(t, ImageSpecPreset, cfg.Build.RuntimeImage.Kind())

	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-alpine", builderImg.AMD64)
	assert.Equal(t, "golang:1.23-alpine", builderImg.ARM64)

	runtimeImg, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "alpine:3.19-amd64", runtimeImg.AMD64)
	assert.Equal(t, "alpine:3.19-arm64", runtimeImg.ARM64)
}

// --- 用例6: 格式1（自动推导）但用户覆盖构建命令 ---

func TestAcceptance_AutoInferImage_ExplicitBuildCommand(t *testing.T) {
	yaml := `
service:
  name: custom-build-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  commands:
    build: "make build"
    pre_build: "make deps"
    post_build: "make test"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate())

	// 镜像自动推导
	builderImg, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-alpine", builderImg.AMD64, "default go version")

	// 构建命令用用户的
	assert.Equal(t, "make build", ResolveBuildCommand(cfg))
}

// ============================================
// 错误用例：应该被 validator 拦截
// ============================================

func TestAcceptance_Error_PresetRef_MissingBaseImages(t *testing.T) {
	yaml := `
service:
  name: bad-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image: "@builders.go_1.23"
  runtime_image: "@runtimes.alpine"
  commands:
    build: "go build"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	err = validator.Validate()
	require.Error(t, err, "preset ref without base_images should fail")
	assert.Contains(t, err.Error(), "base_images is required")
}

func TestAcceptance_Error_PresetRef_NotFound(t *testing.T) {
	yaml := `
base_images:
  builders:
    go_1.22:
      amd64: "golang:1.22"
      arm64: "golang:1.22"
  runtimes:
    alpine:
      amd64: "alpine:3.19"
      arm64: "alpine:3.19"

service:
  name: bad-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image: "@builders.nonexistent"
  runtime_image: "@runtimes.alpine"
  commands:
    build: "go build"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	err = validator.Validate()
	require.Error(t, err, "preset ref to nonexistent preset should fail")
	assert.Contains(t, err.Error(), "not found")
}

func TestAcceptance_Error_PresetRef_WrongCategory(t *testing.T) {
	yaml := `
base_images:
  builders:
    go_1.22:
      amd64: "golang:1.22"
      arm64: "golang:1.22"
  runtimes:
    alpine:
      amd64: "alpine:3.19"
      arm64: "alpine:3.19"

service:
  name: bad-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image: "@runtimes.alpine"
  runtime_image: "@builders.go_1.22"
  commands:
    build: "go build"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	err = validator.Validate()
	require.Error(t, err, "wrong category in preset ref should fail")
	assert.Contains(t, err.Error(), "must reference @builders")
}

func TestAcceptance_Error_PerArch_MissingArm64(t *testing.T) {
	yaml := `
service:
  name: bad-service
  ports:
    - name: http
      port: 8080
      protocol: TCP

language:
  type: go

build:
  builder_image:
    amd64: "golang:1.23"
  runtime_image: "alpine:3.19"
  commands:
    build: "go build"

runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	err = validator.Validate()
	require.Error(t, err, "per-arch with missing arm64 should fail")
	assert.Contains(t, err.Error(), "arm64 image is required")
}

// ============================================
// YAML 序列化 / 反序列化 roundtrip 测试
// ============================================

func TestAcceptance_YAML_Roundtrip_AllFormats(t *testing.T) {
	t.Run("format2 direct string roundtrip", func(t *testing.T) {
		original := `
service:
  name: test
  ports:
    - name: http
      port: 8080
      protocol: TCP
language:
  type: go
build:
  builder_image: "golang:1.23-alpine"
  runtime_image: "alpine:3.19"
  commands:
    build: "go build"
runtime:
  startup:
    command: "./app"
`
		cfg, err := LoadFromBytes([]byte(original))
		require.NoError(t, err)
		assert.Equal(t, ImageSpecDirect, cfg.Build.BuilderImage.Kind())
		assert.Equal(t, "golang:1.23-alpine", cfg.Build.BuilderImage.String())
	})

	t.Run("format3 preset ref roundtrip", func(t *testing.T) {
		original := `
base_images:
  builders:
    go_1.23:
      amd64: "golang:1.23"
      arm64: "golang:1.23"
  runtimes:
    alpine:
      amd64: "alpine:3.19"
      arm64: "alpine:3.19"
service:
  name: test
  ports:
    - name: http
      port: 8080
      protocol: TCP
language:
  type: go
build:
  builder_image: "@builders.go_1.23"
  runtime_image: "@runtimes.alpine"
  commands:
    build: "go build"
runtime:
  startup:
    command: "./app"
`
		cfg, err := LoadFromBytes([]byte(original))
		require.NoError(t, err)
		assert.Equal(t, ImageSpecPreset, cfg.Build.BuilderImage.Kind())
		assert.Equal(t, "@builders.go_1.23", cfg.Build.BuilderImage.String())
	})

	t.Run("format4 per-arch roundtrip", func(t *testing.T) {
		original := `
service:
  name: test
  ports:
    - name: http
      port: 8080
      protocol: TCP
language:
  type: go
build:
  builder_image:
    amd64: "go:amd64"
    arm64: "go:arm64"
  runtime_image:
    amd64: "rt:amd64"
    arm64: "rt:arm64"
  commands:
    build: "go build"
runtime:
  startup:
    command: "./app"
`
		cfg, err := LoadFromBytes([]byte(original))
		require.NoError(t, err)
		assert.Equal(t, ImageSpecPerArch, cfg.Build.BuilderImage.Kind())

		resolved, err := cfg.Build.BuilderImage.Resolve(nil, "builders")
		require.NoError(t, err)
		assert.Equal(t, "go:amd64", resolved.AMD64)
		assert.Equal(t, "go:arm64", resolved.ARM64)
	})

	t.Run("format1 empty roundtrip", func(t *testing.T) {
		original := `
service:
  name: test
  ports:
    - name: http
      port: 8080
      protocol: TCP
language:
  type: go
runtime:
  startup:
    command: "./app"
`
		cfg, err := LoadFromBytes([]byte(original))
		require.NoError(t, err)
		assert.True(t, cfg.Build.BuilderImage.IsEmpty())
		assert.True(t, cfg.Build.RuntimeImage.IsEmpty())
		assert.Equal(t, ImageSpecEmpty, cfg.Build.BuilderImage.Kind())
	})
}

// ============================================
// 多语言默认值版本矩阵测试
// ============================================

func TestAcceptance_DefaultVersionMatrix(t *testing.T) {
	// 验证各语言使用默认版本时的镜像推导结果
	tests := []struct {
		name            string
		langType        string
		wantBuilder     string
		wantRuntime     string
		wantBuildCmdKey string // 构建命令中应包含的关键字
	}{
		{"go defaults", "go", "golang:1.23-alpine", "alpine:3.19", "go build"},
		{"python defaults", "python", "python:3.12-slim", "python:3.12-slim", "cp -r"},
		{"java defaults", "java", "maven:3-eclipse-temurin-21", "eclipse-temurin:21-jre-alpine", "mvn package"},
		{"nodejs defaults", "nodejs", "node:20-alpine", "node:20-alpine", "npm run build"},
		{"rust defaults", "rust", "rust:1.78-alpine", "alpine:3.19", "cargo build"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &ServiceConfig{
				Service: ServiceInfo{
					Name:  "test-" + tt.langType,
					Ports: []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
				},
				Language: LanguageConfig{Type: tt.langType},
				Runtime: RuntimeConfig{
					Startup: StartupConfig{Command: "./app"},
				},
			}
			applyDefaults(cfg)

			// 验证通过
			validator := NewValidator(cfg)
			require.NoError(t, validator.Validate(), "%s should pass validation with all defaults", tt.langType)

			// 镜像推导
			builder, err := ResolveBuilderImageWithDefaults(cfg)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBuilder, builder.AMD64, "%s builder image", tt.langType)
			assert.Equal(t, tt.wantBuilder, builder.ARM64, "%s builder image (arm64 == amd64)", tt.langType)

			runtime, err := ResolveRuntimeImageWithDefaults(cfg)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRuntime, runtime.AMD64, "%s runtime image", tt.langType)

			// 构建命令推导
			buildCmd := ResolveBuildCommand(cfg)
			assert.Contains(t, buildCmd, tt.wantBuildCmdKey, "%s build command", tt.langType)
		})
	}
}

// ============================================
// 优先级测试：显式配置 > 默认推导
// ============================================

func TestAcceptance_Priority_ExplicitOverridesDefault(t *testing.T) {
	t.Run("explicit image overrides language default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Service: ServiceInfo{
				Name:  "test",
				Ports: []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			},
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				BuilderImage: NewImageSpec("custom-golang:1.20"),
				RuntimeImage: NewImageSpec("custom-alpine:3.18"),
			},
			Runtime: RuntimeConfig{
				Startup: StartupConfig{Command: "./app"},
			},
		}
		applyDefaults(cfg)

		builder, err := ResolveBuilderImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "custom-golang:1.20", builder.AMD64, "explicit image should override default")

		runtime, err := ResolveRuntimeImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "custom-alpine:3.18", runtime.AMD64)
	})

	t.Run("explicit build command overrides language default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				Commands: BuildCommandsConfig{Build: "make release"},
			},
		}

		assert.Equal(t, "make release", ResolveBuildCommand(cfg))
	})

	t.Run("per-arch image overrides language default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Service: ServiceInfo{
				Name:  "test",
				Ports: []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			},
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				BuilderImage: NewImageSpecPerArch("private-registry:amd64", "private-registry:arm64"),
				RuntimeImage: NewImageSpec("alpine:3.19"),
			},
			Runtime: RuntimeConfig{
				Startup: StartupConfig{Command: "./app"},
			},
		}
		applyDefaults(cfg)

		validator := NewValidator(cfg)
		require.NoError(t, validator.Validate())

		builder, err := ResolveBuilderImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "private-registry:amd64", builder.AMD64)
		assert.Equal(t, "private-registry:arm64", builder.ARM64)
	})
}

// ============================================
// 边界用例
// ============================================

func TestAcceptance_Edge_NoPortsService(t *testing.T) {
	yaml := `
service:
  name: worker-service
  description: "Background worker with no ports"

language:
  type: go

runtime:
  startup:
    command: "./worker"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	require.NoError(t, validator.Validate(), "service without ports should be valid")

	assert.Empty(t, cfg.Service.Ports)
}

func TestAcceptance_Edge_MinimalGoConfig(t *testing.T) {
	// 最简 Go 配置：仅需 service.name + language.type + runtime.startup.command
	yaml := `
service:
  name: minimal
language:
  type: go
runtime:
  startup:
    command: "./app"
`
	cfg, err := LoadFromBytes([]byte(yaml))
	require.NoError(t, err)

	validator := NewValidator(cfg)
	err = validator.Validate()
	require.NoError(t, err, "minimal Go config should be valid with all defaults inferred")

	// 所有推导正确
	builder, err := ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-alpine", builder.AMD64)

	runtime, err := ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "alpine:3.19", runtime.AMD64)

	buildCmd := ResolveBuildCommand(cfg)
	assert.NotEmpty(t, buildCmd)
	assert.Contains(t, buildCmd, "go build")
}

func TestAcceptance_Edge_BaseImagesOptionalWhenNotUsed(t *testing.T) {
	// base_images 在不使用预设引用时完全可选
	tests := []struct {
		name  string
		yaml  string
		valid bool
	}{
		{
			name: "no base_images, no image refs (auto-infer)",
			yaml: `
service:
  name: test
language:
  type: go
runtime:
  startup:
    command: "./app"
`,
			valid: true,
		},
		{
			name: "no base_images, direct image",
			yaml: `
service:
  name: test
language:
  type: go
build:
  builder_image: "golang:1.23"
  runtime_image: "alpine:3.19"
runtime:
  startup:
    command: "./app"
`,
			valid: true,
		},
		{
			name: "no base_images, per-arch image",
			yaml: `
service:
  name: test
language:
  type: go
build:
  builder_image:
    amd64: "go:amd64"
    arm64: "go:arm64"
  runtime_image:
    amd64: "rt:amd64"
    arm64: "rt:arm64"
runtime:
  startup:
    command: "./app"
`,
			valid: true,
		},
		{
			name: "no base_images, preset ref (should fail)",
			yaml: `
service:
  name: test
language:
  type: go
build:
  builder_image: "@builders.go"
  runtime_image: "@runtimes.alpine"
runtime:
  startup:
    command: "./app"
`,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFromBytes([]byte(tt.yaml))
			require.NoError(t, err)

			validator := NewValidator(cfg)
			err = validator.Validate()

			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), "base_images") || strings.Contains(err.Error(), "preset"),
					"error should mention base_images or preset: %v", err)
			}
		})
	}
}
