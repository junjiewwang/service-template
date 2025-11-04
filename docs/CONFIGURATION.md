# 📖 Configuration Guide

Complete guide to configuring your service using `service.yaml`.

## Table of Contents

- [Configuration Structure](#configuration-structure)
- [Service Configuration](#service-configuration)
- [Language Configuration](#language-configuration)
- [Build Configuration](#build-configuration)
- [Plugin Configuration](#plugin-configuration)
- [Runtime Configuration](#runtime-configuration)
- [Local Development Configuration](#local-development-configuration)
- [CI/CD Path Configuration](#cicd-path-configuration)
- [Variable Substitution](#variable-substitution)

## Configuration Structure

The `service.yaml` file is organized into logical sections:

```yaml
service:          # Basic service information
language:         # Programming language and version
build:            # Build configuration
plugins:          # Third-party plugins
runtime:          # Runtime configuration
local_dev:        # Local development settings
makefile:         # Custom Makefile targets
ci:               # CI/CD path configuration
metadata:         # Template metadata
```

## Service Configuration

Define basic service information and port configuration.

```yaml
service:
  name: my-service              # Service name (required)
  description: "My Service"     # Service description
  ports:                        # Port configuration (required)
    - name: http                # Port name
      port: 8080                # Port number (1-65535)
      protocol: TCP             # Protocol (TCP/UDP)
      expose: true              # Expose in Docker Compose
      description: "HTTP API"   # Port description
  deploy_dir: /usr/local/services  # Deployment directory
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Service name (alphanumeric, hyphens, underscores) |
| `description` | string | ❌ | Service description |
| `ports` | array | ✅ | List of port configurations (at least one) |
| `deploy_dir` | string | ❌ | Deployment directory (default: `/usr/local/services`) |

### Port Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Port name (e.g., http, grpc, metrics) |
| `port` | integer | ✅ | Port number (1-65535) |
| `protocol` | string | ❌ | Protocol (TCP/UDP, default: TCP) |
| `expose` | boolean | ❌ | Expose in Docker Compose (default: true) |
| `description` | string | ❌ | Port description |

## Language Configuration

Specify the programming language and version.

```yaml
language:
  type: go                      # go | python | nodejs | java | rust
  version: "1.23"               # Language version
  config:                       # Language-specific config
    goproxy: "https://goproxy.cn,direct"
    gosumdb: "sum.golang.org"
```

### Supported Languages

| Language | Type Value | Example Versions |
|----------|-----------|------------------|
| Go | `go` | `1.21`, `1.22`, `1.23` |
| Python | `python` | `3.9`, `3.10`, `3.11`, `3.12` |
| Node.js | `nodejs` | `18`, `20`, `21` |
| Java | `java` | `11`, `17`, `21` |
| Rust | `rust` | `1.70`, `1.75` |

### Language-Specific Configuration

#### Go

```yaml
language:
  type: go
  version: "1.23"
  config:
    goproxy: "https://goproxy.cn,direct"
    gosumdb: "sum.golang.org"
    goprivate: "github.com/myorg/*"
```

#### Python

```yaml
language:
  type: python
  version: "3.12"
  config:
    pip_index_url: "https://pypi.org/simple"
    pip_trusted_host: "pypi.org"
```

#### Node.js

```yaml
language:
  type: nodejs
  version: "20"
  config:
    npm_registry: "https://registry.npmjs.org"
```

## Build Configuration

Configure the build process, dependencies, and output.

```yaml
build:
  dependency_files:
    auto_detect: true           # Auto-detect dependency files
  
  builder_image:                # Build stage images
    amd64: "golang:1.23-alpine"
    arm64: "golang:1.23-alpine"
  
  runtime_image:                # Runtime stage images
    amd64: "alpine:3.18"
    arm64: "alpine:3.18"
  
  system_dependencies:          # System packages
    packages:
      - git
      - make
      - ca-certificates
  
  commands:                     # Build commands
    pre_build: |\
      echo "Running tests..."
      go test ./...
    
    build: |\
      CGO_ENABLED=0 go build -ldflags="-s -w" \
        -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} \
        ./cmd/server
    
    post_build: |\
      echo "Build completed successfully"
  
  output_dir: dist              # Build output directory
```

### Dependency Files

The tool can auto-detect dependency files based on language:

| Language | Auto-Detected Files |
|----------|-------------------|
| Go | `go.mod`, `go.sum` |
| Python | `requirements.txt`, `Pipfile`, `pyproject.toml` |
| Node.js | `package.json`, `package-lock.json`, `yarn.lock` |
| Java | `pom.xml`, `build.gradle` |
| Rust | `Cargo.toml`, `Cargo.lock` |

Or specify manually:

```yaml
build:
  dependency_files:
    auto_detect: false
    files:
      - go.mod
      - go.sum
      - vendor/
```

### Build Commands

Build commands support variable substitution and multi-line scripts:

```yaml
build:
  commands:
    pre_build: |\
      # Run linting
      golangci-lint run ./...
      
      # Run tests with coverage
      go test -cover ./...
    
    build: |\
      # Build for Linux
      GOOS=linux GOARCH=${GOARCH} \
      CGO_ENABLED=0 go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} \
        ./cmd/server
      
      # Copy configuration files
      cp -r configs ${BUILD_OUTPUT_DIR}/conf/
    
    post_build: |\
      # Verify binary
      file ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME}
      
      # Print build info
      echo "Build completed at $(date)"
```

## Plugin Configuration

Configure third-party plugins and tools.

```yaml
plugins:
  - name: selfMonitor
    description: "Monitoring tool"
    download_url: "https://example.com/tool.sh"
    install_dir: /opt/monitor
    install_command: |\
      curl -fsSL "${PLUGIN_DOWNLOAD_URL}" | bash
    runtime_env:
      - name: MONITOR_PATH
        value: ${PLUGIN_INSTALL_DIR}
      - name: MONITOR_ENABLED
        value: "true"
    required: true
```

### Plugin Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Plugin name |
| `description` | string | ❌ | Plugin description |
| `download_url` | string | ✅ | Download URL |
| `install_dir` | string | ✅ | Installation directory |
| `install_command` | string | ✅ | Installation script |
| `runtime_env` | array | ❌ | Runtime environment variables |
| `required` | boolean | ❌ | Whether plugin is required (default: false) |

### Plugin Variables

Available in plugin context:

- `${PLUGIN_NAME}` - Plugin name
- `${PLUGIN_DOWNLOAD_URL}` - Download URL
- `${PLUGIN_INSTALL_DIR}` - Installation directory
- `${PLUGIN_WORK_DIR}` - Plugin work directory

## Runtime Configuration

Configure runtime behavior, health checks, and startup.

```yaml
runtime:
  system_dependencies:
    packages:
      - ca-certificates
      - tzdata
  
  healthcheck:
    enabled: true
    type: http                  # http | tcp | exec | custom
    http:
      path: /health
      port: 8080
      timeout: 3
  
  startup:
    command: |\
      #!/bin/sh
      cd ${SERVICE_ROOT}
      exec ./bin/${SERVICE_NAME}
    
    env:
      - name: GO_ENV
        value: production
      - name: LOG_LEVEL
        value: info
```

### Health Check Types

#### HTTP Health Check

```yaml
runtime:
  healthcheck:
    enabled: true
    type: http
    http:
      path: /health
      port: 8080
      method: GET
      timeout: 3
      interval: 30
      retries: 3
```

#### TCP Health Check

```yaml
runtime:
  healthcheck:
    enabled: true
    type: tcp
    tcp:
      port: 8080
      timeout: 3
```

#### Exec Health Check

```yaml
runtime:
  healthcheck:
    enabled: true
    type: exec
    exec:
      command: |\
        #!/bin/sh
        curl -f http://localhost:8080/health || exit 1
```

#### Custom Health Check

```yaml
runtime:
  healthcheck:
    enabled: true
    type: custom
    custom_script: |\
      #!/bin/sh
      # Custom health check logic
      if [ -f /tmp/healthy ]; then
        exit 0
      else
        exit 1
      fi
```

## Local Development Configuration

Configure Docker Compose and Kubernetes for local development.

```yaml
local_dev:
  compose:
    resources:
      limits:
        cpus: "1.0"
        memory: 1G
      reservations:
        cpus: "0.5"
        memory: 512M
    
    volumes:
      - source: ./configs/config.yaml
        target: ${SERVICE_ROOT}/config.yaml
        type: bind
      - source: ./data
        target: ${SERVICE_ROOT}/data
        type: bind
    
    healthcheck:
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
  
  kubernetes:
    enabled: true
    namespace: default
    replicas: 2
    volume_type: configMap      # configMap | persistentVolumeClaim
    
    resources:
      limits:
        cpu: "1000m"
        memory: "1Gi"
      requests:
        cpu: "500m"
        memory: "512Mi"
    
    wait:
      enabled: true
      timeout: 300s
```

### Volume Types

#### Bind Mount (Docker Compose)

```yaml
local_dev:
  compose:
    volumes:
      - source: ./config.yaml
        target: ${SERVICE_ROOT}/config.yaml
        type: bind
```

#### ConfigMap (Kubernetes)

```yaml
local_dev:
  kubernetes:
    enabled: true
    volume_type: configMap
    volumes:
      - source: ./configs/config.yaml
        target: ${SERVICE_ROOT}/config.yaml
```

Automatically generates ConfigMap from source files.

#### PersistentVolumeClaim (Kubernetes)

```yaml
local_dev:
  kubernetes:
    enabled: true
    volume_type: persistentVolumeClaim
    volumes:
      - source: data-pvc
        target: ${SERVICE_ROOT}/data
```

## CI/CD Path Configuration

Customize CI/CD script locations.

```yaml
ci:
  script_dir: bk-ci/tcs              # CI script directory
  build_config_dir: bk-ci/tcs/build  # Build config directory
```

### Default Paths

| Path | Default Value | Description |
|------|--------------|-------------|
| `script_dir` | `bk-ci/tcs` | CI script directory |
| `build_config_dir` | `bk-ci/tcs/build` | Build config directory |

### Affected Files

Changing CI paths affects:

- `build.sh` location
- `build_deps_install.sh` location
- `rt_prepare.sh` location
- `entrypoint.sh` location
- ConfigMap source paths in Kubernetes manifests

## Variable Substitution

The following variables are available in all templates and commands:

### Service Variables

| Variable | Description | Example |
|----------|-------------|------------|
| `${SERVICE_NAME}` | Service name | `my-service` |
| `${SERVICE_PORT}` | Main service port (first port) | `8080` |
| `${SERVICE_ROOT}` | Service root directory | `/usr/local/services/my-service` |
| `${SERVICE_BIN_DIR}` | Binary directory | `/usr/local/services/my-service/bin` |
| `${DEPLOY_DIR}` | Deployment directory | `/usr/local/services` |
| `${CONFIG_DIR}` | Configuration directory | `/usr/local/services/my-service/conf` |

### Build Variables

| Variable | Description | Example |
|----------|-------------|------------|
| `${BUILD_OUTPUT_DIR}` | Build output directory | `/opt/dist` |
| `${PROJECT_ROOT}` | Project root directory | `/opt` |
| `${GOARCH}` | Go architecture (in build context) | `amd64` |
| `${GOOS}` | Go OS (in build context) | `linux` |

### Plugin Variables

| Variable | Description | Example |
|----------|-------------|------------|
| `${PLUGIN_NAME}` | Plugin name | `selfMonitor` |
| `${PLUGIN_INSTALL_DIR}` | Plugin install directory | `/opt/monitor` |
| `${PLUGIN_WORK_DIR}` | Plugin work directory | `/plugins/selfMonitor` |
| `${PLUGIN_DOWNLOAD_URL}` | Plugin download URL | `https://example.com/tool.sh` |

### CI/CD Variables

| Variable | Description | Example |
|----------|-------------|------------|
| `${CI_SCRIPT_DIR}` | CI script directory | `bk-ci/tcs` |
| `${CI_BUILD_CONFIG_DIR}` | Build config directory | `bk-ci/tcs/build` |

## Complete Example

Here's a complete example configuration:

```yaml
service:
  name: my-api-service
  description: "RESTful API Service"
  ports:
    - name: http
      port: 8080
      protocol: TCP
      expose: true
    - name: metrics
      port: 9090
      protocol: TCP
      expose: false

language:
  type: go
  version: "1.23"
  config:
    goproxy: "https://goproxy.cn,direct"

build:
  dependency_files:
    auto_detect: true
  
  builder_image:
    amd64: "golang:1.23-alpine"
    arm64: "golang:1.23-alpine"
  
  runtime_image:
    amd64: "alpine:3.18"
    arm64: "alpine:3.18"
  
  system_dependencies:
    packages:
      - git
      - make
      - ca-certificates
  
  commands:
    pre_build: "go test ./..."
    build: |\
      CGO_ENABLED=0 go build -ldflags="-s -w" \
        -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} \
        ./cmd/server
    post_build: "echo Done"
  
  output_dir: dist

runtime:
  system_dependencies:
    packages:
      - ca-certificates
      - tzdata
  
  healthcheck:
    enabled: true
    type: http
    http:
      path: /health
      port: 8080
      timeout: 3
  
  startup:
    command: |\
      #!/bin/sh
      cd ${SERVICE_ROOT}
      exec ./bin/${SERVICE_NAME}
    
    env:
      - name: GO_ENV
        value: production

local_dev:
  compose:
    resources:
      limits:
        cpus: "1.0"
        memory: 1G
    volumes:
      - source: ./configs/config.yaml
        target: ${SERVICE_ROOT}/config.yaml
        type: bind
  
  kubernetes:
    enabled: true
    namespace: default
    volume_type: configMap

metadata:
  template_version: "2.0.0"
  generator: "svcgen"
```

---

[← Back to README](../README.md)
