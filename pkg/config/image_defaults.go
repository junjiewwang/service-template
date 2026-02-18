package config

import "fmt"

// DefaultBuilderImage 根据语言类型和语言配置推导默认的构建镜像
// 返回的是 multi-arch 公开镜像（Docker Hub），同一地址支持 amd64/arm64
func DefaultBuilderImage(langType string, langCfg *LanguageConfig) (string, error) {
	fn, ok := defaultBuilderImageFuncs[langType]
	if !ok {
		return "", fmt.Errorf("no default builder image for language: %s", langType)
	}
	return fn(langCfg), nil
}

// DefaultRuntimeImage 根据语言类型和语言配置推导默认的运行时镜像
func DefaultRuntimeImage(langType string, langCfg *LanguageConfig) (string, error) {
	fn, ok := defaultRuntimeImageFuncs[langType]
	if !ok {
		return "", fmt.Errorf("no default runtime image for language: %s", langType)
	}
	return fn(langCfg), nil
}

// HasDefaultImages 检查指定语言是否有默认镜像推导支持
func HasDefaultImages(langType string) bool {
	_, hasBuilder := defaultBuilderImageFuncs[langType]
	_, hasRuntime := defaultRuntimeImageFuncs[langType]
	return hasBuilder && hasRuntime
}

// ResolveBuilderImageWithDefaults 解析构建镜像，支持自动推导
// 优先级：用户显式配置 > 按语言推导
func ResolveBuilderImageWithDefaults(cfg *ServiceConfig) (ArchImageConfig, error) {
	if !cfg.Build.BuilderImage.IsEmpty() {
		return cfg.Build.BuilderImage.Resolve(&cfg.BaseImages, "builders")
	}
	// 未指定，按语言推导
	image, err := DefaultBuilderImage(cfg.Language.Type, &cfg.Language)
	if err != nil {
		return ArchImageConfig{}, fmt.Errorf("builder_image not specified and %w", err)
	}
	return ArchImageConfig{AMD64: image, ARM64: image}, nil
}

// ResolveRuntimeImageWithDefaults 解析运行时镜像，支持自动推导
func ResolveRuntimeImageWithDefaults(cfg *ServiceConfig) (ArchImageConfig, error) {
	if !cfg.Build.RuntimeImage.IsEmpty() {
		return cfg.Build.RuntimeImage.Resolve(&cfg.BaseImages, "runtimes")
	}
	// 未指定，按语言推导
	image, err := DefaultRuntimeImage(cfg.Language.Type, &cfg.Language)
	if err != nil {
		return ArchImageConfig{}, fmt.Errorf("runtime_image not specified and %w", err)
	}
	return ArchImageConfig{AMD64: image, ARM64: image}, nil
}

// ============================================
// 默认镜像推导映射表
// ============================================

var defaultBuilderImageFuncs = map[string]func(cfg *LanguageConfig) string{
	"go": func(cfg *LanguageConfig) string {
		version := cfg.GetString("go_version", "1.23")
		return fmt.Sprintf("golang:%s-alpine", version)
	},
	"python": func(cfg *LanguageConfig) string {
		version := cfg.GetString("python_version", "3.12")
		return fmt.Sprintf("python:%s-slim", version)
	},
	"java": func(cfg *LanguageConfig) string {
		buildTool := cfg.GetString("build_tool", "maven")
		jdkVersion := cfg.GetString("jdk_version", "21")
		if buildTool == "gradle" {
			gradleVersion := cfg.GetString("gradle_version", "8")
			return fmt.Sprintf("gradle:%s-jdk%s", gradleVersion, jdkVersion)
		}
		return fmt.Sprintf("maven:3-eclipse-temurin-%s", jdkVersion)
	},
	"nodejs": func(cfg *LanguageConfig) string {
		version := cfg.GetString("node_version", "20")
		return fmt.Sprintf("node:%s-alpine", version)
	},
	"rust": func(cfg *LanguageConfig) string {
		version := cfg.GetString("rust_version", "1.78")
		return fmt.Sprintf("rust:%s-alpine", version)
	},
}

var defaultRuntimeImageFuncs = map[string]func(cfg *LanguageConfig) string{
	"go": func(cfg *LanguageConfig) string {
		// Go 静态编译，不需要语言运行环境
		return "alpine:3.19"
	},
	"python": func(cfg *LanguageConfig) string {
		// Python 运行时需要 Python 环境
		version := cfg.GetString("python_version", "3.12")
		return fmt.Sprintf("python:%s-slim", version)
	},
	"java": func(cfg *LanguageConfig) string {
		// Java 运行时只需要 JRE
		jdkVersion := cfg.GetString("jdk_version", "21")
		return fmt.Sprintf("eclipse-temurin:%s-jre-alpine", jdkVersion)
	},
	"nodejs": func(cfg *LanguageConfig) string {
		// Node.js 运行时需要 Node 环境
		version := cfg.GetString("node_version", "20")
		return fmt.Sprintf("node:%s-alpine", version)
	},
	"rust": func(cfg *LanguageConfig) string {
		// Rust 静态编译，不需要语言运行环境
		return "alpine:3.19"
	},
}
