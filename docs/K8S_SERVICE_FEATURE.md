# K8s Service Manifest 生成功能

## 需求描述

`svcgen generate` 在生成 `compose.yaml` 的同时，额外生成一个 Kubernetes Service manifest（`.tad/k8s-service.yaml`），确保 port name 信息不丢失。

### 问题背景

用户在 `service.yaml` 中定义了带 `name` 的端口：

```yaml
ports:
  - name: http
    port: 8080
    protocol: TCP
    expose: true
  - name: metrics
    port: 9090
    protocol: TCP
    expose: false
```

通过 `compose.yaml` → `kompose convert` 转换为 K8s Service 时，port name 会丢失（变成端口号），导致 **Istio 无法正确识别协议并进行流量路由**。

Istio 依赖 K8s Service 的 port name 前缀（如 `http`、`grpc`、`tcp`）来识别协议。纯数字 name（如 `"8080"`）会导致 Istio 回退到 TCP 透传模式。

### 根因分析

Docker Compose 的 `ports` 语法本身不支持 `name` 语义：

```yaml
# compose.yaml 只能表达端口号，无法表达 name
ports:
  - "8080"
  - "9090"
```

kompose 转换时由于缺少 name 信息，只能用端口号作为 port name。

## 设计方案

### 解决思路

**不替换 kompose 的输出，而是在部署后用 `kubectl patch` 最小化修正 port name**：

| 步骤 | 操作 |
|------|------|
| `k8s-convert` | kompose 正常转换（生成 Deployment、Service 等） |
| `k8s-deploy` | kubectl apply 所有 manifests |
| **patch** | `kubectl patch service` 用 strategic merge 只修改 port name |

### 为什么用 patch 而不是替换整个 Service

- kompose 生成的 Service 可能包含 annotations、labels、selector 等其他重要配置
- 全量替换会丢失这些信息，导致请求不通
- strategic merge patch 按 `port` 号匹配，只修改 `name` 字段，**最小化改动**

### 生成规则

| 规则 | 说明 |
|------|------|
| 输出路径 | `.tad/k8s-service.yaml`（strategic merge patch） |
| 包含所有端口 | kompose 会转换所有 compose ports，patch 需覆盖全部 |
| port name 使用 `service.yaml` 中定义的 `name` | 确保 Istio 协议识别正确 |
| 只包含 `port` + `name` | 最小化 patch，不修改其他字段 |

### 生成示例

`.tad/k8s-service.yaml`（strategic merge patch）：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
    - port: 8080
      name: http
    - port: 9090
      name: metrics
```

Makefile 中执行：
```bash
kubectl patch service my-service -n <namespace> --type=strategic -p "$(cat .tad/k8s-service.yaml)"
```

效果：只修改 Service 中对应 port 的 name 字段，其他字段（selector、annotations、targetPort、protocol 等）完全不动。

## 实施记录

### 新增文件

- `pkg/generator/generators/k8s/service/generator.go` — K8s Service 生成器
- `pkg/generator/generators/k8s/service/generator_test.go` — 单元测试（7 个测试用例）
- `pkg/generator/generators/k8s/service/templates/service.yaml.tmpl` — K8s Service 模板

### 修改文件

- `pkg/generator/generator.go` — 注册 k8s-service 生成器，在 Generate() 流程中调用
- `pkg/generator/generator_test.go` — expectedFiles 中新增 `.tad/k8s-service.yaml`

### 测试覆盖

- ✅ `TestGenerator_Generate` — 生成 patch + 包含所有端口
- ✅ `TestGenerator_Generate_MultiplePorts` — 多端口 patch
- ✅ `TestGenerator_Generate_NoPorts` — 无端口时不输出 ports section
- ✅ `TestGenerator_GetName` — 生成器名称
- ✅ `TestGenerator_Validate` — 配置验证

### .gitignore 影响

`.tad/k8s-service.yaml` 已被 `.tad/` 目录条目覆盖，无需额外添加。

## 使用方式

已集成到 `make k8s-deploy` 流程中，**无需额外操作**：

```bash
# 正常执行 k8s-deploy 即可，流程自动 patch port names
make k8s-deploy
```

### 自动 patch 流程

`make k8s-deploy` 的执行步骤：

1. `k8s-convert` → kompose 正常转换 compose.yaml → K8s manifests
2. `kubectl apply -f K8S_OUTPUT_DIR/` → 部署所有资源
3. **`kubectl patch service` → 用 `.tad/k8s-service.yaml` 修正 port names**
4. 如果 `.tad/k8s-service.yaml` 不存在 → 打印 warning，不阻塞流程

### 手动 patch

如果需要单独修正已部署的 Service：

```bash
kubectl patch service <service-name> -n <namespace> \
  --type=strategic -p "$(cat .tad/k8s-service.yaml)"
```

## 状态

- [x] 问题分析
- [x] 方案设计
- [x] 核心实现
- [x] 单元测试
- [x] 集成测试验证
- [x] 文档记录
