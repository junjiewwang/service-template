package adapter

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"

	// Import all legacy generators
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/build_tools/makefile"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/compose"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/devops"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/dockerfile"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/build"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/build_plugins"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/deps_install"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/entrypoint"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/healthcheck"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/rt_prepare"
)

// LegacyGeneratorAdapter adapts old generators to DDD architecture
type LegacyGeneratorAdapter struct {
	config    *config.ServiceConfig
	outputDir string
	genCtx    *context.GeneratorContext
}

// NewLegacyGeneratorAdapter creates a new adapter
func NewLegacyGeneratorAdapter(cfg *config.ServiceConfig, outputDir string) *LegacyGeneratorAdapter {
	genCtx := context.NewGeneratorContext(cfg, outputDir)
	return &LegacyGeneratorAdapter{
		config:    cfg,
		outputDir: outputDir,
		genCtx:    genCtx,
	}
}

// GenerateDockerfile generates Dockerfile for specific architecture
func (a *LegacyGeneratorAdapter) GenerateDockerfile(arch string) (string, error) {
	creator, exists := core.DefaultRegistry.Get("dockerfile")
	if !exists {
		return "", fmt.Errorf("dockerfile generator not found")
	}

	generator, err := creator(a.genCtx, arch)
	if err != nil {
		return "", fmt.Errorf("failed to create dockerfile generator: %w", err)
	}

	return generator.Generate()
}

// GenerateCompose generates docker-compose.yaml
func (a *LegacyGeneratorAdapter) GenerateCompose() (string, error) {
	return a.generateByType("compose")
}

// GenerateMakefile generates Makefile
func (a *LegacyGeneratorAdapter) GenerateMakefile() (string, error) {
	return a.generateByType("makefile")
}

// GenerateDevOps generates devops.yaml
func (a *LegacyGeneratorAdapter) GenerateDevOps() (string, error) {
	return a.generateByType("devops")
}

// GenerateBuildScript generates build.sh
func (a *LegacyGeneratorAdapter) GenerateBuildScript() (string, error) {
	return a.generateByType("build-script")
}

// GenerateBuildPluginsScript generates build_plugins.sh
func (a *LegacyGeneratorAdapter) GenerateBuildPluginsScript() (string, error) {
	return a.generateByType("build-plugins-script")
}

// GenerateDepsInstallScript generates deps_install.sh
func (a *LegacyGeneratorAdapter) GenerateDepsInstallScript() (string, error) {
	return a.generateByType("deps-install-script")
}

// GenerateEntrypointScript generates entrypoint.sh
func (a *LegacyGeneratorAdapter) GenerateEntrypointScript() (string, error) {
	return a.generateByType("entrypoint-script")
}

// GenerateHealthcheckScript generates healthcheck.sh
func (a *LegacyGeneratorAdapter) GenerateHealthcheckScript() (string, error) {
	return a.generateByType("healthcheck-script")
}

// GenerateRtPrepareScript generates rt_prepare.sh
func (a *LegacyGeneratorAdapter) GenerateRtPrepareScript() (string, error) {
	return a.generateByType("rt-prepare-script")
}

// generateByType generates content using legacy generator by type
func (a *LegacyGeneratorAdapter) generateByType(generatorType string) (string, error) {
	creator, exists := core.DefaultRegistry.Get(generatorType)
	if !exists {
		return "", fmt.Errorf("generator type %s not found", generatorType)
	}

	generator, err := creator(a.genCtx)
	if err != nil {
		return "", fmt.Errorf("failed to create %s generator: %w", generatorType, err)
	}

	return generator.Generate()
}
