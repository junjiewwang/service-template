package generator

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeBuildScript = "build-script"

// BuildScriptTemplateGenerator generates build.sh script
type BuildScriptTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the BuildScript generator
func init() {
	RegisterGenerator(GeneratorTypeBuildScript, createBuildScriptGenerator)
}

// createBuildScriptGenerator is the creator function for BuildScript generator
func createBuildScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewBuildScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewBuildScriptTemplateGenerator creates a new build script generator
func NewBuildScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *BuildScriptTemplateGenerator {
	return &BuildScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeBuildScript,
		},
	}
}

// Generate generates build.sh content
func (g *BuildScriptTemplateGenerator) Generate() (string, error) {
	// Prepare plugin install commands
	type PluginInfo struct {
		Name           string
		DownloadURL    string
		InstallDir     string
		InstallCommand string
		RuntimeEnv     []config.EnvironmentVariable // 新增：运行时环境变量
	}
	var plugins []PluginInfo
	for _, plugin := range g.config.Plugins {
		plugins = append(plugins, PluginInfo{
			Name:           plugin.Name,
			DownloadURL:    plugin.DownloadURL,
			InstallDir:     plugin.InstallDir,
			InstallCommand: plugin.InstallCommand, // Keep original command with variables
			RuntimeEnv:     plugin.RuntimeEnv,     // 新增：运行时环境变量
		})
	}

	vars := map[string]interface{}{
		"SERVICE_NAME":       g.config.Service.Name,
		"DEPLOY_DIR":         g.config.Service.DeployDir,
		"BUILD_COMMAND":      g.config.Build.Commands.Build,
		"PRE_BUILD_COMMAND":  g.config.Build.Commands.PreBuild,
		"POST_BUILD_COMMAND": g.config.Build.Commands.PostBuild,
		"PLUGINS":            plugins,
		"SERVICE_ROOT":       g.config.Service.DeployDir + "/" + g.config.Service.Name,
		"PLUGIN_ROOT_DIR":    "/plugins",                       // 新增：插件根目录
		"GENERATE_SCRIPTS":   g.config.Runtime.GenerateScripts, // 新增：是否生成运行时脚本
	}

	return g.RenderTemplate(buildScriptTemplate, vars)
}

// buildScriptTemplate is the build.sh template
const buildScriptTemplate = `#!/bin/bash

script_path=$(
	cd $(dirname $0)
	pwd
)

PKG_NAME={{ .SERVICE_NAME }}

# 到项目根目录并获取绝对路径
cd $script_path/../../

# ============================================
# Build Output Directory Configuration
# ============================================
#
# In Docker build: BUILD_OUTPUT_DIR=/opt/dist, PROJECT_ROOT=/opt (set by Dockerfile)
# In local build: BUILD_OUTPUT_DIR=${PROJECT_ROOT}/dist (calculated here)
#
# DO NOT hardcode "dist" anywhere in your scripts - always use ${BUILD_OUTPUT_DIR}

# If PROJECT_ROOT not set (local build), calculate it
if [ -z "${PROJECT_ROOT}" ]; then
	PROJECT_ROOT=$(pwd)
fi

# If BUILD_OUTPUT_DIR not set (local build), use PROJECT_ROOT/dist
BUILD_OUTPUT_DIR="${BUILD_OUTPUT_DIR:-${PROJECT_ROOT}/dist}"

export BUILD_OUTPUT_DIR
export PROJECT_ROOT

echo "Project root: ${PROJECT_ROOT}"
echo "Build output directory: ${BUILD_OUTPUT_DIR}"


# ============================================
# Standardized Multi-language Build Support
# ============================================
echo "========================================="
echo "TCS Service Build System"
echo "Service: ${PKG_NAME}"
echo "========================================="

# Create output directories
mkdir -p ${BUILD_OUTPUT_DIR}/bin ${BUILD_OUTPUT_DIR}/conf

# ============================================
# Execute pre-build commands
# ============================================
{{- if .PRE_BUILD_COMMAND }}
echo "Executing pre-build commands..."
{{ .PRE_BUILD_COMMAND }}
[ $? -ne 0 ] && echo "ERROR: Pre-build commands failed" && exit 1
echo "Pre-build completed successfully"
{{- end }}

# ============================================
# Execute build commands
# ============================================
echo "Executing build commands..."
{{ .BUILD_COMMAND }}
[ $? -ne 0 ] && echo "ERROR: Build commands failed" && exit 1
echo "Build completed successfully"

# ============================================
# Execute post-build commands
# ============================================
{{- if .POST_BUILD_COMMAND }}
echo "Executing post-build commands..."
{{ .POST_BUILD_COMMAND }}
[ $? -ne 0 ] && echo "ERROR: Post-build commands failed" && exit 1
echo "Post-build completed successfully"
{{- end }}

# ============================================
# Deploy artifacts
# ============================================
target_path={{ .DEPLOY_DIR }}/${PKG_NAME}
# create service dir
mkdir -p ${target_path}

# Copy build artifacts from BUILD_OUTPUT_DIR directory
if [ ! -d "${BUILD_OUTPUT_DIR}" ]; then
	echo "ERROR: ${BUILD_OUTPUT_DIR}/ directory not found!"
	echo ""
	echo "The build commands must output artifacts to the ${BUILD_OUTPUT_DIR}/ directory."
	echo "Please ensure build commands create and populate ${BUILD_OUTPUT_DIR}/"
	exit 1
fi

echo "Copying artifacts from ${BUILD_OUTPUT_DIR}/ to ${target_path}/"
cp -rf ${BUILD_OUTPUT_DIR}/* ${target_path}/

[ $? -ne 0 ] && exit

# ============================================
# Setup TCE environment and install plugins
# ============================================
{{- if .PLUGINS }}
# 创建插件根目录
mkdir -p {{ .PLUGIN_ROOT_DIR }}

{{- range .PLUGINS }}
# ============================================
# Install plugin: {{ .Name }}
# ============================================
PLUGIN_NAME="{{ .Name }}"
PLUGIN_DOWNLOAD_URL="{{ .DownloadURL }}"
PLUGIN_INSTALL_DIR="{{ .InstallDir }}"
PLUGIN_WORK_DIR="{{ $.PLUGIN_ROOT_DIR }}/{{ .Name }}"

echo "Installing plugin: ${PLUGIN_NAME}"
echo "Download URL: ${PLUGIN_DOWNLOAD_URL}"
echo "Install directory: ${PLUGIN_INSTALL_DIR}"
echo "Plugin work directory: ${PLUGIN_WORK_DIR}"

# 创建插件工作目录
mkdir -p ${PLUGIN_WORK_DIR}

{{- if .InstallCommand }}
# Use custom install command from configuration
echo "Using custom install command..."
{{ .InstallCommand }}
{{- else }}
# Default install command: download and execute script
echo "Using default install command..."
set -o pipefail # Enable pipefail to catch errors in pipe
if ! curl -fsSL ${PLUGIN_DOWNLOAD_URL} | bash -es ${PLUGIN_WORK_DIR}; then
	echo "ERROR: Failed to download or execute plugin script"
	echo "URL: ${PLUGIN_DOWNLOAD_URL}"
	echo "Target directory: ${PLUGIN_WORK_DIR}"
	exit 1
fi
{{- end }}

if [ $? -eq 0 ]; then
	echo "Plugin {{ .Name }} downloaded to work directory successfully"
else
	echo "ERROR: Plugin {{ .Name }} download failed"
	exit 1
fi

{{- if .RuntimeEnv }}
# ============================================
# Setup runtime environment variables for {{ .Name }}
# ============================================
echo "Setting up runtime environment variables for {{ .Name }}..."

# 创建插件环境变量文件
ENV_FILE="${PLUGIN_WORK_DIR}/.env"
cat > ${ENV_FILE} << 'EOF'
{{- range .RuntimeEnv }}
export {{ .Name }}="{{ .Value }}"
{{- end }}
EOF

echo "Environment variables written to ${ENV_FILE}"
{{- end }}

echo "Plugin {{ .Name }} setup completed"
echo ""
{{- end }}
{{- end }}

# Set default TCE_DIR if no plugins defined
TCE_DIR=${TCE_DIR:-/tce}
mkdir -p ${TCE_DIR}

# ============================================
# Generate runtime scripts
# ============================================
echo "Generating runtime scripts..."

# 复制运行时脚本到服务目录
SERVICE_DIR="{{ .SERVICE_ROOT }}"
mkdir -p ${SERVICE_DIR}

if [ -f "${script_path}/entrypoint.sh" ]; then
	cp -f ${script_path}/entrypoint.sh ${SERVICE_DIR}/entrypoint.sh
	chmod +x ${SERVICE_DIR}/entrypoint.sh
	echo "✓ Generated ${SERVICE_DIR}/entrypoint.sh"
fi

if [ -f "${script_path}/healthchk.sh" ]; then
	cp -f ${script_path}/healthchk.sh ${SERVICE_DIR}/healthcheck.sh
	chmod +x ${SERVICE_DIR}/healthcheck.sh
	echo "✓ Generated ${SERVICE_DIR}/healthcheck.sh"
fi


echo "Done!"
`
