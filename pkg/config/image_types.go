package config

import (
	"fmt"
	"strings"
)

// BaseImagesConfig 基础镜像配置（顶层配置）
// 仅在使用 @builders.* / @runtimes.* 预设引用时需要配置
type BaseImagesConfig struct {
	Builders map[string]ArchImageConfig `yaml:"builders,omitempty"` // 构建镜像预设
	Runtimes map[string]ArchImageConfig `yaml:"runtimes,omitempty"` // 运行时镜像预设
}

// IsEmpty 判断是否未配置任何预设
func (b *BaseImagesConfig) IsEmpty() bool {
	return len(b.Builders) == 0 && len(b.Runtimes) == 0
}

// Validate 验证基础镜像配置（仅在有内容时验证）
func (b *BaseImagesConfig) Validate() error {
	// 验证每个预设
	for name, img := range b.Builders {
		if err := img.Validate(); err != nil {
			return fmt.Errorf("base_images.builders.%s: %w", name, err)
		}
	}
	for name, img := range b.Runtimes {
		if err := img.Validate(); err != nil {
			return fmt.Errorf("base_images.runtimes.%s: %w", name, err)
		}
	}

	return nil
}

// GetBuilder 获取构建镜像预设
func (b *BaseImagesConfig) GetBuilder(name string) (ArchImageConfig, error) {
	img, ok := b.Builders[name]
	if !ok {
		available := make([]string, 0, len(b.Builders))
		for k := range b.Builders {
			available = append(available, k)
		}
		return ArchImageConfig{}, fmt.Errorf(
			"builder preset '%s' not found. Available: %v",
			name, available,
		)
	}
	return img, nil
}

// GetRuntime 获取运行时镜像预设
func (b *BaseImagesConfig) GetRuntime(name string) (ArchImageConfig, error) {
	img, ok := b.Runtimes[name]
	if !ok {
		available := make([]string, 0, len(b.Runtimes))
		for k := range b.Runtimes {
			available = append(available, k)
		}
		return ArchImageConfig{}, fmt.Errorf(
			"runtime preset '%s' not found. Available: %v",
			name, available,
		)
	}
	return img, nil
}

// ListBuilders 列出所有构建镜像预设名称
func (b *BaseImagesConfig) ListBuilders() []string {
	names := make([]string, 0, len(b.Builders))
	for name := range b.Builders {
		names = append(names, name)
	}
	return names
}

// ListRuntimes 列出所有运行时镜像预设名称
func (b *BaseImagesConfig) ListRuntimes() []string {
	names := make([]string, 0, len(b.Runtimes))
	for name := range b.Runtimes {
		names = append(names, name)
	}
	return names
}

// ============================================
// ImageSpec - 统一镜像规格类型
// ============================================

// ImageSpecKind 镜像规格的输入格式类型
type ImageSpecKind int

const (
	ImageSpecEmpty   ImageSpecKind = iota // 未指定，由 Resolver 推导
	ImageSpecDirect                       // 直接镜像名，如 "golang:1.23-alpine"
	ImageSpecPreset                       // 预设引用，如 "@builders.go_1.22"
	ImageSpecPerArch                      // 按架构指定，如 {amd64: "xxx", arm64: "yyy"}
)

// ImageSpec 镜像规格，支持多种输入格式
//
// 格式1: 不填（空）          → 由 Resolver 按 language.type 自动推导
// 格式2: 字符串（直接镜像名） → "golang:1.23-alpine"（multi-arch，同一地址填充两个架构）
// 格式3: 字符串（预设引用）   → "@builders.go_1.22"（从 base_images 查找）
// 格式4: 对象（按架构指定）   → {amd64: "xxx", arm64: "yyy"}
type ImageSpec struct {
	raw interface{} // string | ArchImageConfig | nil
}

// NewImageSpec 从字符串创建 ImageSpec（用于代码中构造）
func NewImageSpec(ref string) ImageSpec {
	return ImageSpec{raw: ref}
}

// NewImageSpecPerArch 从架构映射创建 ImageSpec
func NewImageSpecPerArch(amd64, arm64 string) ImageSpec {
	return ImageSpec{raw: ArchImageConfig{AMD64: amd64, ARM64: arm64}}
}

// Kind 返回镜像规格的输入格式类型
func (s *ImageSpec) Kind() ImageSpecKind {
	if s == nil || s.raw == nil {
		return ImageSpecEmpty
	}
	switch v := s.raw.(type) {
	case string:
		if v == "" {
			return ImageSpecEmpty
		}
		if strings.HasPrefix(v, "@") {
			return ImageSpecPreset
		}
		return ImageSpecDirect
	case ArchImageConfig:
		if v.IsEmpty() {
			return ImageSpecEmpty
		}
		return ImageSpecPerArch
	}
	return ImageSpecEmpty
}

// IsEmpty 判断是否未指定
func (s *ImageSpec) IsEmpty() bool {
	return s.Kind() == ImageSpecEmpty
}

// String 返回字符串表示
func (s *ImageSpec) String() string {
	if s == nil || s.raw == nil {
		return "<empty>"
	}
	switch v := s.raw.(type) {
	case string:
		return v
	case ArchImageConfig:
		return fmt.Sprintf("{amd64: %s, arm64: %s}", v.AMD64, v.ARM64)
	}
	return "<unknown>"
}

// Resolve 将 ImageSpec 统一解析为 ArchImageConfig
// baseImages 用于预设引用解析（当 Kind 为 ImageSpecPreset 时需要）
// expectedCategory 用于预设引用的类型校验（"builders" 或 "runtimes"）
func (s *ImageSpec) Resolve(baseImages *BaseImagesConfig, expectedCategory string) (ArchImageConfig, error) {
	switch s.Kind() {

	case ImageSpecDirect:
		// 直接镜像名（如 "golang:1.23-alpine"）
		// Docker Hub 等公开镜像本身是 multi-arch manifest，同一地址填充两个架构
		image := s.raw.(string)
		return ArchImageConfig{
			AMD64: image,
			ARM64: image,
		}, nil

	case ImageSpecPreset:
		// 预设引用（如 "@builders.go_1.22"）
		ref := s.raw.(string)
		category, name, err := parsePresetRef(ref)
		if err != nil {
			return ArchImageConfig{}, fmt.Errorf("invalid preset reference: %w", err)
		}
		if category != expectedCategory {
			return ArchImageConfig{}, fmt.Errorf(
				"image reference must use @%s.* (got: %s)", expectedCategory, ref,
			)
		}
		if baseImages == nil {
			return ArchImageConfig{}, fmt.Errorf(
				"base_images is required for preset reference: %s", ref,
			)
		}
		if category == "builders" {
			return baseImages.GetBuilder(name)
		}
		return baseImages.GetRuntime(name)

	case ImageSpecPerArch:
		// 按架构指定（如 {amd64: "xxx", arm64: "yyy"}）
		return s.raw.(ArchImageConfig), nil

	case ImageSpecEmpty:
		// 未指定，返回空值
		return ArchImageConfig{}, nil
	}

	return ArchImageConfig{}, fmt.Errorf("unknown image spec kind")
}

// Validate 验证 ImageSpec 格式合法性（不做解析，只检查格式）
func (s *ImageSpec) Validate(baseImages *BaseImagesConfig, expectedCategory string) error {
	switch s.Kind() {
	case ImageSpecEmpty:
		return nil
	case ImageSpecDirect:
		image := s.raw.(string)
		if strings.TrimSpace(image) == "" {
			return fmt.Errorf("image name cannot be empty")
		}
		return nil
	case ImageSpecPreset:
		ref := s.raw.(string)
		category, name, err := parsePresetRef(ref)
		if err != nil {
			return err
		}
		if category != expectedCategory {
			return fmt.Errorf("must reference @%s.* (got: %s)", expectedCategory, ref)
		}
		// 验证预设引用存在
		if baseImages == nil {
			return fmt.Errorf("base_images is required for preset reference: %s", ref)
		}
		if category == "builders" {
			if _, err := baseImages.GetBuilder(name); err != nil {
				return err
			}
		} else {
			if _, err := baseImages.GetRuntime(name); err != nil {
				return err
			}
		}
		return nil
	case ImageSpecPerArch:
		arch := s.raw.(ArchImageConfig)
		return arch.Validate()
	}
	return nil
}

// UnmarshalYAML implements custom YAML unmarshaling for ImageSpec
func (s *ImageSpec) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// 尝试解析为 string（直接镜像名 或 预设引用）
	var str string
	if err := unmarshal(&str); err == nil {
		s.raw = str
		return nil
	}

	// 尝试解析为 ArchImageConfig（按架构指定）
	var arch ArchImageConfig
	if err := unmarshal(&arch); err == nil {
		s.raw = arch
		return nil
	}

	return fmt.Errorf("image spec must be a string (e.g. \"golang:1.23-alpine\" or \"@builders.go_1.23\") " +
		"or an object with amd64/arm64 fields")
}

// MarshalYAML implements custom YAML marshaling for ImageSpec
func (s ImageSpec) MarshalYAML() (interface{}, error) {
	if s.raw == nil {
		return nil, nil
	}
	return s.raw, nil
}

// parsePresetRef 解析预设引用格式，返回 (category, name)
func parsePresetRef(ref string) (category string, name string, err error) {
	if !strings.HasPrefix(ref, "@") {
		return "", "", fmt.Errorf(
			"preset reference must start with '@' (got: %s)", ref,
		)
	}

	trimmed := strings.TrimPrefix(ref, "@")
	parts := strings.SplitN(trimmed, ".", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", "", fmt.Errorf(
			"invalid preset reference format: %s (expected: @category.name, e.g. @builders.go_1.21)", ref,
		)
	}

	category = parts[0]
	name = parts[1]

	if category != "builders" && category != "runtimes" {
		return "", "", fmt.Errorf(
			"invalid category in preset reference: %s (must be 'builders' or 'runtimes')", category,
		)
	}

	return category, name, nil
}

