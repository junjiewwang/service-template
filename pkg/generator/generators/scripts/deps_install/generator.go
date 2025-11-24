package deps_install

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/domain/services"
	"github.com/junjiewwang/service-template/pkg/generator/domain/services/languageservice"
)

const GeneratorType = "deps-install-script"

// init registers the deps install script generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates build_deps_install.sh script
type Generator struct {
	core.BaseGenerator
	depSvc  *services.DependencyService
	langSvc *languageservice.LanguageService
}

// New creates a new deps install script generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
		depSvc:        services.NewDependencyService(ctx, engine),
		langSvc:       languageservice.NewLanguageService(ctx),
	}, nil
}

// Generate generates build_deps_install.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Get build dependencies
	buildDeps := g.depSvc.GetBuildDependencies()

	// Get language-specific install command (with custom config and variable substitution)
	depsInstallCmd := g.langSvc.GetDepsInstallCommand(ctx.Config.Language.Type)

	// Get language-specific config using helper methods
	var goProxy, goSumDB string
	if ctx.Config.Language.Type == "go" {
		goProxy = ctx.Config.Language.GetString("goproxy", "")
		goSumDB = ctx.Config.Language.GetString("gosumdb", "")
	}

	// Prepare template data
	data := map[string]interface{}{
		"HasSystemPackages":  g.depSvc.HasSystemPackages(),
		"SystemPackages":     buildDeps.SystemPkgs,
		"HasCustomPackages":  g.depSvc.HasCustomPackages(),
		"CustomPackages":     buildDeps.CustomPkgs,
		"Language":           ctx.Config.Language.Type,
		"DepsInstallCommand": depsInstallCmd,
		"GoProxy":            goProxy,
		"GoSumDB":            goSumDB,
	}

	return g.RenderTemplate(template, data)
}

//go:embed templates/deps_install.sh.tmpl
var template string
