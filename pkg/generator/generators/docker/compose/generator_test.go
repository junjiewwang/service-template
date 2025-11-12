package compose

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Service.Ports = []config.PortConfig{
		{Port: 8080, Protocol: "tcp"},
		{Port: 9090, Protocol: "tcp"},
	}
	cfg.Runtime.Startup.Env = []config.EnvConfig{
		{Name: "ENV", Value: "production"},
	}
	cfg.LocalDev.Compose.Volumes = []config.VolumeConfig{
		{Source: "./data", Target: "${SERVICE_ROOT}/data"},
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
	if !strings.Contains(content, "services:") {
		t.Error("Expected services field not found")
	}
	if !strings.Contains(content, "8080") {
		t.Error("Expected port 8080 not found")
	}
	if !strings.Contains(content, "9090") {
		t.Error("Expected port 9090 not found")
	}
	if !strings.Contains(content, "ENV=production") {
		t.Error("Expected environment variable not found")
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
