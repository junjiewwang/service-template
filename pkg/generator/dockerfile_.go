package generator

import (
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// DockerfileGenerator generates Dockerfiles
type DockerfileGenerator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
}

// NewDockerfileGenerator creates a new Dockerfile generator
func NewDockerfileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *DockerfileGenerator {
	return &DockerfileGenerator{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
	}
}

// Generate generates a Dockerfile for the specified architecture
func (g *DockerfileGenerator) Generate(arch string) (string, error) {
	// Use embedded template
	templateContent := g.getDefaultDockerfileTemplate()

	// Prepare template variables
	vars := g.prepareTemplateVars(arch)

	return g.templateEngine.Render(templateContent, vars)
}

// prepareTemplateVars prepares variables for Dockerfile template
func (g *DockerfileGenerator) prepareTemplateVars(arch string) map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["ARCH"] = arch
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name

	// Images
	if arch == "amd64" {
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
	var plugins []map[string]string
	for _, plugin := range g.config.Plugins {
		pluginVars := g.variables.WithPlugin(plugin)
		plugins = append(plugins, map[string]string{
			"InstallCommand": SubstituteVariables(plugin.InstallCommand, pluginVars.ToMap()),
		})
	}
	vars["PLUGINS"] = plugins

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

// getDefaultDockerfileTemplate returns embedded default Dockerfile template
func (g *DockerfileGenerator) getDefaultDockerfileTemplate() string {
	return `# Auto-generated Dockerfile for {{ .ARCH }}
# Generated at: {{ .GENERATED_AT }}

# ============================================
# Build Arguments (from devops.yaml)
# ============================================
{{- if eq .ARCH "amd64" }}
ARG BUILDER_IMAGE_X86
ARG TLINUX_BASE_IMAGE_X86
ARG TLINUX_TAG_X86
{{- else }}
ARG BUILDER_IMAGE_ARM
ARG TLINUX_BASE_IMAGE_ARM
ARG TLINUX_TAG_ARM
{{- end }}
ARG DEPLOY_DIR

# ============================================
# Build Stage
# ============================================
{{- if eq .ARCH "amd64" }}
FROM ${BUILDER_IMAGE_X86} AS builder
{{- else }}
FROM ${BUILDER_IMAGE_ARM} AS builder
{{- end }}

WORKDIR /build

{{- if .BUILD_DEPS_PACKAGES }}

# Install build dependencies
RUN {{ .PKG_MANAGER }} install -y {{ join " " .BUILD_DEPS_PACKAGES }}
{{- end }}

# Copy dependency files
{{- range .DEPENDENCY_FILES }}
COPY {{ . }} ./
{{- end }}

# Install dependencies
RUN {{ .DEPS_INSTALL_COMMAND }}

# Copy source code
COPY . .

# Copy scripts directory and ensure they are executable
COPY scripts/ ./scripts/
RUN chmod +x scripts/*.sh

{{- if or .PRE_BUILD_COMMAND .BUILD_COMMAND .POST_BUILD_COMMAND }}

# Build commands
RUN ./scripts/build.sh
{{- end }}

# ============================================
# Runtime Stage
# ============================================
{{- if eq .ARCH "amd64" }}
FROM ${TLINUX_BASE_IMAGE_X86}:${TLINUX_TAG_X86}
{{- else }}
FROM ${TLINUX_BASE_IMAGE_ARM}:${TLINUX_TAG_ARM}
{{- end }}

{{- if eq .ARCH "amd64" }}
ARG BUILDER_IMAGE_X86
ARG TLINUX_BASE_IMAGE_X86
ARG TLINUX_TAG_X86
{{- else }}
ARG BUILDER_IMAGE_ARM
ARG TLINUX_BASE_IMAGE_ARM
ARG TLINUX_TAG_ARM
{{- end }}
ARG DEPLOY_DIR

{{- if .RUNTIME_DEPS_PACKAGES }}

# Install runtime dependencies
RUN {{ .PKG_MANAGER }} install -y {{ join " " .RUNTIME_DEPS_PACKAGES }}
{{- end }}

{{- if or .PLUGINS }}

# Install plugins and prepare runtime
RUN ./scripts/rt_prepare.sh
{{- end }}

# Create service directory
RUN mkdir -p ${DEPLOY_DIR}

# Copy build artifacts
COPY --from=builder /build/{{ .BUILD_OUTPUT_DIR }}/ ${DEPLOY_DIR}/

# Copy hooks
COPY hooks/ ${DEPLOY_DIR}/hooks/

# Set working directory
WORKDIR ${DEPLOY_DIR}

{{- if .EXPOSE_PORTS }}

# Expose ports
{{- range .EXPOSE_PORTS }}
EXPOSE {{ . }}
{{- end }}
{{- end }}

{{- if .HEALTHCHECK_ENABLED }}

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=40s \
  CMD /bin/sh ${DEPLOY_DIR}/hooks/healthchk.sh
{{- end }}

# Start command
CMD ["/bin/sh", "${DEPLOY_DIR}/hooks/start.sh"]
`
}

// getDependencyFilesList returns list of dependency files
func (g *DockerfileGenerator) getDependencyFilesList() []string {
	if g.config.Build.DependencyFiles.AutoDetect {
		// Auto-detect based on language
		return getDefaultDependencyFiles(g.config.Language.Type)
	}

	return g.config.Build.DependencyFiles.Files
}

// getDepsInstallCommand generates dependency installation command
func (g *DockerfileGenerator) getDepsInstallCommand() string {
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

	// Default to yum for unknown images
	return "yum"
}
