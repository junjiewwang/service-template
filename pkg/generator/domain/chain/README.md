# Chain Package - Responsibility Chain Infrastructure

## Overview

The `chain` package provides a flexible and extensible responsibility chain infrastructure for processing domain-specific configurations. It implements the Chain of Responsibility pattern combined with the Abstract Factory pattern to enable clean separation of concerns and easy extensibility.

## Core Concepts

### 1. Handler

The base interface for all handlers in the chain. Three specialized handler types are provided:

- **ParserHandler**: Parses configuration from raw data
- **ValidatorHandler**: Validates parsed configuration
- **GeneratorHandler**: Generates files based on configuration

### 2. ProcessingContext

A thread-safe context that flows through the chain, holding:
- Raw configuration data
- Parsed domain models
- Validation errors
- Generated files
- Metadata
- Error aggregation

### 3. DomainFactory

An abstract factory interface for creating domain-specific handlers. Each domain (e.g., service, language, build) implements this interface to provide its own handlers.

### 4. DomainRegistry

Manages domain factory registration and builds chains automatically based on factory priorities.

### 5. PriorityChain

Provides a fluent API for defining domain processing order explicitly:

```go
chain := NewPriorityChain().
    First(serviceFactory).
    Then(languageFactory).
    ThenAll(pluginFactory, runtimeFactory).
    Finally(crossDomainFactory)
```

### 6. DependencyGraph

Manages domain dependencies and performs topological sorting to ensure dependencies are processed before dependents.

## Usage Examples

### Example 1: Using DomainRegistry (Simple)

```go
// Create registry
registry := chain.NewDomainRegistry()

// Register domain factories
registry.RegisterAll(
    service.NewServiceDomainFactory(),
    language.NewLanguageDomainFactory(),
    build.NewBuildDomainFactory(),
)

// Build chains automatically (sorted by priority)
parseChain := registry.BuildParseChain()
validateChain := registry.BuildValidateChain()
generateChain := registry.BuildGenerateChain()

// Execute chains
ctx := chain.NewProcessingContext(context.Background(), rawConfig)
parseChain.Handle(ctx)
validateChain.Handle(ctx)
generateChain.Handle(ctx)
```

### Example 2: Using PriorityChain (Explicit Order)

```go
// Define priority using factory instances
priorityChain := chain.NewPriorityChain().
    First(serviceFactory).              // Service info first
    Then(languageFactory).              // Language config
    Then(buildFactory).                 // Build config
    ThenAll(pluginFactory, runtimeFactory). // Parallel processing
    Finally(crossDomainFactory)         // Cross-domain validation

// Build chains
parseChain := priorityChain.BuildParseChain()
validateChain := priorityChain.BuildValidateChain()
generateChain := priorityChain.BuildGenerateChain()
```

### Example 3: Using DependencyGraph (Complex Dependencies)

```go
// Define dependencies using factory instances
graph := chain.NewDependencyGraph().
    AddNode(serviceFactory).
    AddNode(languageFactory, serviceFactory).
    AddNode(buildFactory, serviceFactory, languageFactory).
    AddNode(pluginFactory, buildFactory).
    AddNode(runtimeFactory, buildFactory).
    AddNode(localdevFactory, pluginFactory, runtimeFactory)

// Validate dependencies
if err := graph.Validate(); err != nil {
    log.Fatal(err)
}

// Build chains (automatically sorted by dependencies)
parseChain, _ := graph.BuildParseChain()
validateChain, _ := graph.BuildValidateChain()
generateChain, _ := graph.BuildGenerateChain()
```

## Implementing a Domain Factory

```go
package service

import "github.com/junjiewwang/service-template/pkg/generator/domain/chain"

// ServiceDomainFactory creates service domain handlers
type ServiceDomainFactory struct {
    *chain.BaseDomainFactory
}

// NewServiceDomainFactory creates a new service domain factory
func NewServiceDomainFactory() chain.DomainFactory {
    return &ServiceDomainFactory{
        BaseDomainFactory: chain.NewBaseDomainFactory("service", 10),
    }
}

// CreateParserHandler creates a service parser handler
func (f *ServiceDomainFactory) CreateParserHandler() chain.ParserHandler {
    return NewServiceParserHandler()
}

// CreateValidatorHandler creates a service validator handler
func (f *ServiceDomainFactory) CreateValidatorHandler() chain.ValidatorHandler {
    return NewServiceValidatorHandler()
}

// CreateGeneratorHandler creates a service generator handler
func (f *ServiceDomainFactory) CreateGeneratorHandler() chain.GeneratorHandler {
    return NewServiceGeneratorHandler()
}
```

## Middleware Support

The package provides built-in middleware for common concerns:

### Logging Middleware

```go
chain := builder.BuildWithLogging()
```

### Metrics Middleware

```go
handler := chain.NewMetricsMiddleware(yourHandler)
// Later retrieve metrics
metrics := handler.GetMetrics()
```

### Recovery Middleware

```go
handler := chain.NewRecoveryMiddleware(yourHandler)
```

### Custom Middleware

```go
customMiddleware := func(h chain.Handler) chain.Handler {
    // Your custom logic
    return wrappedHandler
}

chain := builder.BuildWithMiddleware(customMiddleware)
```

## Key Benefits

1. **Zero Hardcoding**: Domain names come from factory instances via `GetName()`
2. **Type Safety**: Compile-time checking of factory instances
3. **Flexible Ordering**: Support both explicit (PriorityChain) and dependency-based (DependencyGraph) ordering
4. **Easy Extension**: Add new domains by implementing DomainFactory interface
5. **Thread-Safe**: ProcessingContext uses sync.Map and mutex for concurrent access
6. **Testable**: Each component can be tested independently

## Testing

Run tests with coverage:

```bash
go test -v ./pkg/generator/domain/chain/... -cover
```

Current coverage: 72.2%

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ Application Layer                                            │
│ └── Orchestrator (uses DomainRegistry/PriorityChain)        │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Chain Infrastructure (this package)                          │
│ ├── Handler Interfaces                                       │
│ ├── ProcessingContext                                        │
│ ├── ChainBuilder & Executor                                  │
│ ├── DomainFactory & Registry                                 │
│ ├── PriorityChain & DependencyGraph                          │
│ └── Middleware                                               │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Domain Layer                                                 │
│ └── Subdomains (service, language, build, plugin, etc.)     │
│     ├── Factory (implements DomainFactory)                   │
│     ├── Parser (implements ParserHandler)                    │
│     ├── Validator (implements ValidatorHandler)              │
│     └── Generator (implements GeneratorHandler)              │
└─────────────────────────────────────────────────────────────┘
```

## Next Steps

After implementing this infrastructure, the next phase is to:

1. Create subdomain implementations (service, language, build, plugin, runtime, localdev)
2. Implement the application layer orchestrator
3. Integrate with existing generators

See the main project documentation for the complete implementation plan.
