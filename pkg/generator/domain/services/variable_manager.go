package services

import (
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/domain/models"
)

// VariableManager manages variable preparation for generators
// It provides a centralized way to prepare variables with paths configuration
type VariableManager struct {
	ctx         *context.GeneratorContext
	pathsConfig *models.PathsConfig
}

// NewVariableManager creates a new variable manager
func NewVariableManager(ctx *context.GeneratorContext) *VariableManager {
	// Create paths config from context
	pathsConfig := models.NewPathsConfig().
		WithPluginInstallDir(ctx.Config.Plugins.InstallDir).
		WithServiceDeployDir(ctx.Config.Service.DeployDir)

	return &VariableManager{
		ctx:         ctx,
		pathsConfig: pathsConfig,
	}
}

// NewVariableManagerWithPaths creates a variable manager with custom paths config
func NewVariableManagerWithPaths(ctx *context.GeneratorContext, pathsConfig *models.PathsConfig) *VariableManager {
	return &VariableManager{
		ctx:         ctx,
		pathsConfig: pathsConfig,
	}
}

// GetPathsConfig returns the paths configuration
func (m *VariableManager) GetPathsConfig() *models.PathsConfig {
	return m.pathsConfig
}

// GetContext returns the generator context
func (m *VariableManager) GetContext() *context.GeneratorContext {
	return m.ctx
}

// PrepareForDockerfile prepares variables for Dockerfile generation
func (m *VariableManager) PrepareForDockerfile(arch string) *context.VariableComposer {
	return m.ctx.GetVariablePreset().ForDockerfile(arch)
}

// PrepareForCompose prepares variables for docker-compose generation
func (m *VariableManager) PrepareForCompose() *context.VariableComposer {
	return m.ctx.GetVariablePreset().ForCompose()
}

// PrepareForBuildScript prepares variables for build script generation
func (m *VariableManager) PrepareForBuildScript() *context.VariableComposer {
	return m.ctx.GetVariablePreset().ForBuildScript()
}

// PrepareForScript prepares variables for general script generation
func (m *VariableManager) PrepareForScript() *context.VariableComposer {
	return m.ctx.GetVariablePreset().ForScript()
}

// PrepareForMakefile prepares variables for Makefile generation
func (m *VariableManager) PrepareForMakefile() *context.VariableComposer {
	return m.ctx.GetVariablePreset().ForMakefile()
}

// PrepareForDevOps prepares variables for DevOps configuration generation
func (m *VariableManager) PrepareForDevOps() *context.VariableComposer {
	return m.ctx.GetVariablePreset().ForDevOps()
}

// AddPathVariables adds path-related variables to the composer
func (m *VariableManager) AddPathVariables(composer *context.VariableComposer, serviceName string) *context.VariableComposer {
	return composer.
		WithCustom("SERVICE_ROOT", m.pathsConfig.GetServiceRoot(serviceName)).
		WithCustom("SERVICE_BIN_PATH", m.pathsConfig.GetServiceBinPath(serviceName)).
		WithCustom("SERVICE_CONFIG_PATH", m.pathsConfig.GetServiceConfigPath(serviceName)).
		WithCustom("SERVICE_LOG_PATH", m.pathsConfig.GetServiceLogPath(serviceName)).
		WithCustom("SERVICE_DATA_PATH", m.pathsConfig.GetServiceDataPath(serviceName)).
		WithCustom("PLUGIN_INSTALL_DIR", m.pathsConfig.PluginInstallDir)
}

// PrepareWithPaths prepares variables with path information
func (m *VariableManager) PrepareWithPaths(presetFunc func() *context.VariableComposer) *context.VariableComposer {
	composer := presetFunc()
	return m.AddPathVariables(composer, m.ctx.Config.Service.Name)
}
