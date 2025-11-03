package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/utils"
)

// Generator is the main generator that orchestrates all generation tasks
type Generator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
	factory        *GeneratorFactory
	outputDir      string
}

// NewGenerator creates a new generator instance
func NewGenerator(cfg *config.ServiceConfig, outputDir string) *Generator {
	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	return &Generator{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
		factory:        NewGeneratorFactory(cfg, engine, vars),
		outputDir:      outputDir,
	}
}

// Generate generates all project files
func (g *Generator) Generate() error {
	// Update metadata
	g.config.Metadata.GeneratedAt = time.Now().Format(time.RFC3339)

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
		// Use factory to create generator
		generator, err := g.factory.CreateGenerator("dockerfile", arch)
		if err != nil {
			return fmt.Errorf("failed to create dockerfile generator: %w", err)
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
	generator, err := g.factory.CreateGenerator("compose")
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
	generator, err := g.factory.CreateGenerator("makefile")
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
	scripts := []struct {
		name          string
		generatorType string
	}{
		{"bk-ci/tcs/build.sh", "build-script"},
		{"bk-ci/tcs/build_deps_install.sh", "deps-install-script"},
		{"bk-ci/tcs/rt_prepare.sh", "rt-prepare-script"},
		{"bk-ci/tcs/entrypoint.sh", "entrypoint-script"},
		{"bk-ci/tcs/healthchk.sh", "healthcheck-script"},
	}

	for _, script := range scripts {
		generator, err := g.factory.CreateGenerator(script.generatorType)
		if err != nil {
			return fmt.Errorf("failed to create %s generator: %w", script.generatorType, err)
		}

		content, err := generator.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", script.name, err)
		}

		outputPath := filepath.Join(g.outputDir, script.name)
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", script.name, err)
		}

		if err := utils.WriteFile(outputPath, content); err != nil {
			return fmt.Errorf("failed to write %s: %w", script.name, err)
		}

		// Make shell scripts executable
		if filepath.Ext(script.name) == ".sh" {
			if err := os.Chmod(outputPath, 0755); err != nil {
				return fmt.Errorf("failed to make %s executable: %w", script.name, err)
			}
		}

		fmt.Printf("✓ Generated %s\n", script.name)
	}

	return nil
}

// generateDevOps generates DevOps configuration
func (g *Generator) generateDevOps() error {
	generator, err := g.factory.CreateGenerator("devops")
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
