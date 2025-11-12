package dockerfile

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate_AMD64(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Build.BuilderImage.AMD64 = "golang:1.21-alpine"
	cfg.Build.RuntimeImage.AMD64 = "alpine:3.18"

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx, "amd64")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify content
	if len(content) < 100 {
		t.Errorf("Content too short: %d bytes", len(content))
	}
	if !strings.Contains(content, "FROM") {
		t.Error("Expected FROM statement not found")
	}
}

func TestGenerator_Generate_ARM64(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Build.BuilderImage.ARM64 = "golang:1.21-alpine"
	cfg.Build.RuntimeImage.ARM64 = "alpine:3.18"

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx, "arm64")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify content
	if len(content) < 100 {
		t.Errorf("Content too short: %d bytes", len(content))
	}
	if !strings.Contains(content, "FROM") {
		t.Error("Expected FROM statement not found")
	}
}

func TestGenerator_InvalidArch(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	_, err := New(ctx, "invalid")
	if err == nil {
		t.Error("Expected error for invalid architecture")
	}
}

func TestDetectPackageManager(t *testing.T) {
	tests := []struct {
		image string
		want  string
	}{
		{"alpine:latest", "apk"},
		{"debian:bullseye", "apt-get"},
		{"ubuntu:22.04", "apt-get"},
		{"centos:7", "yum"},
		{"tencentos:3", "yum"},
		{"fedora:38", "dnf"},
		{"unknown:latest", "yum"},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			got := detectPackageManager(tt.image)
			if got != tt.want {
				t.Errorf("detectPackageManager(%s) = %s, want %s", tt.image, got, tt.want)
			}
		})
	}
}

func TestGetDefaultDependencyFiles(t *testing.T) {
	tests := []struct {
		language string
		want     []string
	}{
		{"go", []string{"go.mod", "go.sum"}},
		{"python", []string{"requirements.txt"}},
		{"nodejs", []string{"package.json", "package-lock.json"}},
		{"java", []string{"pom.xml"}},
		{"unknown", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			got := getDefaultDependencyFiles(tt.language)
			if len(got) != len(tt.want) {
				t.Errorf("getDefaultDependencyFiles(%s) length = %d, want %d", tt.language, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getDefaultDependencyFiles(%s)[%d] = %s, want %s", tt.language, i, got[i], tt.want[i])
				}
			}
		})
	}
}
