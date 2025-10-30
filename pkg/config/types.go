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
	DependencyFiles    DependencyFilesConfig    `yaml:"dependency_files"`
	BuilderImage       ArchImageConfig          `yaml:"builder_image"`
	RuntimeImage       ArchImageConfig          `yaml:"runtime_image"`
	SystemDependencies SystemDependenciesConfig `yaml:"system_dependencies,omitempty"`
	Commands           BuildCommandsConfig      `yaml:"commands"`
	OutputDir          string                   `yaml:"output_dir"`
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

// SystemDependenciesConfig for system packages
type SystemDependenciesConfig struct {
	Build   PackagesConfig `yaml:"build,omitempty"`
	Runtime PackagesConfig `yaml:"runtime,omitempty"`
}

// PackagesConfig for package lists
type PackagesConfig struct {
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
}

// RuntimeConfig contains runtime settings
type RuntimeConfig struct {
	SystemDependencies SystemDependenciesConfig `yaml:"system_dependencies,omitempty"`
	Healthcheck        HealthcheckConfig        `yaml:"healthcheck"`
	Startup            StartupConfig            `yaml:"startup"`
}

// HealthcheckConfig for health check settings
type HealthcheckConfig struct {
	Enabled      bool             `yaml:"enabled"`
	Type         string           `yaml:"type"` // http | tcp | exec | custom
	HTTP         HTTPHealthConfig `yaml:"http,omitempty"`
	CustomScript string           `yaml:"custom_script,omitempty"`
}

// HTTPHealthConfig for HTTP health checks
type HTTPHealthConfig struct {
	Path    string `yaml:"path"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
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
