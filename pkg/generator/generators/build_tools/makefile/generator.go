package makefile

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "makefile"

// init registers the makefile generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates Makefile
type Generator struct {
	core.BaseGenerator
}

// New creates a new makefile generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates Makefile content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	vars := g.prepareTemplateVars()
	return g.RenderTemplate(template, vars)
}

// prepareTemplateVars prepares variables for Makefile template
func (g *Generator) prepareTemplateVars() map[string]interface{} {
	ctx := g.GetContext()

	// Use preset for Makefile
	composer := ctx.GetVariablePreset().ForMakefile()

	// Add Makefile-specific custom variables
	composer.
		WithCustom("K8S_ENABLED", ctx.Config.LocalDev.Kubernetes.Enabled).
		WithCustom("K8S_NAMESPACE", ctx.Config.LocalDev.Kubernetes.Namespace).
		WithCustom("K8S_OUTPUT_DIR", ctx.Config.LocalDev.Kubernetes.OutputDir).
		WithCustom("K8S_WAIT_ENABLED", ctx.Config.LocalDev.Kubernetes.Wait.Enabled).
		WithCustom("K8S_WAIT_TIMEOUT", ctx.Config.LocalDev.Kubernetes.Wait.Timeout).
		WithCustom("K8S_VOLUME_TYPE", ctx.Config.LocalDev.Kubernetes.VolumeType).
		WithCustom("CUSTOM_TARGETS", ctx.Config.Makefile.CustomTargets)

	return composer.Build()
}

//go:embed templates/makefile.tmpl
var template string
