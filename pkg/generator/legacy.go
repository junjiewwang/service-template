package generator

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// Legacy generator wrappers for backward compatibility with tests

// DockerfileGenerator is a legacy wrapper for DockerfileTemplateGenerator
type DockerfileGenerator struct {
	generator *DockerfileTemplateGenerator
	arch      string
}

// NewDockerfileGenerator creates a new Dockerfile generator (legacy)
func NewDockerfileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *DockerfileGenerator {
	return &DockerfileGenerator{
		generator: NewDockerfileTemplateGenerator(cfg, engine, vars, "amd64"),
		arch:      "amd64",
	}
}

// Generate generates a Dockerfile for the specified architecture (legacy)
func (g *DockerfileGenerator) Generate(arch string) (string, error) {
	// Create a new generator with the specified architecture
	g.generator.arch = arch
	g.generator.name = "dockerfile-" + arch
	return g.generator.Generate()
}

// ComposeGenerator is a legacy wrapper for ComposeTemplateGenerator
type ComposeGenerator struct {
	generator *ComposeTemplateGenerator
}

// NewComposeGenerator creates a new Compose generator (legacy)
func NewComposeGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ComposeGenerator {
	return &ComposeGenerator{
		generator: NewComposeTemplateGenerator(cfg, engine, vars),
	}
}

// Generate generates docker-compose.yaml content (legacy)
func (g *ComposeGenerator) Generate() (string, error) {
	return g.generator.Generate()
}

// MakefileGenerator is a legacy wrapper for MakefileTemplateGenerator
type MakefileGenerator struct {
	generator *MakefileTemplateGenerator
}

// NewMakefileGenerator creates a new Makefile generator (legacy)
func NewMakefileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *MakefileGenerator {
	return &MakefileGenerator{
		generator: NewMakefileTemplateGenerator(cfg, engine, vars),
	}
}

// Generate generates Makefile content (legacy)
func (g *MakefileGenerator) Generate() (string, error) {
	return g.generator.Generate()
}

// DevOpsGenerator is a legacy wrapper for DevOpsTemplateGenerator
type DevOpsGenerator struct {
	generator *DevOpsTemplateGenerator
}

// NewDevOpsGenerator creates a new DevOps generator (legacy)
func NewDevOpsGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *DevOpsGenerator {
	return &DevOpsGenerator{
		generator: NewDevOpsTemplateGenerator(cfg, engine, vars),
	}
}

// Generate generates devops.yaml content (legacy)
func (g *DevOpsGenerator) Generate() (string, error) {
	return g.generator.Generate()
}

// ConfigMapGenerator is a legacy wrapper for ConfigMapTemplateGenerator
type ConfigMapGenerator struct {
	generator *ConfigMapTemplateGenerator
}

// NewConfigMapGenerator creates a new ConfigMap generator (legacy)
func NewConfigMapGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ConfigMapGenerator {
	return &ConfigMapGenerator{
		generator: NewConfigMapTemplateGenerator(cfg, engine, vars),
	}
}

// Generate generates ConfigMap YAML content (legacy)
func (g *ConfigMapGenerator) Generate() (string, error) {
	return g.generator.Generate()
}

// ScriptsGenerator is a legacy wrapper for script template generators
type ScriptsGenerator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
}

// NewScriptsGenerator creates a new scripts generator (legacy)
func NewScriptsGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *ScriptsGenerator {
	return &ScriptsGenerator{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
	}
}

// GenerateBuildScript generates build.sh (legacy)
func (g *ScriptsGenerator) GenerateBuildScript() (string, error) {
	gen := NewBuildScriptTemplateGenerator(g.config, g.templateEngine, g.variables)
	return gen.Generate()
}

// GenerateDepsInstallScript generates build_deps_install.sh (legacy)
func (g *ScriptsGenerator) GenerateDepsInstallScript() (string, error) {
	gen := NewDepsInstallScriptTemplateGenerator(g.config, g.templateEngine, g.variables)
	return gen.Generate()
}

// GenerateRtPrepareScript generates rt_prepare.sh (legacy)
func (g *ScriptsGenerator) GenerateRtPrepareScript() (string, error) {
	gen := NewRtPrepareScriptTemplateGenerator(g.config, g.templateEngine, g.variables)
	return gen.Generate()
}

// GenerateEntrypointScript generates entrypoint.sh (legacy)
func (g *ScriptsGenerator) GenerateEntrypointScript() (string, error) {
	gen := NewEntrypointScriptTemplateGenerator(g.config, g.templateEngine, g.variables)
	return gen.Generate()
}

// GenerateHealthchkScript generates healthchk.sh (legacy)
func (g *ScriptsGenerator) GenerateHealthchkScript() (string, error) {
	gen := NewHealthcheckScriptTemplateGenerator(g.config, g.templateEngine, g.variables)
	return gen.Generate()
}
