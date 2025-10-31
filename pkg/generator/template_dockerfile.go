package generator

import (
	"fmt"
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeDockerfile = "dockerfile"

// DockerfileTemplateGenerator generates Dockerfiles using factory pattern
type DockerfileTemplateGenerator struct {
	BaseTemplateGenerator
	arch string
}

// init registers the Dockerfile generator
func init() {
	RegisterGenerator(GeneratorTypeDockerfile, createDockerfileGenerator)
}

// createDockerfileGenerator is the creator function for Dockerfile generator
func createDockerfileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
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

	return NewDockerfileTemplateGenerator(cfg, engine, vars, arch), nil
}

// NewDockerfileTemplateGenerator creates a new Dockerfile template generator
func NewDockerfileTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, arch string) *DockerfileTemplateGenerator {
	return &DockerfileTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeDockerfile + "-" + arch,
		},
		arch: arch,
	}
}

// Generate generates Dockerfile content
func (g *DockerfileTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for Dockerfile template
func (g *DockerfileTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["ARCH"] = g.arch
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name
	vars["LANGUAGE"] = g.config.Language.Type
	vars["DEPLOY_DIR"] = g.config.Service.DeployDir

	// Images
	if g.arch == "amd64" {
		vars["BUILDER_IMAGE"] = g.config.Build.BuilderImage.AMD64
		vars["RUNTIME_IMAGE"] = g.config.Build.RuntimeImage.AMD64
	} else {
		vars["BUILDER_IMAGE"] = g.config.Build.BuilderImage.ARM64
		vars["RUNTIME_IMAGE"] = g.config.Build.RuntimeImage.ARM64
	}

	// Package manager
	vars["PKG_MANAGER"] = detectPackageManager(vars["BUILDER_IMAGE"].(string))

	// Dependencies
	vars["BUILD_DEPS_PACKAGES"] = g.config.Build.SystemDependencies.Build.Packages
	vars["RUNTIME_DEPS_PACKAGES"] = g.config.Runtime.SystemDependencies.Runtime.Packages

	// Dependency files
	vars["DEPENDENCY_FILES"] = g.getDependencyFilesList()
	vars["DEPS_INSTALL_COMMAND"] = g.getDepsInstallCommand()

	// Build commands
	vars["PRE_BUILD_COMMAND"] = g.config.Build.Commands.PreBuild
	vars["BUILD_COMMAND"] = g.config.Build.Commands.Build
	vars["POST_BUILD_COMMAND"] = g.config.Build.Commands.PostBuild

	// Plugins
	var plugins []map[string]interface{}
	for _, plugin := range g.config.Plugins {
		pluginVars := g.variables.WithPlugin(plugin)
		plugins = append(plugins, map[string]interface{}{
			"InstallCommand": SubstituteVariables(plugin.InstallCommand, pluginVars.ToMap()),
			"Name":           plugin.Name,
			"InstallDir":     plugin.InstallDir,
			"RuntimeEnv":     plugin.RuntimeEnv,
		})
	}
	vars["PLUGINS"] = plugins
	vars["HAS_PLUGINS"] = len(g.config.Plugins) > 0

	// Ports
	var exposePorts []int
	for _, port := range g.config.Service.Ports {
		if port.Expose {
			exposePorts = append(exposePorts, port.Port)
		}
	}
	vars["EXPOSE_PORTS"] = exposePorts

	// Health check
	vars["HEALTHCHECK_ENABLED"] = g.config.Runtime.Healthcheck.Enabled

	return vars
}

// getDependencyFilesList returns list of dependency files
func (g *DockerfileTemplateGenerator) getDependencyFilesList() []string {
	if g.config.Build.DependencyFiles.AutoDetect {
		return getDefaultDependencyFiles(g.config.Language.Type)
	}
	return g.config.Build.DependencyFiles.Files
}

// getDepsInstallCommand generates dependency installation command
func (g *DockerfileTemplateGenerator) getDepsInstallCommand() string {
	switch g.config.Language.Type {
	case "go":
		return "go mod download"
	case "python":
		return "pip install -r requirements.txt"
	case "nodejs":
		return "npm install"
	case "java":
		return "mvn dependency:go-offline"
	default:
		return "echo 'No dependency installation needed'"
	}
}

// getDefaultDependencyFiles returns default dependency files for a language
func getDefaultDependencyFiles(language string) []string {
	switch language {
	case "go":
		return []string{"go.mod", "go.sum"}
	case "python":
		return []string{"requirements.txt"}
	case "nodejs":
		return []string{"package.json", "package-lock.json"}
	case "java":
		return []string{"pom.xml"}
	default:
		return []string{}
	}
}

// detectPackageManager detects the package manager from the image name
func detectPackageManager(image string) string {
	imageLower := strings.ToLower(image)

	if strings.Contains(imageLower, "alpine") {
		return "apk"
	} else if strings.Contains(imageLower, "debian") || strings.Contains(imageLower, "ubuntu") {
		return "apt-get"
	} else if strings.Contains(imageLower, "centos") || strings.Contains(imageLower, "rhel") || strings.Contains(imageLower, "tencentos") {
		return "yum"
	} else if strings.Contains(imageLower, "fedora") {
		return "dnf"
	}

	return "yum"
}

// getTemplate returns the Dockerfile template
func (g *DockerfileTemplateGenerator) getTemplate() string {
	return `# Define build arguments - only set once
{{- if eq .ARCH "amd64" }}
ARG TLINUX_BASE_IMAGE_X86  
ARG TLINUX_TAG_X86
ARG BUILDER_IMAGE_X86
{{- else }}
ARG TLINUX_BASE_IMAGE_ARM  
ARG TLINUX_TAG_ARM
ARG BUILDER_IMAGE_ARM
{{- end }}
ARG DEPLOY_DIR={{ .DEPLOY_DIR }}

# Builder stage
{{- if eq .ARCH "amd64" }}
FROM ${BUILDER_IMAGE_X86} AS builder
{{- else }}
FROM ${BUILDER_IMAGE_ARM} AS builder
{{- end }}

# Use ARG value as ENV in builder stage
ARG DEPLOY_DIR
ENV DEPLOY_DIR=${DEPLOY_DIR}

# Set build output directory (shared between deps_install.sh and build.sh)
ENV BUILD_OUTPUT_DIR=/opt/dist
ENV PROJECT_ROOT=/opt

WORKDIR /opt

# ============================================
# Layer 1: Copy dependency files only (for caching)
# ============================================
# Copy only dependency manifest files to leverage Docker cache
# When these files don't change, Docker will reuse the cached dependency layer
# Detected dependency files for: {{ .LANGUAGE }}
{{- range .DEPENDENCY_FILES }}
COPY {{ . }} ./
{{- end }}

# Copy build scripts needed for dependency installation
COPY bk-ci/tcs/build_deps_install.sh /opt/bk-ci/tcs/

# ============================================
# Layer 2: Install dependencies (cacheable if deps files unchanged)
# ============================================
# This layer will be cached if dependency files haven't changed
# Only source code changes won't invalidate this cache layer
RUN bash -xe /opt/bk-ci/tcs/build_deps_install.sh

# ============================================
# Layer 3: Copy all source code
# ============================================
# Copy remaining source code after dependencies are installed
# This ensures dependency layer cache is preserved when only source code changes
COPY . /opt/

# ============================================
# Layer 4: Build the service
# ============================================
# Build using already installed dependencies
RUN bash -xe /opt/bk-ci/tcs/build.sh

# Runtime stage
{{- if eq .ARCH "amd64" }}
FROM ${TLINUX_BASE_IMAGE_X86}:${TLINUX_TAG_X86}
{{- else }}
FROM ${TLINUX_BASE_IMAGE_ARM}:${TLINUX_TAG_ARM}
{{- end }}

# Use ARG value in runtime stage
ARG DEPLOY_DIR

# Copy runtime preparation script
COPY bk-ci/tcs/rt_prepare.sh /tmp/rt_prepare.sh

# Install runtime dependencies
RUN sh -xe /tmp/rt_prepare.sh && rm -f /tmp/rt_prepare.sh

# Copy built artifacts from builder stage
COPY --from=builder ${DEPLOY_DIR} ${DEPLOY_DIR}

{{- if .HAS_PLUGINS }}
# ============================================
# Copy plugins from builder stage
# ============================================
# Copy all plugins from /plugins directory
COPY --from=builder /plugins /plugins

# Install plugins to their respective directories
RUN sh -xe /plugins/install.sh
{{- end }}

# Set working directory
WORKDIR ${DEPLOY_DIR}

# Set entrypoint
ENTRYPOINT ["/tce/entrypoint.sh"]
`
}
