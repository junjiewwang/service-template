package entrypoint

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
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

	// 准备插件环境变量信息
	var pluginEnvs []map[string]interface{}
	sharedInstallDir := ctx.Config.Plugins.InstallDir
	for _, plugin := range ctx.Config.Plugins.Items {
		if len(plugin.RuntimeEnv) > 0 {
			pluginEnvs = append(pluginEnvs, map[string]interface{}{
				"Name":       plugin.Name,
				"InstallDir": sharedInstallDir, // 使用共享的安装目录
				"RuntimeEnv": plugin.RuntimeEnv,
			})
		}
	}

	// Add script-specific custom variables
	composer.
		WithCustom("PLUGINS_ENV", pluginEnvs).
		WithCustom("HAS_PLUGINS_ENV", len(pluginEnvs) > 0)

	return g.RenderTemplate(template, composer.Build())
}

//go:embed templates/entrypoint.sh.tmpl
var template string
