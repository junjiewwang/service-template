package generator

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeRtPrepareScript = "rt-prepare-script"

// init registers the rt prepare script generator
func init() {
	RegisterGenerator(GeneratorTypeRtPrepareScript, createRtPrepareScriptGenerator)
}

// RtPrepareScriptTemplateGenerator generates rt_prepare.sh script
type RtPrepareScriptTemplateGenerator struct {
	BaseTemplateGenerator
}

// createRtPrepareScriptGenerator is the creator function for RtPrepareScript generator
func createRtPrepareScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewRtPrepareScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewRtPrepareScriptTemplateGenerator creates a new rt prepare script generator
func NewRtPrepareScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *RtPrepareScriptTemplateGenerator {
	return &RtPrepareScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeRtPrepareScript,
		},
	}
}

// Generate generates rt_prepare.sh content
func (g *RtPrepareScriptTemplateGenerator) Generate() (string, error) {
	vars := map[string]interface{}{
		"RUNTIME_DEPS_PACKAGES": g.config.Runtime.SystemDependencies.Packages,
	}
	return g.RenderTemplate(rtPrepareScriptTemplate, vars)
}

//go:embed templates/rt_prepare.sh.tmpl
var rtPrepareScriptTemplate string
