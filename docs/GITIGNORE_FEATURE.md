# .gitignore 自动管理功能

## 需求描述

`svcgen generate` 生成文件后，自动将生成的文件路径写入 `.gitignore`，使用户的仓库中只需要版本控制 `service.yaml`，其余生成产物均被 git 忽略。

## 设计方案

### 配置开关

通过 `service.yaml` 中的 `metadata.manage_gitignore` 字段控制是否自动管理 `.gitignore`：

```yaml
metadata:
  template_version: "2.0.0"
  generator: "tcs-gen"
  # 是否自动管理 .gitignore（将生成的文件添加到 .gitignore）
  # 默认: false
  manage_gitignore: true
```

| 值 | 行为 |
|------|------|
| `false`（默认） | 不写入 `.gitignore`，用户自行管理 |
| `true` | 自动将生成的文件条目添加到 `.gitignore` |

### 核心思路

在 `Generate()` 流程末尾新增 `updateGitignore()` 步骤，使用 **marker block** 方式管理 `.gitignore` 中的生成条目，保留用户自定义的 ignore 规则。

### Marker Block 格式

```gitignore
# >>> svcgen generated - DO NOT EDIT >>>
# Generated files (only service.yaml needs to be tracked)
.env.make
.tad/
Makefile
compose.yaml
# <<< svcgen generated <<<
```

### 行为规则

| 场景 | 行为 |
|------|------|
| `.gitignore` 不存在 | 创建新文件，写入 marker block |
| `.gitignore` 已存在，无 marker block | 追加 marker block 到文件末尾 |
| `.gitignore` 已存在，有 marker block | 替换 marker block 内容，保留用户自定义内容 |
| 多次运行 | 幂等，内容未变化则不写入 |

### 被忽略的条目

| 条目 | 说明 |
|------|------|
| `.tad/` | 包含 devops.yaml、Dockerfile、所有构建/部署脚本 |
| `compose.yaml` | Docker Compose 文件 |
| `Makefile` | 构建入口 |
| `.env.make` | 由 Makefile 从 devops.yaml 解析生成 |

### 自定义 CI 路径处理

当用户配置了 `ci.script_dir` 指向非 `.tad/` 目录时（如 `bk-ci/tcs`），会额外添加该路径到 ignore 列表。如果 `script_dir` 仍在 `.tad/` 下，则不重复添加。

## 实施记录

### 新增文件

- `pkg/generator/gitignore.go` — `.gitignore` 管理核心逻辑
  - `gitignoreEntries()` — 动态收集需忽略的条目
  - `buildGitignoreBlock()` — 构建 marker block 内容
  - `updateGitignore()` — 创建/追加/替换 `.gitignore` 文件
  - `replaceOrAppendBlock()` — marker block 替换或追加逻辑
- `pkg/generator/gitignore_test.go` — 单元测试（13 个测试用例）

### 修改文件

- `pkg/config/types.go` — `MetadataConfig` 新增 `ManageGitignore bool` 字段
- `pkg/generator/generator.go` — 在 `Generate()` 中根据 `ManageGitignore` 配置决定是否调用 `updateGitignore()`
- `pkg/generator/generator_test.go` — `expectedFiles` 中新增 `.gitignore` 验证（启用 `ManageGitignore`）
- `pkg/config/testutil/config_builder.go` — 新增 `WithManageGitignore()` Builder 方法
- `pkg/config/testutil/options.go` — 新增 `WithManageGitignoreOpt()` Option 函数
- `demo-app/service.yaml` — 示例中展示 `manage_gitignore` 配置（注释说明默认 false）

### 测试覆盖

- ✅ `TestBuildGitignoreBlock` — block 构建
- ✅ `TestReplaceOrAppendBlock_NoExisting` — 空文件追加
- ✅ `TestReplaceOrAppendBlock_AppendToExisting` — 已有文件追加
- ✅ `TestReplaceOrAppendBlock_ReplaceExistingBlock` — 替换已有 block
- ✅ `TestReplaceOrAppendBlock_PreservesUserContent` — 保留用户内容
- ✅ `TestGenerator_GitignoreEntries_DefaultPaths` — 默认路径条目
- ✅ `TestGenerator_GitignoreEntries_CustomScriptDir` — 自定义 CI 路径
- ✅ `TestGenerator_GitignoreEntries_ScriptDirUnderTad` — .tad 内路径去重
- ✅ `TestGenerator_UpdateGitignore_CreateNew` — 创建新文件
- ✅ `TestGenerator_UpdateGitignore_AppendToExisting` — 追加到已有文件
- ✅ `TestGenerator_UpdateGitignore_UpdateExisting` — 更新已有 block
- ✅ `TestGenerator_UpdateGitignore_Idempotent` — 幂等性
- ✅ `TestGenerator_ManageGitignore_DefaultFalse` — 默认关闭时不生成 .gitignore
- ✅ `TestGenerator_ManageGitignore_EnabledTrue` — 启用时生成 .gitignore
- ✅ `TestGenerator_Generate` — 集成测试（验证 .gitignore 在生成列表中）

## 状态

- [x] 需求分析
- [x] 方案设计
- [x] 核心实现
- [x] 单元测试
- [x] 集成测试验证
- [x] 文档记录
