package config

// ServiceConfig represents the complete service configuration
type ServiceConfig struct {
	Service  ServiceInfo    `yaml:"service"`
	Language LanguageConfig `yaml:"language"`
	Build    BuildConfig    `yaml:"build"`
	Plugins  []PluginConfig `yaml:"plugins,omitempty"`
	Runtime  RuntimeConfig  `yaml:"runtime"`
	LocalDev LocalDevConfig `yaml:"local_dev"`
	Makefile MakefileConfig `yaml:"makefile,omitempty"`
	Metadata MetadataConfig `yaml:"metadata"`
	CI       CIConfig       `yaml:"ci,omitempty"`
}

// ServiceInfo contains basic service information
type ServiceInfo struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Ports       []PortConfig `yaml:"ports"`
	DeployDir   string       `yaml:"deploy_dir"`
}

// PortConfig defines a service port
type PortConfig struct {
	Name        string `yaml:"name"`
	Port        int    `yaml:"port"`
	Protocol    string `yaml:"protocol"`
	Expose      bool   `yaml:"expose"`
	Description string `yaml:"description,omitempty"`
}

// LanguageConfig contains language-specific settings
type LanguageConfig struct {
	Type    string            `yaml:"type"`
	Version string            `yaml:"version"`
	Config  map[string]string `yaml:"config,omitempty"`
}

// BuildConfig contains build-related settings
type BuildConfig struct {
	DependencyFiles    DependencyFilesConfig         `yaml:"dependency_files"`
	BuilderImage       ArchImageConfig               `yaml:"builder_image"`
	RuntimeImage       ArchImageConfig               `yaml:"runtime_image"`
	SystemDependencies BuildSystemDependenciesConfig `yaml:"system_dependencies,omitempty"`
	Commands           BuildCommandsConfig           `yaml:"commands"`
}

// DependencyFilesConfig for dependency file detection
type DependencyFilesConfig struct {
	AutoDetect bool     `yaml:"auto_detect"`
	Files      []string `yaml:"files,omitempty"`
}

// ArchImageConfig for architecture-specific images
type ArchImageConfig struct {
	AMD64 string `yaml:"amd64"`
	ARM64 string `yaml:"arm64"`
}

// BuildSystemDependenciesConfig for build stage system packages
// Matches structure: build.system_dependencies.packages
type BuildSystemDependenciesConfig struct {
	Packages []string `yaml:"packages,omitempty"`
}

// RuntimeSystemDependenciesConfig for runtime stage system packages
// Matches structure: runtime.system_dependencies.packages
type RuntimeSystemDependenciesConfig struct {
	Packages []string `yaml:"packages,omitempty"`
}

// BuildCommandsConfig for build commands
type BuildCommandsConfig struct {
	PreBuild  string `yaml:"pre_build,omitempty"`
	Build     string `yaml:"build"`
	PostBuild string `yaml:"post_build,omitempty"`
}

// PluginConfig for plugin installation
type PluginConfig struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	DownloadURL    string `yaml:"download_url"`
	InstallDir     string `yaml:"install_dir"`
	InstallCommand string `yaml:"install_command"`
	Required       bool   `yaml:"required"`
	// 运行时环境变量配置
	RuntimeEnv []EnvironmentVariable `yaml:"runtime_env,omitempty"`
}

// EnvironmentVariable represents an environment variable
type EnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// RuntimeConfig contains runtime settings
type RuntimeConfig struct {
	SystemDependencies RuntimeSystemDependenciesConfig `yaml:"system_dependencies,omitempty"`
	Healthcheck        HealthcheckConfig               `yaml:"healthcheck"`
	Startup            StartupConfig                   `yaml:"startup"`
	// 控制是否生成运行时脚本的开关
	GenerateScripts bool `yaml:"generate_scripts,omitempty"`
}

// HealthcheckConfig for health check settings
type HealthcheckConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Type         string `yaml:"type"`                    // default | custom
	CustomScript string `yaml:"custom_script,omitempty"` // Required when type is 'custom'
}

// StartupConfig for startup settings
type StartupConfig struct {
	Command string      `yaml:"command"`
	Env     []EnvConfig `yaml:"env,omitempty"`
}

// EnvConfig for environment variables
type EnvConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// LocalDevConfig for local development settings
type LocalDevConfig struct {
	Compose    ComposeConfig    `yaml:"compose"`
	Kubernetes KubernetesConfig `yaml:"kubernetes"`
	// 支持的变量列表
	SupportedVariables []string `yaml:"supported_variables,omitempty"`
}

// ComposeConfig for Docker Compose settings
type ComposeConfig struct {
	Resources   ResourcesConfig     `yaml:"resources,omitempty"`
	Volumes     []VolumeConfig      `yaml:"volumes,omitempty"`
	Healthcheck ComposeHealthConfig `yaml:"healthcheck,omitempty"`
	Labels      map[string]string   `yaml:"labels,omitempty"`
}

// ResourcesConfig for resource limits
type ResourcesConfig struct {
	Limits       ResourceLimits `yaml:"limits,omitempty"`
	Reservations ResourceLimits `yaml:"reservations,omitempty"`
}

// ResourceLimits for CPU and memory limits
type ResourceLimits struct {
	CPUs   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// VolumeConfig for volume mounts
type VolumeConfig struct {
	Source      string `yaml:"source"`
	Target      string `yaml:"target"`
	Type        string `yaml:"type"` // bind | volume
	Description string `yaml:"description,omitempty"`
}

// ComposeHealthConfig for Docker Compose health check
type ComposeHealthConfig struct {
	Interval    string `yaml:"interval,omitempty"`
	Timeout     string `yaml:"timeout,omitempty"`
	Retries     int    `yaml:"retries,omitempty"`
	StartPeriod string `yaml:"start_period,omitempty"`
}

// KubernetesConfig for Kubernetes settings
type KubernetesConfig struct {
	Enabled    bool            `yaml:"enabled"`
	Namespace  string          `yaml:"namespace"`
	OutputDir  string          `yaml:"output_dir"`
	VolumeType string          `yaml:"volume_type"` // configMap | persistentVolumeClaim | emptyDir | hostPath
	ConfigMap  ConfigMapConfig `yaml:"configmap,omitempty"`
	Wait       WaitConfig      `yaml:"wait,omitempty"`
}

// ConfigMapConfig for ConfigMap settings
type ConfigMapConfig struct {
	AutoDetect bool   `yaml:"auto_detect"`
	Name       string `yaml:"name,omitempty"`
}

// WaitConfig for deployment wait settings
type WaitConfig struct {
	Enabled bool   `yaml:"enabled"`
	Timeout string `yaml:"timeout"`
}

// MakefileConfig for Makefile generation
type MakefileConfig struct {
	CustomTargets []CustomTarget `yaml:"custom_targets,omitempty"`
}

// CustomTarget for custom Make targets
type CustomTarget struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Commands    []string `yaml:"commands"`
}

// MetadataConfig for metadata
type MetadataConfig struct {
	TemplateVersion string `yaml:"template_version"`
	GeneratedAt     string `yaml:"generated_at,omitempty"`
	Generator       string `yaml:"generator"`
}

// CIConfig CI/CD 相关路径配置
type CIConfig struct {
	// CI 脚本目录路径（相对于项目根目录）
	// 默认: bk-ci/tcs
	ScriptDir string `yaml:"script_dir,omitempty"`

	// 构建配置目录（用于 K8s ConfigMap 等）
	// 默认: bk-ci/tcs/build
	BuildConfigDir string `yaml:"build_config_dir,omitempty"`
}
