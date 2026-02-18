package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBuildCommand(t *testing.T) {
	tests := []struct {
		name     string
		langType string
		langCfg  *LanguageConfig
		contains string // 检查命令中包含的关键字
	}{
		{"go", "go", &LanguageConfig{Type: "go"}, "go build"},
		{"python", "python", &LanguageConfig{Type: "python"}, "cp -r"},
		{"java maven", "java", &LanguageConfig{Type: "java"}, "mvn package"},
		{"java gradle", "java", &LanguageConfig{
			Type:   "java",
			Config: map[string]interface{}{"build_tool": "gradle"},
		}, "gradle build"},
		{"nodejs", "nodejs", &LanguageConfig{Type: "nodejs"}, "npm run build"},
		{"rust", "rust", &LanguageConfig{Type: "rust"}, "cargo build"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := DefaultBuildCommand(tt.langType, tt.langCfg)
			assert.NotEmpty(t, cmd)
			assert.Contains(t, cmd, tt.contains)
		})
	}
}

func TestDefaultBuildCommand_Unsupported(t *testing.T) {
	cmd := DefaultBuildCommand("cobol", &LanguageConfig{Type: "cobol"})
	assert.Empty(t, cmd)
}

func TestHasDefaultBuildCommand(t *testing.T) {
	supported := []string{"go", "python", "java", "nodejs", "rust"}
	for _, lang := range supported {
		assert.True(t, HasDefaultBuildCommand(lang), "expected %s to have default build command", lang)
	}

	assert.False(t, HasDefaultBuildCommand("cobol"))
	assert.False(t, HasDefaultBuildCommand(""))
}

func TestResolveBuildCommand(t *testing.T) {
	t.Run("explicit config takes priority", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build: BuildConfig{
				Commands: BuildCommandsConfig{Build: "custom-build-cmd"},
			},
		}
		assert.Equal(t, "custom-build-cmd", ResolveBuildCommand(cfg))
	})

	t.Run("falls back to language default", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "go"},
			Build:    BuildConfig{},
		}
		cmd := ResolveBuildCommand(cfg)
		assert.Contains(t, cmd, "go build")
	})

	t.Run("unsupported language returns empty", func(t *testing.T) {
		cfg := &ServiceConfig{
			Language: LanguageConfig{Type: "cobol"},
			Build:    BuildConfig{},
		}
		assert.Empty(t, ResolveBuildCommand(cfg))
	})
}
