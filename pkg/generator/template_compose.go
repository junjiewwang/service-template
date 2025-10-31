package generator

import (
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

// Generate generates docker-compose.yaml content
func (g *ComposeTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(composeTemplate, vars)
}

// prepareTemplateVars prepares variables for compose template
func (g *ComposeTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name

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
		variableMap["PLUGIN_INSTALL_DIR"] = "/plugins"

		// 支持本地开发配置中定义的其他变量
		if len(g.config.LocalDev.SupportedVariables) > 0 {
			for _, supportedVar := range g.config.LocalDev.SupportedVariables {
				switch supportedVar {
				case "SERVICE_ROOT":
					variableMap[supportedVar] = g.config.Service.DeployDir + "/" + g.config.Service.Name
				case "PLUGIN_INSTALL_DIR":
					variableMap[supportedVar] = "/plugins"
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
	vars["HEALTHCHECK_HTTP_PORT"] = g.config.Runtime.Healthcheck.HTTP.Port
	vars["HEALTHCHECK_HTTP_PATH"] = g.config.Runtime.Healthcheck.HTTP.Path
	vars["HEALTHCHECK_INTERVAL"] = g.config.LocalDev.Compose.Healthcheck.Interval
	vars["HEALTHCHECK_TIMEOUT"] = g.config.LocalDev.Compose.Healthcheck.Timeout
	vars["HEALTHCHECK_RETRIES"] = g.config.LocalDev.Compose.Healthcheck.Retries
	vars["HEALTHCHECK_START_PERIOD"] = g.config.LocalDev.Compose.Healthcheck.StartPeriod

	// Labels
	vars["LABELS"] = g.config.LocalDev.Compose.Labels

	return vars
}

// composeTemplate is the docker-compose.yaml template
const composeTemplate = `# Auto-generated docker-compose.yaml
# Generated at: {{ .GENERATED_AT }}

version: '3.8'

services:
  {{ .SERVICE_NAME }}:
    image: {{ .SERVICE_NAME }}:latest-amd64
    container_name: {{ .SERVICE_NAME }}
{{- if .PORTS }}
    ports:
{{- range .PORTS }}
{{- if eq .Port .TargetPort }}
      - "{{ .Port }}"
{{- else }}
      - "{{ .Port }}:{{ .TargetPort }}"
{{- end }}
{{- end }}
{{- end }}
{{- if .ENV_VARS }}
    environment:
{{- range .ENV_VARS }}
      - {{ .Name }}={{ .Value }}
{{- end }}
{{- end }}
{{- if .VOLUMES }}
    volumes:
{{- range .VOLUMES }}
      - {{ .Source }}:{{ .Target }}
{{- end }}
{{- end }}
{{- if or .LIMITS_CPUS .LIMITS_MEMORY .RESERVATIONS_CPUS .RESERVATIONS_MEMORY }}
    deploy:
      resources:
{{- if or .LIMITS_CPUS .LIMITS_MEMORY }}
        limits:
{{- if .LIMITS_CPUS }}
          cpus: '{{ .LIMITS_CPUS }}'
{{- end }}
{{- if .LIMITS_MEMORY }}
          memory: {{ .LIMITS_MEMORY }}
{{- end }}
{{- end }}
{{- if or .RESERVATIONS_CPUS .RESERVATIONS_MEMORY }}
        reservations:
{{- if .RESERVATIONS_CPUS }}
          cpus: '{{ .RESERVATIONS_CPUS }}'
{{- end }}
{{- if .RESERVATIONS_MEMORY }}
          memory: {{ .RESERVATIONS_MEMORY }}
{{- end }}
{{- end }}
{{- end }}
{{- if .HEALTHCHECK_ENABLED }}
    healthcheck:
{{- if eq .HEALTHCHECK_TYPE "http" }}
      test: ["CMD", "curl", "-f", "http://localhost:{{ .HEALTHCHECK_HTTP_PORT }}{{ .HEALTHCHECK_HTTP_PATH }}"]
{{- else }}
      test: ["CMD", "/bin/sh", "/usr/local/services/${SERVICE_NAME}/hooks/healthchk.sh"]
{{- end }}
{{- if .HEALTHCHECK_INTERVAL }}
      interval: {{ .HEALTHCHECK_INTERVAL }}
{{- end }}
{{- if .HEALTHCHECK_TIMEOUT }}
      timeout: {{ .HEALTHCHECK_TIMEOUT }}
{{- end }}
{{- if .HEALTHCHECK_RETRIES }}
      retries: {{ .HEALTHCHECK_RETRIES }}
{{- end }}
{{- if .HEALTHCHECK_START_PERIOD }}
      start_period: {{ .HEALTHCHECK_START_PERIOD }}
{{- end }}
{{- end }}
{{- if .LABELS }}
    labels:
{{- range $key, $value := .LABELS }}
      {{ $key }}: "{{ $value }}"
{{- end }}
{{- end }}
    restart: unless-stopped
`
