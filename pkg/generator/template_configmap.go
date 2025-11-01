package generator

import (
	_ "embed"
	"path/filepath"
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeConfigMap = "configmap"

// ConfigMapTemplateGenerator generates Kubernetes ConfigMap using factory pattern
type ConfigMapTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the ConfigMap generator
func init() {
	RegisterGenerator(GeneratorTypeConfigMap, createConfigMapGenerator)
}

// createConfigMapGenerator is the creator function for ConfigMap generator
func createConfigMapGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewConfigMapTemplateGenerator(cfg, engine, vars), nil
}

// NewConfigMapTemplateGenerator creates a new ConfigMap template generator
func NewConfigMapTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ConfigMapTemplateGenerator {
	return &ConfigMapTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeConfigMap,
		},
	}
}

//go:embed templates/configmap.yaml.tmpl
var configmapTemplate string

// Generate generates ConfigMap YAML content
func (g *ConfigMapTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for ConfigMap template
func (g *ConfigMapTemplateGenerator) prepareTemplateVars() map[string]interface{} {
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

// isConfigFile checks if a file is a configuration file
func (g *ConfigMapTemplateGenerator) isConfigFile(path string) bool {
	configExts := []string{".yaml", ".yml", ".json", ".toml", ".ini", ".conf", ".config"}
	ext := strings.ToLower(filepath.Ext(path))

	for _, configExt := range configExts {
		if ext == configExt {
			return true
		}
	}

	return false
}

// getTemplate returns the Kubernetes ConfigMap template
func (g *ConfigMapTemplateGenerator) getTemplate() string {
	return configmapTemplate
}
