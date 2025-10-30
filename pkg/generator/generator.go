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
	outputDir      string
}

// NewGenerator creates a new generator instance
func NewGenerator(cfg *config.ServiceConfig, outputDir string) *Generator {
	return &Generator{
		config:         cfg,
		templateEngine: NewTemplateEngine(),
		variables:      NewVariables(cfg),
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

	// Generate hooks
	if err := g.generateHooks(); err != nil {
		return fmt.Errorf("failed to generate hooks: %w", err)
	}

	// Generate DevOps configuration
	if err := g.generateDevOps(); err != nil {
		return fmt.Errorf("failed to generate DevOps configuration: %w", err)
	}

	// Generate ConfigMap if needed
	if g.config.LocalDev.Kubernetes.Enabled && g.config.LocalDev.Kubernetes.ConfigMap.AutoDetect {
		if err := g.generateConfigMap(); err != nil {
			return fmt.Errorf("failed to generate ConfigMap: %w", err)
		}
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
		vars := g.variables.WithArchitecture(arch)

		generator := NewDockerfileGenerator(g.config, g.templateEngine, vars)
		content, err := generator.Generate(arch)
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
	generator := NewComposeGenerator(g.config, g.templateEngine, g.variables)
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
	generator := NewMakefileGenerator(g.config, g.templateEngine, g.variables)
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
	generator := NewScriptsGenerator(g.config, g.templateEngine, g.variables)

	scripts := []struct {
		name   string
		method func() (string, error)
	}{
		{"bk-ci/tcs/build.sh", generator.GenerateBuildScript},
		{"bk-ci/tcs/deps_install.sh", generator.GenerateDepsInstallScript},
		{"bk-ci/tcs/rt_prepare.sh", generator.GenerateRtPrepareScript},
	}

	for _, script := range scripts {
		content, err := script.method()
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

// generateHooks generates hook scripts
func (g *Generator) generateHooks() error {
	hooksDir := filepath.Join(g.outputDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Generate healthcheck script
	healthcheckContent, err := g.generateHealthcheckScript()
	if err != nil {
		return fmt.Errorf("failed to generate healthcheck script: %w", err)
	}

	healthcheckPath := filepath.Join(hooksDir, "healthchk.sh")
	if err := utils.WriteFile(healthcheckPath, healthcheckContent); err != nil {
		return err
	}
	if err := os.Chmod(healthcheckPath, 0755); err != nil {
		return err
	}
	fmt.Println("✓ Generated hooks/healthchk.sh")

	// Generate start script
	startContent, err := g.generateStartScript()
	if err != nil {
		return fmt.Errorf("failed to generate start script: %w", err)
	}

	startPath := filepath.Join(hooksDir, "start.sh")
	if err := utils.WriteFile(startPath, startContent); err != nil {
		return err
	}
	if err := os.Chmod(startPath, 0755); err != nil {
		return err
	}
	fmt.Println("✓ Generated hooks/start.sh")

	return nil
}

// generateHealthcheckScript generates the health check script
func (g *Generator) generateHealthcheckScript() (string, error) {
	vars := g.variables.ToMap()

	var script string
	if g.config.Runtime.Healthcheck.Type == "custom" {
		script = g.config.Runtime.Healthcheck.CustomScript
	} else if g.config.Runtime.Healthcheck.Type == "http" {
		script = fmt.Sprintf(`#!/bin/sh
# Auto-generated health check script

curl -f http://localhost:%d%s || exit 1
`, g.config.Runtime.Healthcheck.HTTP.Port, g.config.Runtime.Healthcheck.HTTP.Path)
	} else {
		script = `#!/bin/sh
# Default health check
exit 0
`
	}

	return SubstituteVariables(script, vars), nil
}

// generateStartScript generates the startup script
func (g *Generator) generateStartScript() (string, error) {
	vars := g.variables.ToMap()
	return SubstituteVariables(g.config.Runtime.Startup.Command, vars), nil
}

// generateDevOps generates DevOps configuration
func (g *Generator) generateDevOps() error {
	generator := NewDevOpsGenerator(g.config, g.templateEngine, g.variables)
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

// generateConfigMap generates Kubernetes ConfigMap
func (g *Generator) generateConfigMap() error {
	generator := NewConfigMapGenerator(g.config, g.templateEngine, g.variables)
	content, err := generator.Generate()
	if err != nil {
		return err
	}

	k8sDir := filepath.Join(g.outputDir, g.config.LocalDev.Kubernetes.OutputDir)
	if err := os.MkdirAll(k8sDir, 0755); err != nil {
		return fmt.Errorf("failed to create k8s directory: %w", err)
	}

	outputPath := filepath.Join(k8sDir, "configmap.yaml")
	if err := utils.WriteFile(outputPath, content); err != nil {
		return err
	}

	fmt.Println("✓ Generated k8s-manifests/configmap.yaml")
	return nil
}
