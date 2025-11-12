package build

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/domain/services"
)

const GeneratorType = "build-script"

// init registers the build script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates build.sh script
type Generator struct {
	core.BaseGenerator
	engine *core.TemplateEngine
}

// New creates a new build script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
		engine:        engine,
	}, nil
}

// Generate generates build.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Use preset for build script
	composer := ctx.GetVariablePreset().ForBuildScript()

	// Use plugin service to process plugins
	pluginService := services.NewPluginService(ctx, g.engine)
	if pluginService.HasPlugins() {
		plugins := pluginService.PrepareForBuildScript()
		composer.Override("PLUGINS", plugins)
	}

	// Add build script specific variables
	composer.
		WithCustom("GENERATE_SCRIPTS", ctx.Config.Runtime.GenerateScripts)

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/build.sh.tmpl
var template string
