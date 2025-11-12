package dockerfile

import (
	_ "embed"
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
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

	// Add Dockerfile-specific custom variables
	composer.
		WithCustom("PKG_MANAGER", detectPackageManager(builderImage)).
		WithCustom("DEPENDENCY_FILES", getDependencyFilesList(ctx.Config)).
		WithCustom("DEPS_INSTALL_COMMAND", getDepsInstallCommand(ctx.Config.Language.Type))

	// Process plugins with install commands
	var plugins []map[string]interface{}
	sharedInstallDir := ctx.Config.Plugins.InstallDir
	for _, plugin := range ctx.Config.Plugins.Items {
		pluginVars := ctx.Variables.WithPlugin(plugin, sharedInstallDir)
		plugins = append(plugins, map[string]interface{}{
			"InstallCommand": core.SubstituteVariables(plugin.InstallCommand, pluginVars.ToMap()),
			"Name":           plugin.Name,
			"InstallDir":     sharedInstallDir,
			"RuntimeEnv":     plugin.RuntimeEnv,
		})
	}

	// Override plugins with processed version
	if len(plugins) > 0 {
		composer.Override("PLUGINS", plugins)
	}

	return composer.Build()
}

//go:embed templates/dockerfile_.tmpl
var template string
