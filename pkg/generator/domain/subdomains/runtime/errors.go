package runtime

import "errors"

var (
	// ErrStartupCommandRequired 启动命令必需
	ErrStartupCommandRequired = errors.New("startup command is required")

	// ErrCustomHealthcheckScriptRequired 自定义健康检查脚本必需
	ErrCustomHealthcheckScriptRequired = errors.New("custom healthcheck script is required when type is custom")

	// ErrInvalidHealthcheckType 无效的健康检查类型
	ErrInvalidHealthcheckType = errors.New("invalid healthcheck type, must be 'default' or 'custom'")
)
