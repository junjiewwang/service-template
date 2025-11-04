package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	outputDir  string
)

var rootCmd = &cobra.Command{
	Use:   "svcgen",
	Short: "Service Template Generator",
	Long: `Service Template Generator - A configuration-driven tool to generate 
service templates with Docker, Kubernetes, and CI/CD configurations.

All templates are embedded in the binary, no external template files needed.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "service.yaml", "Path to service.yaml configuration file")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", cwd, "Output directory for generated files")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("svcgen version 2.0.0")
	},
}
