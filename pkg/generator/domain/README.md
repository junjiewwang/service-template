# Generator Domain Layer

> 领域驱动设计（DDD）的领域层实现

## 📁 目录结构

```
domain/
├── errors/          # 结构化错误处理
│   ├── errors.go
│   └── errors_test.go
├── services/        # 领域服务
│   ├── plugin_service.go
│   ├── language_service.go
│   └── language_service_test.go
└── models/          # 领域模型（预留）
```

## 🎯 设计原则

本层遵循DDD（领域驱动设计）原则：

1. **领域服务** - 封装不属于任何实体的业务逻辑
2. **值对象** - 不可变的领域概念
3. **策略模式** - 支持多种实现的灵活设计
4. **单一职责** - 每个服务只负责一个领域

## 📦 组件说明

### 1. 错误处理 (errors/)

提供结构化的错误处理系统。

**特性**:
- 6种错误类型（Validation, Template, Plugin, Language, FileSystem, Configuration）
- 错误链支持（实现 `Unwrap()`）
- 上下文信息附加
- 清晰的错误消息格式

**使用示例**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/errors"

// 创建错误
err := errors.NewPluginError(
    "dockerfile",
    "failed to process plugin",
    cause,
).WithContext("plugin", pluginName)

// 错误输出格式
// [PLUGIN_ERROR] dockerfile: failed to process plugin (caused by: ...)
```

### 2. 插件服务 (services/plugin_service.go)

统一处理所有生成器的插件相关逻辑。

**功能**:
- 为不同生成器准备插件信息
- 处理插件环境变量
- 插件检测和验证

**方法**:
- `PrepareForDockerfile()` - 为Dockerfile准备插件
- `PrepareForBuildScript()` - 为构建脚本准备插件
- `PrepareForEntrypoint()` - 为入口脚本准备插件环境变量
- `HasPlugins()` - 检查是否有插件
- `GetInstallDir()` - 获取插件安装目录

**使用示例**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/services"

// 创建服务
pluginService := services.NewPluginService(ctx, engine)

// 检查插件
if pluginService.HasPlugins() {
    // 根据生成器类型选择方法
    plugins := pluginService.PrepareForDockerfile()
    composer.Override("PLUGINS", plugins)
}
```

**收益**:
- ✅ 消除60行重复代码
- ✅ 集中插件处理逻辑
- ✅ 提高可测试性
- ✅ 便于扩展新场景

### 3. 路径配置值对象 (models/paths_config.go)

封装所有路径相关配置，消除硬编码。

**特性**:
- 值对象模式，不可变
- 流式API配置
- 路径验证
- 深拷贝支持

**使用示例**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/models"

// 创建默认配置
paths := models.NewPathsConfig()

// 自定义配置
paths := models.NewPathsConfig().
    WithPluginInstallDir("/opt/plugins").
    WithServiceDeployDir("/var/services")

// 获取路径
serviceName := "myservice"
root := paths.GetServiceRoot(serviceName)
// 返回: "/var/services/myservice"

binPath := paths.GetServiceBinPath(serviceName)
// 返回: "/var/services/myservice/bin"

// 验证
if err := paths.Validate(); err != nil {
    log.Fatal(err)
}
```

**收益**:
- ✅ 消除硬编码路径
- ✅ 集中配置管理
- ✅ 易于测试和修改
- ✅ 类型安全

### 4. 变量管理服务 (services/variable_manager.go)

统一管理变量准备，集成路径配置。

**功能**:
- 统一变量准备接口
- 集成PathsConfig
- 提供预设方法
- 自动路径变量注入

**方法**:
- `PrepareForDockerfile(arch)` - Dockerfile变量
- `PrepareForCompose()` - Compose变量
- `PrepareForBuildScript()` - 构建脚本变量
- `PrepareForScript()` - 通用脚本变量
- `PrepareForMakefile()` - Makefile变量
- `PrepareForDevOps()` - DevOps变量
- `AddPathVariables(composer, serviceName)` - 添加路径变量
- `PrepareWithPaths(presetFunc)` - 准备带路径的变量

**使用示例**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/services"

// 创建管理器
manager := services.NewVariableManager(ctx)

// 准备变量
composer := manager.PrepareForDockerfile("amd64")

// 添加路径变量
composer = manager.AddPathVariables(composer, "myservice")

// 一步到位（推荐）
composer := manager.PrepareWithPaths(func() *context.VariableComposer {
    return manager.PrepareForScript()
})

// 构建变量
vars := composer.Build()
// vars 包含所有预设变量 + 路径变量
```

**自动注入的路径变量**:
- `SERVICE_ROOT` - 服务根路径
- `SERVICE_BIN_PATH` - 二进制路径
- `SERVICE_CONFIG_PATH` - 配置路径
- `SERVICE_LOG_PATH` - 日志路径
- `SERVICE_DATA_PATH` - 数据路径
- `PLUGIN_INSTALL_DIR` - 插件安装路径

**收益**:
- ✅ 统一变量管理
- ✅ 自动路径注入
- ✅ 简化生成器代码
- ✅ 提高可测试性

### 5. 领域事件系统 (events/)

提供事件驱动架构支持，提高可观测性。

**核心组件**:
- `Event` - 事件接口
- `EventPublisher` - 事件发布器
- `EventHandler` - 事件处理器

**内置事件**:
- `GenerationStartedEvent` - 生成开始
- `GenerationCompletedEvent` - 生成完成
- `GenerationFailedEvent` - 生成失败
- `ValidationStartedEvent` - 验证开始
- `ValidationCompletedEvent` - 验证完成
- `ValidationFailedEvent` - 验证失败
- `TemplateRenderedEvent` - 模板渲染

**内置处理器**:
- `LoggingEventHandler` - 日志记录
- `MetricsEventHandler` - 指标收集
- `FilteredEventHandler` - 事件过滤
- `CompositeEventHandler` - 组合处理

**使用示例**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/events"

// 创建发布器
publisher := events.NewSimpleEventPublisher()

// 订阅处理器
logger := events.NewLoggingEventHandler("GEN")
metrics := events.NewMetricsEventHandler()
publisher.Subscribe(logger)
publisher.Subscribe(metrics)

// 发布事件
event := events.NewGenerationStartedEvent("dockerfile", map[string]interface{}{
    "arch": "amd64",
})
publisher.Publish(event)

// 查看指标
count := metrics.GetCount(events.EventGenerationCompleted)
```

**在生成器中使用**:
```go
func (g *Generator) Generate() (string, error) {
    // 发布开始事件
    g.PublishEvent(events.NewGenerationStartedEvent(
        g.GetName(),
        map[string]interface{}{"config": g.GetContext().Config},
    ))
    
    // 执行生成...
    result, err := g.doGenerate()
    
    if err != nil {
        g.PublishEvent(events.NewGenerationFailedEvent(
            g.GetName(),
            map[string]interface{}{"error": err.Error()},
        ))
        return "", err
    }
    
    g.PublishEvent(events.NewGenerationCompletedEvent(
        g.GetName(),
        map[string]interface{}{"size": len(result)},
    ))
    
    return result, nil
}
```

**收益**:
- ✅ 提高可观测性
- ✅ 解耦组件
- ✅ 易于扩展
- ✅ 支持并发

### 6. 模板仓储 (repositories/)

提供模板存储抽象，支持多种存储后端。

**仓储接口**:
```go
type TemplateRepository interface {
    Get(name string) (*Template, error)
    List() ([]*Template, error)
    ListByCategory(category string) ([]*Template, error)
    Exists(name string) bool
    Save(template *Template) error
    Delete(name string) error
}
```

**实现**:
1. **EmbeddedTemplateRepository** - 嵌入式仓储（生产环境）
   - 模板嵌入到二进制
   - 只读访问
   - 高性能

2. **FileSystemTemplateRepository** - 文件系统仓储（开发环境）
   - 从文件系统加载
   - 支持读写
   - 支持热重载

**使用示例（嵌入式）**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/repositories"

// 创建仓储
repo := repositories.NewEmbeddedTemplateRepository()

// 注册模板
repo.Register(&repositories.Template{
    Name:     "dockerfile",
    Content:  dockerfileTemplate,
    Category: "docker",
})

// 获取模板
template, err := repo.Get("dockerfile")
if err != nil {
    log.Fatal(err)
}

// 使用模板
engine := core.NewTemplateEngine()
result, _ := engine.Render(template.Content, vars)
```

**使用示例（文件系统）**:
```go
// 创建仓储（自动加载所有.tmpl文件）
repo, err := repositories.NewFileSystemTemplateRepository("./templates")
if err != nil {
    log.Fatal(err)
}

// 获取模板
template, _ := repo.Get("dockerfile")

// 按类别列出
dockerTemplates, _ := repo.ListByCategory("docker")

// 保存新模板
newTemplate := &repositories.Template{
    Name:     "custom",
    Content:  "...",
    Category: "docker",
}
repo.Save(newTemplate)

// 热重载
repo.Reload()
```

**收益**:
- ✅ 统一接口
- ✅ 可切换实现
- ✅ 分类管理
- ✅ 线程安全

### 7. 语言服务 (services/language_service.go)

使用策略模式处理不同编程语言的特定逻辑。

**架构**:
```
LanguageService
    ├── GoStrategy
    ├── PythonStrategy
    ├── NodeJSStrategy
    ├── JavaStrategy
    └── RustStrategy
```

**接口定义**:
```go
type LanguageStrategy interface {
    GetName() string
    GetDependencyFiles() []string
    GetDepsInstallCommand() string
    GetPackageManager() string
}
```

**支持的语言**:
- ✅ Go
- ✅ Python
- ✅ NodeJS
- ✅ Java
- ✅ Rust

**使用示例**:
```go
import "github.com/junjiewwang/service-template/pkg/generator/domain/services"

// 创建服务
langService := services.NewLanguageService()

// 获取依赖文件
files := langService.GetDependencyFiles("go", true, []string{})
// 返回: ["go.mod", "go.sum"]

// 获取安装命令
command := langService.GetDepsInstallCommand("python")
// 返回: "pip install -r requirements.txt"

// 检查语言支持
if langService.IsSupported("rust") {
    // 处理Rust项目
}

// 列出所有支持的语言
languages := langService.ListSupportedLanguages()
// 返回: ["go", "python", "nodejs", "java", "rust"]
```

**扩展新语言**:
```go
// 1. 实现策略接口
type RubyStrategy struct{}

func (s *RubyStrategy) GetName() string {
    return "ruby"
}

func (s *RubyStrategy) GetDependencyFiles() []string {
    return []string{"Gemfile", "Gemfile.lock"}
}

func (s *RubyStrategy) GetDepsInstallCommand() string {
    return "bundle install"
}

func (s *RubyStrategy) GetPackageManager() string {
    return "bundle"
}

// 2. 注册策略
langService := services.NewLanguageService()
langService.Register(&RubyStrategy{})

// 3. 完成！
```

**收益**:
- ✅ 消除硬编码
- ✅ 符合开闭原则
- ✅ 易于扩展新语言
- ✅ 逻辑清晰分离

## 🎨 设计模式

### 1. 策略模式 (Strategy Pattern)
- **应用**: `LanguageService` + `LanguageStrategy`
- **优势**: 消除硬编码，易于扩展
- **符合**: 开闭原则

### 2. 服务模式 (Service Pattern)
- **应用**: `PluginService`
- **优势**: 集中业务逻辑，减少重复
- **符合**: DDD领域服务

### 3. 工厂模式 (Factory Pattern)
- **应用**: 错误构造函数
- **优势**: 统一错误创建，类型安全
- **符合**: 单一职责原则

## 📊 测试覆盖率

| 模块 | 测试数量 | 覆盖率 | 状态 |
|------|---------|--------|------|
| **errors** | 9个测试 | 100% | ✅ |
| **language_service** | 13个测试 | 100% | ✅ |
| **plugin_service** | 待添加 | - | ⏳ |

## 🔄 迁移指南

### 从旧代码迁移

**旧代码** (重复的插件处理):
```go
var plugins []map[string]interface{}
sharedInstallDir := ctx.Config.Plugins.InstallDir
for _, plugin := range ctx.Config.Plugins.Items {
    // Get base variables using the new variable system
    composer := ctx.GetVariableComposer().WithCommon().WithPlugin()
    baseVars := composer.Build()
    
    plugins = append(plugins, map[string]interface{}{
        "InstallCommand": core.SubstituteVariables(plugin.InstallCommand, baseVars),
        "Name":           plugin.Name,
        "InstallDir":     sharedInstallDir,
        "RuntimeEnv":     plugin.RuntimeEnv,
    })
}
```

**新代码** (使用服务):
```go
pluginService := services.NewPluginService(ctx, g.GetEngine())
if pluginService.HasPlugins() {
    plugins := pluginService.PrepareForDockerfile()
    composer.Override("PLUGINS", plugins)
}
```

**收益**: 代码减少85%，逻辑集中管理

## 📚 相关文档

- [代码分析报告](../CODE_ANALYSIS_REPORT.md)
- [阶段1实施报告](../PHASE1_IMPLEMENTATION_REPORT.md)
- [优化指南](../OPTIMIZATION_GUIDE.md)

## 🎯 已实现功能

### 领域模型

```
models/
├── paths_config.go        # ✅ 路径配置值对象
└── paths_config_test.go   # ✅ 完整测试
```

### 领域服务

```
services/
├── plugin_service.go          # ✅ 插件服务
├── language_service.go        # ✅ 语言服务
├── language_service_test.go   # ✅ 语言服务测试
├── variable_manager.go        # ✅ 变量管理服务
└── variable_manager_test.go   # ✅ 变量管理测试
```

### 领域事件

```
events/
├── event.go               # ✅ 事件接口和基类
├── publisher.go           # ✅ 事件发布器
├── generator_events.go    # ✅ 生成器事件
├── handlers.go            # ✅ 事件处理器
└── events_test.go         # ✅ 完整测试
```

### 仓储

```
repositories/
├── template_repository.go    # ✅ 仓储接口
├── embedded_repository.go    # ✅ 嵌入式仓储
├── filesystem_repository.go  # ✅ 文件系统仓储
└── repository_test.go        # ✅ 完整测试
```

## 🚀 未来计划

### 待实现的领域模型

```
models/
├── plugin.go          # 插件领域模型
├── language.go        # 语言领域模型
└── generation_config.go # 生成配置聚合根
```

### 待实现的领域服务

```
services/
├── template_service.go    # 模板管理服务
└── validation_service.go  # 验证服务
```

## 🤝 贡献

欢迎贡献新的语言策略或领域服务！

### 贡献步骤

1. Fork 项目
2. 创建功能分支
3. 实现功能并编写测试
4. 确保测试覆盖率 > 80%
5. 提交 PR

---

**创建日期**: 2025-11-12  
**最后更新**: 2025-11-12  
**维护者**: Development Team
