package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/utils"

	// Import all generators to register them
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/build_tools/makefile"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/compose"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/devops"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/dockerfile"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/build"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/deps_install"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/entrypoint"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/healthcheck"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/rt_prepare"
)

// Generator is the main generator that orchestrates all generation tasks
type Generator struct {
	config    *config.ServiceConfig
	ctx       *context.GeneratorContext
	outputDir string
}

// NewGenerator creates a new generator instance
func NewGenerator(cfg *config.ServiceConfig, outputDir string) *Generator {
	ctx := context.NewGeneratorContext(cfg, outputDir)

	// Update metadata
	cfg.Metadata.GeneratedAt = time.Now().Format(time.RFC3339)

	return &Generator{
		config:    cfg,
		ctx:       ctx,
		outputDir: outputDir,
	}
}

// Generate generates all project files
func (g *Generator) Generate() error {
	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate Dockerfiles
	if err := g.generateDockerfiles(); err != nil {
		return fmt.Errorf("failed to generate Dockerfiles: %w", err)
	}

	// Generate Docker Compose
	if err := g.generateCompose(); err != nil {
		return fmt.Errorf("failed to generate compose.yaml: %w", err)
	}

	// Generate Makefile
	if err := g.generateMakefile(); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}

	// Generate build scripts
	if err := g.generateScripts(); err != nil {
		return fmt.Errorf("failed to generate scripts: %w", err)
	}

	// Generate DevOps configuration
	if err := g.generateDevOps(); err != nil {
		return fmt.Errorf("failed to generate DevOps configuration: %w", err)
	}

	fmt.Println("✓ Project generated successfully!")
	return nil
}

// generateDockerfiles generates Dockerfiles for different architectures
func (g *Generator) generateDockerfiles() error {
	architectures := []string{"amd64", "arm64"}

	// Create .tad/build/{service-name} directory
	dockerfileDir := filepath.Join(g.outputDir, ".tad", "build", g.config.Service.Name)
	if err := os.MkdirAll(dockerfileDir, 0755); err != nil {
		return fmt.Errorf("failed to create dockerfile directory: %w", err)
	}

	for _, arch := range architectures {
		// Create generator using new registry
		creator, exists := core.DefaultRegistry.Get("dockerfile")
		if !exists {
			return fmt.Errorf("generator type dockerfile not found")
		}

		generator, err := creator(g.ctx, arch)
		if err != nil {
			return fmt.Errorf("failed to create dockerfile generator for %s: %w", arch, err)
		}

		content, err := generator.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate Dockerfile for %s: %w", arch, err)
		}

		// Generate Dockerfile with format: Dockerfile.{service-name}.{arch}
		filename := fmt.Sprintf("Dockerfile.%s.%s", g.config.Service.Name, arch)
		outputPath := filepath.Join(dockerfileDir, filename)
		if err := utils.WriteFile(outputPath, content); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		fmt.Printf("✓ Generated .tad/build/%s/%s\n", g.config.Service.Name, filename)
	}

	return nil
}

// generateCompose generates docker-compose.yaml
func (g *Generator) generateCompose() error {
	generator, err := g.createGenerator("compose")
	if err != nil {
		return fmt.Errorf("failed to create compose generator: %w", err)
	}

	content, err := generator.Generate()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(g.outputDir, "compose.yaml")
	if err := utils.WriteFile(outputPath, content); err != nil {
		return err
	}

	fmt.Println("✓ Generated compose.yaml")
	return nil
}

// generateMakefile generates Makefile
func (g *Generator) generateMakefile() error {
	generator, err := g.createGenerator("makefile")
	if err != nil {
		return fmt.Errorf("failed to create makefile generator: %w", err)
	}

	content, err := generator.Generate()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(g.outputDir, "Makefile")
	if err := utils.WriteFile(outputPath, content); err != nil {
		return err
	}

	fmt.Println("✓ Generated Makefile")
	return nil
}

// generateScripts generates build and deployment scripts
func (g *Generator) generateScripts() error {
	// Script types and their output paths
	scripts := map[string]string{
		"build-script":        g.ctx.Paths.CI.GetScriptPath(g.ctx.Paths.CI.BuildScript),
		"deps-install-script": g.ctx.Paths.CI.GetScriptPath(g.ctx.Paths.CI.DepsInstallScript),
		"rt-prepare-script":   g.ctx.Paths.CI.GetScriptPath(g.ctx.Paths.CI.RtPrepareScript),
		"entrypoint-script":   g.ctx.Paths.CI.GetScriptPath(g.ctx.Paths.CI.EntrypointScript),
		"healthcheck-script":  g.ctx.Paths.CI.GetScriptPath(g.ctx.Paths.CI.HealthcheckScript),
	}

	for generatorType, scriptPath := range scripts {
		generator, err := g.createGenerator(generatorType)
		if err != nil {
			return fmt.Errorf("failed to create %s generator: %w", generatorType, err)
		}

		content, err := generator.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", scriptPath, err)
		}

		outputPath := filepath.Join(g.outputDir, scriptPath)
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", scriptPath, err)
		}

		if err := utils.WriteFile(outputPath, content); err != nil {
			return fmt.Errorf("failed to write %s: %w", scriptPath, err)
		}

		// Make shell scripts executable
		if filepath.Ext(scriptPath) == ".sh" {
			if err := os.Chmod(outputPath, 0755); err != nil {
				return fmt.Errorf("failed to make %s executable: %w", scriptPath, err)
			}
		}

		fmt.Printf("✓ Generated %s\n", scriptPath)
	}

	return nil
}

// generateDevOps generates DevOps configuration
func (g *Generator) generateDevOps() error {
	generator, err := g.createGenerator("devops")
	if err != nil {
		return fmt.Errorf("failed to create devops generator: %w", err)
	}

	content, err := generator.Generate()
	if err != nil {
		return err
	}

	tadDir := filepath.Join(g.outputDir, ".tad")
	if err := os.MkdirAll(tadDir, 0755); err != nil {
		return fmt.Errorf("failed to create .tad directory: %w", err)
	}

	outputPath := filepath.Join(tadDir, "devops.yaml")
	if err := utils.WriteFile(outputPath, content); err != nil {
		return err
	}

	fmt.Println("✓ Generated .tad/devops.yaml")
	return nil
}

// createGenerator creates a generator using the new registry
func (g *Generator) createGenerator(generatorType string) (core.Generator, error) {
	creator, exists := core.DefaultRegistry.Get(generatorType)
	if !exists {
		return nil, fmt.Errorf("generator type %s not found (available: %v)",
			generatorType, core.DefaultRegistry.GetAll())
	}

	return creator(g.ctx)
}
