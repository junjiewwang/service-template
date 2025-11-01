package generator

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeMakefile = "makefile"

// MakefileTemplateGenerator generates Makefile using factory pattern
type MakefileTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the Makefile generator
func init() {
	RegisterGenerator(GeneratorTypeMakefile, createMakefileGenerator)
}

// createMakefileGenerator is the creator function for Makefile generator
func createMakefileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewMakefileTemplateGenerator(cfg, engine, vars), nil
}

// NewMakefileTemplateGenerator creates a new Makefile template generator
func NewMakefileTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *MakefileTemplateGenerator {
	return &MakefileTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeMakefile,
		},
	}
}

//go:embed templates/makefile.tmpl
var makefileTemplate string

// Generate generates Makefile content
func (g *MakefileTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for Makefile template
func (g *MakefileTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name
	vars["SERVICE_PORT"] = g.variables.ServicePort
	vars["OUTPUT_DIR"] = g.config.Build.OutputDir

	// Kubernetes config
	vars["K8S_ENABLED"] = g.config.LocalDev.Kubernetes.Enabled
	vars["K8S_NAMESPACE"] = g.config.LocalDev.Kubernetes.Namespace
	vars["K8S_OUTPUT_DIR"] = g.config.LocalDev.Kubernetes.OutputDir
	vars["K8S_WAIT_ENABLED"] = g.config.LocalDev.Kubernetes.Wait.Enabled
	vars["K8S_WAIT_TIMEOUT"] = g.config.LocalDev.Kubernetes.Wait.Timeout
	vars["K8S_VOLUME_TYPE"] = g.config.LocalDev.Kubernetes.VolumeType

	// Custom targets
	vars["CUSTOM_TARGETS"] = g.config.Makefile.CustomTargets

	return vars
}

// getTemplate returns the Makefile template
func (g *MakefileTemplateGenerator) getTemplate() string {
	return makefileTemplate
}
