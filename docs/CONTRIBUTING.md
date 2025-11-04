# 🤝 Contributing Guide

Thank you for your interest in contributing to TCS Service Template Generator!

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code:

- **Be respectful**: Treat everyone with respect and kindness
- **Be collaborative**: Work together to achieve common goals
- **Be inclusive**: Welcome and support people of all backgrounds
- **Be constructive**: Provide helpful feedback and suggestions

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Git
- Docker (for testing generated Dockerfiles)
- kubectl (optional, for testing Kubernetes manifests)

### Finding Issues to Work On

- Check the [Issues](https://github.com/junjiewwang/service-template/issues) page
- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to let others know you're working on it

### Types of Contributions

We welcome various types of contributions:

- 🐛 **Bug fixes**: Fix issues and improve stability
- ✨ **New features**: Add new generators, languages, or capabilities
- 📝 **Documentation**: Improve docs, add examples, fix typos
- 🧪 **Tests**: Add test coverage, improve test quality
- 🎨 **Code quality**: Refactoring, performance improvements
- 🌐 **Translations**: Translate documentation to other languages

## Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/service-template.git
cd service-template

# Add upstream remote
git remote add upstream https://github.com/junjiewwang/service-template.git
```

### 2. Install Dependencies

```bash
# Download Go dependencies
go mod download

# Install development tools
make install-tools
```

### 3. Build the Project

```bash
# Build the binary
go build -o svcgen ./cmd/svcgen

# Or use make
make build
```

### 4. Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./pkg/config/...
```

### 5. Verify Your Setup

```bash
# Generate a test project
cd demo-app
../svcgen generate

# Build the generated Docker image
make docker-build

# Run the service
make docker-up
```

## How to Contribute

### Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the issue
2. **Steps to Reproduce**: Detailed steps to reproduce the bug
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**: OS, Go version, tool version
6. **Configuration**: Your `service.yaml` (sanitized)
7. **Logs**: Relevant error messages or logs

**Example Bug Report**:

```markdown
### Description
Dockerfile generation fails when using custom health check

### Steps to Reproduce
1. Create service.yaml with custom health check
2. Run `svcgen generate`
3. Error occurs

### Expected Behavior
Dockerfile should be generated with custom health check script

### Actual Behavior
Error: "template: dockerfile:45: unexpected EOF"

### Environment
- OS: macOS 14.0
- Go: 1.23.0
- svcgen: 2.0.0

### Configuration
```yaml
runtime:
  healthcheck:
    enabled: true
    type: custom
    custom_script: |
      #!/bin/sh
      exit 0
```

### Logs
```
Error generating Dockerfile: template: dockerfile:45: unexpected EOF
```
```

### Suggesting Features

When suggesting features, please include:

1. **Use Case**: Why is this feature needed?
2. **Proposed Solution**: How should it work?
3. **Alternatives**: Other approaches you considered
4. **Examples**: Example configuration or usage

**Example Feature Request**:

```markdown
### Use Case
I need to generate Helm charts for Kubernetes deployment

### Proposed Solution
Add a new generator type for Helm charts that:
- Generates Chart.yaml
- Generates values.yaml
- Generates templates/

### Alternatives
- Use existing K8s manifests with Helm wrapper
- Use external Helm chart generator

### Example Configuration
```yaml
local_dev:
  kubernetes:
    enabled: true
    helm:
      enabled: true
      chart_version: "1.0.0"
```
```

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines:

- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Follow standard Go naming conventions
- Write clear, self-documenting code

### Code Organization

```
pkg/
├── config/          # Configuration management
│   ├── types.go     # Type definitions
│   ├── loader.go    # Loading logic
│   └── validator.go # Validation logic
├── generator/       # Code generation
│   ├── generator.go # Main orchestrator
│   ├── factory.go   # Factory pattern
│   └── template_*.go # Individual generators
└── utils/           # Utility functions
```

### Naming Conventions

- **Packages**: Short, lowercase, single word (e.g., `config`, `generator`)
- **Files**: Lowercase with underscores (e.g., `template_dockerfile.go`)
- **Types**: PascalCase (e.g., `ServiceConfig`, `TemplateGenerator`)
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase (e.g., `serviceName`, `portNumber`)
- **Constants**: PascalCase or UPPER_CASE (e.g., `DefaultPort`, `MAX_RETRIES`)

### Comments

- Add godoc comments for all exported types and functions
- Use complete sentences with proper punctuation
- Explain "why" not "what" for complex logic

```go
// ServiceConfig represents the complete service configuration loaded from service.yaml.
// It includes all settings needed to generate infrastructure code for a microservice.
type ServiceConfig struct {
    Service  ServiceInfo
    Language LanguageConfig
    Build    BuildConfig
}

// Validate checks if the configuration is valid and returns an error if not.
// It performs comprehensive validation including type checking, required fields,
// and business rule validation.
func (c *ServiceConfig) Validate() error {
    // Implementation
}
```

### Error Handling

- Return errors, don't panic
- Wrap errors with context using `fmt.Errorf`
- Use custom error types for specific errors

```go
// Good
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Bad
if err != nil {
    panic(err)
}
```

## Testing Guidelines

### Test Coverage

- Aim for 80%+ test coverage
- Write tests for all new code
- Update tests when modifying existing code

### Test Organization

```
pkg/
├── config/
│   ├── loader.go
│   ├── loader_test.go      # Tests for loader.go
│   ├── validator.go
│   └── validator_test.go   # Tests for validator.go
```

### Test Naming

```go
// Test function names should be descriptive
func TestServiceConfig_Validate_ValidConfig(t *testing.T) {}
func TestServiceConfig_Validate_MissingServiceName(t *testing.T) {}
func TestServiceConfig_Validate_InvalidPortRange(t *testing.T) {}
```

### Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestValidatePort(t *testing.T) {
    tests := []struct {
        name    string
        port    int
        wantErr bool
    }{
        {"valid port", 8080, false},
        {"port too low", 0, true},
        {"port too high", 70000, true},
        {"minimum port", 1, false},
        {"maximum port", 65535, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validatePort(tt.port)
            if (err != nil) != tt.wantErr {
                t.Errorf("validatePort() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Helpers

Use testify for assertions:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
    result, err := doSomething()
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Integration Tests

Write integration tests for end-to-end workflows:

```go
func TestGenerateAll_Integration(t *testing.T) {
    // Setup
    cfg := loadTestConfig(t)
    gen := NewGenerator(cfg)
    
    // Execute
    err := gen.GenerateAll()
    require.NoError(t, err)
    
    // Verify
    assert.FileExists(t, "Dockerfile.test-service.amd64")
    assert.FileExists(t, "compose.yaml")
    assert.FileExists(t, "Makefile")
}
```

## Pull Request Process

### 1. Create a Branch

```bash
# Update your fork
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b feature/my-feature
# or
git checkout -b fix/my-bugfix
```

### 2. Make Changes

- Write code following coding standards
- Add tests for new functionality
- Update documentation as needed
- Run tests locally

```bash
# Format code
gofmt -w .

# Run linter
golangci-lint run

# Run tests
go test ./...

# Check coverage
go test -cover ./...
```

### 3. Commit Changes

Write clear, descriptive commit messages:

```bash
# Good commit messages
git commit -m "feat: add Helm chart generator"
git commit -m "fix: handle empty health check configuration"
git commit -m "docs: update configuration guide with examples"
git commit -m "test: add tests for validator edge cases"

# Commit message format
# <type>: <description>
#
# Types:
# feat: New feature
# fix: Bug fix
# docs: Documentation changes
# test: Test changes
# refactor: Code refactoring
# style: Code style changes
# chore: Build/tooling changes
```

### 4. Push Changes

```bash
git push origin feature/my-feature
```

### 5. Create Pull Request

1. Go to GitHub and create a pull request
2. Fill out the PR template
3. Link related issues
4. Request review from maintainers

**PR Template**:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Test improvement

## Related Issues
Fixes #123

## Changes Made
- Added X feature
- Fixed Y bug
- Updated Z documentation

## Testing
- [ ] Added unit tests
- [ ] Added integration tests
- [ ] Manually tested
- [ ] All tests pass

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests added and passing
```

### 6. Code Review

- Respond to review comments
- Make requested changes
- Push updates to the same branch

```bash
# Make changes based on review
git add .
git commit -m "address review comments"
git push origin feature/my-feature
```

### 7. Merge

Once approved, maintainers will merge your PR.

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **Major** (X.0.0): Breaking changes
- **Minor** (x.Y.0): New features (backward compatible)
- **Patch** (x.y.Z): Bug fixes

### Release Checklist

1. Update version in code
2. Update CHANGELOG.md
3. Create git tag
4. Build binaries
5. Create GitHub release
6. Update documentation

### Creating a Release

```bash
# Update version
VERSION=2.1.0

# Update CHANGELOG.md
# Add release notes

# Commit changes
git add .
git commit -m "chore: release v${VERSION}"

# Create tag
git tag -a v${VERSION} -m "Release v${VERSION}"

# Push changes and tag
git push origin main
git push origin v${VERSION}

# Build binaries
make build-all

# Create GitHub release
# Upload binaries
# Add release notes
```

## Getting Help

- **Documentation**: Check the [docs](../README.md)
- **Issues**: Search [existing issues](https://github.com/junjiewwang/service-template/issues)
- **Discussions**: Join [GitHub Discussions](https://github.com/junjiewwang/service-template/discussions)
- **Email**: Contact maintainers at junjiewwang@example.com

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- GitHub contributors page

Thank you for contributing! 🎉

---

[← Back to README](../README.md)
