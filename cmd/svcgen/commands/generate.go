package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator"
	"github.com/junjiewwang/service-template/pkg/generator/adapter"
	"github.com/junjiewwang/service-template/pkg/generator/application/orchestrator"
	"github.com/junjiewwang/service-template/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	skipValidation bool
	useLegacy      bool // Flag to use legacy generator
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate project files from service.yaml",
	Long:  `Generates all project files including Dockerfiles, Compose, Makefile, and scripts based on service.yaml configuration.`,
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skip configuration validation")
	generateCmd.Flags().BoolVar(&useLegacy, "use-legacy", false, "Use legacy generator (for backward compatibility)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	fmt.Println("Loading configuration...")

	// Load configuration
	loader := config.NewLoader(configFile)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration unless skipped
	if !skipValidation {
		fmt.Println("Validating configuration...")
		validator := config.NewValidator(cfg)
		if err := validator.Validate(); err != nil {
			return err
		}
		fmt.Println("✓ Configuration is valid")
	}

	// Use legacy generator if flag is set
	if useLegacy {
		return runLegacyGenerate(cfg)
	}

	// Use new DDD-based orchestrator
	return runDDDGenerate(cfg)
}

// runLegacyGenerate uses the old generator (for backward compatibility)
func runLegacyGenerate(cfg *config.ServiceConfig) error {
	fmt.Println("\nGenerating project files (using legacy generator)...")
	gen := generator.NewGenerator(cfg, outputDir)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	fmt.Println("\n✓ All files generated successfully!")
	fmt.Printf("\nOutput directory: %s\n", outputDir)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the generated files")
	fmt.Println("  2. Run 'make docker-build' to build Docker images")
	fmt.Println("  3. Run 'make docker-up' to start services")

	return nil
}

// runDDDGenerate uses the new DDD-based orchestrator
func runDDDGenerate(cfg *config.ServiceConfig) error {
	fmt.Println("\nGenerating project files (using DDD orchestrator)...")

	// Convert ServiceConfig to raw config map
	rawConfig, err := convertConfigToMap(cfg)
	if err != nil {
		return fmt.Errorf("failed to convert config: %w", err)
	}

	// Create and initialize orchestrator
	orch := orchestrator.NewConfigProcessingOrchestrator()
	if err := orch.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize orchestrator: %w", err)
	}

	// Process configuration through DDD chains (Parse -> Validate)
	procCtx, err := orch.Process(context.Background(), rawConfig)
	if err != nil {
		return fmt.Errorf("orchestrator processing failed: %w", err)
	}

	// Check for validation errors
	if procCtx.HasValidationErrors() {
		fmt.Println("\n❌ Validation errors found:")
		for field, errs := range procCtx.GetAllValidationErrors() {
			for _, err := range errs {
				fmt.Printf("  - %s: %v\n", field, err)
			}
		}
		return fmt.Errorf("validation failed")
	}

	fmt.Println("✓ DDD validation passed")

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use legacy adapter to generate actual files with original config
	// (DDD validation passed, now use legacy generators for file generation)
	legacyAdapter := adapter.NewLegacyGeneratorAdapter(cfg, outputDir)

	// Generate Dockerfiles
	if err := generateDockerfiles(legacyAdapter, cfg); err != nil {
		return fmt.Errorf("failed to generate Dockerfiles: %w", err)
	}

	// Generate Docker Compose
	if err := generateCompose(legacyAdapter); err != nil {
		return fmt.Errorf("failed to generate compose.yaml: %w", err)
	}

	// Generate Makefile
	if err := generateMakefile(legacyAdapter); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}

	// Generate scripts
	if err := generateScripts(legacyAdapter, cfg); err != nil {
		return fmt.Errorf("failed to generate scripts: %w", err)
	}

	// Generate DevOps configuration
	if err := generateDevOps(legacyAdapter); err != nil {
		return fmt.Errorf("failed to generate DevOps configuration: %w", err)
	}

	fmt.Println("\n✓ All files generated successfully!")
	fmt.Printf("\nOutput directory: %s\n", outputDir)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the generated files")
	fmt.Println("  2. Run 'make docker-build' to build Docker images")
	fmt.Println("  3. Run 'make docker-up' to start services")

	return nil
}

// convertConfigToMap converts ServiceConfig to map for orchestrator
func convertConfigToMap(cfg *config.ServiceConfig) (map[string]interface{}, error) {
	// Marshal to YAML then unmarshal to map
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return rawConfig, nil
}

// generateDockerfiles generates Dockerfiles for different architectures
func generateDockerfiles(adapter *adapter.LegacyGeneratorAdapter, cfg *config.ServiceConfig) error {
	architectures := []string{"amd64", "arm64"}

	// Create .tad/build/{service-name} directory
	dockerfileDir := filepath.Join(outputDir, ".tad", "build", cfg.Service.Name)
	if err := os.MkdirAll(dockerfileDir, 0755); err != nil {
		return fmt.Errorf("failed to create dockerfile directory: %w", err)
	}

	for _, arch := range architectures {
		content, err := adapter.GenerateDockerfile(arch)
		if err != nil {
			return fmt.Errorf("failed to generate Dockerfile for %s: %w", arch, err)
		}

		filename := fmt.Sprintf("Dockerfile.%s.%s", cfg.Service.Name, arch)
		outputPath := filepath.Join(dockerfileDir, filename)
		if err := utils.WriteFile(outputPath, content); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		fmt.Printf("✓ Generated .tad/build/%s/%s\n", cfg.Service.Name, filename)
	}

	return nil
}

// generateCompose generates docker-compose.yaml
func generateCompose(adapter *adapter.LegacyGeneratorAdapter) error {
	content, err := adapter.GenerateCompose()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "compose.yaml")
	if err := utils.WriteFile(outputPath, content); err != nil {
		return err
	}

	fmt.Println("✓ Generated compose.yaml")
	return nil
}

// generateMakefile generates Makefile
func generateMakefile(adapter *adapter.LegacyGeneratorAdapter) error {
	content, err := adapter.GenerateMakefile()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "Makefile")
	if err := utils.WriteFile(outputPath, content); err != nil {
		return err
	}

	fmt.Println("✓ Generated Makefile")
	return nil
}

// generateScripts generates build and deployment scripts
func generateScripts(adapter *adapter.LegacyGeneratorAdapter, cfg *config.ServiceConfig) error {
	// Use correct path pattern: .tad/build/{service-name}/
	scriptDir := fmt.Sprintf(".tad/build/%s", cfg.Service.Name)

	scripts := map[string]func() (string, error){
		filepath.Join(scriptDir, "build.sh"):              adapter.GenerateBuildScript,
		filepath.Join(scriptDir, "build_deps_install.sh"): adapter.GenerateDepsInstallScript,
		filepath.Join(scriptDir, "rt_prepare.sh"):         adapter.GenerateRtPrepareScript,
		filepath.Join(scriptDir, "entrypoint.sh"):         adapter.GenerateEntrypointScript,
		filepath.Join(scriptDir, "healthchk.sh"):          adapter.GenerateHealthcheckScript,
	}

	// Add build_plugins.sh if plugins are configured
	if len(cfg.Plugins.Items) > 0 {
		scripts[filepath.Join(scriptDir, "build_plugins.sh")] = adapter.GenerateBuildPluginsScript
	}

	for scriptPath, generateFunc := range scripts {
		content, err := generateFunc()
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", scriptPath, err)
		}

		// Skip if generator returned empty content
		if content == "" {
			continue
		}

		outputPath := filepath.Join(outputDir, scriptPath)
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
func generateDevOps(adapter *adapter.LegacyGeneratorAdapter) error {
	content, err := adapter.GenerateDevOps()
	if err != nil {
		return err
	}

	tadDir := filepath.Join(outputDir, ".tad")
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
