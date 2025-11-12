package makefile

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.LocalDev.Kubernetes.Enabled = true
	cfg.LocalDev.Kubernetes.Namespace = "default"
	cfg.Makefile.CustomTargets = []config.CustomTarget{
		{
			Name:        "custom-test",
			Description: "Run custom tests",
			Commands:    []string{"echo 'Running tests'"},
		},
	}

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
	if !strings.Contains(content, ".PHONY:") {
		t.Error("Expected .PHONY directive not found")
	}
	if !strings.Contains(content, "build:") {
		t.Error("Expected build target not found")
	}
	if !strings.Contains(content, "k8s-deploy:") {
		t.Error("Expected k8s-deploy target not found")
	}
	if !strings.Contains(content, "custom-test:") {
		t.Error("Expected custom target not found")
	}
}

func TestGenerator_GetName(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if gen.GetName() != GeneratorType {
		t.Errorf("Expected name %s, got %s", GeneratorType, gen.GetName())
	}
}

func TestGenerator_Validate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if err := gen.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}
