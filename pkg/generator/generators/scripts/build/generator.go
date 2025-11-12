package build

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
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

// PluginInfo holds plugin information for template
type PluginInfo struct {
	Name           string
	DownloadURL    string
	InstallDir     string
	InstallCommand string
	RuntimeEnv     []config.EnvironmentVariable
}

// Generate generates build.sh content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	ctx := g.GetContext()

	// Use preset for build script
	composer := ctx.GetVariablePreset().ForBuildScript()

	// Prepare plugin install commands with processed environment variables
	var plugins []PluginInfo
	sharedInstallDir := ctx.Config.Plugins.InstallDir

	for _, plugin := range ctx.Config.Plugins.Items {
		// Process runtime environment variables - replace ${PLUGIN_INSTALL_DIR} with actual install dir
		processedEnv := make([]config.EnvironmentVariable, len(plugin.RuntimeEnv))
		for i, env := range plugin.RuntimeEnv {
			processedEnv[i] = config.EnvironmentVariable{
				Name: env.Name,
				Value: g.engine.ReplaceVariables(env.Value, map[string]string{
					"PLUGIN_INSTALL_DIR": sharedInstallDir,
				}),
			}
		}

		plugins = append(plugins, PluginInfo{
			Name:           plugin.Name,
			DownloadURL:    plugin.DownloadURL,
			InstallDir:     sharedInstallDir,
			InstallCommand: plugin.InstallCommand,
			RuntimeEnv:     processedEnv,
		})
	}

	// Add build script specific variables
	composer.
		WithCustom("GENERATE_SCRIPTS", ctx.Config.Runtime.GenerateScripts)

	// Override plugins with processed version
	if len(plugins) > 0 {
		composer.Override("PLUGINS", plugins)
	}

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/build.sh.tmpl
var template string
