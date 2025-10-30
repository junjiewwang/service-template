# TCS Service Template Generator

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A powerful, configuration-driven tool to generate service templates with Docker, Kubernetes, and CI/CD configurations.

## 🎯 Features

- ✅ **Single Configuration Source**: All settings in one `service.yaml` file
- ✅ **Auto-Generation**: Docker, K8s configs generated automatically
- ✅ **Multi-Language Support**: Go, Python, Node.js, Java, Rust
- ✅ **Multi-Architecture**: AMD64 and ARM64 support
- ✅ **Multi-Port Support**: Configure multiple service ports
- ✅ **Plugin System**: Flexible plugin installation mechanism
- ✅ **Type-Safe**: Go implementation with compile-time checks
- ✅ **Incremental Updates**: Safe regeneration preserving customizations

## 📦 Installation

### From Source

```bash
git clone https://github.com/junjiewwang/service-template.git
cd service-template
make install
```

### Using Go Install

```bash
go install github.com/junjiewwang/service-template/cmd/tcs-gen@latest
```

## 🚀 Quick Start

### 1. Initialize Configuration

```bash
cd /path/to/your-project
tcs-gen init
```

This creates a `service.yaml` file with example configuration.

### 2. Edit Configuration

Edit `service.yaml` to configure your service:

```yaml
service:
  name: my-service
  description: "My Service"
  ports:
    - name: http
      port: 8080
      protocol: TCP
      expose: true

language:
  type: go
  version: "1.23"

build:
  dependency_files:
    auto_detect: true
  builder_image:
    amd64: "golang:1.23-alpine"
    arm64: "golang:1.23-alpine"
  runtime_image:
    amd64: "alpine:latest"
    arm64: "alpine:latest"
  commands:
    build: "go build -o app ./cmd/server"
  output_dir: dist

runtime:
  healthcheck:
    enabled: true
    type: http
    http:
      path: /health
      port: 8080
  startup:
    command: "./app"

local_dev:
  compose:
    volumes: []
  kubernetes:
    enabled: false
    namespace: default

metadata:
  template_version: "2.0.0"
  generator: "tcs-gen"
```

### 3. Validate Configuration

```bash
tcs-gen validate
```

### 4. Generate Project Files

```bash
tcs-gen generate
```

This generates:
- ✓ `Dockerfile.amd64` and `Dockerfile.arm64`
- ✓ `compose.yaml`
- ✓ `Makefile`
- ✓ `bk-ci/tcs/build.sh`
- ✓ `bk-ci/tcs/deps_install.sh`
- ✓ `bk-ci/tcs/rt_prepare.sh`
- ✓ `.tad/devops.yaml`
- ✓ `hooks/healthchk.sh`
- ✓ `hooks/start.sh`
- ✓ `k8s-manifests/configmap.yaml` (if K8s enabled)

### 5. Build and Run

```bash
# Build Docker image
make docker-build

# Start services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

## 📖 Documentation

### Commands

#### `tcs-gen init`

Initialize a new `service.yaml` configuration file.

```bash
tcs-gen init [--config service.yaml]
```

#### `tcs-gen validate`

Validate the configuration file.

```bash
tcs-gen validate [--config service.yaml]
```

#### `tcs-gen generate`

Generate all project files.

```bash
tcs-gen generate [--config service.yaml] [--output .] [--skip-validation]
```

Options:
- `-c, --config`: Path to service.yaml (default: "service.yaml")
- `-o, --output`: Output directory (default: current directory)
- `--skip-validation`: Skip configuration validation

#### `tcs-gen version`

Print version information.

```bash
tcs-gen version
```

### Configuration Reference

See [service.yaml.example](service.yaml.example) for a complete configuration example with comments.

Key sections:
- **service**: Basic service information (name, ports, deploy directory)
- **language**: Language type and version
- **build**: Build configuration (images, dependencies, commands)
- **plugins**: Plugin installation configuration
- **runtime**: Runtime configuration (healthcheck, startup)
- **local_dev**: Local development settings (compose, kubernetes)
- **makefile**: Custom Makefile targets
- **metadata**: Template metadata

### Variable Substitution

The following variables are available in templates:

- `${SERVICE_NAME}` - Service name
- `${SERVICE_PORT}` - Main service port (first port)
- `${SERVICE_ROOT}` - Service root directory
- `${DEPLOY_DIR}` - Deployment directory
- `${BUILD_OUTPUT_DIR}` - Build output directory
- `${CONFIG_DIR}` - Configuration directory
- `${SERVICE_BIN_DIR}` - Binary directory
- `${PLUGIN_NAME}` - Plugin name (in plugin context)
- `${PLUGIN_INSTALL_DIR}` - Plugin install directory
- `${GOARCH}` - Go architecture (in build context)
- `${GOOS}` - Go OS (in build context)

## 🧪 Testing

Run tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./pkg/config/...
go test ./pkg/generator/...
```

## 🏗️ Project Structure

```
tcs-service-template/
├── cmd/
│   └── tcs-gen/
│       ├── main.go
│       └── commands/
│           ├── root.go
│           ├── init.go
│           ├── validate.go
│           └── generate.go
│
├── pkg/
│   ├── config/
│   │   ├── types.go
│   │   ├── loader.go
│   │   ├── validator.go
│   │   └── *_test.go
│   │
│   ├── generator/
│   │   ├── generator.go
│   │   ├── template.go
│   │   ├── variables.go
│   │   ├── dockerfile.go
│   │   ├── compose.go
│   │   ├── makefile.go
│   │   ├── scripts.go
│   │   ├── configmap.go
│   │   └── *_test.go
│   │
│   └── utils/
│       └── file.go
│
├── service.yaml.example
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Sprig](https://github.com/Masterminds/sprig) - Template functions

## 📞 Support

For issues and questions, please open an issue on GitHub.
