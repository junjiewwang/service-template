package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// ============================================
// ImageSpec Kind / IsEmpty / String
// ============================================

func TestImageSpec_Kind(t *testing.T) {
	tests := []struct {
		name     string
		spec     ImageSpec
		expected ImageSpecKind
	}{
		{"zero value", ImageSpec{}, ImageSpecEmpty},
		{"nil raw", ImageSpec{raw: nil}, ImageSpecEmpty},
		{"empty string", NewImageSpec(""), ImageSpecEmpty},
		{"direct image", NewImageSpec("golang:1.23-alpine"), ImageSpecDirect},
		{"preset builder", NewImageSpec("@builders.go_1.22"), ImageSpecPreset},
		{"preset runtime", NewImageSpec("@runtimes.alpine"), ImageSpecPreset},
		{"per-arch", NewImageSpecPerArch("img:amd64", "img:arm64"), ImageSpecPerArch},
		{"per-arch both empty", NewImageSpecPerArch("", ""), ImageSpecEmpty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.spec.Kind())
		})
	}
}

func TestImageSpec_IsEmpty(t *testing.T) {
	empty1 := ImageSpec{}
	empty2 := NewImageSpec("")
	empty3 := NewImageSpecPerArch("", "")
	notEmpty1 := NewImageSpec("golang:1.23")
	notEmpty2 := NewImageSpec("@builders.go")
	notEmpty3 := NewImageSpecPerArch("a", "b")

	assert.True(t, empty1.IsEmpty())
	assert.True(t, empty2.IsEmpty())
	assert.True(t, empty3.IsEmpty())
	assert.False(t, notEmpty1.IsEmpty())
	assert.False(t, notEmpty2.IsEmpty())
	assert.False(t, notEmpty3.IsEmpty())
}

func TestImageSpec_String(t *testing.T) {
	tests := []struct {
		name     string
		spec     ImageSpec
		expected string
	}{
		{"empty", ImageSpec{}, "<empty>"},
		{"direct", NewImageSpec("golang:1.23"), "golang:1.23"},
		{"preset", NewImageSpec("@builders.go_1.22"), "@builders.go_1.22"},
		{"per-arch", NewImageSpecPerArch("a:amd64", "a:arm64"), "{amd64: a:amd64, arm64: a:arm64}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.spec.String())
		})
	}
}

// ============================================
// ImageSpec Resolve
// ============================================

func TestImageSpec_Resolve_Direct(t *testing.T) {
	spec := NewImageSpec("golang:1.23-alpine")
	result, err := spec.Resolve(nil, "builders")
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-alpine", result.AMD64)
	assert.Equal(t, "golang:1.23-alpine", result.ARM64)
}

func TestImageSpec_Resolve_Preset(t *testing.T) {
	baseImages := &BaseImagesConfig{
		Builders: map[string]ArchImageConfig{
			"go_1.22": {AMD64: "golang:1.22-amd64", ARM64: "golang:1.22-arm64"},
		},
		Runtimes: map[string]ArchImageConfig{
			"alpine": {AMD64: "alpine:amd64", ARM64: "alpine:arm64"},
		},
	}

	t.Run("builder preset", func(t *testing.T) {
		spec := NewImageSpec("@builders.go_1.22")
		result, err := spec.Resolve(baseImages, "builders")
		require.NoError(t, err)
		assert.Equal(t, "golang:1.22-amd64", result.AMD64)
		assert.Equal(t, "golang:1.22-arm64", result.ARM64)
	})

	t.Run("runtime preset", func(t *testing.T) {
		spec := NewImageSpec("@runtimes.alpine")
		result, err := spec.Resolve(baseImages, "runtimes")
		require.NoError(t, err)
		assert.Equal(t, "alpine:amd64", result.AMD64)
		assert.Equal(t, "alpine:arm64", result.ARM64)
	})

	t.Run("category mismatch", func(t *testing.T) {
		spec := NewImageSpec("@builders.go_1.22")
		_, err := spec.Resolve(baseImages, "runtimes")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must use @runtimes.*")
	})

	t.Run("preset not found", func(t *testing.T) {
		spec := NewImageSpec("@builders.nonexistent")
		_, err := spec.Resolve(baseImages, "builders")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("nil baseImages", func(t *testing.T) {
		spec := NewImageSpec("@builders.go_1.22")
		_, err := spec.Resolve(nil, "builders")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "base_images is required")
	})
}

func TestImageSpec_Resolve_PerArch(t *testing.T) {
	spec := NewImageSpecPerArch("custom:amd64", "custom:arm64")
	result, err := spec.Resolve(nil, "builders")
	require.NoError(t, err)
	assert.Equal(t, "custom:amd64", result.AMD64)
	assert.Equal(t, "custom:arm64", result.ARM64)
}

func TestImageSpec_Resolve_Empty(t *testing.T) {
	spec := ImageSpec{}
	result, err := spec.Resolve(nil, "builders")
	require.NoError(t, err)
	assert.True(t, result.IsEmpty())
}

// ============================================
// ImageSpec Validate
// ============================================

func TestImageSpec_Validate(t *testing.T) {
	baseImages := &BaseImagesConfig{
		Builders: map[string]ArchImageConfig{
			"go_1.22": {AMD64: "golang:amd64", ARM64: "golang:arm64"},
		},
		Runtimes: map[string]ArchImageConfig{
			"alpine": {AMD64: "alpine:amd64", ARM64: "alpine:arm64"},
		},
	}

	tests := []struct {
		name     string
		spec     ImageSpec
		base     *BaseImagesConfig
		category string
		wantErr  bool
		errMsg   string
	}{
		{"empty is valid", ImageSpec{}, nil, "builders", false, ""},
		{"direct is valid", NewImageSpec("golang:1.23"), nil, "builders", false, ""},
		{"preset found", NewImageSpec("@builders.go_1.22"), baseImages, "builders", false, ""},
		{"preset category mismatch", NewImageSpec("@builders.go_1.22"), baseImages, "runtimes", true, "must reference @runtimes.*"},
		{"preset not found", NewImageSpec("@builders.nonexistent"), baseImages, "builders", true, "not found"},
		{"preset nil baseImages", NewImageSpec("@builders.go_1.22"), nil, "builders", true, "base_images is required"},
		{"per-arch valid", NewImageSpecPerArch("a:amd64", "a:arm64"), nil, "builders", false, ""},
		{"per-arch missing amd64", NewImageSpecPerArch("", "a:arm64"), nil, "builders", true, "amd64 image is required"},
		{"per-arch missing arm64", NewImageSpecPerArch("a:amd64", ""), nil, "builders", true, "arm64 image is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate(tt.base, tt.category)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ============================================
// ImageSpec YAML Marshal / Unmarshal
// ============================================

func TestImageSpec_UnmarshalYAML_DirectString(t *testing.T) {
	yamlContent := `builder_image: "golang:1.23-alpine"`

	var cfg struct {
		BuilderImage ImageSpec `yaml:"builder_image"`
	}
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)
	require.NoError(t, err)
	assert.Equal(t, ImageSpecDirect, cfg.BuilderImage.Kind())
	assert.Equal(t, "golang:1.23-alpine", cfg.BuilderImage.String())
}

func TestImageSpec_UnmarshalYAML_PresetRef(t *testing.T) {
	yamlContent := `builder_image: "@builders.go_1.22"`

	var cfg struct {
		BuilderImage ImageSpec `yaml:"builder_image"`
	}
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)
	require.NoError(t, err)
	assert.Equal(t, ImageSpecPreset, cfg.BuilderImage.Kind())
	assert.Equal(t, "@builders.go_1.22", cfg.BuilderImage.String())
}

func TestImageSpec_UnmarshalYAML_PerArch(t *testing.T) {
	yamlContent := `
builder_image:
  amd64: "golang:1.23-amd64"
  arm64: "golang:1.23-arm64"
`

	var cfg struct {
		BuilderImage ImageSpec `yaml:"builder_image"`
	}
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)
	require.NoError(t, err)
	assert.Equal(t, ImageSpecPerArch, cfg.BuilderImage.Kind())

	result, err := cfg.BuilderImage.Resolve(nil, "builders")
	require.NoError(t, err)
	assert.Equal(t, "golang:1.23-amd64", result.AMD64)
	assert.Equal(t, "golang:1.23-arm64", result.ARM64)
}

func TestImageSpec_UnmarshalYAML_NotSpecified(t *testing.T) {
	// 未指定字段时，ImageSpec 保持零值
	yamlContent := `other_field: "value"`

	var cfg struct {
		BuilderImage ImageSpec `yaml:"builder_image,omitempty"`
		OtherField   string    `yaml:"other_field"`
	}
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)
	require.NoError(t, err)
	assert.True(t, cfg.BuilderImage.IsEmpty())
	assert.Equal(t, ImageSpecEmpty, cfg.BuilderImage.Kind())
}

func TestImageSpec_MarshalYAML(t *testing.T) {
	t.Run("direct string roundtrip", func(t *testing.T) {
		original := struct {
			Image ImageSpec `yaml:"image"`
		}{Image: NewImageSpec("golang:1.23")}

		data, err := yaml.Marshal(&original)
		require.NoError(t, err)

		var result struct {
			Image ImageSpec `yaml:"image"`
		}
		err = yaml.Unmarshal(data, &result)
		require.NoError(t, err)
		assert.Equal(t, "golang:1.23", result.Image.String())
	})

	t.Run("per-arch roundtrip", func(t *testing.T) {
		original := struct {
			Image ImageSpec `yaml:"image"`
		}{Image: NewImageSpecPerArch("a:amd64", "a:arm64")}

		data, err := yaml.Marshal(&original)
		require.NoError(t, err)

		var result struct {
			Image ImageSpec `yaml:"image"`
		}
		err = yaml.Unmarshal(data, &result)
		require.NoError(t, err)
		assert.Equal(t, ImageSpecPerArch, result.Image.Kind())
		resolved, err := result.Image.Resolve(nil, "builders")
		require.NoError(t, err)
		assert.Equal(t, "a:amd64", resolved.AMD64)
		assert.Equal(t, "a:arm64", resolved.ARM64)
	})

	t.Run("empty marshals to nil", func(t *testing.T) {
		original := struct {
			Image ImageSpec `yaml:"image,omitempty"`
		}{Image: ImageSpec{}}

		data, err := yaml.Marshal(&original)
		require.NoError(t, err)
		assert.Equal(t, "{}\n", string(data))
	})
}

// ============================================
// parsePresetRef
// ============================================

func TestParsePresetRef(t *testing.T) {
	tests := []struct {
		name         string
		ref          string
		wantCategory string
		wantName     string
		wantErr      bool
		errMsg       string
	}{
		{"valid builder", "@builders.go_1.22", "builders", "go_1.22", false, ""},
		{"valid runtime", "@runtimes.alpine_3.19", "runtimes", "alpine_3.19", false, ""},
		{"no @ prefix", "builders.go_1.22", "", "", true, "must start with '@'"},
		{"no dot separator", "@buildersgo", "", "", true, "invalid preset reference format"},
		{"empty name", "@builders.", "", "", true, "invalid preset reference format"},
		{"invalid category", "@images.go_1.22", "", "", true, "invalid category"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, name, err := parsePresetRef(tt.ref)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCategory, category)
				assert.Equal(t, tt.wantName, name)
			}
		})
	}
}

// ============================================
// BaseImagesConfig
// ============================================

func TestBaseImagesConfig_IsEmpty(t *testing.T) {
	assert.True(t, (&BaseImagesConfig{}).IsEmpty())
	assert.False(t, (&BaseImagesConfig{
		Builders: map[string]ArchImageConfig{"go": {AMD64: "a", ARM64: "b"}},
	}).IsEmpty())
}

func TestBaseImagesConfig_GetBuilder(t *testing.T) {
	base := &BaseImagesConfig{
		Builders: map[string]ArchImageConfig{
			"go_1.22": {AMD64: "go:amd64", ARM64: "go:arm64"},
		},
	}

	t.Run("found", func(t *testing.T) {
		img, err := base.GetBuilder("go_1.22")
		require.NoError(t, err)
		assert.Equal(t, "go:amd64", img.AMD64)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := base.GetBuilder("nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestBaseImagesConfig_GetRuntime(t *testing.T) {
	base := &BaseImagesConfig{
		Runtimes: map[string]ArchImageConfig{
			"alpine": {AMD64: "alpine:amd64", ARM64: "alpine:arm64"},
		},
	}

	t.Run("found", func(t *testing.T) {
		img, err := base.GetRuntime("alpine")
		require.NoError(t, err)
		assert.Equal(t, "alpine:arm64", img.ARM64)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := base.GetRuntime("nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
