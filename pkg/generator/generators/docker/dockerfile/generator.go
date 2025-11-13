package dockerfile

import (
	_ "embed"
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/domain/services"
)

const GeneratorType = "dockerfile"

// init registers the dockerfile generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates Dockerfiles
type Generator struct {
	core.BaseGenerator
	arch string
}

// New creates a new dockerfile generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("dockerfile generator requires architecture parameter (amd64 or arm64)")
	}

	arch, ok := options[0].(string)
	if !ok {
		return nil, fmt.Errorf("dockerfile generator architecture parameter must be a string")
	}

	if arch != "amd64" && arch != "arm64" {
		return nil, fmt.Errorf("dockerfile generator architecture must be 'amd64' or 'arm64', got: %s", arch)
	}

	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType+"-"+arch, ctx, engine),
		arch:          arch,
	}, nil
}

// Generate generates Dockerfile content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	vars := g.prepareTemplateVars()
	return g.RenderTemplate(template, vars)
}

// prepareTemplateVars prepares variables for Dockerfile template
func (g *Generator) prepareTemplateVars() map[string]interface{} {
	ctx := g.GetContext()

	// Use preset for Dockerfile with architecture
	composer := ctx.GetVariablePreset().ForDockerfile(g.arch)

	// Get builder image for package manager detection
	builderImage := ""
	if val, ok := composer.Get("BUILDER_IMAGE"); ok {
		builderImage = val.(string)
	}

	// Use language service for language-specific logic
	langService := services.NewLanguageService()

	// Add Dockerfile-specific custom variables
	composer.
		WithCustom("PKG_MANAGER", detectPackageManager(builderImage)).
		WithCustom("DEPENDENCY_FILES", getDependencyFilesList(ctx.Config, ctx.OutputDir)).
		WithCustom("DEPS_INSTALL_COMMAND", langService.GetDepsInstallCommand(ctx.Config.Language.Type))

	// Use plugin service to process plugins
	pluginService := services.NewPluginService(ctx, g.GetEngine())
	hasPlugins := pluginService.HasPlugins()

	composer.WithCustom("HAS_PLUGINS", hasPlugins)

	if hasPlugins {
		plugins := pluginService.PrepareForDockerfile()
		composer.Override("PLUGINS", plugins)

		// Add plugin build script variables
		composer.
			WithCustom("PLUGIN_BUILD_SCRIPT", "build_plugins.sh").
			WithCustom("PLUGIN_BUILD_SCRIPT_CONTAINER_PATH", fmt.Sprintf("%s/build_plugins.sh", ctx.Paths.CI.ContainerScriptDir))
	}

	return composer.Build()
}

//go:embed templates/dockerfile_.tmpl
var template string
