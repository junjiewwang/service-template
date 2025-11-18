package plugin

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

// ParserHandler Plugin子域解析处理器
type ParserHandler struct {
	*chain.BaseHandler
}

// NewParserHandler 创建Plugin解析处理器
func NewParserHandler() chain.ParserHandler {
	return &ParserHandler{
		BaseHandler: chain.NewBaseHandler("plugin-parser"),
	}
}

// Parse implements ParserHandler interface
func (h *ParserHandler) Parse(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理解析逻辑
func (h *ParserHandler) Handle(ctx *chain.ProcessingContext) error {
	rawConfig, ok := ctx.RawConfig["plugins"]
	if !ok {
		return h.CallNext(ctx)
	}

	pluginMap, ok := rawConfig.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid plugins configuration format")
	}

	domain := &PluginDomain{}

	// 解析安装目录
	if installDir, ok := pluginMap["install_dir"].(string); ok {
		domain.InstallDir = installDir
	}

	// 解析插件列表
	if items, ok := pluginMap["items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				plugin := Plugin{}

				if name, ok := itemMap["name"].(string); ok {
					plugin.Name = name
				}
				if desc, ok := itemMap["description"].(string); ok {
					plugin.Description = desc
				}

				// 下载URL可以是字符串或map
				plugin.DownloadURL = itemMap["download_url"]

				if cmd, ok := itemMap["install_command"].(string); ok {
					plugin.InstallCommand = cmd
				}
				if req, ok := itemMap["required"].(bool); ok {
					plugin.Required = req
				}

				// 解析运行时环境变量
				if runtimeEnv, ok := itemMap["runtime_env"].([]interface{}); ok {
					for _, env := range runtimeEnv {
						if envMap, ok := env.(map[string]interface{}); ok {
							envVar := PluginEnvVar{}
							if name, ok := envMap["name"].(string); ok {
								envVar.Name = name
							}
							if value, ok := envMap["value"].(string); ok {
								envVar.Value = value
							}
							plugin.RuntimeEnv = append(plugin.RuntimeEnv, envVar)
						}
					}
				}

				domain.Items = append(domain.Items, plugin)
			}
		}
	}

	ctx.SetDomainModel("plugin", domain)
	return h.CallNext(ctx)
}

// ValidatorHandler Plugin子域校验处理器
type ValidatorHandler struct {
	*chain.BaseHandler
}

// NewValidatorHandler 创建Plugin校验处理器
func NewValidatorHandler() chain.ValidatorHandler {
	return &ValidatorHandler{
		BaseHandler: chain.NewBaseHandler("plugin-validator"),
	}
}

// Validate implements ValidatorHandler interface
func (h *ValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理校验逻辑
func (h *ValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("plugin")
	if !ok {
		return h.CallNext(ctx)
	}

	pluginDomain, ok := domain.(*PluginDomain)
	if !ok {
		return fmt.Errorf("invalid plugin domain model type")
	}

	// 如果有插件，安装目录必需
	if pluginDomain.HasPlugins() && pluginDomain.InstallDir == "" {
		ctx.AddValidationError("plugins.install_dir", ErrPluginInstallDirRequired)
	}

	// 校验每个插件
	for i, plugin := range pluginDomain.Items {
		prefix := fmt.Sprintf("plugins.items[%d]", i)

		if plugin.Name == "" {
			ctx.AddValidationError(prefix+".name", ErrPluginNameRequired)
		}

		if plugin.DownloadURL == nil {
			ctx.AddValidationError(prefix+".download_url", ErrPluginDownloadURLRequired)
		} else {
			// 验证下载URL格式
			switch v := plugin.DownloadURL.(type) {
			case string:
				if v == "" {
					ctx.AddValidationError(prefix+".download_url", ErrPluginDownloadURLRequired)
				}
			case map[string]interface{}:
				if len(v) == 0 {
					ctx.AddValidationError(prefix+".download_url", ErrPluginDownloadURLRequired)
				}
			default:
				ctx.AddValidationError(prefix+".download_url", ErrInvalidDownloadURLFormat)
			}
		}

		if plugin.InstallCommand == "" {
			ctx.AddValidationError(prefix+".install_command", ErrPluginInstallCommandRequired)
		}
	}

	return h.CallNext(ctx)
}

// GeneratorHandler Plugin子域生成处理器
type GeneratorHandler struct {
	*chain.BaseHandler
}

// NewGeneratorHandler 创建Plugin生成处理器
func NewGeneratorHandler() chain.GeneratorHandler {
	return &GeneratorHandler{
		BaseHandler: chain.NewBaseHandler("plugin-generator"),
	}
}

// Generate implements GeneratorHandler interface
func (h *GeneratorHandler) Generate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理生成逻辑
func (h *GeneratorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("plugin")
	if !ok {
		return h.CallNext(ctx)
	}

	pluginDomain, ok := domain.(*PluginDomain)
	if !ok {
		return fmt.Errorf("invalid plugin domain model type")
	}

	if pluginDomain.HasPlugins() {
		// 记录生成的文件
		ctx.AddGeneratedFile("plugins/install.sh", []byte("# Plugin installation script"))

		// 为每个插件生成.env文件
		for _, plugin := range pluginDomain.Items {
			if plugin.HasRuntimeEnv() {
				ctx.AddGeneratedFile(fmt.Sprintf("plugins/%s/.env", plugin.Name),
					[]byte(fmt.Sprintf("# Runtime environment for %s plugin", plugin.Name)))
			}
		}

		// 添加元数据
		ctx.SetMetadata("plugins.count", len(pluginDomain.Items))
		ctx.SetMetadata("plugins.required_count", len(pluginDomain.GetRequiredPlugins()))
		ctx.SetMetadata("plugins.install_dir", pluginDomain.InstallDir)
	}

	return h.CallNext(ctx)
}
