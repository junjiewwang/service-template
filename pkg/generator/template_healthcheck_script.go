package generator

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeHealthcheckScript = "healthcheck-script"

// init registers the healthcheck script generator
func init() {
	RegisterGenerator(GeneratorTypeHealthcheckScript, createHealthcheckScriptGenerator)
}

// HealthcheckScriptTemplateGenerator generates healthchk.sh script
type HealthcheckScriptTemplateGenerator struct {
	BaseTemplateGenerator
	strategy HealthcheckStrategy
}

// createHealthcheckScriptGenerator is the creator function for HealthcheckScript generator
func createHealthcheckScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewHealthcheckScriptTemplateGenerator(cfg, engine, vars)
}

// NewHealthcheckScriptTemplateGenerator creates a new healthcheck script generator
func NewHealthcheckScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) (*HealthcheckScriptTemplateGenerator, error) {
	// Create strategy factory
	factory := NewHealthcheckStrategyFactory(cfg)

	// Create appropriate strategy
	strategy, err := factory.CreateStrategy()
	if err != nil {
		return nil, err
	}

	// Validate strategy configuration
	if err := strategy.Validate(); err != nil {
		return nil, err
	}

	return &HealthcheckScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeHealthcheckScript,
		},
		strategy: strategy,
	}, nil
}

//go:embed templates/healthcheck.sh.tmpl
var healthcheckScriptTemplate string

// Generate generates healthchk.sh content using the selected strategy
func (g *HealthcheckScriptTemplateGenerator) Generate() (string, error) {
	// Prepare template variables
	vars := map[string]interface{}{
		"SERVICE_NAME":  g.config.Service.Name,
		"DEPLOY_DIR":    g.config.Service.DeployDir,
		"CUSTOM_SCRIPT": g.config.Runtime.Healthcheck.CustomScript,
	}

	// Get script template from strategy
	scriptTemplate, err := g.strategy.GenerateScript(vars)
	if err != nil {
		return "", err
	}

	// Render the template with variables
	return g.RenderTemplate(scriptTemplate, vars)
}

// getTemplate returns the healthcheck script template
func (g *HealthcheckScriptTemplateGenerator) getTemplate() string {
	return healthcheckScriptTemplate
}

// GetStrategy returns the current healthcheck strategy
func (g *HealthcheckScriptTemplateGenerator) GetStrategy() HealthcheckStrategy {
	return g.strategy
}
