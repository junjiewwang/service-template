package deps_install

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "deps-install-script"

// init registers the deps install script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates build_deps_install.sh script
type Generator struct {
	core.BaseGenerator
}

// New creates a new deps install script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates build_deps_install.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Use preset for script
	composer := ctx.GetVariablePreset().ForScript()

	// Get language-specific config
	var goProxy, goSumDB string
	if ctx.Config.Language.Type == "go" {
		if goproxy, ok := ctx.Config.Language.Config["goproxy"]; ok {
			goProxy = goproxy
		}
		if gosumdb, ok := ctx.Config.Language.Config["gosumdb"]; ok {
			goSumDB = gosumdb
		}
	}

	// Add script-specific custom variables
	composer.
		WithCustom("BUILD_DEPS_PACKAGES", ctx.Config.Build.SystemDependencies.Packages).
		WithCustom("GO_PROXY", goProxy).
		WithCustom("GO_SUMDB", goSumDB)

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/deps_install.sh.tmpl
var template string
