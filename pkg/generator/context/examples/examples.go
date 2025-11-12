package examples

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
)

// Example 1: Basic Usage - Using Variable Composer
func ExampleBasicUsage() {
	// Create a sample config
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "my-service",
			DeployDir: "/app",
		},
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
	}

	// Create context
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// Compose variables using fluent API
	vars := ctx.GetVariableComposer().
		WithCommon().
		WithBuild().
		Build()

	fmt.Printf("Service Name: %v\n", vars["SERVICE_NAME"])
	fmt.Printf("Total Variables: %d\n", len(vars))
}

// Example 2: Using Presets - Dockerfile Generation
func ExampleDockerfilePreset() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// Use preset for Dockerfile
	vars := ctx.GetVariablePreset().
		ForDockerfile("amd64").
		Build()

	fmt.Printf("Architecture: %v\n", vars["ARCH"])
	fmt.Printf("Builder Image: %v\n", vars["BUILDER_IMAGE"])
	fmt.Printf("Total Variables: %d\n", len(vars))
}

// Example 3: Custom Variables
func ExampleCustomVariables() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// Add custom variables
	vars := ctx.GetVariableComposer().
		WithCommon().
		WithCustom("MY_CUSTOM_VAR", "custom_value").
		WithCustom("ANOTHER_VAR", 123).
		WithCustomMap(map[string]interface{}{
			"VAR1": "value1",
			"VAR2": true,
		}).
		Build()

	fmt.Printf("Custom Var: %v\n", vars["MY_CUSTOM_VAR"])
	fmt.Printf("Another Var: %v\n", vars["ANOTHER_VAR"])
}

// Example 4: Override Variables
func ExampleOverrideVariables() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// Override existing variables
	vars := ctx.GetVariableComposer().
		WithCommon().
		Override("SERVICE_NAME", "overridden-service").
		Build()

	fmt.Printf("Service Name: %v\n", vars["SERVICE_NAME"])
}

// Example 5: Checking Variables
func ExampleCheckingVariables() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	composer := ctx.GetVariableComposer().WithCommon()

	// Check if variable exists
	if composer.Has("SERVICE_NAME") {
		if val, ok := composer.Get("SERVICE_NAME"); ok {
			fmt.Printf("Found SERVICE_NAME: %v\n", val)
		}
	}

	// Check non-existent variable
	if !composer.Has("NON_EXISTENT") {
		fmt.Println("NON_EXISTENT variable not found")
	}
}

// Example 6: Cloning Composer
func ExampleCloningComposer() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// Create original composer
	original := ctx.GetVariableComposer().WithCommon()

	// Clone and modify
	cloned := original.Clone()
	cloned.WithCustom("NEW_VAR", "new_value")

	fmt.Printf("Original has NEW_VAR: %v\n", original.Has("NEW_VAR"))
	fmt.Printf("Cloned has NEW_VAR: %v\n", cloned.Has("NEW_VAR"))
}

// Example 7: Real Generator Implementation
type ExampleGenerator struct {
	ctx *context.GeneratorContext
}

func (g *ExampleGenerator) PrepareVariables() map[string]interface{} {
	// Old way (before flyweight pattern):
	// vars := make(map[string]interface{})
	// vars["SERVICE_NAME"] = g.ctx.Config.Service.Name
	// vars["DEPLOY_DIR"] = g.ctx.Config.Service.DeployDir
	// ... many more lines

	// New way (with flyweight pattern):
	return g.ctx.GetVariableComposer().
		WithCommon().
		WithBuild().
		WithRuntime().
		WithCustom("GENERATOR_SPECIFIC_VAR", "value").
		Build()
}

// Example 8: Using All Categories
func ExampleAllCategories() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// Get all standard variables
	vars := ctx.GetVariableComposer().
		WithAll().
		Build()

	fmt.Printf("Total Variables: %d\n", len(vars))

	// Print some key variables
	fmt.Printf("SERVICE_NAME: %v\n", vars["SERVICE_NAME"])
	fmt.Printf("BUILD_COMMAND: %v\n", vars["BUILD_COMMAND"])
	fmt.Printf("PLUGIN_ROOT_DIR: %v\n", vars["PLUGIN_ROOT_DIR"])
}

// Example 9: Architecture-Specific Variables
func ExampleArchitectureVariables() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	// AMD64
	varsAMD64 := ctx.GetVariableComposer().
		WithBuild().
		WithArchitecture("amd64").
		Build()

	fmt.Printf("AMD64 - GOARCH: %v\n", varsAMD64["GOARCH"])
	fmt.Printf("AMD64 - BUILDER_IMAGE: %v\n", varsAMD64["BUILDER_IMAGE"])

	// ARM64
	varsARM64 := ctx.GetVariableComposer().
		WithBuild().
		WithArchitecture("arm64").
		Build()

	fmt.Printf("ARM64 - GOARCH: %v\n", varsARM64["GOARCH"])
	fmt.Printf("ARM64 - BUILDER_IMAGE: %v\n", varsARM64["BUILDER_IMAGE"])
}

// Example 10: Preset Comparison
func ExamplePresetComparison() {
	cfg := createSampleConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	preset := ctx.GetVariablePreset()

	// Different presets for different scenarios
	dockerfileVars := preset.ForDockerfile("amd64").Build()
	buildScriptVars := preset.ForBuildScript().Build()
	composeVars := preset.ForCompose().Build()

	fmt.Printf("Dockerfile preset: %d variables\n", len(dockerfileVars))
	fmt.Printf("Build script preset: %d variables\n", len(buildScriptVars))
	fmt.Printf("Compose preset: %d variables\n", len(composeVars))
}

// Helper function to create sample config
func createSampleConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "example-service",
			DeployDir: "/app",
			Ports: []config.PortConfig{
				{Port: 8080, Expose: true},
			},
		},
		Language: config.LanguageConfig{
			Type:    "go",
			Version: "1.21",
		},
		Build: config.BuildConfig{
			Commands: config.BuildCommandsConfig{
				Build:     "go build -o bin/app",
				PreBuild:  "go mod download",
				PostBuild: "echo done",
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/app",
			},
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "http",
			},
			GenerateScripts: true,
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/opt/plugins",
		},
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
	}
}
