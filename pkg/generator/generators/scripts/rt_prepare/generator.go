package rt_prepare

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "rt-prepare-script"

// init registers the rt prepare script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates rt_prepare.sh script
type Generator struct {
	core.BaseGenerator
}

// New creates a new rt prepare script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates rt_prepare.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Use preset for script
	composer := ctx.GetVariablePreset().ForScript()

	// Add rt_prepare-specific custom variable
	composer.WithCustom("RUNTIME_DEPS_PACKAGES", ctx.Config.Runtime.SystemDependencies.Packages)

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/rt_prepare.sh.tmpl
var template string
