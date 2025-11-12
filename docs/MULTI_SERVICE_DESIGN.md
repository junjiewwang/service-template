# 多服务配置设计方案

> 支持多服务配置和 YAML 变量引用机制

**设计日期**: 2025-11-12  
**版本**: 1.0  
**状态**: 设计阶段

---

## 📋 目录

1. [需求分析](#需求分析)
2. [设计目标](#设计目标)
3. [方案对比](#方案对比)
4. [推荐方案](#推荐方案)
5. [配置结构设计](#配置结构设计)
6. [变量引用机制](#变量引用机制)
7. [实现细节](#实现细节)
8. [迁移指南](#迁移指南)
9. [最佳实践](#最佳实践)

---

## 🎯 需求分析

### 当前问题

1. **单服务限制**: 当前配置只支持单个服务，无法在一个配置文件中定义多个服务
2. **镜像配置重复**: 每个服务都需要重复配置 `builder_image` 和 `runtime_image`
3. **缺少变量引用**: 无法定义全局变量并在多处引用，导致配置冗余

### 业务场景

#### 场景 1: 微服务项目
```
project/
├── service-a/  # API 服务
├── service-b/  # Worker 服务
├── service-c/  # Admin 服务
└── service.yaml  # 统一配置
```

所有服务使用相同的：
- 构建镜像（builder_image）
- 运行时镜像（runtime_image）
- 语言版本
- 插件配置

#### 场景 2: 多环境部署
```yaml
# 开发环境使用一套镜像
# 生产环境使用另一套镜像
# 通过变量引用统一管理
```

#### 场景 3: 组织级标准化
```yaml
# 组织统一定义标准镜像
# 各团队引用标准配置
# 减少配置错误和不一致
```

---

## 🎯 设计目标

### 核心目标

1. ✅ **支持多服务配置**: 在一个 YAML 文件中定义多个服务
2. ✅ **提取公共配置**: 将 `builder_image`、`runtime_image` 等提升到顶级
3. ✅ **变量引用机制**: 支持 YAML 锚点和别名，或自定义变量系统
4. ✅ **向后兼容**: 保持对现有单服务配置的兼容
5. ✅ **易于理解**: 配置结构清晰，学习成本低

### 非功能目标

- 🔒 **类型安全**: 变量引用有类型检查
- 📝 **良好的错误提示**: 引用不存在的变量时给出清晰错误
- 🧪 **可测试**: 配置加载和变量解析可单独测试
- 📚 **文档完善**: 提供详细的使用文档和示例

---

## 🔄 方案对比

### 方案 1: YAML 锚点和别名（原生 YAML 特性）

#### 优点
- ✅ YAML 原生支持，无需额外实现
- ✅ 标准化，开发者熟悉
- ✅ 工具链支持好（编辑器、linter）

#### 缺点
- ❌ 只能在同一文件内引用
- ❌ 语法相对复杂（`&anchor` 和 `*alias`）
- ❌ 不支持跨文件引用
- ❌ 不支持变量计算和转换

#### 示例
```yaml
# 定义锚点
x-images: &common-images
  builder_image:
    amd64: "mirrors.tencent.com/tcs-infra/builder:amd64"
    arm64: "mirrors.tencent.com/tcs-infra/builder:arm64"
  runtime_image:
    amd64: "mirrors.tencent.com/tencentos/runtime:latest"
    arm64: "mirrors.tencent.com/tencentos/runtime:latest"

# 引用锚点
services:
  - name: service-a
    build:
      <<: *common-images
      commands:
        build: "go build"
```

---

### 方案 2: 自定义变量系统（类似 Helm Values）

#### 优点
- ✅ 灵活强大，支持复杂逻辑
- ✅ 支持变量计算和转换
- ✅ 可以跨文件引用
- ✅ 更好的错误提示

#### 缺点
- ❌ 需要自己实现解析器
- ❌ 增加学习成本
- ❌ 工具链支持需要自己开发

#### 示例
```yaml
# 定义变量
vars:
  common_builder_amd64: "mirrors.tencent.com/tcs-infra/builder:amd64"
  common_builder_arm64: "mirrors.tencent.com/tcs-infra/builder:arm64"
  common_runtime_amd64: "mirrors.tencent.com/tencentos/runtime:latest"
  common_runtime_arm64: "mirrors.tencent.com/tencentos/runtime:latest"

# 引用变量
services:
  - name: service-a
    build:
      builder_image:
        amd64: ${vars.common_builder_amd64}
        arm64: ${vars.common_builder_arm64}
```

---

### 方案 3: 混合方案（推荐）⭐

结合两种方案的优点：
- 使用 **YAML 锚点** 处理结构化配置（如镜像配置）
- 使用 **自定义变量** 处理简单值引用（如路径、版本号）

#### 优点
- ✅ 充分利用 YAML 原生特性
- ✅ 保持灵活性
- ✅ 学习成本适中
- ✅ 实现成本可控

#### 缺点
- ⚠️ 需要理解两种机制

---

## 🏆 推荐方案：混合方案

### 方案概述

采用**三层配置结构**：

```
┌─────────────────────────────────────┐
│  1. Global Config (全局配置)        │
│     - 变量定义 (vars)               │
│     - 默认值 (defaults)             │
│     - 共享配置 (shared)             │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  2. Service Definitions (服务定义)  │
│     - 多个服务配置                  │
│     - 引用全局配置                  │
│     - 服务特定配置                  │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  3. Metadata (元数据)               │
│     - 版本信息                      │
│     - 生成信息                      │
└─────────────────────────────────────┘
```

---

## 📐 配置结构设计

### 新的配置文件结构

```yaml
# ============================================
# Multi-Service Configuration
# Version: 3.0
# ============================================

# ============================================
# 1. 全局变量定义（可选）
# ============================================
vars:
  # 简单值变量
  go_version: "1.23"
  python_version: "3.11"
  deploy_base_dir: "/usr/local/services"
  plugin_install_dir: "/tce"
  
  # 镜像仓库前缀
  image_registry: "mirrors.tencent.com"
  
  # 构建镜像版本
  builder_version: "v1.0.0"
  runtime_version: "latest"

# ============================================
# 2. 默认配置（可选）
# ============================================
defaults:
  # 默认语言配置
  language:
    type: go
    version: ${vars.go_version}
    config:
      goproxy: "https://goproxy.cn,direct"
      gosumdb: "sum.golang.org"
  
  # 默认构建配置
  build:
    builder_image: &default-builder-image
      amd64: "${vars.image_registry}/tcs-infra/tceforqci_x86_go23:${vars.builder_version}"
      arm64: "${vars.image_registry}/tcs-infra/tceforqci_arm_go23:${vars.builder_version}"
    
    runtime_image: &default-runtime-image
      amd64: "${vars.image_registry}/tencentos/tencentos3-minimal:${vars.runtime_version}"
      arm64: "${vars.image_registry}/tencentos/tencentos3-minimal:${vars.runtime_version}"
    
    system_dependencies:
      packages: &default-build-packages
        - git
        - make
        - gcc
  
  # 默认运行时配置
  runtime:
    system_dependencies:
      packages: &default-runtime-packages
        - ca-certificates
        - tzdata
    
    healthcheck:
      enabled: true
      type: default
  
  # 默认插件配置
  plugins: &default-plugins
    install_dir: ${vars.plugin_install_dir}
    items:
      - name: selfMonitor
        description: "TCE Self Monitor Tool"
        download_url: "https://mirrors.tencent.com/repository/generic/selfMonitor/download_tool.sh"
        install_command: |
          curl -fsSL "${PLUGIN_DOWNLOAD_URL}" | bash -s "${PLUGIN_WORK_DIR}"
        runtime_env:
          - name: TCESTAURY_TOOL_PATH
            value: ${PLUGIN_INSTALL_DIR}
        required: true

# ============================================
# 3. 共享配置（使用 YAML 锚点）
# ============================================
shared:
  # 共享的构建镜像配置
  images:
    go_builder: &go-builder-image
      amd64: "${vars.image_registry}/tcs-infra/tceforqci_x86_go23:${vars.builder_version}"
      arm64: "${vars.image_registry}/tcs-infra/tceforqci_arm_go23:${vars.builder_version}"
    
    python_builder: &python-builder-image
      amd64: "${vars.image_registry}/tcs-infra/python_builder:${vars.builder_version}"
      arm64: "${vars.image_registry}/tcs-infra/python_builder:${vars.builder_version}"
    
    common_runtime: &common-runtime-image
      amd64: "${vars.image_registry}/tencentos/tencentos3-minimal:${vars.runtime_version}"
      arm64: "${vars.image_registry}/tencentos/tencentos3-minimal:${vars.runtime_version}"
  
  # 共享的端口配置
  ports:
    http_8080: &port-http-8080
      - name: http
        port: 8080
        protocol: TCP
        expose: true
        description: "HTTP API port"
    
    metrics_9090: &port-metrics-9090
      - name: metrics
        port: 9090
        protocol: TCP
        expose: false
        description: "Prometheus metrics"

# ============================================
# 4. 服务定义（多服务）
# ============================================
services:
  # -------------------- Service A --------------------
  - name: api-service
    description: "API Service"
    
    # 引用共享端口配置
    ports:
      - <<: *port-http-8080
      - <<: *port-metrics-9090
    
    deploy_dir: ${vars.deploy_base_dir}
    
    # 继承默认语言配置（可覆盖）
    language:
      type: go
      version: ${vars.go_version}
      config:
        goproxy: "https://goproxy.cn,direct"
    
    build:
      dependency_files:
        auto_detect: true
      
      # 引用共享镜像配置
      builder_image: *go-builder-image
      runtime_image: *common-runtime-image
      
      # 引用共享系统依赖
      system_dependencies:
        packages: *default-build-packages
      
      commands:
        build: |
          CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/api
    
    # 引用共享插件配置
    plugins: *default-plugins
    
    runtime:
      system_dependencies:
        packages: *default-runtime-packages
      
      healthcheck:
        enabled: true
        type: default
      
      startup:
        command: |
          #!/bin/sh
          cd ${SERVICE_ROOT}
          exec ./bin/${SERVICE_NAME}
    
    local_dev:
      compose:
        resources:
          limits:
            cpus: "0.5"
            memory: 1G
        volumes:
          - source: ./configs/api-config.yaml
            target: ${SERVICE_ROOT}/config.yaml
            type: bind
      
      kubernetes:
        enabled: true
        namespace: default
  
  # -------------------- Service B --------------------
  - name: worker-service
    description: "Background Worker Service"
    
    ports:
      - name: metrics
        port: 9091
        protocol: TCP
        expose: false
    
    deploy_dir: ${vars.deploy_base_dir}
    
    language:
      type: go
      version: ${vars.go_version}
    
    build:
      dependency_files:
        auto_detect: true
      
      # 引用相同的镜像配置
      builder_image: *go-builder-image
      runtime_image: *common-runtime-image
      
      system_dependencies:
        packages: *default-build-packages
      
      commands:
        build: |
          CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/worker
    
    plugins: *default-plugins
    
    runtime:
      system_dependencies:
        packages: *default-runtime-packages
      
      healthcheck:
        enabled: true
        type: default
      
      startup:
        command: |
          #!/bin/sh
          cd ${SERVICE_ROOT}
          exec ./bin/${SERVICE_NAME}
    
    local_dev:
      compose:
        resources:
          limits:
            cpus: "0.25"
            memory: 512M
      kubernetes:
        enabled: true
        namespace: default
  
  # -------------------- Service C (Python) --------------------
  - name: admin-service
    description: "Admin Dashboard Service"
    
    ports:
      - name: http
        port: 8000
        protocol: TCP
        expose: true
    
    deploy_dir: ${vars.deploy_base_dir}
    
    # 使用不同的语言
    language:
      type: python
      version: ${vars.python_version}
      config:
        pip_index_url: "https://mirrors.tencent.com/pypi/simple"
    
    build:
      dependency_files:
        auto_detect: true
      
      # 使用 Python 构建镜像
      builder_image: *python-builder-image
      runtime_image: *common-runtime-image
      
      system_dependencies:
        packages:
          - python3
          - pip
      
      commands:
        build: |
          pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/lib
          cp -r app ${BUILD_OUTPUT_DIR}/
    
    runtime:
      system_dependencies:
        packages:
          - python3
          - ca-certificates
      
      healthcheck:
        enabled: true
        type: custom
        custom_script: |
          curl -f http://localhost:8000/health || exit 1
      
      startup:
        command: |
          #!/bin/sh
          cd ${SERVICE_ROOT}
          export PYTHONPATH=${SERVICE_ROOT}/lib
          exec python3 app/main.py

# ============================================
# 5. 元数据
# ============================================
metadata:
  template_version: "3.0.0"
  generated_at: ""
  generator: "svcgen"
```

---

## 🔧 变量引用机制

### 1. 简单变量引用（自定义实现）

#### 语法
```yaml
${vars.variable_name}
```

#### 支持的位置
- ✅ 字符串值
- ✅ 数组元素
- ✅ 对象属性值
- ❌ 键名（不支持）

#### 示例
```yaml
vars:
  base_dir: "/usr/local"
  service_name: "my-service"

services:
  - name: ${vars.service_name}
    deploy_dir: ${vars.base_dir}/services
```

#### 解析规则
1. **递归解析**: 变量可以引用其他变量
2. **循环检测**: 检测并报错循环引用
3. **类型保持**: 解析后保持原始类型
4. **默认值**: 支持 `${vars.name:default}` 语法

---

### 2. YAML 锚点和别名（原生支持）

#### 语法
```yaml
# 定义锚点
key: &anchor-name
  field: value

# 引用别名
other_key: *anchor-name

# 合并引用
another_key:
  <<: *anchor-name
  additional_field: value
```

#### 使用场景
- ✅ 复杂对象引用（如镜像配置）
- ✅ 数组引用（如端口列表）
- ✅ 配置模板复用

#### 示例
```yaml
# 定义镜像配置锚点
shared:
  images:
    go_builder: &go-builder
      amd64: "builder:amd64"
      arm64: "builder:arm64"

# 引用镜像配置
services:
  - name: service-a
    build:
      builder_image: *go-builder  # 完全引用
  
  - name: service-b
    build:
      builder_image:
        <<: *go-builder  # 合并引用
        amd64: "custom:amd64"  # 覆盖特定字段
```

---

### 3. 变量作用域

```
┌─────────────────────────────────────┐
│  Global Scope (全局作用域)          │
│  - vars.*                           │
│  - defaults.*                       │
│  - shared.*                         │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  Service Scope (服务作用域)         │
│  - 继承全局变量                     │
│  - 可覆盖全局配置                   │
│  - 服务特定变量                     │
└─────────────────────────────────────┘
```

---

## 💻 实现细节

### 1. 配置结构定义

```go
// pkg/config/types.go

// MultiServiceConfig 多服务配置（新）
type MultiServiceConfig struct {
    // 变量定义
    Vars map[string]interface{} `yaml:"vars,omitempty"`
    
    // 默认配置
    Defaults *DefaultsConfig `yaml:"defaults,omitempty"`
    
    // 共享配置（用于 YAML 锚点）
    Shared map[string]interface{} `yaml:"shared,omitempty"`
    
    // 服务列表
    Services []ServiceConfig `yaml:"services"`
    
    // 元数据
    Metadata MetadataConfig `yaml:"metadata"`
}

// DefaultsConfig 默认配置
type DefaultsConfig struct {
    Language *LanguageConfig `yaml:"language,omitempty"`
    Build    *BuildDefaults  `yaml:"build,omitempty"`
    Runtime  *RuntimeConfig  `yaml:"runtime,omitempty"`
    Plugins  *PluginsConfig  `yaml:"plugins,omitempty"`
}

// BuildDefaults 构建默认配置
type BuildDefaults struct {
    BuilderImage       *ArchImageConfig               `yaml:"builder_image,omitempty"`
    RuntimeImage       *ArchImageConfig               `yaml:"runtime_image,omitempty"`
    SystemDependencies *BuildSystemDependenciesConfig `yaml:"system_dependencies,omitempty"`
}

// ServiceConfig 保持现有结构，但字段变为可选
type ServiceConfig struct {
    // 基础信息（必需）
    Name        string `yaml:"name"`
    Description string `yaml:"description,omitempty"`
    
    // 其他字段变为可选，可从 defaults 继承
    Ports       []PortConfig    `yaml:"ports,omitempty"`
    DeployDir   string          `yaml:"deploy_dir,omitempty"`
    Language    *LanguageConfig `yaml:"language,omitempty"`
    Build       *BuildConfig    `yaml:"build,omitempty"`
    Plugins     *PluginsConfig  `yaml:"plugins,omitempty"`
    Runtime     *RuntimeConfig  `yaml:"runtime,omitempty"`
    LocalDev    *LocalDevConfig `yaml:"local_dev,omitempty"`
    Makefile    *MakefileConfig `yaml:"makefile,omitempty"`
    CI          *CIConfig       `yaml:"ci,omitempty"`
}
```

---

### 2. 配置加载器

```go
// pkg/config/loader.go

// ConfigLoader 配置加载器
type ConfigLoader struct {
    configPath string
}

// Load 加载配置（自动检测单服务或多服务）
func (l *ConfigLoader) Load() (interface{}, error) {
    data, err := os.ReadFile(l.configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    // 尝试解析为多服务配置
    var multiConfig MultiServiceConfig
    if err := yaml.Unmarshal(data, &multiConfig); err == nil {
        if len(multiConfig.Services) > 0 {
            // 多服务配置
            return l.processMultiServiceConfig(&multiConfig)
        }
    }
    
    // 回退到单服务配置（向后兼容）
    var singleConfig ServiceConfig
    if err := yaml.Unmarshal(data, &singleConfig); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &singleConfig, nil
}

// processMultiServiceConfig 处理多服务配置
func (l *ConfigLoader) processMultiServiceConfig(config *MultiServiceConfig) (*MultiServiceConfig, error) {
    // 1. 解析变量
    if err := l.resolveVariables(config); err != nil {
        return nil, fmt.Errorf("failed to resolve variables: %w", err)
    }
    
    // 2. 应用默认配置
    if err := l.applyDefaults(config); err != nil {
        return nil, fmt.Errorf("failed to apply defaults: %w", err)
    }
    
    // 3. 验证配置
    if err := l.validate(config); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    return config, nil
}
```

---

### 3. 变量解析器

```go
// pkg/config/variable_resolver.go

// VariableResolver 变量解析器
type VariableResolver struct {
    vars     map[string]interface{}
    resolved map[string]interface{}
    visiting map[string]bool // 用于循环检测
}

// NewVariableResolver 创建变量解析器
func NewVariableResolver(vars map[string]interface{}) *VariableResolver {
    return &VariableResolver{
        vars:     vars,
        resolved: make(map[string]interface{}),
        visiting: make(map[string]bool),
    }
}

// Resolve 解析变量引用
func (r *VariableResolver) Resolve(value interface{}) (interface{}, error) {
    switch v := value.(type) {
    case string:
        return r.resolveString(v)
    case map[string]interface{}:
        return r.resolveMap(v)
    case []interface{}:
        return r.resolveArray(v)
    default:
        return value, nil
    }
}

// resolveString 解析字符串中的变量引用
func (r *VariableResolver) resolveString(s string) (string, error) {
    // 匹配 ${vars.name} 或 ${vars.name:default}
    re := regexp.MustCompile(`\$\{vars\.([^:}]+)(?::([^}]+))?\}`)
    
    result := re.ReplaceAllStringFunc(s, func(match string) string {
        matches := re.FindStringSubmatch(match)
        varName := matches[1]
        defaultValue := matches[2]
        
        // 检查循环引用
        if r.visiting[varName] {
            return fmt.Sprintf("ERROR:CIRCULAR_REFERENCE:%s", varName)
        }
        
        // 获取变量值
        value, exists := r.vars[varName]
        if !exists {
            if defaultValue != "" {
                return defaultValue
            }
            return fmt.Sprintf("ERROR:UNDEFINED_VAR:%s", varName)
        }
        
        // 递归解析
        r.visiting[varName] = true
        resolved, err := r.Resolve(value)
        delete(r.visiting, varName)
        
        if err != nil {
            return fmt.Sprintf("ERROR:%v", err)
        }
        
        return fmt.Sprint(resolved)
    })
    
    // 检查是否有错误
    if strings.Contains(result, "ERROR:") {
        return "", fmt.Errorf("variable resolution failed: %s", result)
    }
    
    return result, nil
}

// resolveMap 解析 map 中的变量引用
func (r *VariableResolver) resolveMap(m map[string]interface{}) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    for k, v := range m {
        resolved, err := r.Resolve(v)
        if err != nil {
            return nil, err
        }
        result[k] = resolved
    }
    return result, nil
}

// resolveArray 解析数组中的变量引用
func (r *VariableResolver) resolveArray(arr []interface{}) ([]interface{}, error) {
    result := make([]interface{}, len(arr))
    for i, v := range arr {
        resolved, err := r.Resolve(v)
        if err != nil {
            return nil, err
        }
        result[i] = resolved
    }
    return result, nil
}
```

---

### 4. 默认配置应用器

```go
// pkg/config/defaults_applier.go

// DefaultsApplier 默认配置应用器
type DefaultsApplier struct {
    defaults *DefaultsConfig
}

// NewDefaultsApplier 创建默认配置应用器
func NewDefaultsApplier(defaults *DefaultsConfig) *DefaultsApplier {
    return &DefaultsApplier{defaults: defaults}
}

// Apply 应用默认配置到服务
func (a *DefaultsApplier) Apply(service *ServiceConfig) error {
    if a.defaults == nil {
        return nil
    }
    
    // 应用语言默认配置
    if service.Language == nil && a.defaults.Language != nil {
        service.Language = a.defaults.Language
    }
    
    // 应用构建默认配置
    if service.Build != nil && a.defaults.Build != nil {
        if service.Build.BuilderImage.AMD64 == "" && a.defaults.Build.BuilderImage != nil {
            service.Build.BuilderImage = *a.defaults.Build.BuilderImage
        }
        if service.Build.RuntimeImage.AMD64 == "" && a.defaults.Build.RuntimeImage != nil {
            service.Build.RuntimeImage = *a.defaults.Build.RuntimeImage
        }
    }
    
    // 应用插件默认配置
    if service.Plugins == nil && a.defaults.Plugins != nil {
        service.Plugins = a.defaults.Plugins
    }
    
    return nil
}
```

---

## 🔄 迁移指南

### 从单服务配置迁移到多服务配置

#### 步骤 1: 提取公共配置

**原配置** (service.yaml):
```yaml
service:
  name: my-service
  
language:
  type: go
  version: "1.23"

build:
  builder_image:
    amd64: "mirrors.tencent.com/builder:amd64"
    arm64: "mirrors.tencent.com/builder:arm64"
  runtime_image:
    amd64: "mirrors.tencent.com/runtime:latest"
    arm64: "mirrors.tencent.com/runtime:latest"
```

**新配置** (services.yaml):
```yaml
# 1. 提取变量
vars:
  image_registry: "mirrors.tencent.com"
  go_version: "1.23"

# 2. 定义默认配置
defaults:
  language:
    type: go
    version: ${vars.go_version}
  
  build:
    builder_image: &default-builder
      amd64: "${vars.image_registry}/builder:amd64"
      arm64: "${vars.image_registry}/builder:arm64"
    runtime_image: &default-runtime
      amd64: "${vars.image_registry}/runtime:latest"
      arm64: "${vars.image_registry}/runtime:latest"

# 3. 定义服务（继承默认配置）
services:
  - name: my-service
    # 其他配置...
```

#### 步骤 2: 添加新服务

```yaml
services:
  - name: my-service
    # 原有配置
  
  - name: new-service
    description: "New Service"
    # 自动继承 defaults 中的配置
    # 只需配置差异部分
    build:
      commands:
        build: "go build ./cmd/new-service"
```

---

### 向后兼容性

#### 单服务配置仍然支持

```yaml
# 旧的单服务配置格式仍然有效
service:
  name: my-service

language:
  type: go

build:
  builder_image:
    amd64: "..."
```

#### 自动检测机制

```go
// 加载器自动检测配置类型
func (l *ConfigLoader) Load() (interface{}, error) {
    // 1. 尝试解析为多服务配置
    // 2. 如果失败，回退到单服务配置
    // 3. 返回对应的配置对象
}
```

---

## 📚 最佳实践

### 1. 变量命名规范

```yaml
vars:
  # ✅ 好的命名：清晰、有意义
  go_version: "1.23"
  image_registry: "mirrors.tencent.com"
  deploy_base_dir: "/usr/local/services"
  
  # ❌ 不好的命名：模糊、缩写
  v: "1.23"
  reg: "mirrors.tencent.com"
  dir: "/usr/local/services"
```

### 2. 锚点命名规范

```yaml
shared:
  images:
    # ✅ 好的锚点名：描述性强
    go_builder: &go-builder-image
      amd64: "..."
    
    # ❌ 不好的锚点名：过于简短
    gb: &gb
      amd64: "..."
```

### 3. 配置组织建议

```yaml
# 推荐的配置组织顺序：
# 1. vars - 变量定义
# 2. defaults - 默认配置
# 3. shared - 共享配置（锚点）
# 4. services - 服务定义
# 5. metadata - 元数据
```

### 4. 何时使用变量 vs 锚点

| 场景 | 推荐方式 | 原因 |
|------|---------|------|
| 简单字符串值 | 变量 `${vars.name}` | 更直观 |
| 复杂对象 | 锚点 `*anchor` | YAML 原生支持 |
| 需要覆盖部分字段 | 锚点 + 合并 `<<: *anchor` | 灵活性高 |
| 跨文件引用 | 变量（未来支持） | 扩展性好 |

### 5. 配置验证

```yaml
# 使用工具验证配置
$ svcgen validate services.yaml

# 输出：
✓ Configuration is valid
✓ Variables resolved: 12
✓ Services defined: 3
  - api-service
  - worker-service
  - admin-service
✓ All services inherit defaults correctly
```

---

## 🧪 测试策略

### 1. 单元测试

```go
// 测试变量解析
func TestVariableResolver_Resolve(t *testing.T) {
    tests := []struct {
        name     string
        vars     map[string]interface{}
        input    string
        expected string
        wantErr  bool
    }{
        {
            name: "simple variable",
            vars: map[string]interface{}{"version": "1.23"},
            input: "go:${vars.version}",
            expected: "go:1.23",
        },
        {
            name: "nested variable",
            vars: map[string]interface{}{
                "base": "mirrors.tencent.com",
                "image": "${vars.base}/builder",
            },
            input: "${vars.image}:latest",
            expected: "mirrors.tencent.com/builder:latest",
        },
        {
            name: "circular reference",
            vars: map[string]interface{}{
                "a": "${vars.b}",
                "b": "${vars.a}",
            },
            input: "${vars.a}",
            wantErr: true,
        },
    }
    // ...
}
```

### 2. 集成测试

```go
// 测试完整配置加载
func TestConfigLoader_LoadMultiService(t *testing.T) {
    loader := NewConfigLoader("testdata/multi-service.yaml")
    config, err := loader.Load()
    
    assert.NoError(t, err)
    assert.NotNil(t, config)
    
    multiConfig := config.(*MultiServiceConfig)
    assert.Len(t, multiConfig.Services, 3)
    
    // 验证变量解析
    assert.Equal(t, "1.23", multiConfig.Services[0].Language.Version)
    
    // 验证默认配置应用
    assert.NotNil(t, multiConfig.Services[0].Build.BuilderImage)
}
```

---

## 📊 实施计划

### 阶段 1: 基础设施 (3-5 天)

- [ ] **任务 1.1**: 定义新的配置结构
  - 创建 `MultiServiceConfig` 类型
  - 创建 `DefaultsConfig` 类型
  - 更新 `ServiceConfig` 使字段可选

- [ ] **任务 1.2**: 实现变量解析器
  - 创建 `VariableResolver`
  - 支持 `${vars.name}` 语法
  - 支持默认值 `${vars.name:default}`
  - 循环引用检测

- [ ] **任务 1.3**: 实现默认配置应用器
  - 创建 `DefaultsApplier`
  - 实现配置继承逻辑
  - 支持部分覆盖

### 阶段 2: 配置加载 (2-3 天)

- [ ] **任务 2.1**: 更新配置加载器
  - 支持多服务配置加载
  - 保持单服务配置兼容
  - 自动检测配置类型

- [ ] **任务 2.2**: 集成 YAML 锚点支持
  - 验证 YAML 库支持锚点
  - 测试锚点和变量混合使用

### 阶段 3: 生成器适配 (3-5 天)

- [ ] **任务 3.1**: 更新生成器接口
  - 支持多服务生成
  - 为每个服务生成独立文件

- [ ] **任务 3.2**: 更新现有生成器
  - Dockerfile 生成器
  - Compose 生成器
  - Makefile 生成器
  - 脚本生成器

### 阶段 4: 测试和文档 (2-3 天)

- [ ] **任务 4.1**: 编写测试
  - 单元测试
  - 集成测试
  - 端到端测试

- [ ] **任务 4.2**: 编写文档
  - 用户指南
  - 迁移指南
  - API 文档
  - 示例配置

---

## 🎯 总结

### 核心特性

1. ✅ **多服务支持**: 一个配置文件定义多个服务
2. ✅ **变量系统**: `${vars.name}` 语法，支持默认值
3. ✅ **YAML 锚点**: 原生支持，用于复杂对象引用
4. ✅ **默认配置**: 减少重复，提高一致性
5. ✅ **向后兼容**: 单服务配置仍然有效

### 技术优势

- 🔒 **类型安全**: Go 结构体保证类型正确
- 📝 **错误提示**: 清晰的变量解析错误信息
- 🧪 **可测试**: 各组件独立可测试
- 📚 **易维护**: 代码结构清晰，职责分明

### 用户体验

- 🎨 **灵活性**: 支持多种引用方式
- 📖 **易学习**: 语法简单，文档完善
- 🚀 **高效率**: 减少配置重复，提高开发效率
- 🔧 **可扩展**: 易于添加新特性

---

**设计完成，等待实施决策** 🎉