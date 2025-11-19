# Generator Package - Quick Reference

## 📁 New Structure Overview

### Core Layer (`core/`)
Foundation for all generators:
- **types.go**: `Generator` interface, `GeneratorCreator` type
- **registry.go**: `Registry` for generator registration
- **base.go**: `BaseGenerator` with common functionality
- **engine.go**: Template rendering engine
- **errors.go**: Standard error types

### Context Layer (`context/`)
Manages generation context:
- **context.go**: `GeneratorContext` - encapsulates all context
- **variables.go**: `Variables` - template variable management
- **paths.go**: `Paths` and `CIPaths` - path management
- **errors.go**: Context-specific errors

### Generators (`generators/`)
Domain-organized generator implementations:
```
generators/
├── docker/          # Container-related
│   ├── dockerfile/
│   ├── compose/
│   └── devops/
├── scripts/         # Script generation
│   ├── build/
│   ├── deps_install/
│   ├── entrypoint/
│   ├── rt_prepare/
│   └── healthcheck/
└── build_tools/     # Build tools
    └── makefile/
```

### Internal Utilities (`internal/`)
Shared utilities:
- **helpers.go**: Text manipulation, formatting
- **validators.go**: Common validation functions
- **testutil/fixtures.go**: Test fixtures

## 🔧 Usage Examples

### Creating a Generator Context
```go
import (
    "github.com/junjiewwang/service-template/pkg/generator/context"
    "github.com/junjiewwang/service-template/pkg/config"
)

// Create context
cfg := &config.ServiceConfig{...}
ctx := context.NewGeneratorContext(cfg, "/output/dir")
```

### Using the Registry
```go
import "github.com/junjiewwang/service-template/pkg/generator/core"

// Register a generator (in init())
func init() {
    core.DefaultRegistry.Register("my-generator", NewMyGenerator)
}

// Get a generator creator
creator, exists := core.DefaultRegistry.Get("my-generator")
if exists {
    gen, err := creator(ctx, options...)
}
```

### Implementing a Generator
```go
package mygenerator

import (
    "github.com/junjiewwang/service-template/pkg/generator/core"
    "github.com/junjiewwang/service-template/pkg/generator/context"
)

type Generator struct {
    core.BaseGenerator
}

func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
    engine := core.NewTemplateEngine()
    return &Generator{
        BaseGenerator: core.NewBaseGenerator("my-generator", ctx, engine),
    }, nil
}

func (g *Generator) Generate() (string, error) {
    // Use the new variable system
    composer := g.GetContext().GetVariableComposer().WithCommon()
    vars := composer.Build()
    return g.RenderTemplate(myTemplate, vars)
}

func (g *Generator) Validate() error {
    return g.BaseGenerator.Validate()
}
```

## 📋 File Organization Pattern

Each generator package should follow this structure:
```
<generator_name>/
├── generator.go         # Generator implementation
├── generator_test.go    # Generator tests
├── helpers.go           # Helper functions (optional)
├── helpers_test.go      # Helper tests (optional)
└── templates/
    └── <name>.tmpl      # Template files
```

## ✅ Checklist for New Generators

1. [ ] Create directory under appropriate domain
2. [ ] Implement `core.Generator` interface
3. [ ] Register in `init()` function
4. [ ] Add template files
5. [ ] Write tests
6. [ ] Document usage

## 🔄 Migration Status

- ✅ Phase 1: Infrastructure created
- ⏳ Phase 2: Migrate generators (TODO)
- ⏳ Phase 3: Refactor tests (TODO)
- ⏳ Phase 4: Cleanup (TODO)

## 📚 Documentation

- [REFACTORING.md](./REFACTORING.md) - Complete refactoring guide
- [PHASE1_SUMMARY.md](./PHASE1_SUMMARY.md) - Phase 1 completion summary

## 🎯 Design Principles

1. **High Cohesion**: Related code stays together
2. **Low Coupling**: Minimal dependencies
3. **Domain-Driven**: Organized by business domain
4. **Extensible**: Easy to add new generators
5. **Testable**: Each component independently testable
6. **Backward Compatible**: Old code continues to work
