package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// PluginBuilder 插件配置构建器
type PluginBuilder struct {
	plugins *config.PluginsConfig
}

// InstallDir 设置插件安装目录
func (b *PluginBuilder) InstallDir(dir string) *PluginBuilder {
	b.plugins.InstallDir = dir
	return b
}

// AddPlugin 添加插件（简化版）
func (b *PluginBuilder) AddPlugin(name, description, url string) *PluginBuilder {
	b.plugins.Items = append(b.plugins.Items, config.PluginConfig{
		Name:           name,
		Description:    description,
		DownloadURL:    config.NewStaticDownloadURL(url),
		InstallCommand: "echo 'Installing " + name + "...'",
		Required:       true,
	})
	return b
}

// AddPluginWithCommand 添加带自定义安装命令的插件
func (b *PluginBuilder) AddPluginWithCommand(name, description, url, installCmd string) *PluginBuilder {
	b.plugins.Items = append(b.plugins.Items, config.PluginConfig{
		Name:           name,
		Description:    description,
		DownloadURL:    config.NewStaticDownloadURL(url),
		InstallCommand: installCmd,
		Required:       true,
	})
	return b
}

// AddPluginConfig 添加完整的插件配置
func (b *PluginBuilder) AddPluginConfig(cfg config.PluginConfig) *PluginBuilder {
	b.plugins.Items = append(b.plugins.Items, cfg)
	return b
}

// AddOptionalPlugin 添加可选插件
func (b *PluginBuilder) AddOptionalPlugin(name, description, url string) *PluginBuilder {
	b.plugins.Items = append(b.plugins.Items, config.PluginConfig{
		Name:           name,
		Description:    description,
		DownloadURL:    config.NewStaticDownloadURL(url),
		InstallCommand: "echo 'Installing " + name + "...'",
		Required:       false,
	})
	return b
}

// SetPlugins 设置插件列表（替换现有）
func (b *PluginBuilder) SetPlugins(plugins []config.PluginConfig) *PluginBuilder {
	b.plugins.Items = plugins
	return b
}

// ClearPlugins 清空插件列表
func (b *PluginBuilder) ClearPlugins() *PluginBuilder {
	b.plugins.Items = []config.PluginConfig{}
	return b
}
