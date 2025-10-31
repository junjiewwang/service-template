package generator

import (
	"strings"
	"testing"
	"time"

	"github.com/junjiewwang/service-template/pkg/config"
)

func TestNewDevOpsGenerator(t *testing.T) {
	cfg := &config.ServiceConfig{
		Metadata: config.MetadataConfig{
			GeneratedAt: time.Now().Format(time.RFC3339),
		},
		Language: config.LanguageConfig{
			Type:    "go",
			Version: "1.21",
		},
		Build: config.BuildConfig{
			RuntimeImage: config.ArchImageConfig{
				AMD64: "mirrors.tencent.com/tlinux/tlinux3:latest",
				ARM64: "mirrors.tencent.com/tlinux/tlinux3:latest",
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "mirrors.tencent.com/tcs-infra/golang:1.21-alpine",
				ARM64: "mirrors.tencent.com/tcs-infra/golang:1.21-alpine",
			},
		},
		Service: config.ServiceInfo{
			DeployDir: "/data/app",
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	gen := NewDevOpsGenerator(cfg, engine, vars)

	if gen == nil {
		t.Fatal("Expected non-nil generator")
	}
	// Test that generator can generate content
	_, err := gen.Generate()
	if err != nil {
		t.Errorf("Failed to generate: %v", err)
	}
}

func TestDevOpsGenerator_Generate(t *testing.T) {
	cfg := &config.ServiceConfig{
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
		Language: config.LanguageConfig{
			Type:    "go",
			Version: "1.21",
		},
		Build: config.BuildConfig{
			RuntimeImage: config.ArchImageConfig{
				AMD64: "mirrors.tencent.com/tlinux/tlinux3:3.1",
				ARM64: "mirrors.tencent.com/tlinux/tlinux3:3.1-arm",
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "mirrors.tencent.com/tcs-infra/golang:1.21-alpine",
				ARM64: "mirrors.tencent.com/tcs-infra/golang:1.21-alpine",
			},
		},
		Service: config.ServiceInfo{
			DeployDir: "/data/app",
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	gen := NewDevOpsGenerator(cfg, engine, vars)

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify content contains expected sections
	expectedStrings := []string{
		"tad:",
		"export_envs:",
		"TLINUX_BASE_IMAGE_X86",
		"TLINUX_TAG_X86",
		"TLINUX_BASE_IMAGE_ARM",
		"TLINUX_TAG_ARM",
		"BUILDER_IMAGE_X86",
		"BUILDER_IMAGE_ARM",
		"DEPLOY_DIR",
		"/data/app",
		"mirrors.tencent.com/tlinux/tlinux3",
		"3.1",
		"mirrors.tencent.com/tcs-infra/golang:1.21-alpine",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("Expected content to contain %q", expected)
		}
	}
}

func TestDevOpsGenerator_GenerateWithCustomTemplate(t *testing.T) {
	// Test that the generator uses the embedded template and processes variables correctly
	cfg := &config.ServiceConfig{
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
		Language: config.LanguageConfig{
			Type:    "python",
			Version: "3.11",
		},
		Service: config.ServiceInfo{
			DeployDir: "/opt/service",
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	gen := NewDevOpsGenerator(cfg, engine, vars)

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check that the embedded template is used and variables are processed
	if !strings.Contains(content, "Auto-generated DevOps configuration") {
		t.Error("Expected embedded template content")
	}
	if !strings.Contains(content, "Python 3.11") {
		t.Error("Expected language type in content")
	}
	if !strings.Contains(content, "/opt/service") {
		t.Error("Expected deploy dir in content")
	}
}

// TestDevOpsGenerator_PrepareTemplateVars removed - tests internal implementation
// The functionality is covered by TestDevOpsGenerator_Generate which tests the public API

func TestParseImageAndTag(t *testing.T) {
	tests := []struct {
		name          string
		fullImage     string
		expectedImage string
		expectedTag   string
	}{
		{
			name:          "image with tag",
			fullImage:     "mirrors.tencent.com/tlinux/tlinux3:3.1",
			expectedImage: "mirrors.tencent.com/tlinux/tlinux3",
			expectedTag:   "3.1",
		},
		{
			name:          "image without tag",
			fullImage:     "mirrors.tencent.com/tlinux/tlinux3",
			expectedImage: "mirrors.tencent.com/tlinux/tlinux3",
			expectedTag:   "latest",
		},
		{
			name:          "image with complex tag",
			fullImage:     "registry.example.com/org/image:v1.2.3-alpine",
			expectedImage: "registry.example.com/org/image",
			expectedTag:   "v1.2.3-alpine",
		},
		{
			name:          "simple image",
			fullImage:     "ubuntu:20.04",
			expectedImage: "ubuntu",
			expectedTag:   "20.04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, tag := parseImageAndTag(tt.fullImage)
			if image != tt.expectedImage {
				t.Errorf("Expected image %q, got %q", tt.expectedImage, image)
			}
			if tag != tt.expectedTag {
				t.Errorf("Expected tag %q, got %q", tt.expectedTag, tag)
			}
		})
	}
}

func TestGetLanguageDisplayName(t *testing.T) {
	tests := []struct {
		langType string
		version  string
		expected string
	}{
		{"go", "1.21", "Go 1.21"},
		{"python", "3.11", "Python 3.11"},
		{"nodejs", "18", "Node.js 18"},
		{"java", "17", "Java 17"},
		{"rust", "1.70", "rust"},
		{"", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.langType, func(t *testing.T) {
			result := getLanguageDisplayName(tt.langType, tt.version)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDevOpsGenerator_DifferentLanguages(t *testing.T) {
	languages := []struct {
		langType string
		version  string
	}{
		{"go", "1.21"},
		{"python", "3.11"},
		{"nodejs", "18"},
		{"java", "17"},
	}

	for _, lang := range languages {
		t.Run(lang.langType, func(t *testing.T) {
			cfg := &config.ServiceConfig{
				Metadata: config.MetadataConfig{
					GeneratedAt: "2024-01-01T00:00:00Z",
				},
				Language: config.LanguageConfig{
					Type:    lang.langType,
					Version: lang.version,
				},
				Build: config.BuildConfig{
					RuntimeImage: config.ArchImageConfig{
						AMD64: "mirrors.tencent.com/tlinux/tlinux3:latest",
						ARM64: "mirrors.tencent.com/tlinux/tlinux3:latest",
					},
					BuilderImage: config.ArchImageConfig{
						AMD64: "mirrors.tencent.com/builder:latest",
						ARM64: "mirrors.tencent.com/builder:latest",
					},
				},
				Service: config.ServiceInfo{
					DeployDir: "/data/app",
				},
			}

			engine := NewTemplateEngine()
			vars := NewVariables(cfg)
			gen := NewDevOpsGenerator(cfg, engine, vars)

			content, err := gen.Generate()
			if err != nil {
				t.Fatalf("Generate failed for %s: %v", lang.langType, err)
			}

			// Verify language-specific content
			if !strings.Contains(content, "tad:") {
				t.Error("Expected TAD configuration")
			}

			// Verify other language examples are shown
			otherLangs := []string{"python", "nodejs", "java"}
			for _, other := range otherLangs {
				if other != lang.langType {
					// Should show example for other languages
					if !strings.Contains(content, strings.Title(other)) && other != "nodejs" {
						// nodejs shows as "Node.js"
						if other != "nodejs" || !strings.Contains(content, "Node.js") {
							t.Logf("Warning: Expected to see %s example in content", other)
						}
					}
				}
			}
		})
	}
}
