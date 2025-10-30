package generator

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// ScriptsGenerator generates build and deployment scripts
type ScriptsGenerator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
}

// NewScriptsGenerator creates a new scripts generator
func NewScriptsGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ScriptsGenerator {
	return &ScriptsGenerator{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
	}
}

// GenerateBuildScript generates build.sh
func (g *ScriptsGenerator) GenerateBuildScript() (string, error) {
	vars := map[string]interface{}{
		"PRE_BUILD":  g.config.Build.Commands.PreBuild,
		"BUILD":      g.config.Build.Commands.Build,
		"POST_BUILD": g.config.Build.Commands.PostBuild,
	}

	template := `#!/bin/bash
# Auto-generated build script
set -e

echo "Starting build process..."
{{- if .PRE_BUILD }}

# Pre-build commands
{{ .PRE_BUILD }}
{{- end }}

# Build commands
{{ .BUILD }}
{{- if .POST_BUILD }}

# Post-build commands
{{ .POST_BUILD }}
{{- end }}

echo "✓ Build completed successfully"
`

	return g.templateEngine.Render(template, vars)
}

// GenerateDepsInstallScript generates deps_install.sh
func (g *ScriptsGenerator) GenerateDepsInstallScript() (string, error) {
	vars := map[string]interface{}{
		"LANGUAGE": g.config.Language.Type,
	}

	template := `#!/bin/bash
# Auto-generated dependency installation script
set -e

{{- if eq .LANGUAGE "go" }}
echo "Installing Go dependencies..."
go mod download
go mod verify
echo "✓ Go dependencies installed"
{{- else if eq .LANGUAGE "python" }}
echo "Installing Python dependencies..."
pip install -r requirements.txt
echo "✓ Python dependencies installed"
{{- else if eq .LANGUAGE "nodejs" }}
echo "Installing Node.js dependencies..."
npm install
echo "✓ Node.js dependencies installed"
{{- else if eq .LANGUAGE "java" }}
echo "Installing Java dependencies..."
mvn dependency:go-offline
echo "✓ Java dependencies installed"
{{- else }}
echo "No dependency installation needed"
{{- end }}
`

	return g.templateEngine.Render(template, vars)
}

// GenerateRtPrepareScript generates rt_prepare.sh
func (g *ScriptsGenerator) GenerateRtPrepareScript() (string, error) {
	// Prepare plugin install commands
	type PluginInfo struct {
		InstallCommand string
	}
	var plugins []PluginInfo
	for _, plugin := range g.config.Plugins {
		pluginVars := g.variables.WithPlugin(plugin)
		installCmd := SubstituteVariables(plugin.InstallCommand, pluginVars.ToMap())
		plugins = append(plugins, PluginInfo{InstallCommand: installCmd})
	}

	vars := map[string]interface{}{
		"PLUGINS": plugins,
	}

	template := `#!/bin/bash
# Auto-generated runtime preparation script
set -e

echo "Preparing runtime environment..."
{{- if .PLUGINS }}

# Install plugins
{{- range .PLUGINS }}
{{ .InstallCommand }}
{{- end }}
{{- end }}

echo "✓ Runtime environment prepared"
`

	return g.templateEngine.Render(template, vars)
}
