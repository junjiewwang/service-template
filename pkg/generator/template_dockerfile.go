package generator

import (
	_ "embed"
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

//go:embed templates/dockerfile_.tmpl
var dockerfileTemplate string

// Generate generates Dockerfile content
func (g *DockerfileTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for Dockerfile template
func (g *DockerfileTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	// 从 variables 获取所有基础变量（包括 CI 路径变量）
	vars := g.variables.ToMap()

	// Basic info (覆盖或添加特定变量)
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
	return dockerfileTemplate
}
