package generator

import (
	_ "embed"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeBuildScript = "build-script"

// BuildScriptTemplateGenerator generates build.sh script
type BuildScriptTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the BuildScript generator
func init() {
	RegisterGenerator(GeneratorTypeBuildScript, createBuildScriptGenerator)
}

// createBuildScriptGenerator is the creator function for BuildScript generator
func createBuildScriptGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewBuildScriptTemplateGenerator(cfg, engine, vars), nil
}

// NewBuildScriptTemplateGenerator creates a new build script generator
func NewBuildScriptTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *BuildScriptTemplateGenerator {
	return &BuildScriptTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeBuildScript,
		},
	}
}

//go:embed templates/build.sh.tmpl
var buildScriptTemplate string

// Generate generates build.sh content
func (g *BuildScriptTemplateGenerator) Generate() (string, error) {
	// Prepare plugin install commands
	type PluginInfo struct {
		Name           string
		DownloadURL    string
		InstallDir     string
		InstallCommand string
		RuntimeEnv     []config.EnvironmentVariable // 新增：运行时环境变量
	}
	var plugins []PluginInfo
	// 获取共享的安装目录
	sharedInstallDir := g.config.Plugins.InstallDir

	for _, plugin := range g.config.Plugins.Items {
		// Process runtime environment variables - replace ${PLUGIN_INSTALL_DIR} with actual install dir
		processedEnv := make([]config.EnvironmentVariable, len(plugin.RuntimeEnv))
		for i, env := range plugin.RuntimeEnv {
			processedEnv[i] = config.EnvironmentVariable{
				Name: env.Name,
				Value: g.templateEngine.ReplaceVariables(env.Value, map[string]string{
					"PLUGIN_INSTALL_DIR": sharedInstallDir,
				}),
			}
		}

		plugins = append(plugins, PluginInfo{
			Name:           plugin.Name,
			DownloadURL:    plugin.DownloadURL,
			InstallDir:     sharedInstallDir,      // 使用共享的安装目录
			InstallCommand: plugin.InstallCommand, // Keep original command with variables
			RuntimeEnv:     processedEnv,          // 使用处理后的环境变量
		})
	}

	// Get CI paths
	ciPaths := NewCIPaths(g.config)

	vars := map[string]interface{}{
		"SERVICE_NAME":       g.config.Service.Name,
		"DEPLOY_DIR":         g.config.Service.DeployDir,
		"BUILD_COMMAND":      g.config.Build.Commands.Build,
		"PRE_BUILD_COMMAND":  g.config.Build.Commands.PreBuild,
		"POST_BUILD_COMMAND": g.config.Build.Commands.PostBuild,
		"PLUGINS":            plugins,
		"SERVICE_ROOT":       g.config.Service.DeployDir + "/" + g.config.Service.Name,
		"PLUGIN_ROOT_DIR":    "/plugins",                       // 新增：插件根目录
		"GENERATE_SCRIPTS":   g.config.Runtime.GenerateScripts, // 新增：是否生成运行时脚本
		"CI_SCRIPT_DIR":      ciPaths.ScriptDir,                // CI 脚本目录（主机相对路径）
		"CI_CONTAINER_DIR":   ciPaths.ContainerScriptDir,       // CI 脚本目录（容器绝对路径）
	}

	return g.RenderTemplate(g.getTemplate(), vars)
}

// getTemplate returns the build.sh template
func (g *BuildScriptTemplateGenerator) getTemplate() string {
	return buildScriptTemplate
}
