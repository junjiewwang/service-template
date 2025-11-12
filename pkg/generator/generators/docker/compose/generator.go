package compose

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "compose"

// init registers the compose generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates docker-compose.yaml
type Generator struct {
	core.BaseGenerator
}

// New creates a new compose generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates docker-compose.yaml content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	vars := g.prepareTemplateVars()
	return g.RenderTemplate(template, vars)
}

// prepareTemplateVars prepares variables for compose template
func (g *Generator) prepareTemplateVars() map[string]interface{} {
	ctx := g.GetContext()

	// Use preset for compose
	composer := ctx.GetVariablePreset().ForCompose()

	// Prepare custom port mappings
	type PortMapping struct {
		Port       int
		TargetPort int
	}
	var ports []PortMapping
	for _, port := range ctx.Config.Service.Ports {
		ports = append(ports, PortMapping{
			Port:       port.Port,
			TargetPort: port.Port,
		})
	}

	// Prepare volumes with variable substitution
	type VolumeMapping struct {
		Source string
		Target string
	}
	var volumes []VolumeMapping
	for _, vol := range ctx.Config.LocalDev.Compose.Volumes {
		// Build variable map for substitution
		variableMap := composer.Build()
		variableMap["PLUGIN_INSTALL_DIR"] = ctx.Config.Plugins.InstallDir

		// Support additional variables from local dev config
		if len(ctx.Config.LocalDev.SupportedVariables) > 0 {
			for _, supportedVar := range ctx.Config.LocalDev.SupportedVariables {
				switch supportedVar {
				case "SERVICE_ROOT":
					variableMap[supportedVar] = ctx.Config.Service.DeployDir + "/" + ctx.Config.Service.Name
				case "PLUGIN_INSTALL_DIR":
					variableMap[supportedVar] = ctx.Config.Plugins.InstallDir
				}
			}
		}

		volumes = append(volumes, VolumeMapping{
			Source: vol.Source,
			Target: core.SubstituteVariables(vol.Target, variableMap),
		})
	}

	// Add compose-specific custom variables
	composer.
		Override("PORTS", ports).
		WithCustom("VOLUMES", volumes).
		WithCustom("LIMITS_CPUS", ctx.Config.LocalDev.Compose.Resources.Limits.CPUs).
		WithCustom("LIMITS_MEMORY", ctx.Config.LocalDev.Compose.Resources.Limits.Memory).
		WithCustom("RESERVATIONS_CPUS", ctx.Config.LocalDev.Compose.Resources.Reservations.CPUs).
		WithCustom("RESERVATIONS_MEMORY", ctx.Config.LocalDev.Compose.Resources.Reservations.Memory).
		WithCustom("HEALTHCHECK_INTERVAL", ctx.Config.LocalDev.Compose.Healthcheck.Interval).
		WithCustom("HEALTHCHECK_TIMEOUT", ctx.Config.LocalDev.Compose.Healthcheck.Timeout).
		WithCustom("HEALTHCHECK_RETRIES", ctx.Config.LocalDev.Compose.Healthcheck.Retries).
		WithCustom("HEALTHCHECK_START_PERIOD", ctx.Config.LocalDev.Compose.Healthcheck.StartPeriod).
		WithCustom("LABELS", ctx.Config.LocalDev.Compose.Labels)

	return composer.Build()
}

//go:embed templates/compose.yaml.tmpl
var template string
