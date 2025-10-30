package generator

import (
	"path/filepath"
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// ConfigMapGenerator generates Kubernetes ConfigMap
type ConfigMapGenerator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
}

// NewConfigMapGenerator creates a new ConfigMap generator
func NewConfigMapGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ConfigMapGenerator {
	return &ConfigMapGenerator{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
	}
}

// Generate generates ConfigMap YAML content
func (g *ConfigMapGenerator) Generate() (string, error) {
	// Use embedded template
	templateContent := g.getDefaultConfigMapTemplate()

	// Prepare template variables
	vars := g.prepareTemplateVars()

	return g.templateEngine.Render(templateContent, vars)
}

// prepareTemplateVars prepares variables for ConfigMap template
func (g *ConfigMapGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name
	vars["NAMESPACE"] = g.config.LocalDev.Kubernetes.Namespace

	// ConfigMap name
	configMapName := g.config.Service.Name + "-config"
	if g.config.LocalDev.Kubernetes.ConfigMap.Name != "" {
		configMapName = g.config.LocalDev.Kubernetes.ConfigMap.Name
	}
	vars["CONFIGMAP_NAME"] = configMapName

	// Extract config files from volumes
	type ConfigFile struct {
		FileName string
		Source   string
	}
	var configFiles []ConfigFile
	for _, vol := range g.config.LocalDev.Compose.Volumes {
		if vol.Type == "bind" && g.isConfigFile(vol.Source) {
			configFiles = append(configFiles, ConfigFile{
				FileName: filepath.Base(vol.Source),
				Source:   vol.Source,
			})
		}
	}
	vars["CONFIG_FILES"] = configFiles

	return vars
}

// getDefaultConfigMapTemplate returns embedded default ConfigMap template
func (g *ConfigMapGenerator) getDefaultConfigMapTemplate() string {
	return `# Auto-generated Kubernetes ConfigMap
# Generated at: {{ .GENERATED_AT }}

apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .CONFIGMAP_NAME }}
  namespace: {{ .NAMESPACE }}
  labels:
    app: {{ .SERVICE_NAME }}
data:
{{- if .CONFIG_FILES }}
{{- range .CONFIG_FILES }}
  {{ .FileName }}: |
    # Config file content should be provided here
    # Source: {{ .Source }}
{{- end }}
{{- else }}
  # No config files detected
  # Add your configuration data here
{{- end }}
`
}

// isConfigFile checks if a file is a configuration file
func (g *ConfigMapGenerator) isConfigFile(path string) bool {
	configExts := []string{".yaml", ".yml", ".json", ".toml", ".ini", ".conf", ".config"}
	ext := strings.ToLower(filepath.Ext(path))

	for _, configExt := range configExts {
		if ext == configExt {
			return true
		}
	}

	return false
}
