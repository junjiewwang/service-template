package generator

import (
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

	return g.RenderTemplate(entrypointScriptTemplate, vars)
}

// Template content
const entrypointScriptTemplate = `#!/bin/sh

echo "========================================="
echo "TCS Service Entrypoint"
echo "Service: {{ .SERVICE_NAME }}"
echo "Deploy Dir: {{ .DEPLOY_DIR }}"
echo "========================================="

# Export service paths as environment variables
export SERVICE_ROOT="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}"
export SERVICE_BIN_DIR="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}/bin"
export SERVICE_NAME="{{ .SERVICE_NAME }}"

echo "Service Root: ${SERVICE_ROOT}"
echo "Service Bin Dir: ${SERVICE_BIN_DIR}"
echo "Service Name: ${SERVICE_NAME}"

{{- if .ENV_VARS }}
# Set environment variables
{{- range .ENV_VARS }}
export {{ .Name }}={{ .Value }}
{{- end }}
{{- end }}

{{- if .HAS_PLUGINS_ENV }}
# ============================================
# Load plugin environment variables
# ============================================
{{- range .PLUGINS_ENV }}
# Load environment variables for plugin: {{ .Name }}
PLUGIN_ENV_FILE="{{ .InstallDir }}/.env"
if [ -f "${PLUGIN_ENV_FILE}" ]; then
    echo "Loading environment variables for {{ .Name }} from ${PLUGIN_ENV_FILE}"
    set -a  # Automatically export all variables
    . "${PLUGIN_ENV_FILE}"
    set +a
    
    # 显示已加载的环境变量
    echo "Environment variables loaded for {{ .Name }}:"
    {{- range .RuntimeEnv }}
    echo "  {{ .Name }}=${{ .Name }}"
    {{- end }}
else
    echo "Warning: Environment file not found for {{ .Name }}: ${PLUGIN_ENV_FILE}"
fi
echo ""
{{- end }}
{{- end }}

# ============================================
# Language-specific start command (CUSTOMIZE THIS SECTION)
# ============================================
# For Go:
{{ .STARTUP_COMMAND }}

# ============================================
`
