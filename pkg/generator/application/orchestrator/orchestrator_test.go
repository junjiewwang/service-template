package orchestrator

import (
	"context"
	"testing"
)

func TestConfigProcessingOrchestrator(t *testing.T) {
	rawConfig := map[string]interface{}{
		"service": map[string]interface{}{
			"name":        "test-service",
			"description": "Test Service",
			"ports": []interface{}{
				map[string]interface{}{
					"name": "http",
					"port": 8080,
				},
			},
		},
		"language": map[string]interface{}{
			"type": "go",
			"config": map[string]interface{}{
				"goproxy": "https://goproxy.cn",
			},
		},
	}

	t.Run("Process with Registry", func(t *testing.T) {
		orchestrator := NewConfigProcessingOrchestrator()

		if err := orchestrator.Initialize(); err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}

		procCtx, err := orchestrator.Process(context.Background(), rawConfig)
		if err != nil {
			t.Fatalf("Process failed: %v", err)
		}

		// Verify service model
		if _, ok := procCtx.GetDomainModel("service"); !ok {
			t.Error("Service model not found")
		}

		// Verify language model
		if _, ok := procCtx.GetDomainModel("language"); !ok {
			t.Error("Language model not found")
		}

		// Verify generated files
		if _, ok := procCtx.GetGeneratedFile("service-metadata.txt"); !ok {
			t.Error("service-metadata.txt not generated")
		}

		if _, ok := procCtx.GetGeneratedFile("language-config.txt"); !ok {
			t.Error("language-config.txt not generated")
		}

		t.Log("✅ Process with Registry completed successfully")
	})

	t.Run("Process with PriorityChain", func(t *testing.T) {
		orchestrator := NewConfigProcessingOrchestrator()

		procCtx, err := orchestrator.ProcessWithPriorityChain(context.Background(), rawConfig)
		if err != nil {
			t.Fatalf("ProcessWithPriorityChain failed: %v", err)
		}

		if _, ok := procCtx.GetDomainModel("service"); !ok {
			t.Error("Service model not found")
		}

		if _, ok := procCtx.GetDomainModel("language"); !ok {
			t.Error("Language model not found")
		}

		t.Log("✅ Process with PriorityChain completed successfully")
	})

	t.Run("Process with DependencyGraph", func(t *testing.T) {
		orchestrator := NewConfigProcessingOrchestrator()

		procCtx, err := orchestrator.ProcessWithDependencyGraph(context.Background(), rawConfig)
		if err != nil {
			t.Fatalf("ProcessWithDependencyGraph failed: %v", err)
		}

		if _, ok := procCtx.GetDomainModel("service"); !ok {
			t.Error("Service model not found")
		}

		if _, ok := procCtx.GetDomainModel("language"); !ok {
			t.Error("Language model not found")
		}

		t.Log("✅ Process with DependencyGraph completed successfully")
	})
}

func TestConfigProcessingOrchestrator_ErrorHandling(t *testing.T) {
	// Invalid config (missing service name)
	invalidConfig := map[string]interface{}{
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

	orchestrator := NewConfigProcessingOrchestrator()
	if err := orchestrator.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	_, err := orchestrator.Process(context.Background(), invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid config")
	}

	t.Logf("✅ Error handling works correctly: %v", err)
}
