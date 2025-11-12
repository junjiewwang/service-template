# Variable Management with Flyweight Pattern

## Overview

The variable management system uses the **Flyweight Pattern** to efficiently manage template variables across all generators. This design eliminates duplication, improves performance, and makes the codebase more maintainable.

## Architecture

### Core Components

1. **VariablePool** - Manages shared variable sets (Flyweight Pool)
2. **SharedVariables** - Immutable variable sets for different categories (Flyweight Objects)
3. **VariableComposer** - Fluent API for composing variables (Client)
4. **VariablePreset** - Pre-configured variable combinations for common scenarios

### Variable Categories

The system organizes variables into 7 categories:

| Category | Description | Example Variables |
|----------|-------------|-------------------|
| `common` | Shared by all generators | SERVICE_NAME, DEPLOY_DIR, GENERATED_AT |
| `build` | Build-related variables | BUILD_COMMAND, BUILDER_IMAGE |
| `runtime` | Runtime-related variables | STARTUP_COMMAND, HEALTHCHECK_ENABLED |
| `plugin` | Plugin-related variables | PLUGIN_ROOT_DIR, PLUGINS |
| `ci-paths` | CI path variables | CI_SCRIPT_DIR, CI_CONTAINER_DIR |
| `service` | Service-related variables | PORTS, SERVICE_PORT |
| `language` | Language-related variables | LANGUAGE, LANGUAGE_VERSION |

## Usage

### Basic Usage

```go
// In your generator
func (g *Generator) Generate() (string, error) {
    ctx := g.GetContext()
    
    // Compose variables using fluent API
    vars := ctx.GetVariableComposer().
        WithCommon().      // Add common variables
        WithBuild().       // Add build variables
        WithRuntime().     // Add runtime variables
        Build()            // Get final map
    
    return g.RenderTemplate(template, vars)
}
```

### Using All Standard Variables

```go
vars := ctx.GetVariableComposer().
    WithAll().  // Adds all 7 categories
    Build()
```

### Adding Architecture-Specific Variables

```go
vars := ctx.GetVariableComposer().
    WithCommon().
    WithBuild().
    WithArchitecture("amd64").  // Adds GOARCH, GOOS, BUILDER_IMAGE, RUNTIME_IMAGE
    Build()
```

### Adding Custom Variables

```go
vars := ctx.GetVariableComposer().
    WithCommon().
    WithCustom("MY_VAR", "my_value").
    WithCustomMap(map[string]interface{}{
        "VAR1": "value1",
        "VAR2": 123,
    }).
    Build()
```

### Overriding Variables

```go
vars := ctx.GetVariableComposer().
    WithCommon().
    Override("SERVICE_NAME", "custom-name").  // Override existing variable
    Build()
```

### Using Presets

Presets provide pre-configured variable combinations for common scenarios:

```go
// For Dockerfile generation
vars := ctx.GetVariablePreset().
    ForDockerfile("amd64").
    Build()

// For build script generation
vars := ctx.GetVariablePreset().
    ForBuildScript().
    Build()

// For docker-compose generation
vars := ctx.GetVariablePreset().
    ForCompose().
    Build()

// For Makefile generation
vars := ctx.GetVariablePreset().
    ForMakefile().
    Build()

// For DevOps configuration
vars := ctx.GetVariablePreset().
    ForDevOps().
    Build()

// For script generation
vars := ctx.GetVariablePreset().
    ForScript().
    Build()
```

### Checking and Getting Variables

```go
composer := ctx.GetVariableComposer().WithCommon()

// Check if variable exists
if composer.Has("SERVICE_NAME") {
    // Get variable value
    if val, ok := composer.Get("SERVICE_NAME"); ok {
        fmt.Println(val)
    }
}

// Get size
size := composer.Size()
```

### Cloning Composer

```go
original := ctx.GetVariableComposer().WithCommon()
cloned := original.Clone()

// Modify clone without affecting original
cloned.WithCustom("NEW_VAR", "value")
```

## Real-World Examples

### Example 1: Dockerfile Generator

```go
func (g *Generator) prepareTemplateVars() map[string]interface{} {
    ctx := g.GetContext()
    
    return ctx.GetVariableComposer().
        WithCommon().
        WithBuild().
        WithRuntime().
        WithPlugin().
        WithCIPaths().
        WithService().
        WithLanguage().
        WithArchitecture(g.arch).
        WithCustom("PKG_MANAGER", detectPackageManager(...)).
        WithCustom("DEPENDENCY_FILES", getDependencyFilesList(...)).
        Build()
}
```

### Example 2: Build Script Generator

```go
func (g *Generator) Generate() (string, error) {
    ctx := g.GetContext()
    
    vars := ctx.GetVariableComposer().
        WithCommon().
        WithBuild().
        WithPlugin().
        WithCIPaths().
        Build()
    
    return g.RenderTemplate(template, vars)
}
```

### Example 3: Compose Generator with Custom Logic

```go
func (g *Generator) prepareTemplateVars() map[string]interface{} {
    ctx := g.GetContext()
    
    // Prepare custom port mappings
    ports := g.preparePortMappings()
    
    // Prepare custom volumes
    volumes := g.prepareVolumes()
    
    return ctx.GetVariableComposer().
        WithCommon().
        WithRuntime().
        WithService().
        WithCustom("PORTS", ports).
        WithCustom("VOLUMES", volumes).
        Build()
}
```

## Benefits

### 1. Eliminates Duplication
- Shared variables are created once and reused
- No more copy-paste of variable definitions

### 2. Improves Performance
- Variable pool caching mechanism
- Thread-safe with read-write locks
- Immutable shared variables prevent accidental modifications

### 3. Enhances Maintainability
- Centralized variable management
- Clear variable categorization
- Easy to track variable usage

### 4. Increases Extensibility
- Add new variable categories easily
- Generators compose only what they need
- Custom variables for generator-specific needs

### 5. Better Readability
- Fluent API clearly expresses intent
- Self-documenting code
- Consistent patterns across generators

### 6. Type Safety
- Variable keys use constants (VarServiceName, etc.)
- Reduces typos and errors

## Migration Guide

### Old Way (Before Flyweight Pattern)

```go
func (g *Generator) prepareTemplateVars() map[string]interface{} {
    ctx := g.GetContext()
    vars := make(map[string]interface{})
    
    // Manually add each variable
    vars["SERVICE_NAME"] = ctx.Config.Service.Name
    vars["DEPLOY_DIR"] = ctx.Config.Service.DeployDir
    vars["GENERATED_AT"] = ctx.Config.Metadata.GeneratedAt
    vars["BUILD_COMMAND"] = ctx.Config.Build.Commands.Build
    // ... many more lines
    
    return vars
}
```

### New Way (With Flyweight Pattern)

```go
func (g *Generator) prepareTemplateVars() map[string]interface{} {
    ctx := g.GetContext()
    
    return ctx.GetVariableComposer().
        WithCommon().
        WithBuild().
        Build()
}
```

## Constants

All variable keys are defined as constants in `context/constants.go`:

```go
const (
    // Service variables
    VarServiceName   = "SERVICE_NAME"
    VarServicePort   = "SERVICE_PORT"
    VarServiceRoot   = "SERVICE_ROOT"
    
    // Build variables
    VarBuildCommand     = "BUILD_COMMAND"
    VarPreBuildCommand  = "PRE_BUILD_COMMAND"
    
    // Plugin variables
    VarPluginRootDir = "PLUGIN_ROOT_DIR"
    
    // ... more constants
)
```

Use these constants instead of string literals to avoid typos.

## Testing

The variable management system is fully tested:

```bash
# Run variable management tests
go test ./pkg/generator/context/... -v

# Run all generator tests
go test ./pkg/generator/... -v
```

## Performance Characteristics

- **Memory**: Shared variables are cached, reducing memory usage
- **CPU**: Variable creation happens once per category
- **Concurrency**: Thread-safe with RWMutex
- **Scalability**: O(1) lookup for cached variables

## Best Practices

1. **Use Presets When Possible**: They provide optimal variable combinations
2. **Compose Only What You Need**: Don't use `WithAll()` if you only need a few categories
3. **Use Constants**: Always use `VarServiceName` instead of `"SERVICE_NAME"`
4. **Add Custom Variables Last**: Compose shared variables first, then add custom ones
5. **Don't Modify Shared Variables**: They are immutable by design
6. **Use Override Sparingly**: Only override when absolutely necessary

## Troubleshooting

### Variable Not Found

```go
// Check if variable exists
composer := ctx.GetVariableComposer().WithCommon()
if !composer.Has("MY_VAR") {
    // Add it
    composer.WithCustom("MY_VAR", "value")
}
```

### Wrong Variable Value

```go
// Override the value
vars := ctx.GetVariableComposer().
    WithCommon().
    Override("SERVICE_NAME", "correct-name").
    Build()
```

### Need All Variables

```go
// Use WithAll() for debugging
vars := ctx.GetVariableComposer().WithAll().Build()
for k, v := range vars {
    fmt.Printf("%s = %v\n", k, v)
}
```

## Future Enhancements

Potential improvements for the future:

1. **Variable Validation**: Validate required variables before rendering
2. **Variable Tracking**: Track variable sources for debugging
3. **More Presets**: Add presets for more scenarios
4. **Variable Documentation**: Auto-generate variable documentation
5. **Variable Inheritance**: Support variable inheritance between categories

## References

- Design Pattern: [Flyweight Pattern](https://refactoring.guru/design-patterns/flyweight)
- Code Location: `pkg/generator/context/`
- Tests: `pkg/generator/context/*_test.go`
