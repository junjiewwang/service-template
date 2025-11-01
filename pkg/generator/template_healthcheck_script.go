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
}

// createHealthcheckScriptGenerator is the creator function for HealthcheckScript generator
func createHealthcheckScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewHealthcheckScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewHealthcheckScriptTemplateGenerator creates a new healthcheck script generator
func NewHealthcheckScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *HealthcheckScriptTemplateGenerator {
	return &HealthcheckScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeHealthcheckScript,
		},
	}
}

//go:embed templates/healthcheck.sh.tmpl
var healthcheckScriptTemplate string

// Generate generates healthchk.sh content
func (g *HealthcheckScriptTemplateGenerator) Generate() (string, error) {
	vars := map[string]interface{}{
		"SERVICE_NAME": g.config.Service.Name,
		"DEPLOY_DIR":   g.config.Service.DeployDir,
	}
	return g.RenderTemplate(g.getTemplate(), vars)
}

// getTemplate returns the healthcheck script template
func (g *HealthcheckScriptTemplateGenerator) getTemplate() string {
	return healthcheckScriptTemplate
}
