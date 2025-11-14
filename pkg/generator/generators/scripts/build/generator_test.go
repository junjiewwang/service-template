package build

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Build.Commands.Build = "go build -o bin/app"
	cfg.Build.Commands.PreBuild = "echo 'pre-build'"
	cfg.Build.Commands.PostBuild = "echo 'post-build'"
	cfg.Plugins.InstallDir = "/opt/plugins"
	cfg.Plugins.Items = []config.PluginConfig{
		{
			Name:           "test-plugin",
			DownloadURL:    config.NewStaticDownloadURL("https://example.com/plugin.tar.gz"),
			InstallCommand: "tar -xzf plugin.tar.gz -C ${PLUGIN_INSTALL_DIR}",
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
	if !strings.Contains(content, "#!/bin/bash") && !strings.Contains(content, "#!/bin/sh") {
		t.Error("Expected shebang not found")
	}
	if !strings.Contains(content, "go build -o bin/app") {
		t.Error("Expected build command not found")
	}
	if !strings.Contains(content, "pre-build") {
		t.Error("Expected pre-build command not found")
	}
	if !strings.Contains(content, "post-build") {
		t.Error("Expected post-build command not found")
	}
	// Plugin logic has been moved to build_plugins.sh, so build.sh should NOT contain plugin names
	if strings.Contains(content, "test-plugin") {
		t.Error("build.sh should not contain plugin logic (moved to build_plugins.sh)")
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
