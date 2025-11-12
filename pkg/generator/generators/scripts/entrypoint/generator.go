package entrypoint

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/domain/services"
)

const GeneratorType = "entrypoint-script"

// init registers the entrypoint script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates entrypoint.sh script
type Generator struct {
	core.BaseGenerator
}

// New creates a new entrypoint script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates entrypoint.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Use preset for script
	composer := ctx.GetVariablePreset().ForScript()

	// Use plugin service to prepare plugin environment variables
	pluginService := services.NewPluginService(ctx, g.GetEngine())
	pluginEnvs := pluginService.PrepareForEntrypoint()

	// Add script-specific custom variables
	composer.
		WithCustom("PLUGINS_ENV", pluginEnvs).
		WithCustom("HAS_PLUGINS_ENV", len(pluginEnvs) > 0)

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/entrypoint.sh.tmpl
var template string
