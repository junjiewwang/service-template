package generators

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"

	// Import all generators to register them
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/build_tools/makefile"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/compose"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/devops"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/docker/dockerfile"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/build"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/deps_install"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/entrypoint"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/healthcheck"
	_ "github.com/junjiewwang/service-template/pkg/generator/generators/scripts/rt_prepare"
)

func TestAllGeneratorsRegistered(t *testing.T) {
	expectedGenerators := []string{
		"rt-prepare-script",
		"deps-install-script",
		"entrypoint-script",
		"build-script",
		"healthcheck-script",
		"makefile",
		"compose",
		"devops",
		"dockerfile",
	}

	registered := core.DefaultRegistry.GetAll()
	registeredMap := make(map[string]bool)
	for _, name := range registered {
		registeredMap[name] = true
	}

	for _, expected := range expectedGenerators {
		if !registeredMap[expected] {
			t.Errorf("Generator %s is not registered", expected)
		}
	}

	t.Logf("Total registered generators: %d", len(registered))
	t.Logf("Registered generators: %v", registered)
}

func TestAllGeneratorsCanCreate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	tests := []struct {
		name    string
		genType string
		options []interface{}
	}{
		{"rt-prepare", "rt-prepare-script", nil},
		{"deps-install", "deps-install-script", nil},
		{"entrypoint", "entrypoint-script", nil},
		{"build", "build-script", nil},
		{"healthcheck", "healthcheck-script", nil},
		{"makefile", "makefile", nil},
		{"compose", "compose", nil},
		{"devops", "devops", nil},
		{"dockerfile-amd64", "dockerfile", []interface{}{"amd64"}},
		{"dockerfile-arm64", "dockerfile", []interface{}{"arm64"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator, exists := core.DefaultRegistry.Get(tt.genType)
			if !exists {
				t.Fatalf("Generator %s not found in registry", tt.genType)
			}

			gen, err := creator(ctx, tt.options...)
			if err != nil {
				t.Fatalf("Failed to create generator %s: %v", tt.genType, err)
			}

			if gen == nil {
				t.Fatalf("Generator %s is nil", tt.genType)
			}

			// Test Generate
			content, err := gen.Generate()
			if err != nil {
				t.Fatalf("Failed to generate content for %s: %v", tt.genType, err)
			}

			if content == "" {
				t.Errorf("Generated content for %s is empty", tt.genType)
			}

			t.Logf("Generator %s created successfully, content length: %d", tt.name, len(content))
		})
	}
}

func TestGeneratorValidation(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	creator, _ := core.DefaultRegistry.Get("rt-prepare-script")
	gen, _ := creator(ctx)

	if err := gen.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}
