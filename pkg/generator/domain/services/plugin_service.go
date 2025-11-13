package services

import (
	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

// PluginService handles plugin-related business logic
type PluginService struct {
	ctx    *context.GeneratorContext
	engine *core.TemplateEngine
}

// NewPluginService creates a new plugin service
func NewPluginService(ctx *context.GeneratorContext, engine *core.TemplateEngine) *PluginService {
	return &PluginService{
		ctx:    ctx,
		engine: engine,
	}
}

// PluginInfo represents plugin information (domain model)
type PluginInfo struct {
	Name           string
	Description    string
	DownloadURL    string
	InstallDir     string
	InstallCommand string
	RuntimeEnv     []config.EnvironmentVariable
}

// PrepareForDockerfile prepares plugin information for Dockerfile
func (s *PluginService) PrepareForDockerfile() []map[string]interface{} {
	var plugins []map[string]interface{}
	sharedInstallDir := s.ctx.Config.Plugins.InstallDir

	for _, plugin := range s.ctx.Config.Plugins.Items {
		pluginVars := s.ctx.Variables.WithPlugin(plugin, sharedInstallDir)
		plugins = append(plugins, map[string]interface{}{
			"InstallCommand": core.SubstituteVariables(plugin.InstallCommand, pluginVars.ToMap()),
			"Name":           plugin.Name,
			"InstallDir":     sharedInstallDir,
			"RuntimeEnv":     plugin.RuntimeEnv,
		})
	}

	return plugins
}

// PrepareForBuildScript prepares plugin information for build script
func (s *PluginService) PrepareForBuildScript() []PluginInfo {
	var plugins []PluginInfo
	sharedInstallDir := s.ctx.Config.Plugins.InstallDir

	for _, plugin := range s.ctx.Config.Plugins.Items {
		// Process runtime environment variables
		processedEnv := s.processRuntimeEnv(plugin.RuntimeEnv, sharedInstallDir)

		plugins = append(plugins, PluginInfo{
			Name:           plugin.Name,
			Description:    plugin.Description,
			DownloadURL:    plugin.DownloadURL,
			InstallDir:     sharedInstallDir,
			InstallCommand: plugin.InstallCommand,
			RuntimeEnv:     processedEnv,
		})
	}

	return plugins
}

// PrepareForEntrypoint prepares plugin environment variables for entrypoint script
func (s *PluginService) PrepareForEntrypoint() []map[string]interface{} {
	var pluginEnvs []map[string]interface{}
	sharedInstallDir := s.ctx.Config.Plugins.InstallDir

	for _, plugin := range s.ctx.Config.Plugins.Items {
		if len(plugin.RuntimeEnv) > 0 {
			pluginEnvs = append(pluginEnvs, map[string]interface{}{
				"Name":       plugin.Name,
				"InstallDir": sharedInstallDir,
				"RuntimeEnv": plugin.RuntimeEnv,
			})
		}
	}

	return pluginEnvs
}

// processRuntimeEnv processes runtime environment variables
func (s *PluginService) processRuntimeEnv(envVars []config.EnvironmentVariable, installDir string) []config.EnvironmentVariable {
	processed := make([]config.EnvironmentVariable, len(envVars))

	for i, env := range envVars {
		processed[i] = config.EnvironmentVariable{
			Name: env.Name,
			Value: s.engine.ReplaceVariables(env.Value, map[string]string{
				"PLUGIN_INSTALL_DIR": installDir,
			}),
		}
	}

	return processed
}

// HasPlugins checks if there are any plugins configured
func (s *PluginService) HasPlugins() bool {
	return len(s.ctx.Config.Plugins.Items) > 0
}

// GetInstallDir returns the plugin installation directory
func (s *PluginService) GetInstallDir() string {
	return s.ctx.Config.Plugins.InstallDir
}

// GetPluginCount returns the number of configured plugins
func (s *PluginService) GetPluginCount() int {
	return len(s.ctx.Config.Plugins.Items)
}
