package config

import (
	"fmt"
	"strings"
)

// ServiceConfig represents the complete service configuration
type ServiceConfig struct {
	// 基础镜像配置（顶层，与 service 同级）
	BaseImages BaseImagesConfig `yaml:"base_images"`

	Service  ServiceInfo    `yaml:"service"`
	Language LanguageConfig `yaml:"language"`
	Build    BuildConfig    `yaml:"build"`
	Plugins  PluginsConfig  `yaml:"plugins,omitempty"`
	Runtime  RuntimeConfig  `yaml:"runtime"`
	LocalDev LocalDevConfig `yaml:"local_dev"`
	Makefile MakefileConfig `yaml:"makefile,omitempty"`
	Metadata MetadataConfig `yaml:"metadata"`
	CI       CIConfig       `yaml:"ci,omitempty"`
}

// NewImageResolver 创建镜像解析器
func (s *ServiceConfig) NewImageResolver() *ImageResolver {
	return NewImageResolver(s)
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
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config,omitempty"`
}

// GetString gets a string value from config with default
func (l *LanguageConfig) GetString(key string, defaultValue string) string {
	if l.Config == nil {
		return defaultValue
	}
	if val, ok := l.Config[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

// GetInt gets an int value from config with default
func (l *LanguageConfig) GetInt(key string, defaultValue int) int {
	if l.Config == nil {
		return defaultValue
	}
	if val, ok := l.Config[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return defaultValue
}

// GetBool gets a bool value from config with default
func (l *LanguageConfig) GetBool(key string, defaultValue bool) bool {
	if l.Config == nil {
		return defaultValue
	}
	if val, ok := l.Config[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

// GetStringSlice gets a string slice from config
func (l *LanguageConfig) GetStringSlice(key string) []string {
	if l.Config == nil {
		return nil
	}
	if val, ok := l.Config[key]; ok {
		switch v := val.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

// BuildConfig contains build-related settings
type BuildConfig struct {
	DependencyFiles DependencyFilesConfig `yaml:"dependency_files"`
	// 使用 ImageRef 类型强制引用
	BuilderImage ImageRef            `yaml:"builder_image"`
	RuntimeImage ImageRef            `yaml:"runtime_image"`
	Dependencies DependenciesConfig  `yaml:"dependencies"`
	Commands     BuildCommandsConfig `yaml:"commands"`
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

// Validate 验证架构镜像配置
func (a *ArchImageConfig) Validate() error {
	if a.AMD64 == "" {
		return fmt.Errorf("amd64 image is required")
	}
	if a.ARM64 == "" {
		return fmt.Errorf("arm64 image is required")
	}
	return nil
}

// GetByArch 根据架构获取镜像
func (a *ArchImageConfig) GetByArch(arch string) (string, error) {
	arch = normalizeArch(arch)

	switch arch {
	case "amd64":
		if a.AMD64 == "" {
			return "", fmt.Errorf("amd64 image not configured")
		}
		return a.AMD64, nil
	case "arm64":
		if a.ARM64 == "" {
			return "", fmt.Errorf("arm64 image not configured")
		}
		return a.ARM64, nil
	default:
		return "", fmt.Errorf("unsupported architecture: %s (supported: amd64, arm64)", arch)
	}
}

// SupportedArchs 返回支持的架构列表
func (a *ArchImageConfig) SupportedArchs() []string {
	archs := make([]string, 0, 2)
	if a.AMD64 != "" {
		archs = append(archs, "amd64")
	}
	if a.ARM64 != "" {
		archs = append(archs, "arm64")
	}
	return archs
}

// normalizeArch 标准化架构名称
func normalizeArch(arch string) string {
	arch = strings.ToLower(strings.TrimSpace(arch))
	switch arch {
	case "x86_64", "x64":
		return "amd64"
	case "aarch64":
		return "arm64"
	default:
		return arch
	}
}

// DependenciesConfig for build stage dependencies
// Matches structure: build.dependencies
type DependenciesConfig struct {
	SystemPkgs []string        `yaml:"system_pkgs,omitempty"`
	CustomPkgs []CustomPackage `yaml:"custom_pkgs,omitempty"`
}

// CustomPackage for custom package installation
type CustomPackage struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description,omitempty"`
	InstallCommand string `yaml:"install_command"`
	Required       bool   `yaml:"required"`
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

// PluginsConfig for plugins configuration
type PluginsConfig struct {
	// 所有插件共用的安装目录
	InstallDir string `yaml:"install_dir"`
	// 插件列表
	Items []PluginConfig `yaml:"items,omitempty"`
}

// PluginConfig for plugin installation
type PluginConfig struct {
	Name           string            `yaml:"name"`
	Description    string            `yaml:"description"`
	DownloadURL    DownloadURLConfig `yaml:"download_url"`
	InstallCommand string            `yaml:"install_command"`
	Required       bool              `yaml:"required"`
	// 运行时环境变量配置
	RuntimeEnv []EnvironmentVariable `yaml:"runtime_env,omitempty"`
}

// DownloadURLConfig supports both static string URL and architecture-specific URL mapping
type DownloadURLConfig struct {
	value interface{} // string or map[string]string
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
}

// ComposeConfig for Docker Compose settings
type ComposeConfig struct {
	Resources   ResourcesConfig     `yaml:"resources,omitempty"`
	Volumes     []VolumeConfig      `yaml:"volumes,omitempty"`
	Environment []EnvConfig         `yaml:"environment,omitempty"` // Compose-specific environment variables
	Entrypoint  []string            `yaml:"entrypoint,omitempty"`  // Override container entrypoint
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
	Enabled    bool       `yaml:"enabled"`
	Namespace  string     `yaml:"namespace"`
	OutputDir  string     `yaml:"output_dir"`
	VolumeType string     `yaml:"volume_type"` // configMap | persistentVolumeClaim | emptyDir | hostPath
	Wait       WaitConfig `yaml:"wait,omitempty"`
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
	// 默认: .tad/build/{service-name}
	ScriptDir string `yaml:"script_dir,omitempty"`

	// 构建配置目录（用于 K8s ConfigMap 等）
	// 默认: {script_dir}/build
	BuildConfigDir string `yaml:"build_config_dir,omitempty"`

	// 配置模板目录（用于用户自定义配置模板）
	// 默认: {script_dir}/config_template
	ConfigTemplateDir string `yaml:"config_template_dir,omitempty"`
}

// DownloadURLConfig methods

// UnmarshalYAML implements custom YAML unmarshaling for DownloadURLConfig
func (d *DownloadURLConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to unmarshal as string
	var str string
	if err := unmarshal(&str); err == nil {
		if str == "" {
			return fmt.Errorf("download_url cannot be empty string")
		}
		d.value = str
		return nil
	}

	// Try to unmarshal as map
	var m map[string]string
	if err := unmarshal(&m); err == nil {
		if len(m) == 0 {
			return fmt.Errorf("download_url map cannot be empty")
		}
		d.value = m
		return nil
	}

	return fmt.Errorf("download_url must be either a string or a map[string]string")
}

// MarshalYAML implements custom YAML marshaling for DownloadURLConfig
func (d DownloadURLConfig) MarshalYAML() (interface{}, error) {
	return d.value, nil
}

// IsStatic returns true if the download URL is a static string
func (d *DownloadURLConfig) IsStatic() bool {
	_, ok := d.value.(string)
	return ok
}

// IsArchMapping returns true if the download URL is an architecture mapping
func (d *DownloadURLConfig) IsArchMapping() bool {
	_, ok := d.value.(map[string]string)
	return ok
}

// GetStaticURL returns the static URL string
func (d *DownloadURLConfig) GetStaticURL() (string, error) {
	if str, ok := d.value.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("download_url is not a static string")
}

// GetArchURLs returns the architecture-specific URL mapping
func (d *DownloadURLConfig) GetArchURLs() (map[string]string, error) {
	if m, ok := d.value.(map[string]string); ok {
		return m, nil
	}
	return nil, fmt.Errorf("download_url is not an architecture mapping")
}

// IsEmpty returns true if the download URL is not configured
func (d *DownloadURLConfig) IsEmpty() bool {
	return d.value == nil
}

// String returns a string representation for logging
func (d *DownloadURLConfig) String() string {
	if d.IsStatic() {
		url, _ := d.GetStaticURL()
		return url
	}
	if d.IsArchMapping() {
		urls, _ := d.GetArchURLs()
		return fmt.Sprintf("arch_mapping(%d archs)", len(urls))
	}
	return "<empty>"
}

// NewStaticDownloadURL creates a DownloadURLConfig with a static URL
func NewStaticDownloadURL(url string) DownloadURLConfig {
	return DownloadURLConfig{value: url}
}

// NewArchMappingDownloadURL creates a DownloadURLConfig with architecture mapping
func NewArchMappingDownloadURL(urls map[string]string) DownloadURLConfig {
	return DownloadURLConfig{value: urls}
}
