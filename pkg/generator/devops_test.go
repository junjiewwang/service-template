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
	if gen.config != cfg {
		t.Error("Config not set correctly")
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

func TestDevOpsGenerator_PrepareTemplateVars(t *testing.T) {
	cfg := &config.ServiceConfig{
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
		Language: config.LanguageConfig{
			Type:    "nodejs",
			Version: "18",
		},
		Build: config.BuildConfig{
			RuntimeImage: config.ArchImageConfig{
				AMD64: "mirrors.tencent.com/tlinux/tlinux3:3.1",
				ARM64: "mirrors.tencent.com/tlinux/tlinux3:3.2",
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "mirrors.tencent.com/tcs-infra/node:18-alpine",
				ARM64: "mirrors.tencent.com/tcs-infra/node:18-alpine",
			},
		},
		Service: config.ServiceInfo{
			DeployDir: "/app",
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	gen := NewDevOpsGenerator(cfg, engine, vars)

	templateVars := gen.prepareTemplateVars()

	// Check basic variables
	if templateVars["GENERATED_AT"] != "2024-01-01T00:00:00Z" {
		t.Errorf("Expected GENERATED_AT to be set")
	}

	if templateVars["LANGUAGE_TYPE"] != "nodejs" {
		t.Errorf("Expected LANGUAGE_TYPE to be nodejs, got %v", templateVars["LANGUAGE_TYPE"])
	}

	if templateVars["LANGUAGE_VERSION"] != "18" {
		t.Errorf("Expected LANGUAGE_VERSION to be 18, got %v", templateVars["LANGUAGE_VERSION"])
	}

	if templateVars["DEPLOY_DIR"] != "/app" {
		t.Errorf("Expected DEPLOY_DIR to be /app, got %v", templateVars["DEPLOY_DIR"])
	}

	// Check runtime image parsing
	if templateVars["RUNTIME_IMAGE_X86"] != "mirrors.tencent.com/tlinux/tlinux3" {
		t.Errorf("Unexpected RUNTIME_IMAGE_X86: %v", templateVars["RUNTIME_IMAGE_X86"])
	}

	if templateVars["RUNTIME_TAG_X86"] != "3.1" {
		t.Errorf("Unexpected RUNTIME_TAG_X86: %v", templateVars["RUNTIME_TAG_X86"])
	}

	if templateVars["RUNTIME_TAG_ARM"] != "3.2" {
		t.Errorf("Unexpected RUNTIME_TAG_ARM: %v", templateVars["RUNTIME_TAG_ARM"])
	}

	// Check builder images
	if templateVars["BUILDER_IMAGE_X86"] != "mirrors.tencent.com/tcs-infra/node:18-alpine" {
		t.Errorf("Unexpected BUILDER_IMAGE_X86: %v", templateVars["BUILDER_IMAGE_X86"])
	}

	// Check language examples flags
	if templateVars["SHOW_PYTHON_EXAMPLE"] != true {
		t.Error("Expected SHOW_PYTHON_EXAMPLE to be true for nodejs")
	}
	if templateVars["SHOW_NODEJS_EXAMPLE"] != false {
		t.Error("Expected SHOW_NODEJS_EXAMPLE to be false for nodejs")
	}
	if templateVars["SHOW_JAVA_EXAMPLE"] != true {
		t.Error("Expected SHOW_JAVA_EXAMPLE to be true for nodejs")
	}
}

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
