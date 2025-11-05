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
	v.validateService()
	v.validateLanguage()
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

	if v.config.Service.DeployDir == "" {
		v.errors = append(v.errors, "service.deploy_dir is required")
	}
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

	if v.config.Language.Version == "" {
		v.errors = append(v.errors, "language.version is required")
	}
}

func (v *Validator) validateBuild() {
	if v.config.Build.BuilderImage.AMD64 == "" {
		v.errors = append(v.errors, "build.builder_image.amd64 is required")
	}
	if v.config.Build.BuilderImage.ARM64 == "" {
		v.errors = append(v.errors, "build.builder_image.arm64 is required")
	}

	if v.config.Build.RuntimeImage.AMD64 == "" {
		v.errors = append(v.errors, "build.runtime_image.amd64 is required")
	}
	if v.config.Build.RuntimeImage.ARM64 == "" {
		v.errors = append(v.errors, "build.runtime_image.arm64 is required")
	}

	if v.config.Build.Commands.Build == "" {
		v.errors = append(v.errors, "build.commands.build is required")
	}
}

func (v *Validator) validatePlugins() {
	for i, plugin := range v.config.Plugins {
		if plugin.Name == "" {
			v.errors = append(v.errors, fmt.Sprintf("plugins[%d].name is required", i))
		}
		if plugin.DownloadURL == "" {
			v.errors = append(v.errors, fmt.Sprintf("plugins[%d].download_url is required", i))
		}
		if plugin.InstallDir == "" {
			v.errors = append(v.errors, fmt.Sprintf("plugins[%d].install_dir is required", i))
		}
		if plugin.InstallCommand == "" {
			v.errors = append(v.errors, fmt.Sprintf("plugins[%d].install_command is required", i))
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
