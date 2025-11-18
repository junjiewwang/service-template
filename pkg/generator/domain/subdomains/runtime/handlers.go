package runtime

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

// ParserHandler Runtime子域解析处理器
type ParserHandler struct {
	*chain.BaseHandler
}

// NewParserHandler 创建Runtime解析处理器
func NewParserHandler() chain.ParserHandler {
	return &ParserHandler{
		BaseHandler: chain.NewBaseHandler("runtime-parser"),
	}
}

// Parse implements ParserHandler interface
func (h *ParserHandler) Parse(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理解析逻辑
func (h *ParserHandler) Handle(ctx *chain.ProcessingContext) error {
	rawConfig, ok := ctx.RawConfig["runtime"]
	if !ok {
		return h.CallNext(ctx)
	}

	runtimeMap, ok := rawConfig.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid runtime configuration format")
	}

	domain := &RuntimeDomain{}

	// 解析系统依赖
	if sysDeps, ok := runtimeMap["system_dependencies"].(map[string]interface{}); ok {
		domain.SystemDependencies = &SystemDependencies{}
		if packages, ok := sysDeps["packages"].([]interface{}); ok {
			for _, pkg := range packages {
				if pkgStr, ok := pkg.(string); ok {
					domain.SystemDependencies.Packages = append(domain.SystemDependencies.Packages, pkgStr)
				}
			}
		}
	}

	// 解析健康检查配置
	if healthcheck, ok := runtimeMap["healthcheck"].(map[string]interface{}); ok {
		domain.Healthcheck = &HealthcheckConfig{}
		if enabled, ok := healthcheck["enabled"].(bool); ok {
			domain.Healthcheck.Enabled = enabled
		}
		if hcType, ok := healthcheck["type"].(string); ok {
			domain.Healthcheck.Type = hcType
		}
		if customScript, ok := healthcheck["custom_script"].(string); ok {
			domain.Healthcheck.CustomScript = customScript
		}
	}

	// 解析启动配置
	if startup, ok := runtimeMap["startup"].(map[string]interface{}); ok {
		domain.Startup = &StartupConfig{}
		if command, ok := startup["command"].(string); ok {
			domain.Startup.Command = command
		}
		if env, ok := startup["env"].([]interface{}); ok {
			for _, e := range env {
				if envMap, ok := e.(map[string]interface{}); ok {
					envVar := EnvVar{}
					if name, ok := envMap["name"].(string); ok {
						envVar.Name = name
					}
					if value, ok := envMap["value"].(string); ok {
						envVar.Value = value
					}
					domain.Startup.Env = append(domain.Startup.Env, envVar)
				}
			}
		}
	}

	ctx.SetDomainModel("runtime", domain)
	return h.CallNext(ctx)
}

// ValidatorHandler Runtime子域校验处理器
type ValidatorHandler struct {
	*chain.BaseHandler
}

// NewValidatorHandler 创建Runtime校验处理器
func NewValidatorHandler() chain.ValidatorHandler {
	return &ValidatorHandler{
		BaseHandler: chain.NewBaseHandler("runtime-validator"),
	}
}

// Validate implements ValidatorHandler interface
func (h *ValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理校验逻辑
func (h *ValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("runtime")
	if !ok {
		return h.CallNext(ctx)
	}

	runtimeDomain, ok := domain.(*RuntimeDomain)
	if !ok {
		return fmt.Errorf("invalid runtime domain model type")
	}

	// 校验启动命令
	if runtimeDomain.Startup == nil || runtimeDomain.Startup.Command == "" {
		ctx.AddValidationError("runtime.startup.command", ErrStartupCommandRequired)
	}

	// 校验健康检查配置
	if runtimeDomain.IsHealthcheckEnabled() {
		if runtimeDomain.Healthcheck.Type != "" &&
			runtimeDomain.Healthcheck.Type != "default" &&
			runtimeDomain.Healthcheck.Type != "custom" {
			ctx.AddValidationError("runtime.healthcheck.type", ErrInvalidHealthcheckType)
		}

		// 如果是自定义健康检查，必须提供脚本
		if runtimeDomain.IsCustomHealthcheck() && runtimeDomain.Healthcheck.CustomScript == "" {
			ctx.AddValidationError("runtime.healthcheck.custom_script", ErrCustomHealthcheckScriptRequired)
		}
	}

	return h.CallNext(ctx)
}

// GeneratorHandler Runtime子域生成处理器
type GeneratorHandler struct {
	*chain.BaseHandler
}

// NewGeneratorHandler 创建Runtime生成处理器
func NewGeneratorHandler() chain.GeneratorHandler {
	return &GeneratorHandler{
		BaseHandler: chain.NewBaseHandler("runtime-generator"),
	}
}

// Generate implements GeneratorHandler interface
func (h *GeneratorHandler) Generate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理生成逻辑
func (h *GeneratorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("runtime")
	if !ok {
		return h.CallNext(ctx)
	}

	runtimeDomain, ok := domain.(*RuntimeDomain)
	if !ok {
		return fmt.Errorf("invalid runtime domain model type")
	}

	// 记录生成的文件
	ctx.AddGeneratedFile("entrypoint.sh", []byte("# Service entrypoint script"))
	ctx.AddGeneratedFile("rt_prepare.sh", []byte("# Runtime preparation script"))

	if runtimeDomain.IsHealthcheckEnabled() {
		ctx.AddGeneratedFile("healthcheck.sh", []byte("# Health check script"))
	}

	// 添加元数据
	ctx.SetMetadata("runtime.has_system_deps", runtimeDomain.HasSystemDependencies())
	ctx.SetMetadata("runtime.healthcheck_enabled", runtimeDomain.IsHealthcheckEnabled())
	ctx.SetMetadata("runtime.healthcheck_type", func() string {
		if runtimeDomain.Healthcheck != nil {
			return runtimeDomain.Healthcheck.Type
		}
		return "default"
	}())

	return h.CallNext(ctx)
}
