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

不替换 kompose 流程，而是**补充**一个直接生成的 K8s Service manifest：

| 资源 | 生成方式 |
|------|---------|
| Deployment / PVC 等 | 继续用 kompose 转换 |
| **Service** | svcgen 直接生成（带正确 port name） |

### 生成规则

| 规则 | 说明 |
|------|------|
| 输出路径 | `.tad/k8s-service.yaml` |
| 只包含 `expose: true` 的端口 | 非暴露端口不写入 K8s Service |
| port name 使用 `service.yaml` 中定义的 `name` | 确保 Istio 协议识别正确 |
| protocol 统一转大写 | K8s 规范要求（TCP/UDP/SCTP） |

### 生成示例

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
  labels:
    app: my-service
spec:
  selector:
    app: my-service
  ports:
    - name: http
      port: 8080
      targetPort: 8080
      protocol: TCP
```

## 实施记录

### 新增文件

- `pkg/generator/generators/k8s/service/generator.go` — K8s Service 生成器
- `pkg/generator/generators/k8s/service/generator_test.go` — 单元测试（7 个测试用例）
- `pkg/generator/generators/k8s/service/templates/service.yaml.tmpl` — K8s Service 模板

### 修改文件

- `pkg/generator/generator.go` — 注册 k8s-service 生成器，在 Generate() 流程中调用
- `pkg/generator/generator_test.go` — expectedFiles 中新增 `.tad/k8s-service.yaml`

### 测试覆盖

- ✅ `TestGenerator_Generate` — 基本生成 + expose 过滤
- ✅ `TestGenerator_Generate_MultiplePorts` — 多端口混合（expose/non-expose）
- ✅ `TestGenerator_Generate_AllPortsExposed` — 所有端口暴露
- ✅ `TestGenerator_Generate_NoExposedPorts` — 无暴露端口时不输出 ports section
- ✅ `TestGenerator_Generate_ProtocolUppercase` — 小写协议自动转大写
- ✅ `TestGenerator_GetName` — 生成器名称
- ✅ `TestGenerator_Validate` — 配置验证

### .gitignore 影响

`.tad/k8s-service.yaml` 已被 `.tad/` 目录条目覆盖，无需额外添加。

## 使用方式

生成后，部署时使用 `.tad/k8s-service.yaml` 替代 kompose 生成的 Service 资源：

```bash
# 生成所有产物
svcgen generate

# 部署时使用 svcgen 生成的 Service（而非 kompose 转换的）
kubectl apply -f .tad/k8s-service.yaml
```

## 状态

- [x] 问题分析
- [x] 方案设计
- [x] 核心实现
- [x] 单元测试
- [x] 集成测试验证
- [x] 文档记录
