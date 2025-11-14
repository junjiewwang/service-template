package context

// Variable category constants - 变量类别常量
const (
	CategoryCommon   = "common"   // Common variables shared by all generators
	CategoryBuild    = "build"    // Build-related variables
	CategoryRuntime  = "runtime"  // Runtime-related variables
	CategoryPlugin   = "plugin"   // Plugin-related variables
	CategoryCIPaths  = "ci-paths" // CI path variables
	CategoryService  = "service"  // Service-related variables
	CategoryLanguage = "language" // Language-related variables
)

// Path constants - 路径常量
const (
	// Plugin paths
	DefaultPluginRootDir = "/plugins" // 插件根目录

	// Container paths
	ContainerProjectRoot = "/opt" // 容器内项目根目录

	// CI paths patterns
	DefaultCIScriptDirPattern       = ".tad/build/%s"      // %s = service-name
	DefaultCIBuildConfigDirPattern  = "%s/build"           // %s = script_dir
	DefaultConfigTemplateDirPattern = "%s/config_template" // %s = script_dir
)

// Script file names - 脚本文件名常量
const (
	BuildScriptName        = "build.sh"
	DepsInstallScriptName  = "build_deps_install.sh"
	RtPrepareScriptName    = "rt_prepare.sh"
	EntrypointScriptName   = "entrypoint.sh"
	HealthcheckScriptName  = "healthchk.sh"
	BuildPluginsScriptName = "build_plugins.sh"
)

// Directory names - 目录名常量
const (
	ConfigDirName = "configs"
	BinDirName    = "bin"
	LogDirName    = "logs"
	DataDirName   = "data"
)

// Template variable keys - 模板变量键常量
const (
	// Service variables
	VarServiceName   = "SERVICE_NAME"
	VarServicePort   = "SERVICE_PORT"
	VarServiceRoot   = "SERVICE_ROOT"
	VarDeployDir     = "DEPLOY_DIR"
	VarConfigDir     = "CONFIG_DIR"
	VarServiceBinDir = "SERVICE_BIN_DIR"

	// Plugin variables
	VarPluginName        = "PLUGIN_NAME"
	VarPluginDescription = "PLUGIN_DESCRIPTION"
	VarPluginDownloadURL = "PLUGIN_DOWNLOAD_URL"
	VarPluginInstallDir  = "PLUGIN_INSTALL_DIR"
	VarPluginRootDir     = "PLUGIN_ROOT_DIR"

	// Build variables
	VarBuildCommand     = "BUILD_COMMAND"
	VarPreBuildCommand  = "PRE_BUILD_COMMAND"
	VarPostBuildCommand = "POST_BUILD_COMMAND"
	VarBuildOutputDir   = "BUILD_OUTPUT_DIR"
	VarProjectRoot      = "PROJECT_ROOT"

	// Language variables
	VarLanguage        = "LANGUAGE"
	VarLanguageVersion = "LANGUAGE_VERSION"

	// Architecture variables
	VarGOARCH = "GOARCH"
	VarGOOS   = "GOOS"

	// Metadata variables
	VarGeneratedAt = "GENERATED_AT"

	// CI path variables
	VarCIScriptDir         = "CI_SCRIPT_DIR"
	VarCIBuildConfigDir    = "CI_BUILD_CONFIG_DIR"
	VarCIConfigTemplateDir = "CI_CONFIG_TEMPLATE_DIR"
	VarCIContainerDir      = "CI_CONTAINER_DIR"
)

// GetDefaultPaths returns default path values
func GetDefaultPaths() map[string]string {
	return map[string]string{
		"plugin_root_dir":        DefaultPluginRootDir,
		"container_project_root": ContainerProjectRoot,
	}
}
