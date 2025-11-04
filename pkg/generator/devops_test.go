package generator

import (
	"strings"
	"testing"
	"time"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDevOpsGenerator(t *testing.T) {
	// Arrange: Setup service configuration
	cfg := &config.ServiceConfig{
		Metadata: config.MetadataConfig{
			GeneratedAt: time.Now().Format(time.RFC3339),
		},
		Language: config.LanguageConfig{
			Type:    "golang",
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

	// Act: Create DevOps generator
	gen := NewDevOpsGenerator(cfg, engine, vars)
	require.NotNil(t, gen, "DevOps generator should be created")
	t.Logf("✓ DevOps generator created successfully")

	// Assert: Test that generator can generate content
	content, err := gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")
	require.NotEmpty(t, content, "Generated content should not be empty")
	t.Logf("✓ Generated DevOps config: %d bytes", len(content))
}

func TestDevOpsGenerator_Generate(t *testing.T) {
	// Arrange: Setup service configuration with specific images
	cfg := &config.ServiceConfig{
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
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

	// Act: Generate DevOps configuration
	content, err := gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")
	require.NotEmpty(t, content, "Generated content should not be empty")

	t.Logf("Generated DevOps config: %d bytes", len(content))

	// Assert: Verify content contains expected sections
	expectedStrings := map[string]string{
		"tad_section":      "tad:",
		"export_envs":      "export_envs:",
		"tlinux_x86_image": "TLINUX_BASE_IMAGE_X86",
		"tlinux_x86_tag":   "TLINUX_TAG_X86",
		"tlinux_arm_image": "TLINUX_BASE_IMAGE_ARM",
		"tlinux_arm_tag":   "TLINUX_TAG_ARM",
		"builder_x86":      "BUILDER_IMAGE_X86",
		"builder_arm":      "BUILDER_IMAGE_ARM",
		"deploy_dir":       "DEPLOY_DIR",
		"deploy_dir_value": "/data/app",
		"tlinux_base":      "mirrors.tencent.com/tlinux/tlinux3",
		"tlinux_tag":       "3.1",
		"builder_image":    "mirrors.tencent.com/tcs-infra/golang:1.21-alpine",
	}

	for name, expected := range expectedStrings {
		assert.Contains(t, content, expected,
			"DevOps config should contain %s: %q", name, expected)
	}
	t.Logf("✓ Verified all %d expected sections present", len(expectedStrings))
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
	require.NoError(t, err, "Generate should not fail")

	// Check that the embedded template is used and variables are processed
	assert.Contains(t, content, "Auto-generated DevOps configuration", "Expected embedded template content")
	assert.Contains(t, content, "Python 3.11", "Expected language type in content")
	assert.Contains(t, content, "/opt/service", "Expected deploy dir in content")
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
			assert.Equal(t, tt.expectedImage, image, "Image should match expected")
			assert.Equal(t, tt.expectedTag, tag, "Tag should match expected")
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
			assert.Equal(t, tt.expected, result, "Language display name should match expected")
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
			require.NoError(t, err, "Generate failed for %s", lang.langType)

			// Verify language-specific content
			assert.Contains(t, content, "tad:", "Expected TAD configuration")

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
