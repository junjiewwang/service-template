package crossdomain

import "errors"

var (
	// ErrLanguageBuildMismatch 语言与构建配置不匹配
	ErrLanguageBuildMismatch = errors.New("language configuration does not match build configuration")

	// ErrPluginRuntimeMismatch 插件与运行时配置不匹配
	ErrPluginRuntimeMismatch = errors.New("plugin configuration does not match runtime configuration")

	// ErrMissingRequiredDomain 缺少必需的子域配置
	ErrMissingRequiredDomain = errors.New("missing required domain configuration")
)
