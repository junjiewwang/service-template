package generator

import (
	_ "embed"

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

//go:embed templates/deps_install.sh.tmpl
var depsInstallScriptTemplate string

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
	return depsInstallScriptTemplate
}
