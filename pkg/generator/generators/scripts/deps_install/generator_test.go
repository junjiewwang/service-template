package deps_install

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Language.Type = "go"
	cfg.Language.Config = map[string]interface{}{
		"goproxy": "https://goproxy.cn",
		"gosumdb": "sum.golang.org",
	}
	cfg.Build.Dependencies.SystemPkgs = []string{"git", "make"}

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
	if !strings.Contains(content, "#!/bin/bash") && !strings.Contains(content, "#!/bin/sh") {
		t.Error("Expected shebang not found")
	}
	if !strings.Contains(content, "git") {
		t.Error("Expected package git not found")
	}
	if !strings.Contains(content, "make") {
		t.Error("Expected package make not found")
	}
	if !strings.Contains(content, "https://goproxy.cn") {
		t.Error("Expected GOPROXY not found")
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

func TestGenerator_WithCustomDepsInstallCommand(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Language.Type = "python"
	cfg.Language.Config = map[string]interface{}{
		"deps_install_command": "pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/lib",
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

	// Verify custom command is applied with variable substitution
	if !strings.Contains(content, "pip install -r requirements.txt -t /opt/dist/lib") {
		t.Error("Expected custom deps_install_command with variable substitution not found")
	}
}

func TestGenerator_DefaultDepsInstallCommand(t *testing.T) {
	tests := []struct {
		name            string
		language        string
		expectedCommand string
	}{
		{
			name:            "go default command",
			language:        "go",
			expectedCommand: "go mod download",
		},
		{
			name:            "python default command",
			language:        "python",
			expectedCommand: "pip install -r requirements.txt",
		},
		{
			name:            "nodejs default command",
			language:        "nodejs",
			expectedCommand: "npm install",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testutil.NewTestConfig()
			cfg.Language.Type = tt.language

			ctx := context.NewGeneratorContext(cfg, "/tmp/output")
			gen, err := New(ctx)
			if err != nil {
				t.Fatalf("Failed to create generator: %v", err)
			}

			content, err := gen.Generate()
			if err != nil {
				t.Fatalf("Failed to generate: %v", err)
			}

			// Verify default command is used
			if !strings.Contains(content, tt.expectedCommand) {
				t.Errorf("Expected default command '%s' not found in generated content", tt.expectedCommand)
			}
		})
	}
}
