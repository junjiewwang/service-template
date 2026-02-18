package config

// DefaultBuildCommand 根据语言类型和语言配置推导默认的构建命令
// 返回空字符串表示该语言没有合理的默认构建命令
func DefaultBuildCommand(langType string, langCfg *LanguageConfig) string {
	fn, ok := defaultBuildCommandFuncs[langType]
	if !ok {
		return ""
	}
	return fn(langCfg)
}

// HasDefaultBuildCommand 检查指定语言是否有默认构建命令
func HasDefaultBuildCommand(langType string) bool {
	_, ok := defaultBuildCommandFuncs[langType]
	return ok
}

// ResolveBuildCommand 解析构建命令，支持自动推导
// 优先级：用户显式配置 > 按语言推导
func ResolveBuildCommand(cfg *ServiceConfig) string {
	if cfg.Build.Commands.Build != "" {
		return cfg.Build.Commands.Build
	}
	return DefaultBuildCommand(cfg.Language.Type, &cfg.Language)
}

// ============================================
// 默认构建命令推导映射表
// ============================================

var defaultBuildCommandFuncs = map[string]func(cfg *LanguageConfig) string{
	"go": func(cfg *LanguageConfig) string {
		// Go 标准构建：静态编译，输出到 ${BUILD_OUTPUT_DIR}/bin/
		return `CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/server`
	},
	"python": func(cfg *LanguageConfig) string {
		// Python 拷贝源码到构建输出目录
		return `cp -r . ${BUILD_OUTPUT_DIR}/`
	},
	"java": func(cfg *LanguageConfig) string {
		buildTool := cfg.GetString("build_tool", "maven")
		if buildTool == "gradle" {
			return `gradle build -x test && cp build/libs/*.jar ${BUILD_OUTPUT_DIR}/app.jar`
		}
		return `mvn package -DskipTests && cp target/*.jar ${BUILD_OUTPUT_DIR}/app.jar`
	},
	"nodejs": func(cfg *LanguageConfig) string {
		// Node.js 拷贝源码和 node_modules 到构建输出目录
		return `npm run build 2>/dev/null || true && cp -r . ${BUILD_OUTPUT_DIR}/`
	},
	"rust": func(cfg *LanguageConfig) string {
		// Rust 静态编译，输出到 ${BUILD_OUTPUT_DIR}/bin/
		return `cargo build --release && cp target/release/${SERVICE_NAME} ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME}`
	},
}
