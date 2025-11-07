package commands

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate service.yaml configuration",
	Long:  `Validates the service.yaml configuration file for correctness and completeness.`,
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Load configuration
	loader := config.NewLoader(configFile)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	validator := config.NewValidator(cfg)
	if err := validator.Validate(); err != nil {
		return err
	}

	fmt.Println("✓ Configuration is valid")
	fmt.Printf("\nService: %s\n", cfg.Service.Name)
	fmt.Printf("Language: %s %s\n", cfg.Language.Type, cfg.Language.Version)
	fmt.Printf("Ports: %d configured\n", len(cfg.Service.Ports))
	if len(cfg.Plugins.Items) > 0 {
		fmt.Printf("Plugins: %d configured\n", len(cfg.Plugins.Items))
	}

	return nil
}
