package generator

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeCompose = "compose"

// ComposeTemplateGenerator generates docker-compose.yaml using factory pattern
type ComposeTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the Compose generator
func init() {
	RegisterGenerator(GeneratorTypeCompose, createComposeGenerator)
}

// createComposeGenerator is the creator function for Compose generator
func createComposeGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewComposeTemplateGenerator(cfg, engine, vars), nil
}

// NewComposeTemplateGenerator creates a new Compose template generator
func NewComposeTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ComposeTemplateGenerator {
	return &ComposeTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeCompose,
		},
	}
}

//go:embed templates/compose.yaml.tmpl
var composeTemplate string

// Generate generates docker-compose.yaml content
func (g *ComposeTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for compose template
func (g *ComposeTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name
	vars["SERVICE_ROOT"] = g.config.Service.DeployDir + "/" + g.config.Service.Name

	// Ports
	type PortMapping struct {
		Port       int
		TargetPort int
	}
	var ports []PortMapping
	for _, port := range g.config.Service.Ports {
		ports = append(ports, PortMapping{
			Port:       port.Port,
			TargetPort: port.Port,
		})
	}
	vars["PORTS"] = ports

	// Environment variables
	vars["ENV_VARS"] = g.config.Runtime.Startup.Env

	// Volumes
	type VolumeMapping struct {
		Source string
		Target string
	}
	var volumes []VolumeMapping
	for _, vol := range g.config.LocalDev.Compose.Volumes {
		// 扩展变量映射以支持SERVICE_ROOT和PLUGIN_INSTALL_DIR
		variableMap := g.variables.ToMap()
		variableMap["SERVICE_ROOT"] = g.config.Service.DeployDir + "/" + g.config.Service.Name
		variableMap["PLUGIN_INSTALL_DIR"] = g.config.Plugins.InstallDir

		// 支持本地开发配置中定义的其他变量
		if len(g.config.LocalDev.SupportedVariables) > 0 {
			for _, supportedVar := range g.config.LocalDev.SupportedVariables {
				switch supportedVar {
				case "SERVICE_ROOT":
					variableMap[supportedVar] = g.config.Service.DeployDir + "/" + g.config.Service.Name
				case "PLUGIN_INSTALL_DIR":
					variableMap[supportedVar] = g.config.Plugins.InstallDir
				}
			}
		}

		volumes = append(volumes, VolumeMapping{
			Source: vol.Source,
			Target: SubstituteVariables(vol.Target, variableMap),
		})
	}
	vars["VOLUMES"] = volumes

	// Resources
	vars["LIMITS_CPUS"] = g.config.LocalDev.Compose.Resources.Limits.CPUs
	vars["LIMITS_MEMORY"] = g.config.LocalDev.Compose.Resources.Limits.Memory
	vars["RESERVATIONS_CPUS"] = g.config.LocalDev.Compose.Resources.Reservations.CPUs
	vars["RESERVATIONS_MEMORY"] = g.config.LocalDev.Compose.Resources.Reservations.Memory

	// Health check
	vars["HEALTHCHECK_ENABLED"] = g.config.Runtime.Healthcheck.Enabled
	vars["HEALTHCHECK_TYPE"] = g.config.Runtime.Healthcheck.Type
	vars["HEALTHCHECK_INTERVAL"] = g.config.LocalDev.Compose.Healthcheck.Interval
	vars["HEALTHCHECK_TIMEOUT"] = g.config.LocalDev.Compose.Healthcheck.Timeout
	vars["HEALTHCHECK_RETRIES"] = g.config.LocalDev.Compose.Healthcheck.Retries
	vars["HEALTHCHECK_START_PERIOD"] = g.config.LocalDev.Compose.Healthcheck.StartPeriod

	// Labels
	vars["LABELS"] = g.config.LocalDev.Compose.Labels

	return vars
}

// getTemplate returns the docker-compose.yaml template
func (g *ComposeTemplateGenerator) getTemplate() string {
	return composeTemplate
}
