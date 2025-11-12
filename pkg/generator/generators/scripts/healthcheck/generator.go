package healthcheck

import (
	_ "embed"
	"fmt"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "healthcheck-script"

// init registers the healthcheck script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates healthchk.sh script
type Generator struct {
	core.BaseGenerator
	strategy Strategy
}

// New creates a new healthcheck script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()

	// Create strategy factory
	factory := NewStrategyFactory(ctx.Config)

	// Create appropriate strategy
	strategy, err := factory.CreateStrategy()
	if err != nil {
		return nil, err
	}

	// Validate strategy configuration
	if err := strategy.Validate(); err != nil {
		return nil, err
	}

	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
		strategy:      strategy,
	}, nil
}

// Generate generates healthchk.sh content using the selected strategy
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Use preset for script
	composer := ctx.GetVariablePreset().ForScript()

	// Add healthcheck-specific custom variable
	composer.WithCustom("CUSTOM_SCRIPT", ctx.Config.Runtime.Healthcheck.CustomScript)

	vars := composer.Build()

	// Get script template from strategy
	scriptTemplate, err := g.strategy.GenerateScript(vars)
	if err != nil {
		return "", err
	}

	// Render the template with variables
	return g.RenderTemplate(scriptTemplate, vars)
}

// GetStrategy returns the current healthcheck strategy
func (g *Generator) GetStrategy() Strategy {
	return g.strategy
}

// Strategy defines the interface for health check strategies
type Strategy interface {
	// GetType returns the strategy type
	GetType() string
	// GenerateScript generates the health check script content
	GenerateScript(vars map[string]interface{}) (string, error)
	// Validate validates the strategy configuration
	Validate() error
}

// StrategyFactory creates health check strategies
type StrategyFactory struct {
	config *config.ServiceConfig
}

// NewStrategyFactory creates a new factory
func NewStrategyFactory(cfg *config.ServiceConfig) *StrategyFactory {
	return &StrategyFactory{
		config: cfg,
	}
}

// CreateStrategy creates a health check strategy based on configuration
func (f *StrategyFactory) CreateStrategy() (Strategy, error) {
	if !f.config.Runtime.Healthcheck.Enabled {
		return NewDefaultStrategy(f.config), nil
	}

	switch f.config.Runtime.Healthcheck.Type {
	case "default", "":
		return NewDefaultStrategy(f.config), nil
	case "custom":
		return NewCustomStrategy(f.config), nil
	default:
		return nil, fmt.Errorf("unsupported healthcheck type: %s (valid: default, custom)", f.config.Runtime.Healthcheck.Type)
	}
}

// DefaultStrategy implements default process-based health check
type DefaultStrategy struct {
	config *config.ServiceConfig
}

// NewDefaultStrategy creates a new default strategy
func NewDefaultStrategy(cfg *config.ServiceConfig) *DefaultStrategy {
	return &DefaultStrategy{
		config: cfg,
	}
}

// GetType returns the strategy type
func (s *DefaultStrategy) GetType() string {
	return "default"
}

// GenerateScript generates default health check script
func (s *DefaultStrategy) GenerateScript(vars map[string]interface{}) (string, error) {
	script := `#!/bin/sh

# Export service paths as environment variables
export SERVICE_ROOT="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}"
export SERVICE_BIN_DIR="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}/bin"
export SERVICE_NAME="{{ .SERVICE_NAME }}"

# Default healthcheck: check if service process is running
ps=$(ls -l /proc/*/exe 2>/dev/null | grep "${SERVICE_NAME}" | grep -v grep)

# abnormal
[[ "$ps" == "" ]] && exit 1

# normal
exit 0
`
	return script, nil
}

// Validate validates the default strategy configuration
func (s *DefaultStrategy) Validate() error {
	// Default strategy has no additional validation requirements
	return nil
}

// CustomStrategy implements custom user-defined health check
type CustomStrategy struct {
	config *config.ServiceConfig
}

// NewCustomStrategy creates a new custom strategy
func NewCustomStrategy(cfg *config.ServiceConfig) *CustomStrategy {
	return &CustomStrategy{
		config: cfg,
	}
}

// GetType returns the strategy type
func (s *CustomStrategy) GetType() string {
	return "custom"
}

// GenerateScript generates custom health check script
func (s *CustomStrategy) GenerateScript(vars map[string]interface{}) (string, error) {
	if s.config.Runtime.Healthcheck.CustomScript == "" {
		return "", fmt.Errorf("custom_script is required for custom healthcheck type")
	}

	script := `#!/bin/sh

# Export service paths as environment variables
export SERVICE_ROOT="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}"
export SERVICE_BIN_DIR="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}/bin"
export SERVICE_NAME="{{ .SERVICE_NAME }}"

# Custom healthcheck script
{{ .CUSTOM_SCRIPT }}
`
	return script, nil
}

// Validate validates the custom strategy configuration
func (s *CustomStrategy) Validate() error {
	if s.config.Runtime.Healthcheck.CustomScript == "" {
		return fmt.Errorf("runtime.healthcheck.custom_script is required when type is 'custom'")
	}
	return nil
}

//go:embed templates/healthcheck.sh.tmpl
var template string
