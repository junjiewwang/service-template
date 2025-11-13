package build_plugins

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/domain/services"
)

const GeneratorType = "build-plugins-script"

// init registers the build plugins script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates build_plugins.sh script
type Generator struct {
	core.BaseGenerator
	engine *core.TemplateEngine
}

// New creates a new build plugins script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
		engine:        engine,
	}, nil
}

// Generate generates build_plugins.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Check if plugins exist
	pluginService := services.NewPluginService(ctx, g.engine)
	if !pluginService.HasPlugins() {
		// Return empty string if no plugins, caller should skip file generation
		return "", nil
	}

	// Use preset for build script
	composer := ctx.GetVariablePreset().ForBuildScript()

	// Prepare plugin data
	plugins := pluginService.PrepareForBuildScript()
	composer.Override("PLUGINS", plugins)

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/build_plugins.sh.tmpl
var template string
