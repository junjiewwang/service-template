package generator

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeRtPrepareScript = "rt-prepare-script"

// init registers the rt prepare script generator
func init() {
	RegisterGenerator(GeneratorTypeRtPrepareScript, createRtPrepareScriptGenerator)
}

// RtPrepareScriptTemplateGenerator generates rt_prepare.sh script
type RtPrepareScriptTemplateGenerator struct {
	BaseTemplateGenerator
}

// createRtPrepareScriptGenerator is the creator function for RtPrepareScript generator
func createRtPrepareScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewRtPrepareScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewRtPrepareScriptTemplateGenerator creates a new rt prepare script generator
func NewRtPrepareScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *RtPrepareScriptTemplateGenerator {
	return &RtPrepareScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeRtPrepareScript,
		},
	}
}

// Generate generates rt_prepare.sh content
func (g *RtPrepareScriptTemplateGenerator) Generate() (string, error) {
	vars := map[string]interface{}{
		"RUNTIME_DEPS_PACKAGES": g.config.Runtime.SystemDependencies.Runtime.Packages,
	}
	return g.RenderTemplate(rtPrepareScriptTemplate, vars)
}

// Template content
const rtPrepareScriptTemplate = `#!/bin/sh
# rt_prepare.sh - Runtime preparation script
# This script is used to install runtime dependencies and prepare the environment
# for the application to run properly.

set -e

# ============================================
# Runtime Dependencies Installation
# ============================================
echo "========================================="
echo "TCS Runtime Preparation"
echo "========================================="

# Detect architecture
ARCH=$(uname -m)
echo "Detected architecture: ${ARCH}"
echo ""

# ============================================
# Check and Install Required Tools
# ============================================
# Define required tools here (space-separated list)
# Example: REQUIRED_TOOLS="curl tar wget"
REQUIRED_TOOLS="{{ join " " .RUNTIME_DEPS_PACKAGES }}"

if [ -n "$REQUIRED_TOOLS" ]; then
	echo "Checking required tools..."

	for tool in $REQUIRED_TOOLS; do
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
else
	echo "No required tools specified, skipping tool installation"
	echo ""
fi

# ============================================
# Architecture-specific Dependencies
# ============================================
case "${ARCH}" in
x86_64 | amd64)
	echo "Installing x86_64/amd64 specific dependencies..."
	# Add x86_64 specific packages here
	;;
aarch64 | arm64)
	echo "Installing aarch64/arm64 specific dependencies..."
	# Add ARM64 specific packages here
	;;
*)
	echo "WARNING: Unknown architecture ${ARCH}, skipping architecture-specific dependencies"
	;;
esac

# ============================================
# Verify Essential Dependencies
# ============================================
echo "Verifying essential dependencies..."

# Example: Check if required commands are available
# command -v python3 >/dev/null 2>&1 || { echo "ERROR: python3 not found"; exit 1; }


echo "Runtime preparation completed successfully"
`
