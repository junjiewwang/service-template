package devops

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Build.BuilderImage.AMD64 = "golang:1.21"
	cfg.Build.BuilderImage.ARM64 = "golang:1.21"
	cfg.Build.RuntimeImage.AMD64 = "alpine:3.18"
	cfg.Build.RuntimeImage.ARM64 = "alpine:3.18"
	cfg.Language.Type = "go"
	cfg.Language.Version = "1.21"

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
	if !strings.Contains(content, "tad:") {
		t.Error("Expected tad section not found")
	}
	if !strings.Contains(content, "BUILDER_IMAGE_X86") {
		t.Error("Expected BUILDER_IMAGE_X86 not found")
	}
	if !strings.Contains(content, "BUILDER_IMAGE_ARM") {
		t.Error("Expected BUILDER_IMAGE_ARM not found")
	}
	if !strings.Contains(content, "golang:1.21") {
		t.Error("Expected builder image not found")
	}
	if !strings.Contains(content, "alpine") {
		t.Error("Expected runtime image not found")
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

func TestParseImageAndTag(t *testing.T) {
	tests := []struct {
		name      string
		fullImage string
		wantImage string
		wantTag   string
	}{
		{"with tag", "alpine:3.18", "alpine", "3.18"},
		{"without tag", "alpine", "alpine", "latest"},
		{"complex tag", "golang:1.21-alpine", "golang", "1.21-alpine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImage, gotTag := parseImageAndTag(tt.fullImage)
			if gotImage != tt.wantImage {
				t.Errorf("parseImageAndTag(%s) image = %s, want %s", tt.fullImage, gotImage, tt.wantImage)
			}
			if gotTag != tt.wantTag {
				t.Errorf("parseImageAndTag(%s) tag = %s, want %s", tt.fullImage, gotTag, tt.wantTag)
			}
		})
	}
}

func TestGetLanguageDisplayName(t *testing.T) {
	tests := []struct {
		langType string
		version  string
		want     string
	}{
		{"go", "1.21", "Go 1.21"},
		{"python", "3.11", "Python 3.11"},
		{"nodejs", "18", "Node.js 18"},
		{"java", "17", "Java 17"},
		{"unknown", "1.0", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.langType, func(t *testing.T) {
			got := getLanguageDisplayName(tt.langType, tt.version)
			if got != tt.want {
				t.Errorf("getLanguageDisplayName(%s, %s) = %s, want %s", tt.langType, tt.version, got, tt.want)
			}
		})
	}
}
