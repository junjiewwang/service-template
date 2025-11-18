package localdev

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

// ParserHandler LocalDev子域解析处理器
type ParserHandler struct {
	*chain.BaseHandler
}

// NewParserHandler 创建LocalDev解析处理器
func NewParserHandler() chain.ParserHandler {
	return &ParserHandler{
		BaseHandler: chain.NewBaseHandler("localdev-parser"),
	}
}

// Parse implements ParserHandler interface
func (h *ParserHandler) Parse(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理解析逻辑
func (h *ParserHandler) Handle(ctx *chain.ProcessingContext) error {
	rawConfig, ok := ctx.RawConfig["local_dev"]
	if !ok {
		return h.CallNext(ctx)
	}

	localDevMap, ok := rawConfig.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid local_dev configuration format")
	}

	domain := &LocalDevDomain{}

	// 解析Compose配置
	if compose, ok := localDevMap["compose"].(map[string]interface{}); ok {
		domain.Compose = &ComposeConfig{}

		// 解析资源配置
		if resources, ok := compose["resources"].(map[string]interface{}); ok {
			domain.Compose.Resources = &ResourcesConfig{}

			if limits, ok := resources["limits"].(map[string]interface{}); ok {
				domain.Compose.Resources.Limits = &ResourceLimit{}
				if cpus, ok := limits["cpus"].(string); ok {
					domain.Compose.Resources.Limits.CPUs = cpus
				}
				if memory, ok := limits["memory"].(string); ok {
					domain.Compose.Resources.Limits.Memory = memory
				}
			}

			if reservations, ok := resources["reservations"].(map[string]interface{}); ok {
				domain.Compose.Resources.Reservations = &ResourceLimit{}
				if cpus, ok := reservations["cpus"].(string); ok {
					domain.Compose.Resources.Reservations.CPUs = cpus
				}
				if memory, ok := reservations["memory"].(string); ok {
					domain.Compose.Resources.Reservations.Memory = memory
				}
			}
		}

		// 解析环境变量
		if environment, ok := compose["environment"].([]interface{}); ok {
			for _, env := range environment {
				if envMap, ok := env.(map[string]interface{}); ok {
					envVar := EnvVar{}
					if name, ok := envMap["name"].(string); ok {
						envVar.Name = name
					}
					if value, ok := envMap["value"].(string); ok {
						envVar.Value = value
					}
					domain.Compose.Environment = append(domain.Compose.Environment, envVar)
				}
			}
		}

		// 解析Entrypoint
		if entrypoint, ok := compose["entrypoint"].([]interface{}); ok {
			for _, ep := range entrypoint {
				if epStr, ok := ep.(string); ok {
					domain.Compose.Entrypoint = append(domain.Compose.Entrypoint, epStr)
				}
			}
		}

		// 解析卷挂载
		if volumes, ok := compose["volumes"].([]interface{}); ok {
			for _, vol := range volumes {
				if volMap, ok := vol.(map[string]interface{}); ok {
					volume := VolumeMount{}
					if source, ok := volMap["source"].(string); ok {
						volume.Source = source
					}
					if target, ok := volMap["target"].(string); ok {
						volume.Target = target
					}
					if volType, ok := volMap["type"].(string); ok {
						volume.Type = volType
					}
					if desc, ok := volMap["description"].(string); ok {
						volume.Description = desc
					}
					domain.Compose.Volumes = append(domain.Compose.Volumes, volume)
				}
			}
		}

		// 解析健康检查
		if healthcheck, ok := compose["healthcheck"].(map[string]interface{}); ok {
			domain.Compose.Healthcheck = &ComposeHealthcheck{}
			if interval, ok := healthcheck["interval"].(string); ok {
				domain.Compose.Healthcheck.Interval = interval
			}
			if timeout, ok := healthcheck["timeout"].(string); ok {
				domain.Compose.Healthcheck.Timeout = timeout
			}
			if retries, ok := healthcheck["retries"].(int); ok {
				domain.Compose.Healthcheck.Retries = retries
			}
			if startPeriod, ok := healthcheck["start_period"].(string); ok {
				domain.Compose.Healthcheck.StartPeriod = startPeriod
			}
		}

		// 解析标签
		if labels, ok := compose["labels"].(map[string]interface{}); ok {
			domain.Compose.Labels = make(map[string]string)
			for key, value := range labels {
				if valueStr, ok := value.(string); ok {
					domain.Compose.Labels[key] = valueStr
				}
			}
		}
	}

	// 解析Kubernetes配置
	if k8s, ok := localDevMap["kubernetes"].(map[string]interface{}); ok {
		domain.Kubernetes = &KubernetesConfig{}
		if enabled, ok := k8s["enabled"].(bool); ok {
			domain.Kubernetes.Enabled = enabled
		}
		if namespace, ok := k8s["namespace"].(string); ok {
			domain.Kubernetes.Namespace = namespace
		}
		if outputDir, ok := k8s["output_dir"].(string); ok {
			domain.Kubernetes.OutputDir = outputDir
		}
		if volumeType, ok := k8s["volume_type"].(string); ok {
			domain.Kubernetes.VolumeType = volumeType
		}

		if wait, ok := k8s["wait"].(map[string]interface{}); ok {
			domain.Kubernetes.Wait = &K8sWaitConfig{}
			if enabled, ok := wait["enabled"].(bool); ok {
				domain.Kubernetes.Wait.Enabled = enabled
			}
			if timeout, ok := wait["timeout"].(string); ok {
				domain.Kubernetes.Wait.Timeout = timeout
			}
		}
	}

	ctx.SetDomainModel("localdev", domain)
	return h.CallNext(ctx)
}

// ValidatorHandler LocalDev子域校验处理器
type ValidatorHandler struct {
	*chain.BaseHandler
}

// NewValidatorHandler 创建LocalDev校验处理器
func NewValidatorHandler() chain.ValidatorHandler {
	return &ValidatorHandler{
		BaseHandler: chain.NewBaseHandler("localdev-validator"),
	}
}

// Validate implements ValidatorHandler interface
func (h *ValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理校验逻辑
func (h *ValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("localdev")
	if !ok {
		return h.CallNext(ctx)
	}

	localDevDomain, ok := domain.(*LocalDevDomain)
	if !ok {
		return fmt.Errorf("invalid localdev domain model type")
	}

	// Compose配置必需
	if localDevDomain.Compose == nil {
		ctx.AddValidationError("local_dev.compose", ErrComposeConfigRequired)
		return h.CallNext(ctx)
	}

	// 校验卷挂载
	if localDevDomain.Compose.HasVolumes() {
		for i, vol := range localDevDomain.Compose.Volumes {
			prefix := fmt.Sprintf("local_dev.compose.volumes[%d]", i)

			if vol.Source == "" {
				ctx.AddValidationError(prefix+".source", ErrVolumeSourceRequired)
			}
			if vol.Target == "" {
				ctx.AddValidationError(prefix+".target", ErrVolumeTargetRequired)
			}
			if vol.Type != "" && vol.Type != "bind" && vol.Type != "volume" {
				ctx.AddValidationError(prefix+".type", ErrInvalidVolumeType)
			}
		}
	}

	// 校验Kubernetes配置
	if localDevDomain.IsK8sEnabled() {
		if localDevDomain.Kubernetes.VolumeType != "" {
			validTypes := map[string]bool{
				"configMap":             true,
				"persistentVolumeClaim": true,
				"emptyDir":              true,
				"hostPath":              true,
			}
			if !validTypes[localDevDomain.Kubernetes.VolumeType] {
				ctx.AddValidationError("local_dev.kubernetes.volume_type", ErrInvalidK8sVolumeType)
			}
		}
	}

	return h.CallNext(ctx)
}

// GeneratorHandler LocalDev子域生成处理器
type GeneratorHandler struct {
	*chain.BaseHandler
}

// NewGeneratorHandler 创建LocalDev生成处理器
func NewGeneratorHandler() chain.GeneratorHandler {
	return &GeneratorHandler{
		BaseHandler: chain.NewBaseHandler("localdev-generator"),
	}
}

// Generate implements GeneratorHandler interface
func (h *GeneratorHandler) Generate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理生成逻辑
func (h *GeneratorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("localdev")
	if !ok {
		return h.CallNext(ctx)
	}

	localDevDomain, ok := domain.(*LocalDevDomain)
	if !ok {
		return fmt.Errorf("invalid localdev domain model type")
	}

	// 记录生成的文件
	ctx.AddGeneratedFile("compose.yaml", []byte("# Docker Compose configuration"))
	ctx.AddGeneratedFile("Makefile", []byte("# Build and deployment automation"))

	if localDevDomain.IsK8sEnabled() {
		ctx.AddGeneratedFile("k8s-manifests/README.md", []byte("# Kubernetes manifests directory"))
	}

	// 添加元数据
	ctx.SetMetadata("localdev.has_volumes", localDevDomain.Compose.HasVolumes())
	ctx.SetMetadata("localdev.has_custom_entrypoint", localDevDomain.Compose.HasEntrypoint())
	ctx.SetMetadata("localdev.k8s_enabled", localDevDomain.IsK8sEnabled())
	if localDevDomain.Compose.HasVolumes() {
		ctx.SetMetadata("localdev.volume_count", len(localDevDomain.Compose.Volumes))
	}

	return h.CallNext(ctx)
}
