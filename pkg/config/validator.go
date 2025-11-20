package config

import (
	"fmt"
	"strings"
)

// Validator validates service configuration
type Validator struct {
	config *ServiceConfig
	errors []string
}

// NewValidator creates a new configuration validator
func NewValidator(config *ServiceConfig) *Validator {
	return &Validator{
		config: config,
		errors: []string{},
	}
}

// Validate performs comprehensive validation of the configuration
func (v *Validator) Validate() error {
	// 1. 验证基础镜像配置（必须先验证，因为后续会引用）
	v.validateBaseImages()

	// 2. 验证基础信息
	v.validateService()
	v.validateLanguage()

	// 3. 验证镜像引用（依赖 base_images）
	v.validateImageReferences()

	// 4. 验证其他配置
	v.validateBuild()
	v.validatePlugins()
	v.validateRuntime()
	v.validateLocalDev()

	if len(v.errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(v.errors, "\n  - "))
	}

	return nil
}

func (v *Validator) validateService() {
	if v.config.Service.Name == "" {
		v.errors = append(v.errors, "service.name is required")
	}

	if len(v.config.Service.Ports) == 0 {
		v.errors = append(v.errors, "service.ports must have at least one port")
	}

	for i, port := range v.config.Service.Ports {
		if port.Name == "" {
			v.errors = append(v.errors, fmt.Sprintf("service.ports[%d].name is required", i))
		}
		if port.Port <= 0 || port.Port > 65535 {
			v.errors = append(v.errors, fmt.Sprintf("service.ports[%d].port must be between 1 and 65535", i))
		}
		if port.Protocol == "" {
			v.errors = append(v.errors, fmt.Sprintf("service.ports[%d].protocol is required", i))
		}
	}

	// deploy_dir has a default value, so no validation needed
}

func (v *Validator) validateLanguage() {
	validLanguages := map[string]bool{
		"go":     true,
		"python": true,
		"nodejs": true,
		"java":   true,
		"rust":   true,
	}

	if v.config.Language.Type == "" {
		v.errors = append(v.errors, "language.type is required")
	} else if !validLanguages[v.config.Language.Type] {
		v.errors = append(v.errors, fmt.Sprintf("language.type '%s' is not supported (valid: go, python, nodejs, java, rust)", v.config.Language.Type))
	}

	// Config is optional, no validation needed
}

// validateBaseImages 验证基础镜像配置
func (v *Validator) validateBaseImages() {
	if err := v.config.BaseImages.Validate(); err != nil {
		v.errors = append(v.errors, fmt.Sprintf("base_images: %v", err))
	}
}

// validateImageReferences 验证镜像引用
func (v *Validator) validateImageReferences() {
	// 验证 builder_image 引用格式
	if err := v.config.Build.BuilderImage.Validate(); err != nil {
		v.errors = append(v.errors, fmt.Sprintf("build.builder_image: %v", err))
		return // 格式错误，后续验证无意义
	}

	// 验证 runtime_image 引用格式
	if err := v.config.Build.RuntimeImage.Validate(); err != nil {
		v.errors = append(v.errors, fmt.Sprintf("build.runtime_image: %v", err))
		return
	}

	// 验证引用是否存在且类型正确
	resolver := NewImageResolver(v.config)

	// 验证 builder_image
	if !v.config.Build.BuilderImage.IsBuilder() {
		v.errors = append(v.errors, fmt.Sprintf(
			"build.builder_image: must reference @builders.* (got: %s)",
			v.config.Build.BuilderImage,
		))
	} else {
		if _, err := resolver.ResolveBuilderImage(); err != nil {
			v.errors = append(v.errors, fmt.Sprintf("build.builder_image: %v", err))
		}
	}

	// 验证 runtime_image
	if !v.config.Build.RuntimeImage.IsRuntime() {
		v.errors = append(v.errors, fmt.Sprintf(
			"build.runtime_image: must reference @runtimes.* (got: %s)",
			v.config.Build.RuntimeImage,
		))
	} else {
		if _, err := resolver.ResolveRuntimeImage(); err != nil {
			v.errors = append(v.errors, fmt.Sprintf("build.runtime_image: %v", err))
		}
	}
}

func (v *Validator) validateBuild() {
	// 镜像验证已在 validateImageReferences 中完成

	if v.config.Build.Commands.Build == "" {
		v.errors = append(v.errors, "build.commands.build is required")
	}
}

func (v *Validator) validatePlugins() {
	// 如果有插件配置，验证 install_dir
	if len(v.config.Plugins.Items) > 0 {
		if v.config.Plugins.InstallDir == "" {
			v.errors = append(v.errors, "plugins.install_dir is required when plugins are configured")
		}
	}

	// 验证每个插件
	for i, plugin := range v.config.Plugins.Items {
		if plugin.Name == "" {
			v.errors = append(v.errors, fmt.Sprintf("plugins.items[%d].name is required", i))
		}

		// 验证 download_url
		if plugin.DownloadURL.IsEmpty() {
			v.errors = append(v.errors, fmt.Sprintf("plugins.items[%d].download_url is required", i))
		} else if plugin.DownloadURL.IsArchMapping() {
			// 如果是架构映射，验证架构键的合法性
			urls, _ := plugin.DownloadURL.GetArchURLs()

			validArchs := map[string]bool{
				"x86_64":  true,
				"amd64":   true,
				"aarch64": true,
				"arm64":   true,
				"default": true,
			}

			for arch, url := range urls {
				if !validArchs[arch] {
					v.errors = append(v.errors, fmt.Sprintf(
						"plugins.items[%d].download_url: unsupported architecture '%s'. "+
							"Supported: x86_64, amd64, aarch64, arm64, default",
						i, arch,
					))
				}
				if url == "" {
					v.errors = append(v.errors, fmt.Sprintf(
						"plugins.items[%d].download_url: URL for architecture '%s' cannot be empty",
						i, arch,
					))
				}
			}
		}

		if plugin.InstallCommand == "" {
			v.errors = append(v.errors, fmt.Sprintf("plugins.items[%d].install_command is required", i))
		}
	}
}

func (v *Validator) validateRuntime() {
	// Validate healthcheck configuration
	if v.config.Runtime.Healthcheck.Enabled {
		// Validate healthcheck type
		validTypes := map[string]bool{
			"default": true,
			"custom":  true,
			"":        true, // Empty defaults to "default"
		}

		hcType := v.config.Runtime.Healthcheck.Type
		if hcType == "" {
			hcType = "default" // Set default value
		}

		if !validTypes[hcType] {
			v.errors = append(v.errors, fmt.Sprintf("runtime.healthcheck.type '%s' is not valid (valid: default, custom)", hcType))
		}

		// Validate custom healthcheck requirements
		if hcType == "custom" {
			if v.config.Runtime.Healthcheck.CustomScript == "" {
				v.errors = append(v.errors, "runtime.healthcheck.custom_script is required when type is 'custom'")
			}
		}
	}

	// Validate startup command
	if v.config.Runtime.Startup.Command == "" {
		v.errors = append(v.errors, "runtime.startup.command is required")
	}
}

func (v *Validator) validateLocalDev() {
	if v.config.LocalDev.Kubernetes.Enabled {
		validVolumeTypes := map[string]bool{
			"configMap":             true,
			"persistentVolumeClaim": true,
			"emptyDir":              true,
			"hostPath":              true,
		}

		if v.config.LocalDev.Kubernetes.VolumeType != "" && !validVolumeTypes[v.config.LocalDev.Kubernetes.VolumeType] {
			v.errors = append(v.errors, fmt.Sprintf("local_dev.kubernetes.volume_type '%s' is not valid", v.config.LocalDev.Kubernetes.VolumeType))
		}

		if v.config.LocalDev.Kubernetes.Namespace == "" {
			v.errors = append(v.errors, "local_dev.kubernetes.namespace is required when kubernetes is enabled")
		}
	}

	for i, vol := range v.config.LocalDev.Compose.Volumes {
		if vol.Source == "" {
			v.errors = append(v.errors, fmt.Sprintf("local_dev.compose.volumes[%d].source is required", i))
		}
		if vol.Target == "" {
			v.errors = append(v.errors, fmt.Sprintf("local_dev.compose.volumes[%d].target is required", i))
		}
		if vol.Type != "bind" && vol.Type != "volume" {
			v.errors = append(v.errors, fmt.Sprintf("local_dev.compose.volumes[%d].type must be 'bind' or 'volume'", i))
		}
	}
}
