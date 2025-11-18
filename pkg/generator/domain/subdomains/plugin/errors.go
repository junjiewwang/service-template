package plugin

import "errors"

var (
	// ErrPluginNameRequired 插件名称必需
	ErrPluginNameRequired = errors.New("plugin name is required")

	// ErrPluginDownloadURLRequired 插件下载URL必需
	ErrPluginDownloadURLRequired = errors.New("plugin download URL is required")

	// ErrPluginInstallCommandRequired 插件安装命令必需
	ErrPluginInstallCommandRequired = errors.New("plugin install command is required")

	// ErrPluginInstallDirRequired 插件安装目录必需
	ErrPluginInstallDirRequired = errors.New("plugin install directory is required")

	// ErrInvalidDownloadURLFormat 无效的下载URL格式
	ErrInvalidDownloadURLFormat = errors.New("invalid download URL format")
)
