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
	Name              string
	Description       string
	DownloadURL       string // For template rendering (static URL or placeholder)
	URLResolverScript string // Shell script to resolve URL at runtime
	InstallDir        string
	InstallCommand    string
	RuntimeEnv        []config.EnvironmentVariable
}

// PrepareForDockerfile prepares plugin information for Dockerfile
func (s *PluginService) PrepareForDockerfile() []map[string]interface{} {
	var plugins []map[string]interface{}
	sharedInstallDir := s.ctx.Config.Plugins.InstallDir

	// Get base variables using the new variable system
	composer := s.ctx.GetVariableComposer().WithCommon().WithPlugin()
	baseVars := composer.Build()

	for _, plugin := range s.ctx.Config.Plugins.Items {
		// Add plugin-specific variables
		pluginVars := make(map[string]interface{})
		for k, v := range baseVars {
			pluginVars[k] = v
		}
		pluginVars[context.VarPluginName] = plugin.Name
		pluginVars[context.VarPluginDescription] = plugin.Description
		pluginVars[context.VarPluginDownloadURL] = s.resolveDownloadURL(plugin.DownloadURL)
		pluginVars[context.VarPluginInstallDir] = sharedInstallDir

		plugins = append(plugins, map[string]interface{}{
			"InstallCommand": core.SubstituteVariables(plugin.InstallCommand, pluginVars),
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

		// Generate URL resolver script
		urlResolverScript := s.GenerateURLResolverScript(plugin.DownloadURL)

		// Get download URL for display (static or placeholder)
		downloadURL := s.resolveDownloadURL(plugin.DownloadURL)

		plugins = append(plugins, PluginInfo{
			Name:              plugin.Name,
			Description:       plugin.Description,
			DownloadURL:       downloadURL,
			URLResolverScript: urlResolverScript,
			InstallDir:        sharedInstallDir,
			InstallCommand:    plugin.InstallCommand,
			RuntimeEnv:        processedEnv,
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

// resolveDownloadURL resolves the download URL based on configuration type
// For static URL: returns the URL directly
// For arch mapping: returns a placeholder that will be resolved at runtime
func (s *PluginService) resolveDownloadURL(urlConfig config.DownloadURLConfig) string {
	if urlConfig.IsStatic() {
		url, _ := urlConfig.GetStaticURL()
		return url
	}
	// For arch mapping, return a placeholder
	// The actual URL will be resolved in the shell script at runtime
	return "${PLUGIN_DOWNLOAD_URL}"
}

// GenerateURLResolverScript generates shell script to resolve download URL
// This is used in build_plugins.sh template
func (s *PluginService) GenerateURLResolverScript(urlConfig config.DownloadURLConfig) string {
	if urlConfig.IsStatic() {
		// For static URL, just echo the URL
		url, _ := urlConfig.GetStaticURL()
		return "PLUGIN_DOWNLOAD_URL=\"" + url + "\""
	}

	// For arch mapping, generate case statement
	urls, _ := urlConfig.GetArchURLs()
	script := `# Detect architecture and set download URL
ARCH=$(uname -m)
case "${ARCH}" in
`

	// Normalize architecture names and generate case branches
	archMap := s.normalizeArchMapping(urls)

	for arch, url := range archMap {
		if arch == "default" {
			continue // Handle default at the end
		}
		script += "  " + arch + ")\n"
		script += "    PLUGIN_DOWNLOAD_URL=\"" + url + "\"\n"
		script += "    ;;\n"
	}

	// Add default case
	if defaultURL, ok := urls["default"]; ok {
		script += "  *)\n"
		script += "    PLUGIN_DOWNLOAD_URL=\"" + defaultURL + "\"\n"
		script += "    ;;\n"
	} else {
		script += "  *)\n"
		script += "    echo \"ERROR: Unsupported architecture ${ARCH}\"\n"
		script += "    exit 1\n"
		script += "    ;;\n"
	}

	script += "esac"
	return script
}

// normalizeArchMapping normalizes architecture names to standard forms
// Maps common aliases to standard architecture names
func (s *PluginService) normalizeArchMapping(urls map[string]string) map[string]string {
	normalized := make(map[string]string)

	for arch, url := range urls {
		switch arch {
		case "x86_64", "amd64":
			// Combine x86_64 and amd64 into one case
			if existing, ok := normalized["x86_64|amd64"]; !ok || existing == "" {
				normalized["x86_64|amd64"] = url
			}
		case "aarch64", "arm64":
			// Combine aarch64 and arm64 into one case
			if existing, ok := normalized["aarch64|arm64"]; !ok || existing == "" {
				normalized["aarch64|arm64"] = url
			}
		case "default":
			normalized["default"] = url
		default:
			normalized[arch] = url
		}
	}

	return normalized
}
