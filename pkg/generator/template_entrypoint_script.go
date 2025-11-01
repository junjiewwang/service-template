package generator

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeEntrypointScript = "entrypoint-script"

// init registers the entrypoint script generator
func init() {
	RegisterGenerator(GeneratorTypeEntrypointScript, createEntrypointScriptGenerator)
}

// EntrypointScriptTemplateGenerator generates entrypoint.sh script
type EntrypointScriptTemplateGenerator struct {
	BaseTemplateGenerator
}

// createEntrypointScriptGenerator is the creator function for EntrypointScript generator
func createEntrypointScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewEntrypointScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewEntrypointScriptTemplateGenerator creates a new entrypoint script generator
func NewEntrypointScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *EntrypointScriptTemplateGenerator {
	return &EntrypointScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeEntrypointScript,
		},
	}
}

//go:embed templates/entrypoint.sh.tmpl
var entrypointScriptTemplate string

// Generate generates entrypoint.sh content
func (g *EntrypointScriptTemplateGenerator) Generate() (string, error) {
	// 准备插件环境变量信息
	var pluginEnvs []map[string]interface{}
	for _, plugin := range g.config.Plugins {
		if len(plugin.RuntimeEnv) > 0 {
			pluginEnvs = append(pluginEnvs, map[string]interface{}{
				"Name":       plugin.Name,
				"InstallDir": plugin.InstallDir,
				"RuntimeEnv": plugin.RuntimeEnv,
			})
		}
	}

	vars := map[string]interface{}{
		"SERVICE_NAME":    g.config.Service.Name,
		"DEPLOY_DIR":      g.config.Service.DeployDir,
		"STARTUP_COMMAND": g.config.Runtime.Startup.Command,
		"ENV_VARS":        g.config.Runtime.Startup.Env,
		"PLUGINS_ENV":     pluginEnvs,          // 新增：插件环境变量
		"HAS_PLUGINS_ENV": len(pluginEnvs) > 0, // 新增：是否有插件环境变量
	}

	return g.RenderTemplate(g.getTemplate(), vars)
}

// getTemplate returns the entrypoint script template
func (g *EntrypointScriptTemplateGenerator) getTemplate() string {
	return entrypointScriptTemplate
}
