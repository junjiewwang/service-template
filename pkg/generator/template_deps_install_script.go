package generator

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeDepsInstallScript = "deps-install-script"

// DepsInstallScriptTemplateGenerator generates build_deps_install.sh script
type DepsInstallScriptTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the DepsInstallScript generator
func init() {
	RegisterGenerator(GeneratorTypeDepsInstallScript, createDepsInstallScriptGenerator)
}

// createDepsInstallScriptGenerator is the creator function for DepsInstallScript generator
func createDepsInstallScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewDepsInstallScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewDepsInstallScriptTemplateGenerator creates a new deps install script generator
func NewDepsInstallScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *DepsInstallScriptTemplateGenerator {
	return &DepsInstallScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeDepsInstallScript,
		},
	}
}

// Generate generates build_deps_install.sh content
func (g *DepsInstallScriptTemplateGenerator) Generate() (string, error) {
	// Get language-specific config
	var goProxy, goSumDB string
	if g.config.Language.Type == "go" {
		if goproxy, ok := g.config.Language.Config["goproxy"]; ok {
			goProxy = goproxy
		}
		if gosumdb, ok := g.config.Language.Config["gosumdb"]; ok {
			goSumDB = gosumdb
		}
	}

	vars := map[string]interface{}{
		"LANGUAGE":            g.config.Language.Type,
		"BUILD_DEPS_PACKAGES": g.config.Build.SystemDependencies.Build.Packages,
		"GO_PROXY":            goProxy,
		"GO_SUMDB":            goSumDB,
	}

	return g.RenderTemplate(g.getTemplate(), vars)
}

// depsInstallScriptTemplate is the build_deps_install.sh template
func (g *DepsInstallScriptTemplateGenerator) getTemplate() string {
	return `#!/bin/bash
# deps_install.sh - Dependency installation script
# This script is used to install build-time dependencies in a separate Docker layer
# to leverage Docker's layer caching mechanism for faster builds.
#
# IMPORTANT: This script is executed BEFORE build.sh in the Dockerfile
#
# Environment Variables (automatically set by Dockerfile):
#   - BUILD_OUTPUT_DIR: Absolute path to build output directory (e.g., /opt/dist)
#   - PROJECT_ROOT: Absolute path to project root (e.g., /opt)
#
# Usage Scenarios:
#   1. Install language-specific dependencies (npm install, pip install, go mod download, etc.)
#   2. Download external tools or binaries
#   3. Set up build environment
#
# DO NOT:
#   - Perform actual compilation/build (that's build.sh's job)
#   - Copy/move source files (Dockerfile handles this)
#
# ============================================
# Script Setup
# ============================================

set -e # Exit on error

cd "${PROJECT_ROOT}"

echo "========================================="
echo "TCS Dependency Installation"
echo "Project root: ${PROJECT_ROOT}"
echo "Build output: ${BUILD_OUTPUT_DIR}"
echo "========================================="

# ============================================
# Check and Install Required Tools
# ============================================
echo "Checking required tools..."

# Check and install required tools
{{- if .BUILD_DEPS_PACKAGES }}
for tool in {{ join " " .BUILD_DEPS_PACKAGES }}; do
{{- else }}
for tool in ; do
{{- end }}
	if ! command -v "$tool" >/dev/null 2>&1; then
		echo "Installing $tool..."
		if command -v apt-get >/dev/null 2>&1; then
			apt-get update -qq && apt-get install -y "$tool"
		elif command -v yum >/dev/null 2>&1; then
			yum install -y "$tool"
		elif command -v apk >/dev/null 2>&1; then
			apk add --no-cache "$tool"
		elif command -v dnf >/dev/null 2>&1; then
			dnf install -y "$tool"
		elif command -v zypper >/dev/null 2>&1; then
			zypper install -y "$tool"
		else
			echo "ERROR: No package manager found. Please install $tool manually"
			exit 1
		fi
		command -v "$tool" >/dev/null 2>&1 || {
			echo "ERROR: Failed to install $tool"
			exit 1
		}
		echo "✓ $tool installed successfully"
	else
		echo "✓ $tool is already installed"
	fi
done

echo "All required tools are available"
echo ""

# ============================================
# Ensure Build Output Directory Exists
# ============================================
echo "Ensuring build output directory exists..."
mkdir -p "${BUILD_OUTPUT_DIR}"
echo "✓ Build output directory ready: ${BUILD_OUTPUT_DIR}"
echo ""

# ============================================
# Detect Project Type and Install Dependencies
# ============================================

# Example 1: Go Project
{{- if eq .LANGUAGE "go" }}
if [ -f "go.mod" ]; then
{{- if .GO_PROXY }}
	go env -w GOPROXY="{{ .GO_PROXY }}"
{{- end }}
{{- if .GO_SUMDB }}
	#Set Go checksum database
	go env -w GOSUMDB="{{ .GO_SUMDB }}"
{{- end }}
	echo "Detected Go project, downloading dependencies..."
	go mod download
fi
{{- end }}

# Example 2: Node.js Project
{{- if eq .LANGUAGE "nodejs" }}
# if [ -f "package.json" ]; then
#     echo "Detected Node.js project, installing dependencies..."
#     npm install --production
#     # Or use yarn
#     # yarn install --production
# fi
{{- end }}

# Example 3: Python Project
{{- if eq .LANGUAGE "python" }}
# if [ -f "requirements.txt" ]; then
#     echo "Detected Python project, installing dependencies..."
#     # Install to system site-packages
#     pip install -r requirements.txt
#
#     # Or install to custom directory (useful for packaging)
#     # mkdir -p "${BUILD_OUTPUT_DIR}/bin"
#     # pip install -r requirements.txt -t "${BUILD_OUTPUT_DIR}/bin/"
# fi
{{- end }}

# Example 4: Java Maven Project
{{- if eq .LANGUAGE "java" }}
# if [ -f "pom.xml" ]; then
#     echo "Detected Maven project, downloading dependencies..."
#     mvn dependency:go-offline
# fi
{{- end }}


# ============================================
# Custom Dependency Installation
# ============================================
# Add your custom dependency installation logic here
# Examples:
#   - Download external binaries
#   - Install system packages (if needed in builder image)
#   - Set up build tools

echo "Dependency installation completed successfully"
`
}
