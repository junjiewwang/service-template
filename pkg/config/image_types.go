package config

import (
	"fmt"
	"strings"
)

// BaseImagesConfig 基础镜像配置（顶层配置）
type BaseImagesConfig struct {
	Builders map[string]ArchImageConfig `yaml:"builders"` // 构建镜像预设
	Runtimes map[string]ArchImageConfig `yaml:"runtimes"` // 运行时镜像预设
}

// Validate 验证基础镜像配置
func (b *BaseImagesConfig) Validate() error {
	if len(b.Builders) == 0 {
		return fmt.Errorf("base_images.builders cannot be empty")
	}
	if len(b.Runtimes) == 0 {
		return fmt.Errorf("base_images.runtimes cannot be empty")
	}

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

// ImageRef 镜像引用类型
type ImageRef string

// Validate 验证引用格式
func (r ImageRef) Validate() error {
	if r == "" {
		return fmt.Errorf("image reference cannot be empty")
	}

	ref := string(r)
	if !strings.HasPrefix(ref, "@") {
		return fmt.Errorf(
			"image reference must start with '@' (got: %s). "+
				"Example: @builders.go_1.21 or @runtimes.alpine_3.18",
			ref,
		)
	}

	// 移除 @ 前缀
	ref = strings.TrimPrefix(ref, "@")

	// 验证格式：category.name
	parts := strings.SplitN(ref, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf(
			"invalid reference format: @%s (expected: @category.name). "+
				"Example: @builders.go_1.21",
			ref,
		)
	}

	category := parts[0]
	name := parts[1]

	if category != "builders" && category != "runtimes" {
		return fmt.Errorf(
			"invalid category: %s (must be 'builders' or 'runtimes')",
			category,
		)
	}

	if name == "" {
		return fmt.Errorf("preset name cannot be empty")
	}

	return nil
}

// Parse 解析引用，返回 (category, name)
func (r ImageRef) Parse() (category string, name string, err error) {
	if err := r.Validate(); err != nil {
		return "", "", err
	}

	ref := strings.TrimPrefix(string(r), "@")
	parts := strings.SplitN(ref, ".", 2)
	return parts[0], parts[1], nil
}

// String 返回字符串表示
func (r ImageRef) String() string {
	return string(r)
}

// IsBuilder 判断是否引用构建镜像
func (r ImageRef) IsBuilder() bool {
	category, _, err := r.Parse()
	return err == nil && category == "builders"
}

// IsRuntime 判断是否引用运行时镜像
func (r ImageRef) IsRuntime() bool {
	category, _, err := r.Parse()
	return err == nil && category == "runtimes"
}

// UnmarshalYAML implements custom YAML unmarshaling for ImageRef
func (r *ImageRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return fmt.Errorf("image reference must be a string: %w", err)
	}
	*r = ImageRef(str)
	return nil
}

// MarshalYAML implements custom YAML marshaling for ImageRef
func (r ImageRef) MarshalYAML() (interface{}, error) {
	return string(r), nil
}

// ImageResolver 镜像解析器
type ImageResolver struct {
	cfg *ServiceConfig
}

// NewImageResolver 创建解析器
func NewImageResolver(cfg *ServiceConfig) *ImageResolver {
	return &ImageResolver{cfg: cfg}
}

// ResolveBuilderImage 解析构建镜像
func (r *ImageResolver) ResolveBuilderImage() (ArchImageConfig, error) {
	category, name, err := r.cfg.Build.BuilderImage.Parse()
	if err != nil {
		return ArchImageConfig{}, fmt.Errorf(
			"invalid builder_image reference: %w", err,
		)
	}

	if category != "builders" {
		return ArchImageConfig{}, fmt.Errorf(
			"builder_image must reference @builders.* (got: @%s.%s)",
			category, name,
		)
	}

	img, err := r.cfg.BaseImages.GetBuilder(name)
	if err != nil {
		return ArchImageConfig{}, fmt.Errorf(
			"failed to resolve builder_image: %w", err,
		)
	}

	return img, nil
}

// ResolveRuntimeImage 解析运行时镜像
func (r *ImageResolver) ResolveRuntimeImage() (ArchImageConfig, error) {
	category, name, err := r.cfg.Build.RuntimeImage.Parse()
	if err != nil {
		return ArchImageConfig{}, fmt.Errorf(
			"invalid runtime_image reference: %w", err,
		)
	}

	if category != "runtimes" {
		return ArchImageConfig{}, fmt.Errorf(
			"runtime_image must reference @runtimes.* (got: @%s.%s)",
			category, name,
		)
	}

	img, err := r.cfg.BaseImages.GetRuntime(name)
	if err != nil {
		return ArchImageConfig{}, fmt.Errorf(
			"failed to resolve runtime_image: %w", err,
		)
	}

	return img, nil
}

// GetBuilderImageForArch 获取指定架构的构建镜像
func (r *ImageResolver) GetBuilderImageForArch(arch string) (string, error) {
	images, err := r.ResolveBuilderImage()
	if err != nil {
		return "", err
	}

	img, err := images.GetByArch(arch)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get builder image for %s: %w", arch, err,
		)
	}

	return img, nil
}

// GetRuntimeImageForArch 获取指定架构的运行时镜像
func (r *ImageResolver) GetRuntimeImageForArch(arch string) (string, error) {
	images, err := r.ResolveRuntimeImage()
	if err != nil {
		return "", err
	}

	img, err := images.GetByArch(arch)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get runtime image for %s: %w", arch, err,
		)
	}

	return img, nil
}

// MustResolveBuilderImage 解析构建镜像（panic on error）
func (r *ImageResolver) MustResolveBuilderImage() ArchImageConfig {
	img, err := r.ResolveBuilderImage()
	if err != nil {
		panic(fmt.Errorf("failed to resolve builder image: %w", err))
	}
	return img
}

// MustResolveRuntimeImage 解析运行时镜像（panic on error）
func (r *ImageResolver) MustResolveRuntimeImage() ArchImageConfig {
	img, err := r.ResolveRuntimeImage()
	if err != nil {
		panic(fmt.Errorf("failed to resolve runtime image: %w", err))
	}
	return img
}

// GetBuilderRef 获取构建镜像引用信息
func (r *ImageResolver) GetBuilderRef() (category, name string, err error) {
	return r.cfg.Build.BuilderImage.Parse()
}

// GetRuntimeRef 获取运行时镜像引用信息
func (r *ImageResolver) GetRuntimeRef() (category, name string, err error) {
	return r.cfg.Build.RuntimeImage.Parse()
}
