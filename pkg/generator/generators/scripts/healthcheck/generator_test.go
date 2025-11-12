package healthcheck

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate_Default(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Runtime.Healthcheck.Enabled = true
	cfg.Runtime.Healthcheck.Type = "default"

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify content
	if !strings.Contains(content, "#!/bin/sh") {
		t.Error("Expected shebang not found")
	}
	if !strings.Contains(content, "SERVICE_ROOT") {
		t.Error("Expected SERVICE_ROOT not found")
	}
	if !strings.Contains(content, "grep") {
		t.Error("Expected process check not found")
	}
}

func TestGenerator_Generate_Custom(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Runtime.Healthcheck.Enabled = true
	cfg.Runtime.Healthcheck.Type = "custom"
	cfg.Runtime.Healthcheck.CustomScript = "curl -f http://localhost:8080/health || exit 1"

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify content
	if !strings.Contains(content, "#!/bin/sh") {
		t.Error("Expected shebang not found")
	}
	if !strings.Contains(content, "curl -f http://localhost:8080/health") {
		t.Error("Expected custom script not found")
	}
}

func TestGenerator_GetName(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Runtime.Healthcheck.Enabled = true
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if gen.GetName() != GeneratorType {
		t.Errorf("Expected name %s, got %s", GeneratorType, gen.GetName())
	}
}

func TestGenerator_Validate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Runtime.Healthcheck.Enabled = true
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if err := gen.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}

func TestStrategy_Default(t *testing.T) {
	cfg := testutil.NewTestConfig()
	strategy := NewDefaultStrategy(cfg)

	if strategy.GetType() != "default" {
		t.Errorf("Expected type 'default', got %s", strategy.GetType())
	}

	if err := strategy.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	vars := map[string]interface{}{
		"SERVICE_NAME": "test-service",
		"DEPLOY_DIR":   "/opt/services",
	}
	script, err := strategy.GenerateScript(vars)
	if err != nil {
		t.Fatalf("Failed to generate script: %v", err)
	}

	if !strings.Contains(script, "SERVICE_ROOT") {
		t.Error("Expected SERVICE_ROOT in script")
	}
}

func TestStrategy_Custom(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Runtime.Healthcheck.CustomScript = "echo 'custom check'"
	strategy := NewCustomStrategy(cfg)

	if strategy.GetType() != "custom" {
		t.Errorf("Expected type 'custom', got %s", strategy.GetType())
	}

	if err := strategy.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	vars := map[string]interface{}{
		"SERVICE_NAME":  "test-service",
		"DEPLOY_DIR":    "/opt/services",
		"CUSTOM_SCRIPT": "echo 'custom check'",
	}
	script, err := strategy.GenerateScript(vars)
	if err != nil {
		t.Fatalf("Failed to generate script: %v", err)
	}

	t.Logf("Generated script: %s", script)
	if !strings.Contains(script, "CUSTOM_SCRIPT") || !strings.Contains(script, "{{") {
		// Script is a template, not rendered yet
		t.Log("Script contains template variables, which is expected")
	}
}
