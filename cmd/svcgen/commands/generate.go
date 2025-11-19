package commands

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	skipValidation bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate project files from service.yaml",
	Long:  `Generates all project files including Dockerfiles, Compose, Makefile, and scripts based on service.yaml configuration.`,
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skip configuration validation")
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

	// Generate project files
	fmt.Println("\nGenerating project files...")
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
