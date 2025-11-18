package crossdomain

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/build"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/language"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/plugin"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/runtime"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/service"
)

// ValidatorHandler 跨域校验处理器
type ValidatorHandler struct {
	*chain.BaseHandler
}

// NewValidatorHandler 创建跨域校验处理器
func NewValidatorHandler() chain.ValidatorHandler {
	return &ValidatorHandler{
		BaseHandler: chain.NewBaseHandler("crossdomain-validator"),
	}
}

// Validate implements ValidatorHandler interface
func (h *ValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理跨域校验逻辑
func (h *ValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	// 1. 校验必需的子域是否存在
	h.validateRequiredDomains(ctx)

	// 2. 校验语言与构建配置的一致性
	h.validateLanguageBuildConsistency(ctx)

	// 3. 校验插件与运行时的一致性
	h.validatePluginRuntimeConsistency(ctx)

	// 4. 校验服务端口与健康检查的一致性
	h.validateServiceHealthcheckConsistency(ctx)

	return h.CallNext(ctx)
}

// validateRequiredDomains 校验必需的子域
func (h *ValidatorHandler) validateRequiredDomains(ctx *chain.ProcessingContext) {
	// Service子域是必需的
	if _, ok := ctx.GetDomainModel("service"); !ok {
		ctx.AddValidationError("service",
			fmt.Errorf("%w: service configuration is required", ErrMissingRequiredDomain))
	}

	// Runtime子域是必需的
	if _, ok := ctx.GetDomainModel("runtime"); !ok {
		ctx.AddValidationError("runtime",
			fmt.Errorf("%w: runtime configuration is required", ErrMissingRequiredDomain))
	}
}

// validateLanguageBuildConsistency 校验语言与构建配置的一致性
func (h *ValidatorHandler) validateLanguageBuildConsistency(ctx *chain.ProcessingContext) {
	langModel, hasLang := ctx.GetDomainModel("language")
	buildModel, hasBuild := ctx.GetDomainModel("build")

	if !hasLang || !hasBuild {
		return
	}

	langDomain, ok := langModel.(*language.LanguageConfig)
	if !ok {
		return
	}

	buildDomain, ok := buildModel.(*build.BuildDomain)
	if !ok {
		return
	}

	// 检查构建镜像是否适合该语言
	// 这里可以添加更复杂的校验逻辑
	// 例如：Go语言应该使用Go构建镜像
	if langDomain.Type != "" && len(buildDomain.BuilderImage) > 0 {
		// 简单的示例校验：确保构建命令不为空
		if buildDomain.Commands == nil || buildDomain.Commands.Build == "" {
			ctx.AddValidationError("build.commands.build",
				fmt.Errorf("%w: build command is required for language %s",
					ErrLanguageBuildMismatch, langDomain.Type))
		}
	}
}

// validatePluginRuntimeConsistency 校验插件与运行时的一致性
func (h *ValidatorHandler) validatePluginRuntimeConsistency(ctx *chain.ProcessingContext) {
	pluginModel, hasPlugin := ctx.GetDomainModel("plugin")
	runtimeModel, hasRuntime := ctx.GetDomainModel("runtime")

	if !hasPlugin || !hasRuntime {
		return
	}

	pluginDomain, ok := pluginModel.(*plugin.PluginDomain)
	if !ok {
		return
	}

	runtimeDomain, ok := runtimeModel.(*runtime.RuntimeDomain)
	if !ok {
		return
	}

	// 如果有必需的插件，确保运行时配置正确
	requiredPlugins := pluginDomain.GetRequiredPlugins()
	if len(requiredPlugins) > 0 {
		// 确保运行时有启动命令
		if runtimeDomain.Startup == nil || runtimeDomain.Startup.Command == "" {
			ctx.AddValidationError("runtime.startup.command",
				fmt.Errorf("%w: startup command is required when using plugins",
					ErrPluginRuntimeMismatch))
		}
	}
}

// validateServiceHealthcheckConsistency 校验服务端口与健康检查的一致性
func (h *ValidatorHandler) validateServiceHealthcheckConsistency(ctx *chain.ProcessingContext) {
	serviceModel, hasService := ctx.GetDomainModel("service")
	runtimeModel, hasRuntime := ctx.GetDomainModel("runtime")

	if !hasService || !hasRuntime {
		return
	}

	serviceDomain, ok := serviceModel.(*service.ServiceConfig)
	if !ok {
		return
	}

	runtimeDomain, ok := runtimeModel.(*runtime.RuntimeDomain)
	if !ok {
		return
	}

	// 如果服务有暴露的端口，建议启用健康检查
	if serviceDomain.HasExposedPorts() && !runtimeDomain.IsHealthcheckEnabled() {
		// 这只是一个警告，不是错误
		ctx.SetMetadata("warning.healthcheck",
			"Service has exposed ports but healthcheck is not enabled. Consider enabling healthcheck for better reliability.")
	}
}
