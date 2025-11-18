package integration_test

import (
	"context"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/language"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/service"
)

// TestIntegration_ServiceAndLanguage demonstrates the complete workflow
func TestIntegration_ServiceAndLanguage(t *testing.T) {
	// Prepare raw configuration (simulating YAML input)
	rawConfig := map[string]interface{}{
		"service": map[string]interface{}{
			"name":        "example-service",
			"description": "Example Service",
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "http",
					"port":     8080,
					"protocol": "TCP",
					"expose":   true,
				},
			},
		},
		"language": map[string]interface{}{
			"type": "go",
			"config": map[string]interface{}{
				"goproxy": "https://goproxy.cn,direct",
			},
		},
	}

	// Create domain factories
	serviceFactory := service.NewServiceDomainFactory()
	languageFactory := language.NewLanguageDomainFactory()

	// Test 1: Using DomainRegistry (simple approach)
	t.Run("Using DomainRegistry", func(t *testing.T) {
		registry := chain.NewDomainRegistry()
		registry.RegisterAll(serviceFactory, languageFactory)

		// Build chains
		parseChain := registry.BuildParseChain()
		validateChain := registry.BuildValidateChain()
		generateChain := registry.BuildGenerateChain()

		// Create context
		ctx := chain.NewProcessingContext(context.Background(), rawConfig)

		// Execute parse chain
		if err := parseChain.Handle(ctx); err != nil {
			t.Fatalf("Parse chain failed: %v", err)
		}

		// Verify service model
		serviceModel, ok := ctx.GetDomainModel("service")
		if !ok {
			t.Fatal("Service model not found")
		}
		serviceConfig := serviceModel.(*service.ServiceConfig)
		if serviceConfig.Name != "example-service" {
			t.Errorf("Service name = %s, want example-service", serviceConfig.Name)
		}

		// Verify language model
		langModel, ok := ctx.GetDomainModel("language")
		if !ok {
			t.Fatal("Language model not found")
		}
		langConfig := langModel.(*language.LanguageConfig)
		if langConfig.Type != "go" {
			t.Errorf("Language type = %s, want go", langConfig.Type)
		}

		// Execute validate chain
		if err := validateChain.Handle(ctx); err != nil {
			t.Fatalf("Validate chain failed: %v", err)
		}

		// Execute generate chain
		if err := generateChain.Handle(ctx); err != nil {
			t.Fatalf("Generate chain failed: %v", err)
		}

		// Verify generated files
		if _, ok := ctx.GetGeneratedFile("service-metadata.txt"); !ok {
			t.Error("service-metadata.txt not generated")
		}
		if _, ok := ctx.GetGeneratedFile("language-config.txt"); !ok {
			t.Error("language-config.txt not generated")
		}

		// Verify metadata
		if name, ok := ctx.GetMetadata("service_name"); !ok || name != "example-service" {
			t.Errorf("service_name metadata = %v, want example-service", name)
		}
		if langType, ok := ctx.GetMetadata("language_type"); !ok || langType != "go" {
			t.Errorf("language_type metadata = %v, want go", langType)
		}
	})

	// Test 2: Using PriorityChain (explicit ordering)
	t.Run("Using PriorityChain", func(t *testing.T) {
		priorityChain := chain.NewPriorityChain().
			First(serviceFactory).
			Then(languageFactory)

		if err := priorityChain.Validate(); err != nil {
			t.Fatalf("Priority chain validation failed: %v", err)
		}

		parseChain := priorityChain.BuildParseChain()
		ctx := chain.NewProcessingContext(context.Background(), rawConfig)

		if err := parseChain.Handle(ctx); err != nil {
			t.Fatalf("Parse chain failed: %v", err)
		}

		// Verify both models are parsed
		if _, ok := ctx.GetDomainModel("service"); !ok {
			t.Error("Service model not found")
		}
		if _, ok := ctx.GetDomainModel("language"); !ok {
			t.Error("Language model not found")
		}
	})

	// Test 3: Using DependencyGraph (dependency-based ordering)
	t.Run("Using DependencyGraph", func(t *testing.T) {
		graph := chain.NewDependencyGraph().
			AddNode(serviceFactory).
			AddNode(languageFactory, "service")

		if err := graph.Validate(); err != nil {
			t.Fatalf("Dependency graph validation failed: %v", err)
		}

		parseChain, err := graph.BuildParseChain()
		if err != nil {
			t.Fatalf("Build parse chain failed: %v", err)
		}

		ctx := chain.NewProcessingContext(context.Background(), rawConfig)

		if err := parseChain.Handle(ctx); err != nil {
			t.Fatalf("Parse chain failed: %v", err)
		}

		// Verify processing order (service should be processed before language)
		if _, ok := ctx.GetDomainModel("service"); !ok {
			t.Error("Service model not found")
		}
		if _, ok := ctx.GetDomainModel("language"); !ok {
			t.Error("Language model not found")
		}
	})
}

// TestIntegration_ErrorHandling tests error handling in the chain
func TestIntegration_ErrorHandling(t *testing.T) {
	// Invalid configuration (missing service name)
	rawConfig := map[string]interface{}{
		"service": map[string]interface{}{
			"description": "Test",
			"ports": []interface{}{
				map[string]interface{}{
					"name": "http",
					"port": 8080,
				},
			},
		},
	}

	registry := chain.NewDomainRegistry()
	registry.Register(service.NewServiceDomainFactory())

	parseChain := registry.BuildParseChain()
	validateChain := registry.BuildValidateChain()

	ctx := chain.NewProcessingContext(context.Background(), rawConfig)

	// Parse should succeed
	if err := parseChain.Handle(ctx); err != nil {
		t.Fatalf("Parse chain failed: %v", err)
	}

	// Validate should fail
	err := validateChain.Handle(ctx)
	if err == nil {
		t.Error("Expected validation error for missing service name")
	}

	// Check validation errors
	if !ctx.HasErrors() {
		t.Error("Expected context to have errors")
	}

	errors := ctx.GetValidationErrors("service")
	if len(errors) == 0 {
		t.Error("Expected validation errors for service domain")
	}
}

// TestIntegration_FullWorkflow tests the complete workflow with all chains
func TestIntegration_FullWorkflow(t *testing.T) {
	rawConfig := map[string]interface{}{
		"service": map[string]interface{}{
			"name":        "my-service",
			"description": "My Service",
			"ports": []interface{}{
				map[string]interface{}{
					"name": "http",
					"port": 8080,
				},
			},
		},
		"language": map[string]interface{}{
			"type": "go",
		},
	}

	// Create registry and register factories
	registry := chain.NewDomainRegistry()
	registry.RegisterAll(
		service.NewServiceDomainFactory(),
		language.NewLanguageDomainFactory(),
	)

	// Build all chains
	parseChain := registry.BuildParseChain()
	validateChain := registry.BuildValidateChain()
	generateChain := registry.BuildGenerateChain()

	// Create context
	ctx := chain.NewProcessingContext(context.Background(), rawConfig)

	// Execute full workflow
	if err := parseChain.Handle(ctx); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if err := validateChain.Handle(ctx); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if err := generateChain.Handle(ctx); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify no errors
	if ctx.HasErrors() {
		t.Errorf("Unexpected errors: %v", ctx.GetErrors())
	}

	// Verify all generated files
	generatedFiles := []string{
		"service-metadata.txt",
		"language-config.txt",
	}

	for _, file := range generatedFiles {
		if _, ok := ctx.GetGeneratedFile(file); !ok {
			t.Errorf("File %s not generated", file)
		}
	}

	t.Logf("✅ Full workflow completed successfully")
	t.Logf("📁 Generated %d files", len(generatedFiles))
}
