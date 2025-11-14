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
