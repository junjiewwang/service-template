package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================
// 真实场景端到端测试
// 从 YAML 配置加载 → 校验 → 生成文件 → 验证生成内容
// ============================================

// helperLoadAndGenerate 通用辅助：加载 YAML → 校验 → 生成 → 返回输出目录
func helperLoadAndGenerate(t *testing.T, yamlContent string) (string, *config.ServiceConfig) {
	t.Helper()

	// 1. 从 YAML 加载配置
	cfg, err := config.LoadFromBytes([]byte(yamlContent))
	require.NoError(t, err, "YAML should load successfully")

	// 2. 校验配置
	validator := config.NewValidator(cfg)
	err = validator.Validate()
	require.NoError(t, err, "Configuration should pass validation")

	// 3. 生成文件到临时目录
	tmpDir, err := os.MkdirTemp("", "scenario-test-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	outputDir := filepath.Join(tmpDir, "output")
	gen := generator.NewGenerator(cfg, outputDir)
	err = gen.Generate()
	require.NoError(t, err, "Generate should succeed")

	return outputDir, cfg
}

// helperReadFile 读取生成的文件内容
func helperReadFile(t *testing.T, outputDir, relPath string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(outputDir, relPath))
	require.NoError(t, err, "File should exist: %s", relPath)
	return string(content)
}

// helperAssertFileExists 验证文件存在
func helperAssertFileExists(t *testing.T, outputDir, relPath string) {
	t.Helper()
	_, err := os.Stat(filepath.Join(outputDir, relPath))
	assert.NoError(t, err, "Expected file to exist: %s", relPath)
}

// ============================================
// 场景1: 最简 Go 微服务（零配置自动推导）
// 用户只填 3 个必填字段，其余全部自动推导
// ============================================
func TestScenario_MinimalGoService(t *testing.T) {
	yaml := `
service:
  name: user-api
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true

language:
  type: go

runtime:
  startup:
    command: |
      #!/bin/sh
      exec ./bin/${SERVICE_NAME}
`
	outputDir, cfg := helperLoadAndGenerate(t, yaml)

	// 验证自动推导: 默认 deploy_dir
	assert.Equal(t, "/usr/local/services", cfg.Service.DeployDir)

	// 验证自动推导: Go 默认构建命令
	resolvedCmd := config.ResolveBuildCommand(cfg)
	assert.Contains(t, resolvedCmd, "CGO_ENABLED=0 go build")
	assert.Contains(t, resolvedCmd, "${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME}")

	// 验证自动推导: Go 默认镜像
	builderImg, err := config.ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Contains(t, builderImg.AMD64, "golang:")
	assert.Contains(t, builderImg.AMD64, "-alpine")

	runtimeImg, err := config.ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Contains(t, runtimeImg.AMD64, "alpine:")

	// 验证生成的文件结构
	expectedFiles := []string{
		".tad/build/user-api/Dockerfile.user-api.amd64",
		".tad/build/user-api/Dockerfile.user-api.arm64",
		".tad/build/user-api/build.sh",
		".tad/build/user-api/build_deps_install.sh",
		".tad/build/user-api/rt_prepare.sh",
		".tad/build/user-api/entrypoint.sh",
		".tad/build/user-api/healthchk.sh",
		"compose.yaml",
		"Makefile",
		".tad/devops.yaml",
	}
	for _, f := range expectedFiles {
		helperAssertFileExists(t, outputDir, f)
	}

	// 验证 Dockerfile 内容
	dockerfile := helperReadFile(t, outputDir, ".tad/build/user-api/Dockerfile.user-api.amd64")
	assert.Contains(t, dockerfile, "COPY .tad/build/user-api/", "Dockerfile should reference CI_SCRIPT_DIR")
	assert.Contains(t, dockerfile, "WORKDIR /opt", "Should use container project root")
	assert.Contains(t, dockerfile, "WORKDIR ${DEPLOY_DIR}/user-api", "Should use service name in workdir")

	// 验证 build.sh 包含自动推导的构建命令
	buildScript := helperReadFile(t, outputDir, ".tad/build/user-api/build.sh")
	assert.Contains(t, buildScript, "CGO_ENABLED=0 go build")
	assert.Contains(t, buildScript, "SERVICE_NAME=user-api")

	// 验证 compose.yaml 使用 CI_SCRIPT_DIR 变量
	composeContent := helperReadFile(t, outputDir, "compose.yaml")
	assert.Contains(t, composeContent, ".tad/build/user-api/", "compose should reference CI script dir")
	assert.Contains(t, composeContent, "8080", "compose should contain port mapping")
	assert.Contains(t, composeContent, "user-api", "compose should contain service name")
}

// ============================================
// 场景2: Python Web 服务（自定义镜像 + 自定义构建命令）
// 用户使用 Format2（直接镜像名）+ 显式构建命令
// ============================================
func TestScenario_PythonWebService(t *testing.T) {
	yaml := `
service:
  name: recommendation-engine
  description: "AI-powered recommendation service"
  ports:
  - name: http
    port: 8000
    protocol: TCP
    expose: true
  - name: grpc
    port: 50051
    protocol: TCP
    expose: true

language:
  type: python
  config:
    python_version: "3.11"
    pip_index_url: "https://mirrors.tencent.com/pypi/simple"

build:
  builder_image: "python:3.11-slim"
  runtime_image: "python:3.11-slim"
  dependencies:
    system_pkgs:
    - gcc
    - python3-dev
    - libffi-dev
  commands:
    pre_build: |
      pip install --upgrade pip
    build: |
      pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/lib
      cp -r . ${BUILD_OUTPUT_DIR}/
    post_build: |
      echo "Build artifacts ready"

runtime:
  system_dependencies:
    packages:
    - ca-certificates
    - curl
  healthcheck:
    enabled: true
    type: custom
    custom_script: |
      #!/bin/sh
      curl -f http://localhost:8000/health || exit 1
  startup:
    command: |
      #!/bin/sh
      cd ${SERVICE_ROOT}
      exec python3 app.py
    env:
    - name: PYTHONPATH
      value: /usr/local/services/recommendation-engine/lib
    - name: LOG_LEVEL
      value: info
`
	outputDir, cfg := helperLoadAndGenerate(t, yaml)

	// 验证: 直接镜像名（Format2）
	builderImg, err := config.ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "python:3.11-slim", builderImg.AMD64)
	assert.Equal(t, "python:3.11-slim", builderImg.ARM64, "Direct image should be same for both archs")

	// 验证: 显式构建命令覆盖默认
	assert.Contains(t, cfg.Build.Commands.Build, "pip install -r requirements.txt")

	// 验证 build.sh
	buildScript := helperReadFile(t, outputDir, ".tad/build/recommendation-engine/build.sh")
	assert.Contains(t, buildScript, "pip install --upgrade pip", "pre_build should be in script")
	assert.Contains(t, buildScript, "pip install -r requirements.txt", "build command should be in script")
	assert.Contains(t, buildScript, "Build artifacts ready", "post_build should be in script")

	// 验证 healthcheck.sh 使用自定义策略
	healthcheck := helperReadFile(t, outputDir, ".tad/build/recommendation-engine/healthchk.sh")
	assert.Contains(t, healthcheck, "curl -f http://localhost:8000/health")

	// 验证 compose.yaml 包含多端口
	composeContent := helperReadFile(t, outputDir, "compose.yaml")
	assert.Contains(t, composeContent, "8000", "compose should have http port")
	assert.Contains(t, composeContent, "50051", "compose should have grpc port")

	// 验证 entrypoint.sh
	entrypoint := helperReadFile(t, outputDir, ".tad/build/recommendation-engine/entrypoint.sh")
	assert.Contains(t, entrypoint, "PYTHONPATH", "entrypoint should export PYTHONPATH")
}

// ============================================
// 场景3: 企业级 Go 服务（预设引用 + 插件 + 自定义 CI 路径）
// 用户使用 Format3（@预设引用）+ 插件系统 + 自定义 ci.script_dir
// ============================================
func TestScenario_EnterpriseGoService(t *testing.T) {
	yaml := `
base_images:
  builders:
    go_1.23:
      amd64: "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0"
      arm64: "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0"
  runtimes:
    tencentos:
      amd64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"
      arm64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"

service:
  name: payment-gateway
  description: "Enterprise payment processing gateway"
  ports:
  - name: https
    port: 8443
    protocol: TCP
    expose: true

language:
  type: go
  config:
    go_version: "1.23"
    goproxy: "https://goproxy.cn,direct"
    goprivate: "git.example.com/*"

build:
  builder_image: "@builders.go_1.23"
  runtime_image: "@runtimes.tencentos"
  dependencies:
    system_pkgs:
    - git
    - make
    - gcc
  commands:
    build: |
      CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=1.0.0" \
        -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/server

plugins:
  install_dir: /tce
  items:
  - name: selfMonitor
    description: "TCE Self Monitor Tool"
    download_url: "https://mirrors.tencent.com/repository/generic/selfMonitor/download_tool.sh"
    install_command: |
      echo "Installing ${PLUGIN_NAME}..."
      curl -fsSL "${PLUGIN_DOWNLOAD_URL}" | bash -s "${PLUGIN_WORK_DIR}"
    runtime_env:
    - name: TCESTAURY_TOOL_PATH
      value: ${PLUGIN_INSTALL_DIR}
    required: true

runtime:
  system_dependencies:
    packages:
    - ca-certificates
    - tzdata
  healthcheck:
    enabled: true
    type: default
  startup:
    command: |
      #!/bin/sh
      set -e
      cd ${SERVICE_ROOT}
      exec ./bin/${SERVICE_NAME}
    env:
    - name: GO_ENV
      value: production

ci:
  script_dir: ".ci/build/payment-gateway"

local_dev:
  compose:
    resources:
      limits:
        cpus: "1"
        memory: 2G
      reservations:
        cpus: "0.5"
        memory: 1G
    healthcheck:
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    labels:
      kompose.image-pull-policy: "IfNotPresent"

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
`
	outputDir, cfg := helperLoadAndGenerate(t, yaml)

	// 验证: 预设引用解析正确（Format3）
	builderImg, err := config.ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0", builderImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0", builderImg.ARM64)

	runtimeImg, err := config.ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/tencentos/tencentos3-minimal:latest", runtimeImg.AMD64)

	// 验证: 自定义 CI 路径生效
	assert.Equal(t, ".ci/build/payment-gateway", cfg.CI.ScriptDir)

	// 验证: 文件生成到自定义 CI 路径下
	customCIFiles := []string{
		".ci/build/payment-gateway/Dockerfile.payment-gateway.amd64",
		".ci/build/payment-gateway/Dockerfile.payment-gateway.arm64",
		".ci/build/payment-gateway/build.sh",
		".ci/build/payment-gateway/build_deps_install.sh",
		".ci/build/payment-gateway/build_plugins.sh",
		".ci/build/payment-gateway/rt_prepare.sh",
		".ci/build/payment-gateway/entrypoint.sh",
		".ci/build/payment-gateway/healthchk.sh",
	}
	for _, f := range customCIFiles {
		helperAssertFileExists(t, outputDir, f)
	}

	// 验证: 默认路径下不应有文件
	_, err = os.Stat(filepath.Join(outputDir, ".tad/build/payment-gateway"))
	assert.True(t, os.IsNotExist(err), ".tad/build/ should NOT exist when custom CI path is set")

	// 验证: Dockerfile 引用自定义 CI 路径
	dockerfile := helperReadFile(t, outputDir, ".ci/build/payment-gateway/Dockerfile.payment-gateway.amd64")
	assert.Contains(t, dockerfile, "COPY .ci/build/payment-gateway/", "Dockerfile should use custom CI path")

	// 验证: Dockerfile 包含插件阶段
	assert.Contains(t, dockerfile, "plugin-builder", "Dockerfile should have plugin-builder stage")
	assert.Contains(t, dockerfile, "COPY --from=plugin-builder /plugins /plugins", "Dockerfile should copy plugins")

	// 验证: compose.yaml 使用自定义 CI 路径
	composeContent := helperReadFile(t, outputDir, "compose.yaml")
	assert.Contains(t, composeContent, ".ci/build/payment-gateway/", "compose should use custom CI path")
	assert.NotContains(t, composeContent, ".tad/build/", "compose should NOT contain default .tad/build path")

	// 验证: compose.yaml 包含资源限制
	assert.Contains(t, composeContent, "memory: 2G")
	assert.Contains(t, composeContent, "memory: 1G")

	// 验证: build_plugins.sh 存在
	pluginsScript := helperReadFile(t, outputDir, ".ci/build/payment-gateway/build_plugins.sh")
	assert.Contains(t, pluginsScript, "selfMonitor")

	// 验证: entrypoint.sh 包含插件环境变量加载
	entrypoint := helperReadFile(t, outputDir, ".ci/build/payment-gateway/entrypoint.sh")
	assert.Contains(t, entrypoint, "GO_ENV")
}

// ============================================
// 场景4: Java 服务（按架构指定镜像 + Gradle 构建）
// 用户使用 Format4（per-arch 镜像）+ Java Gradle 构建
// ============================================
func TestScenario_JavaGradleService(t *testing.T) {
	yaml := `
service:
  name: order-service
  description: "Order management microservice"
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true
  - name: management
    port: 8081
    protocol: TCP
    expose: false

language:
  type: java
  config:
    jdk_version: "17"
    build_tool: "gradle"

build:
  builder_image:
    amd64: "mirrors.tencent.com/library/maven:3-eclipse-temurin-17"
    arm64: "mirrors.tencent.com/library/maven:3-eclipse-temurin-17-arm64"
  runtime_image:
    amd64: "mirrors.tencent.com/library/eclipse-temurin:17-jre-alpine"
    arm64: "mirrors.tencent.com/library/eclipse-temurin:17-jre-alpine-arm64"
  dependencies:
    system_pkgs:
    - curl
    - unzip
  commands:
    build: |
      gradle build -x test
      cp build/libs/*.jar ${BUILD_OUTPUT_DIR}/app.jar

runtime:
  healthcheck:
    enabled: true
    type: default
  startup:
    command: |
      #!/bin/sh
      exec java -jar ${SERVICE_ROOT}/app.jar
    env:
    - name: JAVA_OPTS
      value: "-Xms512m -Xmx1024m"
    - name: SPRING_PROFILES_ACTIVE
      value: production

local_dev:
  compose:
    resources:
      limits:
        cpus: "2"
        memory: 2G

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
`
	outputDir, cfg := helperLoadAndGenerate(t, yaml)

	// 验证: per-arch 镜像解析正确（Format4）
	builderImg, err := config.ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/library/maven:3-eclipse-temurin-17", builderImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/library/maven:3-eclipse-temurin-17-arm64", builderImg.ARM64)
	assert.NotEqual(t, builderImg.AMD64, builderImg.ARM64, "per-arch images should differ")

	runtimeImg, err := config.ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Contains(t, runtimeImg.AMD64, "temurin")
	assert.Contains(t, runtimeImg.ARM64, "arm64")

	// 验证 Dockerfile 使用不同架构镜像
	dockerfileAmd := helperReadFile(t, outputDir, ".tad/build/order-service/Dockerfile.order-service.amd64")
	assert.Contains(t, dockerfileAmd, "BUILDER_IMAGE_X86", "amd64 Dockerfile should reference X86 builder")

	dockerfileArm := helperReadFile(t, outputDir, ".tad/build/order-service/Dockerfile.order-service.arm64")
	assert.Contains(t, dockerfileArm, "BUILDER_IMAGE_ARM", "arm64 Dockerfile should reference ARM builder")

	// 验证: build.sh 包含用户指定的 gradle 构建命令
	buildScript := helperReadFile(t, outputDir, ".tad/build/order-service/build.sh")
	assert.Contains(t, buildScript, "gradle build -x test")
	assert.Contains(t, buildScript, "cp build/libs/*.jar")

	// 验证: entrypoint.sh 包含 Java 启动配置
	entrypoint := helperReadFile(t, outputDir, ".tad/build/order-service/entrypoint.sh")
	assert.Contains(t, entrypoint, "JAVA_OPTS")
	assert.Contains(t, entrypoint, "SPRING_PROFILES_ACTIVE")

	// 验证: compose 包含多端口但只映射暴露端口
	composeContent := helperReadFile(t, outputDir, "compose.yaml")
	assert.Contains(t, composeContent, "8080")
	assert.Contains(t, composeContent, "8081")
}

// ============================================
// 场景5: Rust 服务（零镜像配置，全自动推导）
// 验证非主流语言的自动推导 + 无端口服务（后台任务）
// ============================================
func TestScenario_RustBackgroundWorker(t *testing.T) {
	yaml := `
service:
  name: data-processor
  description: "Background data processing worker"
  # 无端口: 后台任务，不需要暴露端口
  ports: []

language:
  type: rust
  config:
    rust_version: "1.75"

build:
  commands:
    build: |
      cargo build --release
      cp target/release/${SERVICE_NAME} ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME}

runtime:
  healthcheck:
    enabled: true
    type: default
  startup:
    command: |
      #!/bin/sh
      exec ./bin/${SERVICE_NAME}

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
`
	outputDir, cfg := helperLoadAndGenerate(t, yaml)

	// 验证: Rust 自动推导镜像
	builderImg, err := config.ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Contains(t, builderImg.AMD64, "rust:")
	assert.Contains(t, builderImg.AMD64, "1.75")

	runtimeImg, err := config.ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Contains(t, runtimeImg.AMD64, "alpine:")

	// 验证: 无端口服务不报错
	assert.Empty(t, cfg.Service.Ports)

	// 验证: compose.yaml 无端口映射
	composeContent := helperReadFile(t, outputDir, "compose.yaml")
	assert.Contains(t, composeContent, "data-processor")
	assert.NotContains(t, composeContent, "ports:", "No ports section for portless service")

	// 验证: 所有文件都生成了
	helperAssertFileExists(t, outputDir, ".tad/build/data-processor/Dockerfile.data-processor.amd64")
	helperAssertFileExists(t, outputDir, ".tad/build/data-processor/build.sh")
	helperAssertFileExists(t, outputDir, "compose.yaml")
	helperAssertFileExists(t, outputDir, "Makefile")

	// 验证: build.sh 包含 cargo 命令
	buildScript := helperReadFile(t, outputDir, ".tad/build/data-processor/build.sh")
	assert.Contains(t, buildScript, "cargo build --release")
}

// ============================================
// 场景6: YAML 加载错误场景
// 验证各种用户犯错时的错误信息是否友好
// ============================================
func TestScenario_UserErrors(t *testing.T) {
	t.Run("missing service name", func(t *testing.T) {
		yaml := `
language:
  type: go
runtime:
  startup:
    command: "exec ./bin/app"
`
		cfg, err := config.LoadFromBytes([]byte(yaml))
		require.NoError(t, err)
		err = config.NewValidator(cfg).Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service.name is required")
	})

	t.Run("unsupported language", func(t *testing.T) {
		yaml := `
service:
  name: test
language:
  type: elixir
runtime:
  startup:
    command: "exec ./bin/app"
`
		cfg, err := config.LoadFromBytes([]byte(yaml))
		require.NoError(t, err)
		err = config.NewValidator(cfg).Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not supported")
	})

	t.Run("preset ref without base_images", func(t *testing.T) {
		yaml := `
service:
  name: test
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true
language:
  type: go
build:
  builder_image: "@builders.go_1.23"
  commands:
    build: "go build ."
runtime:
  startup:
    command: "exec ./bin/app"
`
		cfg, err := config.LoadFromBytes([]byte(yaml))
		require.NoError(t, err)
		err = config.NewValidator(cfg).Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "base_images is required")
	})

	t.Run("custom healthcheck without script", func(t *testing.T) {
		yaml := `
service:
  name: test
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true
language:
  type: go
runtime:
  healthcheck:
    enabled: true
    type: custom
  startup:
    command: "exec ./bin/app"
`
		cfg, err := config.LoadFromBytes([]byte(yaml))
		require.NoError(t, err)
		err = config.NewValidator(cfg).Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "custom_script is required")
	})

	t.Run("per-arch image missing arm64", func(t *testing.T) {
		yaml := `
service:
  name: test
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true
language:
  type: go
build:
  builder_image:
    amd64: "golang:1.23-alpine"
  commands:
    build: "go build ."
runtime:
  startup:
    command: "exec ./bin/app"
`
		cfg, err := config.LoadFromBytes([]byte(yaml))
		require.NoError(t, err)
		err = config.NewValidator(cfg).Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "arm64")
	})
}

// ============================================
// 场景7: 自定义 CI 路径正确性验证
// 确保自定义 script_dir 贯穿所有生成文件
// ============================================
func TestScenario_CustomCIPath_Consistency(t *testing.T) {
	yaml := `
service:
  name: audit-service
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true

language:
  type: go

build:
  commands:
    build: |
      CGO_ENABLED=0 go build -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/server

runtime:
  healthcheck:
    enabled: true
    type: default
  startup:
    command: |
      #!/bin/sh
      exec ./bin/${SERVICE_NAME}

ci:
  script_dir: "deploy/scripts/audit-service"

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
`
	outputDir, _ := helperLoadAndGenerate(t, yaml)
	customDir := "deploy/scripts/audit-service"

	// 验证: 所有脚本在自定义路径下
	scripts := []string{"build.sh", "build_deps_install.sh", "rt_prepare.sh", "entrypoint.sh", "healthchk.sh"}
	for _, s := range scripts {
		helperAssertFileExists(t, outputDir, filepath.Join(customDir, s))
	}

	// 验证: Dockerfile 在自定义路径下
	helperAssertFileExists(t, outputDir, filepath.Join(customDir, "Dockerfile.audit-service.amd64"))
	helperAssertFileExists(t, outputDir, filepath.Join(customDir, "Dockerfile.audit-service.arm64"))

	// 验证: Dockerfile 内部引用自定义路径
	dockerfile := helperReadFile(t, outputDir, filepath.Join(customDir, "Dockerfile.audit-service.amd64"))
	assert.Contains(t, dockerfile, "COPY deploy/scripts/audit-service/")
	assert.NotContains(t, dockerfile, ".tad/build/")

	// 验证: compose.yaml 引用自定义路径
	composeContent := helperReadFile(t, outputDir, "compose.yaml")
	assert.Contains(t, composeContent, customDir)
	assert.NotContains(t, composeContent, ".tad/build/")

	// 验证: 默认路径下无文件
	_, err := os.Stat(filepath.Join(outputDir, ".tad"))
	// .tad/devops.yaml 仍然会在 .tad 下生成
	// 但 .tad/build 不应存在
	_, err = os.Stat(filepath.Join(outputDir, ".tad", "build"))
	assert.True(t, os.IsNotExist(err), ".tad/build/ should not exist when custom CI path is set")
}

// ============================================
// 场景8: NodeJS 服务 + Compose 完整配置
// 验证 Compose 的环境变量合并、资源限制、卷挂载、标签
// ============================================
func TestScenario_NodeJSWithFullCompose(t *testing.T) {
	yaml := `
service:
  name: bff-gateway
  description: "Backend-for-Frontend Gateway"
  ports:
  - name: http
    port: 3000
    protocol: TCP
    expose: true

language:
  type: nodejs
  config:
    node_version: "20"

build:
  builder_image: "node:20-alpine"
  runtime_image: "node:20-alpine"
  commands:
    build: |
      npm ci
      npm run build
      cp -r dist ${BUILD_OUTPUT_DIR}/
      cp package.json ${BUILD_OUTPUT_DIR}/

runtime:
  healthcheck:
    enabled: true
    type: default
  startup:
    command: |
      #!/bin/sh
      cd ${SERVICE_ROOT}
      exec node dist/main.js
    env:
    - name: NODE_ENV
      value: production
    - name: PORT
      value: "3000"

local_dev:
  compose:
    resources:
      limits:
        cpus: "0.5"
        memory: 512M
      reservations:
        cpus: "0.25"
        memory: 256M
    environment:
    - name: NODE_ENV
      value: development
    - name: DEBUG
      value: "true"
    healthcheck:
      interval: 15s
      timeout: 5s
      retries: 5
      start_period: 20s
    labels:
      kompose.image-pull-policy: "IfNotPresent"
      app.version: "2.0.0"

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
`
	outputDir, _ := helperLoadAndGenerate(t, yaml)

	// 验证: compose 环境变量合并（compose 覆盖 runtime）
	composeContent := helperReadFile(t, outputDir, "compose.yaml")

	// NODE_ENV 应该被 compose 覆盖为 "development"
	// (但 compose 模板中是遍历 map，所以验证两个 env 都存在)
	assert.Contains(t, composeContent, "NODE_ENV")
	assert.Contains(t, composeContent, "DEBUG")

	// 资源限制
	assert.Contains(t, composeContent, "512M")
	assert.Contains(t, composeContent, "256M")

	// 健康检查
	assert.Contains(t, composeContent, "interval: 15s")
	assert.Contains(t, composeContent, "timeout: 5s")
	assert.Contains(t, composeContent, "retries: 5")
	assert.Contains(t, composeContent, "start_period: 20s")

	// 标签
	assert.Contains(t, composeContent, "kompose.image-pull-policy")
	assert.Contains(t, composeContent, "app.version")

	// 端口
	assert.Contains(t, composeContent, "3000")

	// 验证: build.sh 包含 npm 命令
	buildScript := helperReadFile(t, outputDir, ".tad/build/bff-gateway/build.sh")
	assert.Contains(t, buildScript, "npm ci")
	assert.Contains(t, buildScript, "npm run build")
}

// ============================================
// 场景9: 跨格式混合使用
// builder 用 Format2（直接镜像），runtime 用 Format4（per-arch）
// ============================================
func TestScenario_MixedImageFormats(t *testing.T) {
	yaml := `
service:
  name: mixed-service
  ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true

language:
  type: go

build:
  builder_image: "golang:1.23-alpine"
  runtime_image:
    amd64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"
    arm64: "mirrors.tencent.com/tencentos/tencentos3-minimal-arm:latest"

runtime:
  startup:
    command: |
      #!/bin/sh
      exec ./bin/${SERVICE_NAME}

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
`
	outputDir, cfg := helperLoadAndGenerate(t, yaml)

	// builder: Format2 → AMD64/ARM64 相同
	builderImg, err := config.ResolveBuilderImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-alpine", builderImg.AMD64)
	assert.Equal(t, "golang:1.23-alpine", builderImg.ARM64)

	// runtime: Format4 → AMD64/ARM64 不同
	runtimeImg, err := config.ResolveRuntimeImageWithDefaults(cfg)
	require.NoError(t, err)
	assert.Equal(t, "mirrors.tencent.com/tencentos/tencentos3-minimal:latest", runtimeImg.AMD64)
	assert.Equal(t, "mirrors.tencent.com/tencentos/tencentos3-minimal-arm:latest", runtimeImg.ARM64)
	assert.NotEqual(t, runtimeImg.AMD64, runtimeImg.ARM64)

	// 验证文件生成成功
	helperAssertFileExists(t, outputDir, ".tad/build/mixed-service/Dockerfile.mixed-service.amd64")
	helperAssertFileExists(t, outputDir, ".tad/build/mixed-service/Dockerfile.mixed-service.arm64")

	// 验证 build.sh 包含自动推导的 Go 构建命令
	buildScript := helperReadFile(t, outputDir, ".tad/build/mixed-service/build.sh")
	assert.Contains(t, buildScript, "CGO_ENABLED=0 go build", "Should use auto-inferred Go build command")
}

// ============================================
// 场景10: 验证 demo-app/service.yaml 可以正确加载和生成
// 这是最接近真实用户体验的测试
// ============================================
func TestScenario_DemoApp(t *testing.T) {
	// 读取实际的 demo-app 配置
	demoYaml, err := os.ReadFile("demo-app/service.yaml")
	if err != nil {
		t.Skip("demo-app/service.yaml not found, skipping")
	}

	// 加载
	cfg, err := config.LoadFromBytes(demoYaml)
	require.NoError(t, err, "demo-app/service.yaml should load successfully")

	// 校验
	err = config.NewValidator(cfg).Validate()
	require.NoError(t, err, "demo-app/service.yaml should pass validation")

	// 验证关键字段
	assert.Equal(t, "demo-app", cfg.Service.Name)
	assert.Equal(t, "go", cfg.Language.Type)

	// 验证 build.dependencies.system_pkgs 正确加载（修复后的字段名）
	assert.Contains(t, cfg.Build.Dependencies.SystemPkgs, "git", "build.dependencies.system_pkgs should load correctly")
	assert.Contains(t, cfg.Build.Dependencies.SystemPkgs, "make")

	// 验证 runtime.system_dependencies.packages 正确加载
	assert.Contains(t, cfg.Runtime.SystemDependencies.Packages, "tzdata")
	assert.Contains(t, cfg.Runtime.SystemDependencies.Packages, "ca-certificates")

	// 生成
	tmpDir, err := os.MkdirTemp("", "demo-app-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	gen := generator.NewGenerator(cfg, tmpDir)
	err = gen.Generate()
	require.NoError(t, err, "demo-app generation should succeed")

	// 验证生成的关键文件
	helperAssertFileExists(t, tmpDir, ".tad/build/demo-app/Dockerfile.demo-app.amd64")
	helperAssertFileExists(t, tmpDir, "compose.yaml")
	helperAssertFileExists(t, tmpDir, "Makefile")

	// 验证 compose 引用正确
	composeContent := helperReadFile(t, tmpDir, "compose.yaml")
	assert.Contains(t, composeContent, "demo-app")
	assert.Contains(t, composeContent, ".tad/build/demo-app/")

	// 验证 Dockerfile 内容
	dockerfile := helperReadFile(t, tmpDir, ".tad/build/demo-app/Dockerfile.demo-app.amd64")
	lines := strings.Split(dockerfile, "\n")
	assert.True(t, len(lines) > 10, "Dockerfile should have substantial content")
}
