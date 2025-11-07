# 🏗️ Architecture & Design

Deep dive into the architecture, design patterns, and technical decisions.

## Table of Contents

- [System Architecture](#system-architecture)
- [Design Patterns](#design-patterns)
- [Component Details](#component-details)
- [Data Flow](#data-flow)
- [Extension Points](#extension-points)

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     SvcGen CLI                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   init   │  │ validate │  │ generate │  │ version  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Configuration Layer                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  service.yaml (Single Source of Truth)              │   │
│  │  - Service Info    - Build Config                   │   │
│  │  - Language        - Runtime Config                 │   │
│  │  - Plugins         - Local Dev                      │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                 │
│                            ▼                                 │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Validator (Comprehensive Validation)               │   │
│  │  - Type checking   - Required fields                │   │
│  │  - Port ranges     - Image validation               │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Generator Layer                           │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  GeneratorFactory (Factory Pattern)                 │   │
│  │  - Registry-based generator creation                │   │
│  │  - Dynamic generator instantiation                  │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                 │
│         ┌──────────────────┼──────────────────┐             │
│         ▼                  ▼                  ▼             │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐          │
│  │Dockerfile│      │ Compose  │      │ Makefile │          │
│  │Generator │      │Generator │      │Generator │          │
│  └──────────┘      └──────────┘      └──────────┘          │
│         │                  │                  │             │
│         ▼                  ▼                  ▼             │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐          │
│  │  Script  │      │  DevOps  │      │ConfigMap │          │
│  │Generator │      │Generator │      │Generator │          │
│  └──────────┘      └──────────┘      └──────────┘          │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Template Engine                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Go Templates + Sprig Functions                     │   │
│  │  - Variable substitution                            │   │
│  │  - Conditional rendering                            │   │
│  │  - Custom functions (indent, join, etc.)           │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Output Files                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │Dockerfile│  │ compose  │  │ Makefile │  │  Scripts │   │
│  │  (x2)    │  │   .yaml  │  │          │  │  (x5)    │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│  ┌──────────┐  ┌──────────┐                                │
│  │   K8s    │  │  DevOps  │                                │
│  │Manifests │  │  Config  │                                │
│  └──────────┘  └──────────┘                                │
└─────────────────────────────────────────────────────────────┘
```

### Layered Architecture

The system follows a clean layered architecture:

1. **CLI Layer** - User interaction and command routing
2. **Configuration Layer** - YAML parsing and validation
3. **Generator Layer** - Code generation orchestration
4. **Template Layer** - Template rendering and processing
5. **Output Layer** - File writing and organization

## Design Patterns

### 1. Factory Pattern 🏭

**Purpose**: Create generators dynamically based on type.

**Implementation**: `GeneratorFactory` with registry system.

**Code Location**: [`pkg/generator/factory.go`](../pkg/generator/factory.go)

```go
// Generator registry
var generatorRegistry = make(map[GeneratorType]GeneratorCreator)

// Register a generator
func RegisterGenerator(genType GeneratorType, creator GeneratorCreator) {
    generatorRegistry[genType] = creator
}

// Create generator dynamically
func (f *GeneratorFactory) CreateGenerator(genType GeneratorType, arch string) (TemplateGenerator, error) {
    creator, exists := generatorRegistry[genType]
    if !exists {
        return nil, fmt.Errorf("unknown generator type: %s", genType)
    }
    return creator(f.config, f.engine, f.variables, arch)
}

// Register at init time
func init() {
    RegisterGenerator(GeneratorTypeDockerfile, createDockerfileGenerator)
    RegisterGenerator(GeneratorTypeMakefile, createMakefileGenerator)
    // ... more generators
}
```

**Benefits**:
- ✅ Easy to add new generator types
- ✅ Loose coupling between generator creation and usage
- ✅ Centralized generator management
- ✅ Runtime generator selection

### 2. Template Method Pattern 📝

**Purpose**: Define common template rendering logic in base class.

**Implementation**: `BaseTemplateGenerator` with shared methods.

**Code Location**: [`pkg/generator/factory.go`](../pkg/generator/factory.go)

```go
type BaseTemplateGenerator struct {
    config         *config.ServiceConfig
    templateEngine *TemplateEngine
    variables      *Variables
}

// Shared template rendering logic
func (g *BaseTemplateGenerator) RenderTemplate(template string, vars map[string]interface{}) (string, error) {
    return g.templateEngine.Render(template, vars)
}

// Each generator implements Generate()
type TemplateGenerator interface {
    Generate() (string, error)
    GetOutputPath() string
}
```

**Benefits**:
- ✅ Code reuse across all generators
- ✅ Consistent template processing
- ✅ Easy to add common functionality
- ✅ Reduced duplication

### 3. Strategy Pattern 🎯

**Purpose**: Support different health check strategies without code changes.

**Implementation**: Health check type selection in configuration.

**Code Location**: [`pkg/generator/template_healthcheck.go`](../pkg/generator/template_healthcheck.go)

```go
// Health check configuration
type HealthcheckConfig struct {
    Enabled      bool
    Type         string  // http | tcp | exec | custom
    HTTP         *HTTPHealthcheck
    TCP          *TCPHealthcheck
    Exec         *ExecHealthcheck
    CustomScript string
}

// Template renders different strategies
{{- if eq .HEALTHCHECK_TYPE "http" }}
# HTTP health check
curl -f http://localhost:{{.HTTP_PORT}}{{.HTTP_PATH}} || exit 1
{{- else if eq .HEALTHCHECK_TYPE "tcp" }}
# TCP health check
nc -z localhost {{.TCP_PORT}} || exit 1
{{- else if eq .HEALTHCHECK_TYPE "exec" }}
# Exec health check
{{ .EXEC_COMMAND }}
{{- else if eq .HEALTHCHECK_TYPE "custom" }}
# Custom health check
{{ .CUSTOM_SCRIPT }}
{{- end }}
```

**Benefits**:
- ✅ Flexible health check configuration
- ✅ Easy to add new strategies
- ✅ No code changes required
- ✅ User-defined custom strategies

### 4. Builder Pattern 🔨

**Purpose**: Construct complex variable objects step by step.

**Implementation**: `Variables` struct with fluent API.

**Code Location**: [`pkg/generator/variables.go`](../pkg/generator/variables.go)

```go
type Variables struct {
    ServiceName    string
    ServicePort    int
    DeployDir      string
    // ... more fields
    CIPaths        *CIPaths
}

// Build variables from config
func NewVariables(cfg *config.ServiceConfig) *Variables {
    vars := &Variables{
        ServiceName: cfg.Service.Name,
        ServicePort: cfg.Service.Ports[0].Port,
        DeployDir:   cfg.Service.DeployDir,
        // ... initialize all fields
    }
    
    // Initialize CI paths
    vars.CIPaths = NewCIPaths(cfg.CI)
    
    return vars
}

// Convert to map for template rendering
func (v *Variables) ToMap() map[string]interface{} {
    m := map[string]interface{}{
        "SERVICE_NAME": v.ServiceName,
        "SERVICE_PORT": v.ServicePort,
        // ... all variables
    }
    
    // Merge CI paths
    for k, v := range v.CIPaths.ToMap() {
        m[k] = v
    }
    
    return m
}
```

**Benefits**:
- ✅ Clean variable construction
- ✅ Easy to extend with new variables
- ✅ Type-safe variable access
- ✅ Centralized variable management

### 5. Centralized Configuration Management ⚙️

**Purpose**: Single source of truth for all CI/CD paths.

**Implementation**: `CIPaths` structure.

**Code Location**: [`pkg/generator/paths.go`](../pkg/generator/paths.go)

```go
// Default path patterns
const (
    DefaultCIScriptDirPattern       = ".tad/build/%s"      // %s = service-name
    DefaultCIBuildConfigDirPattern  = "%s/build"           // %s = script_dir
    DefaultConfigTemplateDirPattern = "%s/config_template" // %s = script_dir
)

type CIPaths struct {
    ScriptDir         string
    BuildConfigDir    string
    ConfigTemplateDir string
}

// Create with defaults or custom config
func NewCIPaths(cfg *config.ServiceConfig) *CIPaths {
    serviceName := cfg.Service.Name
    defaultScriptDir := fmt.Sprintf(DefaultCIScriptDirPattern, serviceName)
    
    paths := &CIPaths{
        ScriptDir:      defaultScriptDir,
        BuildConfigDir: DefaultCIBuildConfigDir,
    }
    
    if cfg != nil && cfg.ScriptDir != "" {
        paths.ScriptDir = cfg.ScriptDir
    }
    if cfg != nil && cfg.BuildConfigDir != "" {
        paths.BuildConfigDir = cfg.BuildConfigDir
    }
    
    return paths
}

// Get specific paths
func (p *CIPaths) GetBuildScriptPath() string {
    return filepath.Join(p.ScriptDir, "build.sh")
}
```

**Benefits**:
- ✅ Single source of truth for paths
- ✅ Easy to customize paths
- ✅ Consistent path usage across all generators
- ✅ No hardcoded paths in templates

### 6. Registry Pattern 📋

**Purpose**: Register and discover generators at runtime.

**Implementation**: Global generator registry.

```go
var generatorRegistry = make(map[GeneratorType]GeneratorCreator)

func RegisterGenerator(genType GeneratorType, creator GeneratorCreator) {
    generatorRegistry[genType] = creator
}

func init() {
    // Auto-registration at startup
    RegisterGenerator(GeneratorTypeDockerfile, createDockerfileGenerator)
    RegisterGenerator(GeneratorTypeMakefile, createMakefileGenerator)
    // ... more registrations
}
```

**Benefits**:
- ✅ Automatic generator discovery
- ✅ Plugin-like extensibility
- ✅ No central registration file needed
- ✅ Easy to add third-party generators

## Component Details

### CLI Layer (`cmd/svcgen/`)

**Responsibilities**:
- Parse command-line arguments
- Route commands to appropriate handlers
- Display user-friendly output
- Handle errors gracefully

**Key Components**:
- `main.go` - Entry point
- `commands/root.go` - Root command setup
- `commands/init.go` - Initialize configuration
- `commands/validate.go` - Validate configuration
- `commands/generate.go` - Generate code

### Configuration Layer (`pkg/config/`)

**Responsibilities**:
- Load YAML configuration
- Parse into strongly-typed structures
- Validate configuration completeness
- Provide default values

**Key Components**:
- `types.go` - Configuration structures
- `loader.go` - YAML loading logic
- `validator.go` - Validation rules

**Validation Rules**:
```go
func (v *Validator) validateService() {
    // Service name required
    if v.config.Service.Name == "" {
        v.addError("service.name is required")
    }
    
    // At least one port required
    if len(v.config.Service.Ports) == 0 {
        v.addError("at least one port must be configured")
    }
    
    // Validate port ranges
    for _, port := range v.config.Service.Ports {
        if port.Port < 1 || port.Port > 65535 {
            v.addError("port must be between 1 and 65535")
        }
    }
}
```

### Generator Layer (`pkg/generator/`)

**Responsibilities**:
- Orchestrate code generation
- Manage generator lifecycle
- Coordinate template rendering
- Write output files

**Key Components**:
- `generator.go` - Main orchestrator
- `factory.go` - Generator factory
- `template.go` - Template engine wrapper
- `variables.go` - Variable management
- `paths.go` - Path management
- `template_*.go` - Individual generators

**Generator Interface**:
```go
type TemplateGenerator interface {
    Generate() (string, error)
    GetOutputPath() string
}
```

### Template Layer (`pkg/generator/templates/`)

**Responsibilities**:
- Define file templates
- Support variable substitution
- Enable conditional rendering
- Provide template functions

**Template Features**:
- Go `text/template` syntax
- Sprig function library
- Custom helper functions
- Embedded in binary

**Example Template**:
```go
//go:embed templates/dockerfile.tmpl
var dockerfileTemplate string

func (g *DockerfileTemplateGenerator) getTemplate() string {
    return dockerfileTemplate
}
```

## Data Flow

### Generation Flow

```
1. User runs: svcgen generate
                    │
                    ▼
2. CLI parses arguments
                    │
                    ▼
3. Load service.yaml
                    │
                    ▼
4. Validate configuration
                    │
                    ▼
5. Create GeneratorFactory
                    │
                    ▼
6. Initialize Variables
                    │
                    ▼
7. For each generator type:
   ├─ Create generator instance
   ├─ Render template with variables
   ├─ Write output file
   └─ Report progress
                    │
                    ▼
8. Display summary
```

### Variable Resolution Flow

```
1. Load ServiceConfig from YAML
                    │
                    ▼
2. Create Variables struct
   ├─ Extract service info
   ├─ Extract language config
   ├─ Extract build config
   ├─ Initialize CI paths
   └─ Calculate derived values
                    │
                    ▼
3. Convert to map for templates
   ├─ Add service variables
   ├─ Add build variables
   ├─ Add plugin variables
   └─ Add CI path variables
                    │
                    ▼
4. Merge with generator-specific vars
                    │
                    ▼
5. Render template
```

## Extension Points

### Adding a New Generator

1. **Create generator file**: `pkg/generator/template_myfile.go`

```go
package generator

import "embed"

//go:embed templates/myfile.tmpl
var myfileTemplate string

type MyFileTemplateGenerator struct {
    BaseTemplateGenerator
}

func NewMyFileTemplateGenerator(
    cfg *config.ServiceConfig,
    engine *TemplateEngine,
    vars *Variables,
) *MyFileTemplateGenerator {
    return &MyFileTemplateGenerator{
        BaseTemplateGenerator: BaseTemplateGenerator{
            config:         cfg,
            templateEngine: engine,
            variables:      vars,
        },
    }
}

func (g *MyFileTemplateGenerator) Generate() (string, error) {
    vars := g.variables.ToMap()
    // Add generator-specific variables
    vars["MY_VAR"] = "value"
    
    return g.RenderTemplate(g.getTemplate(), vars)
}

func (g *MyFileTemplateGenerator) GetOutputPath() string {
    return "myfile.txt"
}

func (g *MyFileTemplateGenerator) getTemplate() string {
    return myfileTemplate
}
```

2. **Create template**: `pkg/generator/templates/myfile.tmpl`

```
# My File
Service: {{ .SERVICE_NAME }}
Port: {{ .SERVICE_PORT }}
```

3. **Register generator**: `pkg/generator/factory.go`

```go
const (
    GeneratorTypeMyFile GeneratorType = "myfile"
)

func init() {
    RegisterGenerator(GeneratorTypeMyFile, createMyFileGenerator)
}

func createMyFileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, arch string) (TemplateGenerator, error) {
    return NewMyFileTemplateGenerator(cfg, engine, vars), nil
}
```

4. **Add to generator list**: `pkg/generator/generator.go`

```go
func (g *Generator) GenerateAll() error {
    generators := []TemplateGenerator{
        // ... existing generators
        g.factory.CreateGenerator(GeneratorTypeMyFile, ""),
    }
    // ... generate all
}
```

### Adding a New Language

1. **Update language types**: `pkg/config/types.go`

```go
const (
    LanguageTypeGo     = "go"
    LanguageTypePython = "python"
    LanguageTypeMyLang = "mylang"  // Add new language
)
```

2. **Add validation**: `pkg/config/validator.go`

```go
func (v *Validator) validateLanguage() {
    validLanguages := []string{"go", "python", "nodejs", "java", "rust", "mylang"}
    // ... validation logic
}
```

3. **Update templates**: Add language-specific logic in templates

```
{{- if eq .LANGUAGE_TYPE "mylang" }}
# MyLang specific configuration
{{- end }}
```

### Adding Custom Template Functions

1. **Add function to template engine**: `pkg/generator/template.go`

```go
func (e *TemplateEngine) Render(templateStr string, data map[string]interface{}) (string, error) {
    funcMap := sprig.TxtFuncMap()
    
    // Add custom functions
    funcMap["myCustomFunc"] = func(s string) string {
        return strings.ToUpper(s)
    }
    
    tmpl, err := template.New("template").Funcs(funcMap).Parse(templateStr)
    // ... render
}
```

2. **Use in templates**:

```
{{ myCustomFunc .SERVICE_NAME }}
```

---

[← Back to README](../README.md)
