package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================
// DefaultBuilderImage / DefaultRuntimeImage
// ============================================

func TestDefaultBuilderImage(t *testing.T) {
	tests := []struct {
		name     string
		langType string
		langCfg  *LanguageConfig
		expected string
	}{
		{"go default", "go", &LanguageConfig{Type: "go"}, "golang:1.23-alpine"},
		{"go custom version", "go", &LanguageConfig{
			Type:   "go",
			Config: map[string]interface{}{"go_version": "1.21"},
		}, "golang:1.21-alpine"},
		{"python default", "python", &LanguageConfig{Type: "python"}, "python:3.12-slim"},
		{"python custom version", "python", &LanguageConfig{
			Type:   "python",
			Config: map[string]interface{}{"python_version": "3.11"},
		}, "python:3.11-slim"},
		{"java maven default", "java", &LanguageConfig{Type: "java"}, "maven:3-eclipse-temurin-21"},
		{"java gradle", "java", &LanguageConfig{
			Type:   "java",
			Config: map[string]interface{}{"build_tool": "gradle", "gradle_version": "8", "jdk_version": "17"},
		}, "gradle:8-jdk17"},
		{"nodejs default", "nodejs", &LanguageConfig{Type: "nodejs"}, "node:20-alpine"},
		{"nodejs custom version", "nodejs", &LanguageConfig{
			Type:   "nodejs",
			Config: map[string]interface{}{"node_version": "18"},
		}, "node:18-alpine"},
		{"rust default", "rust", &LanguageConfig{Type: "rust"}, "rust:1.78-alpine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DefaultBuilderImage(tt.langType, tt.langCfg)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultBuilderImage_UnsupportedLanguage(t *testing.T) {
	_, err := DefaultBuilderImage("cobol", &LanguageConfig{Type: "cobol"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no default builder image")
}

func TestDefaultRuntimeImage(t *testing.T) {
	tests := []struct {
		name     string
		langType string
		langCfg  *LanguageConfig
		expected string
	}{
		{"go runtime", "go", &LanguageConfig{Type: "go"}, "alpine:3.19"},
		{"python runtime", "python", &LanguageConfig{Type: "python"}, "python:3.12-slim"},
		{"java runtime", "java", &LanguageConfig{Type: "java"}, "eclipse-temurin:21-jre-alpine"},
		{"java custom jdk", "java", &LanguageConfig{
			Type:   "java",
			Config: map[string]interface{}{"jdk_version": "17"},
		}, "eclipse-temurin:17-jre-alpine"},
		{"nodejs runtime", "nodejs", &LanguageConfig{Type: "nodejs"}, "node:20-alpine"},
		{"rust runtime", "rust", &LanguageConfig{Type: "rust"}, "alpine:3.19"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DefaultRuntimeImage(tt.langType, tt.langCfg)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultRuntimeImage_UnsupportedLanguage(t *testing.T) {
	_, err := DefaultRuntimeImage("cobol", &LanguageConfig{Type: "cobol"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no default runtime image")
}

// ============================================
// HasDefaultImages
// ============================================

func TestHasDefaultImages(t *testing.T) {
	supported := []string{"go", "python", "java", "nodejs", "rust"}
	for _, lang := range supported {
		assert.True(t, HasDefaultImages(lang), "expected %s to have default images", lang)
	}

	unsupported := []string{"cobol", "ruby", "php", ""}
	for _, lang := range unsupported {
		assert.False(t, HasDefaultImages(lang), "expected %s NOT to have default images", lang)
	}
}

// ============================================
// ResolveBuilderImageWithDefaults / ResolveRuntimeImageWithDefaults
// ============================================

func TestResolveBuilderImageWithDefaults(t *testing.T) {
	t.Run("explicit direct image takes priority", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				BuilderImage: NewImageSpec("custom-golang:1.20"),
			},
		}
		result, err := ResolveBuilderImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "custom-golang:1.20", result.AMD64)
		assert.Equal(t, "custom-golang:1.20", result.ARM64)
	})

	t.Run("explicit preset reference", func(t *testing.T) {
		cfg := &ServiceConfig{
			BaseImages: BaseImagesConfig{
				Builders: map[string]ArchImageConfig{
					"go_1.22": {AMD64: "go:amd64", ARM64: "go:arm64"},
				},
			},
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				BuilderImage: NewImageSpec("@builders.go_1.22"),
			},
		}
		result, err := ResolveBuilderImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "go:amd64", result.AMD64)
		assert.Equal(t, "go:arm64", result.ARM64)
	})

	t.Run("empty falls back to language default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build:    BuildConfig{},
		}
		result, err := ResolveBuilderImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "golang:1.23-alpine", result.AMD64)
		assert.Equal(t, "golang:1.23-alpine", result.ARM64)
	})

	t.Run("empty with custom go version", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{
				Type:   "go",
				Config: map[string]interface{}{"go_version": "1.21"},
			},
			Build: BuildConfig{},
		}
		result, err := ResolveBuilderImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "golang:1.21-alpine", result.AMD64)
	})

	t.Run("unsupported language no default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "cobol"},
			Build:    BuildConfig{},
		}
		_, err := ResolveBuilderImageWithDefaults(cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no default builder image")
	})
}

func TestResolveRuntimeImageWithDefaults(t *testing.T) {
	t.Run("explicit direct image takes priority", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				RuntimeImage: NewImageSpec("custom-alpine:3.18"),
			},
		}
		result, err := ResolveRuntimeImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "custom-alpine:3.18", result.AMD64)
		assert.Equal(t, "custom-alpine:3.18", result.ARM64)
	})

	t.Run("empty falls back to language default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "python"},
			Build:    BuildConfig{},
		}
		result, err := ResolveRuntimeImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "python:3.12-slim", result.AMD64)
	})

	t.Run("per-arch explicit", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				RuntimeImage: NewImageSpecPerArch("custom:amd64", "custom:arm64"),
			},
		}
		result, err := ResolveRuntimeImageWithDefaults(cfg)
		require.NoError(t, err)
		assert.Equal(t, "custom:amd64", result.AMD64)
		assert.Equal(t, "custom:arm64", result.ARM64)
	})
}
