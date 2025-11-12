package entrypoint

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Runtime.Startup.Command = "./bin/myapp"
	cfg.Runtime.Startup.Env = []config.EnvConfig{
		{Name: "ENV", Value: "production"},
	}
	cfg.Plugins.InstallDir = "/opt/plugins"
	cfg.Plugins.Items = []config.PluginConfig{
		{
			Name: "test-plugin",
			RuntimeEnv: []config.EnvironmentVariable{
				{Name: "PLUGIN_PATH", Value: "${PLUGIN_INSTALL_DIR}/bin"},
			},
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
	if !strings.Contains(content, "#!/bin/sh") {
		t.Error("Expected shebang not found")
	}
	if !strings.Contains(content, "./bin/myapp") {
		t.Error("Expected startup command not found")
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
